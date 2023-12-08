package mq

import (
	"context"
	"encoding/json"
	"github.com/THK-IM/THK-IM-Base-Server/conf"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type KafkaPublisher struct {
	writer *kafka.Writer
	logger *logrus.Entry
}

func (k KafkaPublisher) Pub(id string, msg map[string]interface{}) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	kMsg := kafka.Message{
		Key:   []byte(id),
		Value: value,
	}
	ctx := context.Background()
	err = k.writer.WriteMessages(ctx, kMsg)
	return err
}

func NewKafkaPublisher(config *conf.Publisher, clientId string, logger *logrus.Entry) Publisher {
	brokers := strings.Split(config.KafkaPublisher.Brokers, ",")
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        config.Topic,
		Balancer:     &kafka.RoundRobin{},
		RequiredAcks: kafka.RequiredAcks(config.KafkaPublisher.RequireAck),
		BatchSize:    config.KafkaPublisher.BatchSize,
		BatchTimeout: 50 * time.Millisecond,
		Async:        config.KafkaPublisher.Async,
	}
	return &KafkaPublisher{
		writer: writer,
		logger: logger.WithField("search_index", "kafka_publisher"),
	}
}
