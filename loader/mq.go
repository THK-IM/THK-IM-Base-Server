package loader

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/mq"
)

func LoadPublishers(pubConfigs []*conf.Publisher, nodeId int64, logger *logrus.Entry) map[string]mq.Publisher {
	clientId := fmt.Sprintf("%d", nodeId)
	publisherMap := make(map[string]mq.Publisher, 0)
	for _, pubConfig := range pubConfigs {
		if pubConfig.RedisPublisher != nil {
			client := LoadRedis(pubConfig.RedisPublisher.RedisSource)
			publisherMap[pubConfig.Topic] = mq.NewRedisPublisher(pubConfig, clientId, logger, client)
		} else if pubConfig.KafkaPublisher != nil {
			publisherMap[pubConfig.Topic] = mq.NewKafkaPublisher(pubConfig, clientId, logger)
		}
	}
	return publisherMap
}

func LoadSubscribers(subConfigs []*conf.Subscriber, nodeId int64, logger *logrus.Entry) map[string]mq.Subscriber {
	clientId := fmt.Sprintf("%d", nodeId)
	subscriberMap := make(map[string]mq.Subscriber, 0)
	for _, subConfig := range subConfigs {
		if subConfig.RedisSubscriber != nil {
			client := LoadRedis(subConfig.RedisSubscriber.RedisSource)
			subscriberMap[subConfig.Topic] = mq.NewRedisSubscribe(subConfig, clientId, logger, client)
		} else if subConfig.KafkaSubscriber != nil {
			subscriberMap[subConfig.Topic] = mq.NewKafkaSubscriber(subConfig, clientId, logger)
		}
	}
	return subscriberMap
}
