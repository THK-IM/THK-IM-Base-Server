package event

import (
	"encoding/json"
)

const (
	// ServerEventTypeKey 服务端事件类型Key
	ServerEventTypeKey = "ServerEventType"
	// ServerEventBodyKey 服务端事件Body Key
	ServerEventBodyKey = "ServerEventBody"
	// ServerEventUserOnline 服务端事件:用户上线
	ServerEventUserOnline = "user_online_event"
)

type (
	OnlineBody struct {
		NodeId     int64  `json:"node_id"`
		Online     bool   `json:"online"`
		UserId     int64  `json:"user_id"`
		ConnId     int64  `json:"conn_id"`
		Platform   string `json:"platform"`
		OnLineTime int64  `json:"on_line_time"`
		Token      string `json:"token"`
	}
)

func BuildUserOnlineEvent(nodeId int64, online bool, uid, connId, onLineTime int64, platform, token string) (map[string]interface{}, error) {
	onlineBody := &OnlineBody{
		NodeId:     nodeId,
		Online:     online,
		UserId:     uid,
		ConnId:     connId,
		Platform:   platform,
		OnLineTime: onLineTime,
		Token:      token,
	}
	b, err := json.Marshal(onlineBody)
	if err != nil {
		return nil, err
	}
	m := make(map[string]interface{}, 0)
	m[ServerEventTypeKey] = ServerEventUserOnline
	m[ServerEventBodyKey] = string(b)
	return m, nil
}

func ParserOnlineBody(jsonStr string) *OnlineBody {
	body := &OnlineBody{}
	if err := json.Unmarshal([]byte(jsonStr), body); err != nil {
		return nil
	} else {
		return body
	}
}
