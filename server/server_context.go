package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/loader"
	"github.com/thk-im/thk-im-base-server/locker"
	"github.com/thk-im/thk-im-base-server/metric"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-base-server/mq"
	"github.com/thk-im/thk-im-base-server/object"
	"github.com/thk-im/thk-im-base-server/snowflake"
	"github.com/thk-im/thk-im-base-server/websocket"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Context struct {
	startTime       int64
	nodeId          int64
	metricService   *metric.Service
	config          *conf.Config
	logger          *logrus.Entry
	redisCache      *redis.Client
	lockerFactory   locker.Factory
	database        *gorm.DB
	snowflakeNode   *snowflake.Node
	httpEngine      *gin.Engine
	objectStorage   object.Storage
	websocketServer websocket.Server
	publisherMap    map[string]mq.Publisher
	subscriberMap   map[string]mq.Subscriber
	SdkMap          map[string]interface{}
	ModelMap        map[string]interface{}
}

func (app *Context) SupportLanguage() []language.Tag {
	return []language.Tag{
		language.Chinese,
		language.SimplifiedChinese,
		language.TraditionalChinese,
	}
}

func (app *Context) StartTime() int64 {
	return app.startTime
}

func (app *Context) NodeId() int64 {
	return app.nodeId
}

func (app *Context) Config() *conf.Config {
	return app.config
}

func (app *Context) RedisCache() *redis.Client {
	return app.redisCache
}

func (app *Context) Database() *gorm.DB {
	return app.database
}

func (app *Context) SnowflakeNode() *snowflake.Node {
	return app.snowflakeNode
}

func (app *Context) HttpEngine() *gin.Engine {
	return app.httpEngine
}

func (app *Context) Logger() *logrus.Entry {
	return app.logger
}

func (app *Context) NewLocker(key string, waitMs int, timeoutMs int) locker.Locker {
	return app.lockerFactory.NewLocker(key, waitMs, timeoutMs)
}

func (app *Context) MetricService() *metric.Service {
	return app.metricService
}

func (app *Context) WebsocketServer() websocket.Server {
	return app.websocketServer
}

func (app *Context) ObjectStorage() object.Storage {
	return app.objectStorage
}

func (app *Context) ServerEventPublisher() mq.Publisher {
	return app.publisherMap["server_event"]
}

func (app *Context) MsgPusherPublisher() mq.Publisher {
	return app.publisherMap["push_msg"]
}

func (app *Context) MsgSaverPublisher() mq.Publisher {
	return app.publisherMap["save_msg"]
}

func (app *Context) MsgPusherSubscriber() mq.Subscriber {
	return app.subscriberMap["push_msg"]
}

func (app *Context) MsgSaverSubscriber() mq.Subscriber {
	return app.subscriberMap["save_msg"]
}

func (app *Context) ServerEventSubscriber() mq.Subscriber {
	return app.subscriberMap["server_event"]
}

func (app *Context) Init(config *conf.Config) {
	logger := loader.LoadLogger(config.Name, config.Logger)
	redisCache := loader.LoadRedis(config.RedisSource)
	nodeId, startTime := loader.LoadNodeId(config, redisCache)
	snowflakeNode, err := snowflake.NewNode(nodeId)
	if err != nil {
		panic(err)
	}
	gin.SetMode(config.Mode)
	httpEngine := gin.Default()
	claimsMiddleware := middleware.Claims()
	httpEngine.Use(claimsMiddleware)
	app.httpEngine = httpEngine
	app.config = config
	app.logger = logger
	app.redisCache = redisCache
	app.nodeId = nodeId
	app.startTime = startTime
	app.snowflakeNode = snowflakeNode

	if config.MysqlSource != nil {
		app.database = loader.LoadMysql(logger, config.MysqlSource)
	}

	if config.MsgQueue.Publishers != nil {
		app.publisherMap = loader.LoadPublishers(config.MsgQueue.Publishers, nodeId, logger)
	} else {
		app.publisherMap = make(map[string]mq.Publisher, 0)
	}
	if config.MsgQueue.Subscribers != nil {
		app.subscriberMap = loader.LoadSubscribers(config.MsgQueue.Subscribers, nodeId, logger)
	} else {
		app.subscriberMap = make(map[string]mq.Subscriber, 0)
	}

	if redisCache != nil {
		app.lockerFactory = locker.NewRedisLockerFactory(redisCache, logger)
	}

	if config.ObjectStorage != nil {
		app.objectStorage = object.NewMinioStorage(logger, config.ObjectStorage)
	}
	if config.WebSocket != nil {
		app.websocketServer = websocket.NewServer(config.WebSocket, logger, httpEngine, snowflakeNode, config.Mode)
	}

	if config.Metric != nil {
		app.metricService = metric.NewService(config.Name, nodeId, logger)
	}
}

func (app *Context) StartServe() {
	address := fmt.Sprintf("%s:%s", app.config.Host, app.config.Port)
	server := http.Server{
		Addr:    address,
		Handler: app.httpEngine,
	}
	go func() {
		app.logger.Infof("%s server start at: %s", app.config.Name, address)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			app.logger.Errorf("%s server start error: %v", app.config.Name, err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, channel := context.WithTimeout(context.Background(), 20*time.Second)
	defer channel()
	if err := server.Shutdown(ctx); err != nil {
		app.logger.Errorf("%s server start error: %v", app.config.Name, err)
	}
	app.logger.Infof("%s server end at: %v", app.config.Name, time.Now().UTC())
}
