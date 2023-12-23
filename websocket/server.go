package websocket

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/middleware"
	"github.com/thk-im/thk-im-base-server/snowflake"
	"golang.org/x/net/websocket"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

const (
	UidKey = "uid"
)

type OnClientConnected func(client Client)
type OnClientClosed func(client Client)
type OnClientMsgReceived func(msg string, client Client)
type UidGetter func(claim dto.ThkClaims) (uid int64, err error)

type Server interface {
	Init() error
	Clients() map[int64][]Client
	ClientCount() int64
	SetUidGetter(g UidGetter)
	SetOnClientConnected(f OnClientConnected)
	SetOnClientClosed(f OnClientClosed)
	SetOnClientMsgReceived(r OnClientMsgReceived)
	GetUserClient(uid int64) []Client
	AddClient(uid int64, client Client) (err error)
	RemoveClient(uid int64, reason string, client Client) error
	SendMessage(uid int64, msg string) (err error)
	SendMessageToUsers(uIds []int64, msg string) (err error)
}

type WsServer struct {
	g                   *gin.Engine
	mode                string
	conf                *conf.WebSocket
	mutex               *sync.RWMutex
	logger              *logrus.Entry // 日志打印
	curCount            *atomic.Int64
	OnClientMsgReceived OnClientMsgReceived
	snowflakeNode       *snowflake.Node
	UidGetter           UidGetter
	userClients         map[int64][]Client
	OnClientConnected   OnClientConnected
	OnClientClosed      OnClientClosed
}

func NewServer(conf *conf.WebSocket, logger *logrus.Entry, g *gin.Engine, snowflakeNode *snowflake.Node, mode string) *WsServer {
	curCount := &atomic.Int64{}
	curCount.Store(0)
	mutex := &sync.RWMutex{}
	return &WsServer{
		g:             g,
		mode:          mode,
		logger:        logger,
		conf:          conf,
		curCount:      curCount,
		mutex:         mutex,
		snowflakeNode: snowflakeNode,
		userClients:   make(map[int64][]Client, 0),
	}
}

func (server *WsServer) SetUidGetter(getter UidGetter) {
	server.UidGetter = getter
}

func (server *WsServer) SetOnClientConnected(f OnClientConnected) {
	server.OnClientConnected = f
}
func (server *WsServer) SetOnClientClosed(f OnClientClosed) {
	server.OnClientClosed = f
}

func (server *WsServer) AddClient(uid int64, client Client) (err error) {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	clients, ok := server.userClients[uid]
	if !ok {
		clients = make([]Client, 0)
	}
	clients = append(clients, client)
	server.userClients[uid] = clients

	server.curCount.Add(1)
	if server.OnClientConnected != nil {
		server.OnClientConnected(client)
	}
	return
}

func (server *WsServer) RemoveClient(uid int64, reason string, client Client) (err error) {
	server.mutex.Lock()
	clients, ok := server.userClients[uid]
	if ok {
		if len(clients) == 1 {
			if clients[0].Info().Id == client.Info().Id {
				delete(server.userClients, uid)
				server.curCount.Add(-1)
			}
		} else {
			for i := 0; i < len(clients); i++ {
				if clients[i].Info().Id == client.Info().Id {
					newClients := append(clients[:i], clients[i+1:]...)
					server.userClients[uid] = newClients
					server.curCount.Add(-1)
					break
				}
			}
		}
	}
	server.mutex.Unlock()
	ok, err = client.Close(reason)
	if err == nil && ok {
		if server.OnClientConnected != nil {
			server.OnClientClosed(client)
		}
	}
	return err
}

func (server *WsServer) GetUserClient(uid int64) []Client {
	server.mutex.Lock()
	clients, ok := server.userClients[uid]
	server.mutex.Unlock()
	if ok == false {
		return []Client{}
	} else {
		return clients
	}
}

