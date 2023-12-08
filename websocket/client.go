package websocket

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

type ClientInfo struct {
	Id              int64  // 唯一id
	UId             int64  // 用户id
	FirstOnLineTime int64  // 首次上线时间 毫秒
	LastOnlineTime  int64  // 最近心跳时间 毫秒
	Platform        string // 客户端平台 "android/ios/web/windows"
}

type Client interface {
	Info() *ClientInfo
	SetLastOnlineTime(mill int64)
	AcceptMessage()
	WriteMessage(msg string) error
	Close(reason string) (bool, error)
}

type WsClient struct {
	isClosed bool
	server   *WsServer
	logger   *logrus.Entry // 日志打印
	ws       *websocket.Conn
	info     *ClientInfo
	locker   *sync.Mutex
}

func (w *WsClient) LastOnlineTime() int64 {
	w.locker.Lock()
	defer w.locker.Unlock()
	return w.info.LastOnlineTime
}

func (w *WsClient) SetLastOnlineTime(mill int64) {
	w.locker.Lock()
	defer w.locker.Unlock()
	w.info.LastOnlineTime = mill
}

func (w *WsClient) FirstOnlineTime() int64 {
	w.locker.Lock()
	defer w.locker.Unlock()
	return w.info.FirstOnLineTime
}

func (w *WsClient) WriteMessage(msg string) error {
	if w.IsClosed() {
		err := w.server.RemoveClient(w.info.UId, "websocket closed", w)
		if err != nil {
			w.logger.Error(err)
			return err
		}
		return nil
	}
	w.locker.Lock()
	defer w.locker.Unlock()
	return websocket.Message.Send(w.ws, msg)
}

func (w *WsClient) IsClosed() bool {
	w.locker.Lock()
	defer w.locker.Unlock()
	return w.isClosed
}

func (w *WsClient) AcceptMessage() {
	w.read()
}

func (w *WsClient) read() {
	for {
		if w.IsClosed() {
			break
		}
		reply := ""
		if e := websocket.Message.Receive(w.ws, &reply); e == nil {
			if w.server.OnClientMsgReceived != nil {
				go w.server.OnClientMsgReceived(reply, w)
			} else {
				w.logger.Error(w.info, "onMsg handler is nil")
			}
		} else {
			w.logger.Error(e)
			if err := w.server.RemoveClient(w.info.UId, e.Error(), w); err != nil {
				w.logger.Error(w.info, err)
			}
			break
		}
	}
}

func (w *WsClient) Close(reason string) (bool, error) {
	w.locker.Lock()
	defer w.locker.Unlock()
	w.logger.Tracef("client: %v, close reason: %s", w.info, reason)
	if !w.isClosed {
		err := w.ws.Close()
		return true, err
	} else {
		return false, nil
	}
}

func (w *WsClient) Info() *ClientInfo {
	return w.info
}

func NewClient(ws *websocket.Conn, id, uId int64, platform string, server *WsServer) Client {
	onLineTime := time.Now().UnixMilli()
	info := ClientInfo{
		Id:              id,
		UId:             uId,
		FirstOnLineTime: onLineTime,
		LastOnlineTime:  onLineTime,
		Platform:        platform,
	}
	return &WsClient{
		server:   server,
		logger:   server.logger.WithField("uid", uId),
		ws:       ws,
		info:     &info,
		isClosed: false,
		locker:   &sync.Mutex{},
	}
}
