package pool

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/filter"
	"golang.org/x/exp/slices"
	"strconv"
	"time"
)

type (
	RedisRecommendPool struct {
		key       string // 推荐池key
		client    *redis.Client
		logger    *logrus.Entry
		bitFilter filter.BitFilter
	}
)

func (r RedisRecommendPool) Add(elements ...*Element) error {
	members := make([]redis.Z, 0)
	for _, element := range elements {
		member := redis.Z{
			Score:  element.Score,
			Member: element.Id,
		}
		members = append(members, member)
	}
	_, err := r.client.ZAdd(context.Background(), r.key, members...).Result()
	return err
}

func (r RedisRecommendPool) Remove(elements ...*Element) error {
	args := make([]interface{}, 0)
	for _, element := range elements {
		args = append(args, element.Id)
	}
	_, err := r.client.ZRem(context.Background(), r.key, args...).Result()
	return err
}

func (r RedisRecommendPool) RemoveByScore(score float64) (int64, error) {
	return r.client.ZRemRangeByScore(
		context.Background(), r.key,
		"-inf",
		fmt.Sprintf("(%f", score),
	).Result()
}

func (r RedisRecommendPool) ElementCountByRange(score1, score2 float64) (int64, error) {
	return r.client.ZCount(
		context.Background(), r.key,
		fmt.Sprintf("%f", score1),
		fmt.Sprintf("(%f", score2),
	).Result()
}

func (r RedisRecommendPool) ElementCount() (int64, error) {
	return r.client.ZCard(context.Background(), r.key).Result()
}

func (r RedisRecommendPool) FetchElements(uId int64, strategies []*RecommendStrategy) ([]Element, error) {
	elements := make([]Element, 0)
	for _, strategy := range strategies {
		if strategy.Type == StrategyRandom {
			remainCount := int64(strategy.Count)
			for fetchTimes := 0; fetchTimes < strategy.RepeatRetryTimes+1; fetchTimes++ {
				subElements, errFetch := r.fetchMemberByRandom(5 * remainCount)
				if errFetch != nil {
					return nil, errFetch
				}
				if len(subElements) == 0 {
					break
				}
				remainCount, errFetch = r.checkElements(&elements, subElements, uId, remainCount)
				if errFetch != nil {
					return nil, errFetch
				}
				if remainCount <= 0 {
					break
				}
			}
		} else if strategy.Type == StrategyScore {
			offset := int64(0)
			remainCount := int64(strategy.Count)
			for fetchTimes := 0; fetchTimes < strategy.RepeatRetryTimes+1; fetchTimes++ {
				subElements, errFetch := r.fetchMembersByScore(5*remainCount, offset)
				offset += 5 * remainCount
				if errFetch != nil {
					return nil, errFetch
				}
				if len(subElements) == 0 {
					break
				}
				remainCount, errFetch = r.checkElements(&elements, subElements, uId, remainCount)
				if errFetch != nil {
					return nil, errFetch
				}
				if remainCount <= 0 {
					break
				}
			}
		}
	}
	return elements, nil
}

func (r RedisRecommendPool) checkElements(elements *[]Element, fetchElements []Element, uId, count int64) (int64, error) {
	added := int64(0)
	ids := make([]uint32, 0)
	tmpElements := make([]Element, 0)
	for _, element := range fetchElements {
		if !slices.Contains(*elements, element) { // 去重
			ids = append(ids, uint32(element.Id))
			tmpElements = append(tmpElements, element)
		}
	}
	userKey := r.userRecordKey(uId)
	contains, err := r.bitFilter.Contains(userKey, ids...)
	if err != nil {
		return count, err
	}
	for index, c := range contains {
		if c == false {
			*elements = append(*elements, tmpElements[index])
			added++
			if added >= count {
				break
			}
		}
	}
	return count - added, nil
}

func (r RedisRecommendPool) fetchMemberByRandom(count int64) ([]Element, error) {
	elements := make([]Element, 0)
	results, err := r.client.ZRandMemberWithScores(context.Background(), r.key, int(count)).Result()
	if err != nil {
		return elements, err
	}
	for _, result := range results {
		id, errId := strconv.ParseInt(result.Member.(string), 10, 64)
		if errId != nil {
			return elements, errId
		}
		element := Element{
			Id:    id,
			Score: result.Score,
		}
		elements = append(elements, element)
	}
	return elements, nil
}

func (r RedisRecommendPool) fetchMembersByScore(count int64, offset int64) ([]Element, error) {
	elements := make([]Element, 0)
	opt := &redis.ZRangeBy{
		Min:    "-inf",
		Max:    "+inf",
		Count:  count,
		Offset: offset,
	}
	results, err := r.client.ZRangeByScoreWithScores(context.Background(), r.key, opt).Result()
	if err != nil {
		return elements, err
	}
	for _, result := range results {
		id, errId := strconv.ParseInt(result.Member.(string), 10, 64)
		if errId != nil {
			return elements, errId
		}
		element := Element{
			Id:    id,
			Score: result.Score,
		}
		elements = append(elements, element)
	}
	return elements, nil
}

func (r RedisRecommendPool) userRecordKey(uId int64) string {
	userKey := fmt.Sprintf("%s:%d", r.key, uId)
	return userKey
}

func (r RedisRecommendPool) ClearUserRecord(uId int64) error {
	userKey := r.userRecordKey(uId)
	return r.bitFilter.Delete(userKey)
}

func (r RedisRecommendPool) AddUserRecord(uId int64, ex time.Duration, elementIds ...int64) error {
	userKey := r.userRecordKey(uId)
	poss := make([]uint32, len(elementIds))
	for _, id := range elementIds {
		poss = append(poss, uint32(id))
	}
	return r.bitFilter.AddPos(userKey, ex, poss...)
}

func NewRedisRecommendPool(client *redis.Client, logger *logrus.Entry, key string, maxRecordCount uint32) RecommendPool {
	redisPool := &RedisRecommendPool{
		key:       key,
		client:    client,
		logger:    logger.WithField("module", "redis-recommend-pool"),
		bitFilter: filter.NewRedisBitFilter(client, logger, maxRecordCount),
	}
	return redisPool
}
