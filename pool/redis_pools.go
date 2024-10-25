package pool

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
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
	return nil
}

func (r RedisRecommendPool) Remove(elements ...*Element) error {
	return nil
}

func (r RedisRecommendPool) MemberCount() (int64, error) {
	return 0, nil
}

func (r RedisRecommendPool) FetchMembers(uId int64, strategies []*Strategy) ([]*Element, error) {
	return nil, nil
}

func NewRedisRecommendPool(client *redis.Client, logger *logrus.Entry, key, allowRepetition bool) RecommendPool {
	redisPool := &RedisRecommendPool{
		client:          client,
		logger:          logger.WithField("module", "redis-recommend-pool"),
		allowRepetition: allowRepetition,
	}
	return redisPool
}
