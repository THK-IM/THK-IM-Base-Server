package dto

type GetSessionMessageReq struct {
	SId    int64  `json:"s_id" form:"id"`
	CTime  int64  `json:"c_time" form:"c_time"`
	Offset int    `json:"offset" form:"offset"`
	Count  int    `json:"count" form:"count"`
	MsgIds string `json:"msg_ids" form:"msg_ids"`
}

type DelSessionMessageReq struct {
	SId      int64   `json:"s_id"`
	MsgIds   []int64 `json:"msg_ids"`
	TimeFrom int64   `json:"time_from"`
	TimeTo   int64   `json:"time_to"`
}

type GetMessageReq struct {
	UId    int64 `json:"u_id" form:"u_id"`
	Offset int   `json:"offset" form:"offset"`
	Count  int   `json:"count" form:"count"`
	CTime  int64 `json:"c_time" form:"c_time"`
}

type GetMessageRes struct {
	Data interface{} `json:"data"`
}

type DeleteMessageReq struct {
	UId        int64   `json:"u_id" binding:"required"`
	SId        int64   `json:"s_id" binding:"required"`
	MessageIds []int64 `json:"msg_ids"`
	TimeFrom   *int64  `json:"time_from"`
	TimeTo     *int64  `json:"time_to"`
}

type PushMessageReq struct {
	UIds        []int64 `json:"u_ids" binding:"required"`
	Type        int     `json:"type"`
	Body        string  `json:"body" binding:"required"`
	OfflinePush bool    `json:"offline_push"`
}

type PushMessageRes struct {
	OnlineUIds  []int64 `json:"online_ids,omitempty"`
	OfflineUIds []int64 `json:"offline_ids,omitempty"`
}

type Message struct {
	CId     int64   `json:"c_id"` // 消息客户端id
	SId     int64   `json:"s_id"`
	MsgId   int64   `json:"msg_id"` // 消息服务端id
	Type    int     `json:"type"`
	FUid    int64   `json:"f_u_id"`
	CTime   int64   `json:"c_time"`
	Body    string  `json:"body"`
	Status  *int    `json:"status,omitempty"`
	RMsgId  *int64  `json:"r_msg_id,omitempty"`
	AtUsers *string `json:"at_users,omitempty"`
	ExtData *string `json:"ext_data,omitempty"`
}

type SendMessageReq struct {
	CId       int64   `json:"c_id" binding:"required"`
	SId       int64   `json:"s_id" binding:"required"`
	Type      int     `json:"type" binding:"required"`
	FUid      int64   `json:"f_u_id"`
	CTime     int64   `json:"c_time" binding:"required"`
	Body      string  `json:"body" binding:"required"`
	RMsgId    *int64  `json:"r_msg_id,omitempty"`
	AtUsers   *string `json:"at_users,omitempty"`
	Receivers []int64 `json:"receivers,omitempty"`
	ExtData   *string `json:"ext_data,omitempty"`
}

type SendMessageRes struct {
	MsgId      int64   `json:"msg_id"`
	CreateTime int64   `json:"c_time"`
	OnlineIds  []int64 `json:"online_ids,omitempty"`
	OfflineIds []int64 `json:"offline_ids,omitempty"`
}

type KickUserReq struct {
	UId int64 `json:"u_id" binding:"required"`
}
