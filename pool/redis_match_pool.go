package pool

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type (
	RedisMatchPool struct {
		key    string // 推荐池key
		client *redis.Client
		logger *logrus.Entry
	}
)

func (r RedisMatchPool) Clear() error {
	return r.client.Del(context.Background(), r.key).Err()
}

func (r RedisMatchPool) Add(ids ...string) (int64, error) {
	return r.client.SAdd(context.Background(), r.key, ids).Result()
}

func (r RedisMatchPool) Remove(ids ...string) (int64, error) {
	return r.client.SRem(context.Background(), r.key, ids).Result()
}

func (r RedisMatchPool) Contain(id string) (bool, error) {
	return r.client.SIsMember(context.Background(), r.key, id).Result()
}

func (r RedisMatchPool) Count() (int64, error) {
	return r.client.SCard(context.Background(), r.key).Result()
}

func (r RedisMatchPool) FetchOne() (*string, error) {
	id, err := r.client.SPop(context.Background(), r.key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
	}
	return &id, err
}

func (r RedisMatchPool) Match(forId string, retryTimes int, f MatchFunction) (matchedId *string, err error) {
	putBack := false
	passedIds := make([]string, 0)
	times := 0
	for matchedId == nil {
		id, errPop := r.client.SPop(context.Background(), r.key).Result()
		if errPop != nil {
			if errors.Is(errPop, redis.Nil) {
				return nil, nil
			} else {
				return nil, errPop
			}
		}
		matchedId, putBack, err = f(forId, id)
		if err != nil {
			return nil, err
		}
		if putBack {
			passedIds = append(passedIds, id)
			_, errBack := r.Add(id)
			if errBack != nil {
				return matchedId, errBack
			}
		}
		if times >= retryTimes {
			break
		}
		times++
	}
	return matchedId, err
}

func NewRedisMatchPool(client *redis.Client, logger *logrus.Entry, key string) MatchPool {
	return RedisMatchPool{
		key:    key,
		client: client,
		logger: logger.WithField("module", "redis-match-pool"),
	}
}
