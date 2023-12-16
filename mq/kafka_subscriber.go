package mq

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"net"
	"strconv"
	"strings"
	"time"
)

type kafkaSubscribe struct {
	reader            *kafka.Reader
	logger            *logrus.Entry
	numPartitions     int
	replicationFactor int
	groupId           *string
	topic             string
	brokers           []string
}

func (k *kafkaSubscribe) Sub(onReceived OnMessageReceived) {
	k.createTopicIfNeed()
	if k.groupId == nil {
		go k.subscribe(onReceived)
	} else {
		go k.subscribeGroup(onReceived)
	}
}

func (k *kafkaSubscribe) createTopicIfNeed() {
	if len(k.brokers) == 0 {
		panic(errors.New("brokers length is 0"))
	}
	conn, err := kafka.Dial("tcp", k.brokers[0])
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	partitions, e := conn.ReadPartitions()
	if e != nil {
		panic(e)
	}

	existed := false
	for _, p := range partitions {
		if p.Topic == k.topic {
			existed = true
			break
		}
	}

	if !existed {
		var controller kafka.Broker
		controller, err = conn.Controller()
		if err != nil {
			panic(err.Error())
		}
		var connLeader *kafka.Conn
		connLeader, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
		if err != nil {
			panic(err.Error())
		}
		defer func() {
			_ = connLeader.Close()
		}()

		topicConfigs := []kafka.TopicConfig{
			{
				Topic:             k.topic,
				NumPartitions:     k.numPartitions,
				ReplicationFactor: k.replicationFactor,
			},
		}

		err = connLeader.CreateTopics(topicConfigs...)
		if err != nil {
			panic(err.Error())
		}
	}
}

func (k *kafkaSubscribe) subscribe(received OnMessageReceived) {
	if e := k.reader.SetOffsetAt(context.Background(), time.Now()); e != nil {
		panic(e)
	}
	for {
		m, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			k.logger.Errorf("subscribe %v", err)
			break
		}
		msg := make(map[string]interface{})
		if err = json.Unmarshal(m.Value, &msg); err == nil {
			if err = received(msg); err != nil {
				k.logger.Errorf("subscribe %v", err)
			}
		} else {
			k.logger.Errorf("subscribe %v", err)
		}
	}

	if err := k.reader.Close(); err != nil {
		k.logger.Error(err)
	}
}

func (k *kafkaSubscribe) subscribeGroup(received OnMessageReceived) {
	ctx := context.Background()
	for {
		m, err := k.reader.FetchMessage(ctx)
		if err != nil {
			break
		}
		msg := make(map[string]interface{})
		if err = json.Unmarshal(m.Value, &msg); err == nil {
			if err = received(msg); err != nil {
				k.logger.Errorf("subscribeGroup %v, %v", msg, err)
			} else {
				if err = k.reader.CommitMessages(ctx, m); err != nil {
					k.logger.Errorf("subscribeGroup %v, %v", m, err)
				}
			}
		} else {
			k.logger.Errorf("subscribeGroup %v, %v", m.Value, err)
		}
	}
}

func NewKafkaSubscriber(config *conf.Subscriber, nodeId string, logger *logrus.Entry) Subscriber {
	brokers := strings.Split(config.KafkaSubscriber.Brokers, ",")
	if config.Group == nil {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:       brokers,
			Topic:         config.Topic,
			Partition:     config.KafkaSubscriber.Partition,
			QueueCapacity: 1,
			MaxWait:       50 * time.Millisecond,
			MinBytes:      1,    // 1KB
			MaxBytes:      10e6, // 10MB
		})
		return &kafkaSubscribe{
			topic:             config.Topic,
			brokers:           brokers,
			numPartitions:     config.KafkaSubscriber.NumPartitions,
			replicationFactor: config.KafkaSubscriber.ReplicationFactor,
			reader:            r,
			logger:            logger.WithField("search_index", "kafka_subscribe"),
		}
	} else {
		r := kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          config.Topic,
			GroupID:        *config.Group,
			CommitInterval: time.Second,
			StartOffset:    kafka.LastOffset,
			QueueCapacity:  1,
			MaxWait:        50 * time.Millisecond,
			MinBytes:       1,    // 1KB
			MaxBytes:       10e6, // 10MB
		})
		return &kafkaSubscribe{
			topic:             config.Topic,
			brokers:           brokers,
			numPartitions:     config.KafkaSubscriber.NumPartitions,
			replicationFactor: config.KafkaSubscriber.ReplicationFactor,
			groupId:           config.Group,
			reader:            r,
			logger:            logger.WithField("search_index", "kafka_group_subscribe"),
		}
	}
}
