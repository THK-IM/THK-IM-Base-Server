package server

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/loader"
	"github.com/thk-im/thk-im-base-server/locker"
	"github.com/thk-im/thk-im-base-server/metric"
	"github.com/thk-im/thk-im-base-server/model"
	"github.com/thk-im/thk-im-base-server/mq"
	"github.com/thk-im/thk-im-base-server/object"
	"github.com/thk-im/thk-im-base-server/rpc"
	"github.com/thk-im/thk-im-base-server/websocket"
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
	objectStorage   object.Storage
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

func (app *Context) SessionModel() model.SessionModel {
	return app.modelMap["session"].(model.SessionModel)
}

func (app *Context) SessionMessageModel() model.SessionMessageModel {
	return app.modelMap["session_message"].(model.SessionMessageModel)
}

func (app *Context) SessionUserModel() model.SessionUserModel {
	return app.modelMap["session_user"].(model.SessionUserModel)
}

func (app *Context) UserMessageModel() model.UserMessageModel {
	return app.modelMap["user_message"].(model.UserMessageModel)
}

func (app *Context) UserSessionModel() model.UserSessionModel {
	return app.modelMap["user_session"].(model.UserSessionModel)
}

func (app *Context) ObjectModel() model.ObjectModel {
	return app.modelMap["object"].(model.ObjectModel)
}

func (app *Context) SessionObjectModel() model.SessionObjectModel {
	return app.modelMap["session_object"].(model.SessionObjectModel)
}

func (app *Context) UserOnlineStatusModel() model.UserOnlineStatusModel {
	return app.modelMap["user_online_status"].(model.UserOnlineStatusModel)
}

func (app *Context) RpcMsgApi() rpc.MsgApi {
	api, ok := app.rpcMap["msg-api"].(rpc.MsgApi)
	if ok {
		return api
	} else {
		return nil
	}
}

func (app *Context) RpcUserApi() rpc.UserApi {
	api, ok := app.rpcMap["user-api"].(rpc.UserApi)
	if ok {
		return api
	} else {
		return nil
	}
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
		if config.Models != nil {
			app.modelMap = loader.LoadModels(config.Models, app.database, logger, snowflakeNode)
		}
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

	if config.ObjectStorage != nil {
		app.objectStorage = object.NewMinioStorage(logger, config.ObjectStorage)
	}
	if config.WebSocket != nil {
		app.websocketServer = websocket.NewServer(config.WebSocket, logger, httpEngine, snowflakeNode, config.Mode)
	}

	if config.Metric != nil {
		app.metricServer = metric.NewService(config.Name, nodeId, logger)
	}
}

func (app *Context) StartServe() {
	address := fmt.Sprintf("%s:%s", app.config.Host, app.config.Port)
	if e := app.httpEngine.Run(address); e != nil {
		panic(e)
	} else {
		app.logger.Infof("%s server start at: %s", app.config.Name, address)
	}
}
