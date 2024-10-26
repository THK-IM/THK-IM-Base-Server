package pool

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Factory interface {

	// NewMatchPool(key string) MatchPool

	// NewRecommendPool 新建推荐池
	// @param key 池子key
	// @param userMaxRecordCount 用户推荐记录最大保存数量 如设置 8*1024 每个用户占用2KB
	NewRecommendPool(key string, userMaxRecordCount uint32) RecommendPool
}

type RedisPoolFactory struct {
	client *redis.Client
	logger *logrus.Entry
}

func (r RedisPoolFactory) NewRecommendPool(key string, userMaxRecordCount uint32) RecommendPool {
	return NewRedisRecommendPool(r.client, r.logger, key, userMaxRecordCount)
}

func NewRedisPoolFactory(client *redis.Client, logger *logrus.Entry) Factory {
	return &RedisPoolFactory{
		client: client,
		logger: logger,
	}
}