func (server *WsServer) SendMessage(uid int64, msg string) (err error) {
	server.mutex.RLock()
	clients, ok := server.userClients[uid]
	server.mutex.RUnlock()
	if ok {
		for _, c := range clients {
			if e := c.WriteMessage(msg); e != nil {
				server.logger.Errorf("client: %s, err, %s", c.Info(), err.Error())
			}
		}
	}
	return nil
}

func (server *WsServer) SendMessageToUsers(uIds []int64, msg string) (err error) {
	server.mutex.RLock()
	allClients := make([]Client, 0)
	for _, uid := range uIds {
		clients, ok := server.userClients[uid]
		if ok {
			allClients = append(allClients, clients...)
		}
	}
	server.mutex.RUnlock()
	server.logger.Info("SendMessageToUsers", uIds, len(allClients))
	for _, c := range allClients {
		e := c.WriteMessage(msg)
		if e != nil {
			server.logger.WithFields(logrus.Fields(c.Claims())).Errorf("client: %v, err, %s", c.Info(), err.Error())
		}
	}
	return nil
}

func (server *WsServer) Init() error {
	ws := websocket.Server{
		Handshake: func(c *websocket.Config, r *http.Request) error {
			return nil
		},
		Handler: server.onNewConn,
	}
	server.g.GET(server.conf.Uri, func(ctx *gin.Context) {
		err := server.getToken(ctx)
		if err != nil {
			ctx.AbortWithStatus(http.StatusForbidden)
		} else {
			ws.ServeHTTP(ctx.Writer, ctx.Request)
		}
	})
	return nil
}

func (server *WsServer) Clients() map[int64][]Client {
	return server.userClients
}

func (server *WsServer) ClientCount() int64 {
	return server.curCount.Load()
}

func (server *WsServer) SetOnClientMsgReceived(r OnClientMsgReceived) {
	server.OnClientMsgReceived = r
}

func (server *WsServer) onNewConn(ws *websocket.Conn) {
	claims := dto.ThkClaims{}
	claims.PutValue(dto.TraceID, ws.Request().Header.Get(dto.TraceID))
	claims.PutValue(dto.Language, ws.Request().Header.Get(dto.Language))
	claims.PutValue(dto.ClientVersion, ws.Request().Header.Get(dto.ClientPlatform))
	claims.PutValue(dto.ClientOriginIP, ws.Request().Header.Get(dto.ClientPlatform))
	claims.PutValue(dto.ClientPlatform, ws.Request().Header.Get(dto.ClientPlatform))

	if server.curCount.Load() >= server.conf.MaxClient {
		_ = ws.Close()
		server.logger.WithFields(logrus.Fields(claims)).Infof("client count reach max count %d", server.conf.MaxClient)
		return
	}

	uid := ws.Request().Header.Get(UidKey)
	uId, err := strconv.Atoi(uid)
	if err != nil {
		_ = ws.Close()
		server.logger.WithFields(logrus.Fields(claims)).Infof("uid: %s is invaild", uid)
		return
	}

	id := server.snowflakeNode.Generate()
	client := NewClient(ws, int64(id), int64(uId), claims, server)
	err = server.AddClient(int64(uId), client)
	if err != nil {
		server.logger.Error(err)
	} else {
		client.AcceptMessage()
	}
}

func (server *WsServer) getToken(ctx *gin.Context) error {
	claims := ctx.MustGet(middleware.ClaimsKey).(dto.ThkClaims)
	uid, err := server.UidGetter(claims)
	if err == nil {
		ctx.Request.Header.Set(dto.TraceID, claims.GetTraceId())
		ctx.Request.Header.Set(dto.Language, claims.GetLanguage())
		ctx.Request.Header.Set(dto.ClientVersion, claims.GetClientVersion())
		ctx.Request.Header.Set(dto.ClientOriginIP, claims.GetClientOriginIP())
		ctx.Request.Header.Set(dto.ClientPlatform, claims.GetClientPlatform())
		ctx.Request.Header.Set(UidKey, fmt.Sprintf("%d", uid))
	}
	return err

}
