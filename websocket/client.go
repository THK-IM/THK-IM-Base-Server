package websocket

import (
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/dto"
	"golang.org/x/net/websocket"
	"sync"
	"time"
)

type ClientInfo struct {
	Id              int64 // 唯一id
	UId             int64 // 用户id
	FirstOnLineTime int64 // 首次上线时间 毫秒
	LastOnlineTime  int64 // 最近心跳时间 毫秒
}

type Client interface {
	Info() *ClientInfo
	SetLastOnlineTime(mill int64)
	AcceptMessage()
	WriteMessage(msg string) error
	Close(reason string) (bool, error)
	Claims() dto.ThkClaims
}

type WsClient struct {
	isClosed bool
	server   *WsServer
	ws       *websocket.Conn
	locker   *sync.Mutex
	logger   *logrus.Entry // 日志打印
	info     *ClientInfo   // 用户信息
	claims   dto.ThkClaims // 用户数据
}

func (w *WsClient) Claims() dto.ThkClaims {
	return w.claims
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
			w.logger.Errorf("WriteMessage RemoveClient: %v %v", w.info.UId, err)
			return err
		} else {
			w.logger.Infof("WriteMessage RemoveClient: %v %v", w.info.UId, "success")
			return nil
		}
	}
	w.locker.Lock()
	defer w.locker.Unlock()
	encryptMsg := msg
	if w.server.crypto != nil {
		uri := w.server.conf.Uri
		encryptMessage, errEncrypt := w.server.crypto.EncryptUriBody(uri, []byte(msg))
		if errEncrypt != nil {
			return errEncrypt
		}
		encryptMsg = encryptMessage
	}
	err := websocket.Message.Send(w.ws, encryptMsg)
	if err != nil {
		w.logger.Errorf("WriteMessage: %v %v", msg, err)
		return err
	} else {
		w.logger.Infof("WriteMessage: %v %v", msg, "success")
		return nil
	}
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
			go w.server.OnClientMsg(w, reply)
		} else {
			w.logger.Errorf("read message error %v %v ", w.info, e)
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
	w.logger.Warnf("Close client: %v, reason: %s", w.info, reason)
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

func NewClient(ws *websocket.Conn, id, uId int64, claims dto.ThkClaims, server *WsServer) Client {
	onLineTime := time.Now().UnixMilli()
	info := ClientInfo{
		Id:              id,
		UId:             uId,
		FirstOnLineTime: onLineTime,
		LastOnlineTime:  onLineTime,
	}
	return &WsClient{
		server:   server,
		logger:   server.logger.WithFields(logrus.Fields(claims)).WithField("uid", uId),
		ws:       ws,
		info:     &info,
		isClosed: false,
		locker:   &sync.Mutex{},
		claims:   claims,
	}
}
