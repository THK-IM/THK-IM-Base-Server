package pool

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"math"
)

type (
	RedisRecommendPool struct {
		client          *redis.Client
		logger          *logrus.Entry
		key             string // 推荐池key
		allowRepetition bool   // 允许重复
		eliminateScore  int64  // 淘汰分数
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

func (r RedisRecommendPool) MemberCountByRange(score1, score2 float64) (int64, error) {
	return r.client.ZCount(
		context.Background(), r.key,
		fmt.Sprintf("%f", score1),
		fmt.Sprintf("%f", score2),
	).Result()
}

func (r RedisRecommendPool) MemberCount() (int64, error) {
	return r.client.ZCount(
		context.Background(), r.key,
		fmt.Sprintf("%f", 0-math.MaxFloat64),
		fmt.Sprintf("%f", math.MaxFloat64),
	).Result()
}

func (r RedisRecommendPool) FetchMembers(uId int64, strategies []*Strategy) ([]*Element, error) {
	elements := make([]*Element, 0)
	for _, strategy := range strategies {
		if subElements, errFetch := r.fetchMembersByStrategy(strategy); errFetch != nil {
			return nil, errFetch
		} else {
			elements = append(elements, subElements...)
		}
	}
	return elements, nil
}

func (r RedisRecommendPool) fetchMembersByStrategy(strategy *Strategy) ([]*Element, error) {
	elements := make([]*Element, 0)
	if strategy == nil {
		return elements, nil
	}
	if strategy.Type == StrategyRandom {
		return nil, nil
	} else if strategy.Type == StrategyScoreAsc {
		return nil, nil
	} else if strategy.Type == StrategyScoreDesc {
		return nil, nil
	}
	return elements, nil
}

func NewRedisRecommendPool(client *redis.Client, logger *logrus.Entry, key, allowRepetition bool) RecommendPool {
	redisPool := &RedisRecommendPool{
		client:          client,
		logger:          logger.WithField("module", "redis-recommend-pool"),
		allowRepetition: allowRepetition,
	}
	return redisPool
}
