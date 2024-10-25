package filter

import (
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type (
	Factory interface {
		NewBloomFilter(params *BloomFilterParams) BloomFilter
		NewBitFilter(maxBit uint32) BitFilter
	}

	RedisFilterFactory struct {
		client *redis.Client
		logger *logrus.Entry
	}
)

func NewRedisFactory(client *redis.Client, logger *logrus.Entry) Factory {
	return &RedisFilterFactory{
		client: client,
		logger: logger,
	}
}
