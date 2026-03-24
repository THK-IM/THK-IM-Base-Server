package loader

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoLogger struct {
	logger *logrus.Entry
}

func (l *MongoLogger) CommandStarted(_ context.Context, evt *event.CommandStartedEvent) {
	l.logger.Tracef("mongo start: %s %s", evt.CommandName, evt.Command)
}

func (l *MongoLogger) CommandSucceeded(_ context.Context, evt *event.CommandSucceededEvent) {
	l.logger.Tracef("mongo success: %s duration=%dms", evt.CommandName, evt.Duration.Microseconds())
}

func (l *MongoLogger) CommandFailed(_ context.Context, evt *event.CommandFailedEvent) {
	l.logger.Errorf("mongo failed: %s duration=%dms err=%v", evt.CommandName, evt.Duration.Microseconds(), evt.Failure)
}

func LoadMongo(entry *logrus.Entry, source *conf.MongoSource) *mongo.Client {
	if source == nil {
		return nil
	}

	mLogger := &MongoLogger{
		logger: entry,
	}

	monitor := &event.CommandMonitor{
		Started:   mLogger.CommandStarted,
		Succeeded: mLogger.CommandSucceeded,
		Failed:    mLogger.CommandFailed,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().
		ApplyURI(fmt.Sprintf("%s%s", source.Endpoint, source.Uri)).
		SetMaxPoolSize(uint64(source.MaxOpenConn)).
		SetMinPoolSize(uint64(source.MaxIdleConn)).
		SetMaxConnIdleTime(time.Duration(source.ConnMaxIdleTime) * time.Second).
		SetMonitor(monitor)

	client, err := mongo.Connect(clientOpts)
	if err != nil {
		panic(err)
	}

	// ping 检查连接
	if errPing := client.Ping(ctx, nil); errPing != nil {
		panic(errPing)
	}

	entry.Info("MongoDB connected")

	return client
}
