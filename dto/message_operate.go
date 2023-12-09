package dto

type AckUserMessagesReq struct {
	UId    int64   `json:"u_id"`
	SId    int64   `json:"s_id" binding:"required"`
	MsgIds []int64 `json:"msg_ids" binding:"required"`
}

type ReadUserMessageReq struct {
	UId    int64   `json:"u_id"`
	SId    int64   `json:"s_id" binding:"required"`
	MsgIds []int64 `json:"msg_ids" binding:"required"`
}

type RevokeUserMessageReq struct {
	UId   int64 `json:"u_id"`
	SId   int64 `json:"s_id" binding:"required"`
	MsgId int64 `json:"msg_id" binding:"required"`
}

type ReeditUserMessageReq struct {
	UId     int64  `json:"u_id"`
	SId     int64  `json:"s_id" binding:"required"`
	MsgId   int64  `json:"msg_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ForwardUserMessageReq struct {
	SendMessageReq
	ForwardSId       int64   `json:"fwd_s_id" binding:"required"`
	ForwardFromUIds  []int64 `json:"fwd_from_u_ids" binding:"required"`
	ForwardClientIds []int64 `json:"fwd_client_ids" binding:"required"`
}
