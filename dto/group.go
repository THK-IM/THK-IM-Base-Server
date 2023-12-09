package dto

type CreateGroupReq struct {
	UId     int64   `json:"u_id"`
	Members []int64 `json:"members"`
}

type CreateGroupRes struct {
	GroupId     int64        `json:"group_id"`
	UserSession *UserSession `json:"user_session"`
}

type DeleteGroupReq struct {
	UId     int64 `json:"u_id"`
	GroupId int64 `json:"group_id"`
}

type TransferGroupReq struct {
	UId     int64 `json:"u_id"`
	ToUId   int64 `json:"to_u_id"`
	GroupId int64 `json:"group_id"`
}

type GroupInviteReq struct {
	UId       int64 `json:"u_id"`
	GroupId   int64 `json:"group_id"`
	InviteUId int64 `json:"invite_u_id"`
}

type GroupJoinApplyReq struct {
	UId     int64 `json:"u_id"`
	GroupId int64 `json:"group_id"`
}

type GroupJoinReq struct {
	UId     int64 `json:"u_id"`
	GroupId int64 `json:"group_id"`
}
