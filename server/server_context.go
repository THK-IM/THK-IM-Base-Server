package server

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Base-Server/conf"
	"github.com/THK-IM/THK-IM-Base-Server/loader"
	"github.com/THK-IM/THK-IM-Base-Server/locker"
	"github.com/THK-IM/THK-IM-Base-Server/metric"
	"github.com/THK-IM/THK-IM-Base-Server/mq"
	"github.com/THK-IM/THK-IM-Base-Server/websocket"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Context struct {
	startTime       int64
	nodeId          int64
	metricServer    *metric.Service
	config          *conf.Config
	logger          *logrus.Entry
	redisCache      *redis.Client
	lockerFactory   locker.Factory
	database        *gorm.DB
	snowflakeNode   *snowflake.Node
	httpEngine      *gin.Engine
	websocketServer websocket.Server
	rpcMap          map[string]interface{}
	modelMap        map[string]interface{}
	publisherMap    map[string]mq.Publisher
	subscriberMap   map[string]mq.Subscriber
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

func (app *Context) AddRpc(name string, rpc interface{}) {
	app.rpcMap[name] = rpc
}

func (app *Context) AddModel(name string, model interface{}) {
	app.modelMap[name] = model
}

func (app *Context) NewLocker(key string, waitMs int, timeoutMs int) locker.Locker {
	return app.lockerFactory.NewLocker(key, waitMs, timeoutMs)
}

func (app *Context) Init(config *conf.Config) {
	logger := loader.LoadLogger(config.Name, config.Logger)
	redisCache := loader.LoadRedis(config.RedisSource)
	nodeId, startTime := loader.LoadNodeId(config, redisCache)
	snowflakeNode, err := snowflake.NewNode(nodeId)
	if err != nil {
		panic(err)
	}
	httpEngine := gin.New()
	app.logger = logger
	app.redisCache = redisCache
	app.nodeId = nodeId
	app.startTime = startTime
	app.snowflakeNode = snowflakeNode

	if config.MysqlSource != nil {
		app.database = loader.LoadMysql(logger, config.MysqlSource)
	}

	app.rpcMap = make(map[string]interface{}, 0)
	app.modelMap = make(map[string]interface{}, 0)

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

	if config.WebSocket != nil {
		app.websocketServer = websocket.NewServer(config.WebSocket, logger, httpEngine, snowflakeNode, config.Mode)
	}

	if config.Metric != nil {
		app.metricServer = metric.NewService(config.Name, nodeId, logger)
	}
}

func (app *Context) Start() {
	address := fmt.Sprintf("%s:%s", app.config.Host, app.config.Port)
	if e := app.httpEngine.Run(address); e != nil {
		panic(e)
	} else {
		app.logger.Infof("%s server start at: %s", app.config.Name, address)
	}
}
