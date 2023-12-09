package mq

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
)

type (
	RedisPublisher struct {
		stream   string
		clientId string
		maxLen   int64
		client   *redis.Client
		logger   *logrus.Entry
	}
)

func (d RedisPublisher) Pub(id string, msg map[string]interface{}) error {
	ctx := context.Background()
	err := d.client.XAdd(ctx, &redis.XAddArgs{
		Stream: d.stream,
		MaxLen: d.maxLen,
		ID:     "",
		Values: msg,
	}).Err()
	return err
}

func NewRedisPublisher(config *conf.Publisher, clientId string, logger *logrus.Entry, client *redis.Client) Publisher {
	return RedisPublisher{
		maxLen:   config.RedisPublisher.MaxQueueLen,
		stream:   config.Topic,
		clientId: clientId,
		logger:   logger,
		client:   client,
	}
}
