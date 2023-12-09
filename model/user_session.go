package model

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	// MutedAllBitInUserSessionStatus 全员被禁言标志位
	MutedAllBitInUserSessionStatus = 1 << 0
	// MutedSingleBitInUserSessionStatus 用户被禁言标志位
	MutedSingleBitInUserSessionStatus = 1 << 1
	// RejectBitInUserSessionStatus 拒收标志位
	RejectBitInUserSessionStatus = 1 << 0
	// SilenceBitInUserSessionStatus 静音标志位
	SilenceBitInUserSessionStatus = 1 << 1
)

type (
	UserSession struct {
		SessionId  int64   `gorm:"session_id" json:"session_id"`
		UserId     int64   `gorm:"user_id" json:"user_id"`
		ParentId   int64   `gorm:"parent_id" json:"parent_id"`
		Type       int     `gorm:"type" json:"type"`
		EntityId   int64   `gorm:"entity_id" json:"entity_id"`
		Name       string  `gorm:"name" json:"name"`
		Remark     string  `gorm:"remark" json:"remark"`
		ExtData    *string `json:"ext_data" json:"ext_data"`
		Top        int64   `gorm:"top" json:"top"`
		Role       int     `gorm:"role" json:"role"`
		Mute       int     `gorm:"mute" json:"mute"`
		Status     int     `gorm:"status" json:"status"`
		CreateTime int64   `gorm:"create_time" json:"create_time"`
		UpdateTime int64   `gorm:"update_time" json:"update_time"`
		Deleted    int8    `gorm:"deleted" json:"deleted"`
	}

	UserSessionModel interface {
		FindUserSessionByEntityId(userId, entityId int64, sessionType int, containDeleted bool) (*UserSession, error)
		UpdateUserSession(userIds []int64, sessionId int64, sessionName, sessionRemark, mute, extData *string, top *int64, status, role *int, parentId *int64) error
		FindEntityIdsInUserSession(userId, sessionId int64) []int64
		GetUserSessions(userId, mTime int64, offset, count int) ([]*UserSession, error)
		GetUserSession(userId, sessionId int64) (*UserSession, error)
		GenUserSessionTableName(userId int64) string
	}

	defaultUserSessionModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultUserSessionModel) FindUserSessionByEntityId(userId, entityId int64, sessionType int, containDeleted bool) (*UserSession, error) {
	userSession := &UserSession{}
	sqlStr := "select * from " + d.GenUserSessionTableName(userId) + " where user_id = ? and entity_id = ? and type = ?"
	if !containDeleted {
		sqlStr += " and deleted = 0"
	}
	err := d.db.Raw(sqlStr, userId, entityId, sessionType).Scan(&userSession).Error
	return userSession, err
}

func (d defaultUserSessionModel) UpdateUserSession(userIds []int64, sessionId int64, sessionName, sessionRemark, mute, extData *string, top *int64, status, role *int, parentId *int64) (err error) {
	// 分表uid数组
	sharedUIds := make(map[int64][]int64, 0)
	for _, uId := range userIds {
		share := uId % d.shards
		if sharedUIds[share] == nil {
			sharedUIds[share] = make([]int64, 0)
		}
		sharedUIds[share] = append(sharedUIds[share], uId)
	}

	tx := d.db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	for k, v := range sharedUIds {
		if sessionName == nil && sessionRemark == nil && top == nil && status == nil && mute == nil && role == nil {
			continue
		}
		sqlBuffer := bytes.Buffer{}
		sqlBuffer.WriteString(fmt.Sprintf("update %s set ", d.GenUserSessionTableName(k)))
		if sessionName != nil {
			sqlBuffer.WriteString(fmt.Sprintf("name = '%s', ", *sessionName))
		}
		if sessionRemark != nil {
			sqlBuffer.WriteString(fmt.Sprintf("remark = '%s', ", *sessionRemark))
		}
		if top != nil {
			sqlBuffer.WriteString(fmt.Sprintf("top = %d, ", *top))
		}
		if status != nil {
			sqlBuffer.WriteString(fmt.Sprintf("status = %d, ", *status))
		}
		if mute != nil {
			sqlBuffer.WriteString(fmt.Sprintf("mute = %s, ", *mute))
		}
		if extData != nil {
			sqlBuffer.WriteString(fmt.Sprintf("ext_data = %s, ", *extData))
		}
		if role != nil {
			sqlBuffer.WriteString(fmt.Sprintf("role = %d, ", *role))
		}
		if parentId != nil {
			sqlBuffer.WriteString(fmt.Sprintf("parent_id = %d, ", *parentId))
		}
		sqlBuffer.WriteString(fmt.Sprintf("update_time = %d ", time.Now().UnixMilli()))
		sqlBuffer.WriteString("where session_id = ? and user_id in ? ")
		err = tx.Exec(sqlBuffer.String(), sessionId, v).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (d defaultUserSessionModel) FindEntityIdsInUserSession(userId, sessionId int64) []int64 {
	entityIds := make([]int64, 0)
	sqlStr := fmt.Sprintf("select entity_id from %s where user_id = ? and session_id = ? and deleted = 0", d.GenUserSessionTableName(userId))
	_ = d.db.Raw(sqlStr, userId, sessionId).Scan(&entityIds).Error
	return entityIds
}

func (d defaultUserSessionModel) GetUserSessions(userId, mTime int64, offset, count int) ([]*UserSession, error) {
	userSessions := make([]*UserSession, 0)
	sqlStr := "select * from " + d.GenUserSessionTableName(userId) + " where user_id = ? and update_time <= ? limit ? offset ?"
	err := d.db.Raw(sqlStr, userId, mTime, count, offset).Scan(&userSessions).Error
	if err != nil {
		return nil, err
	}
	return userSessions, nil
}

func (d defaultUserSessionModel) GetUserSession(userId, sessionId int64) (*UserSession, error) {
	userSession := &UserSession{}
	sqlStr := "select * from " + d.GenUserSessionTableName(userId) + " where user_id = ? and session_id = ?"
	err := d.db.Raw(sqlStr, userId, sessionId).Scan(userSession).Error
	if err != nil {
		return nil, err
	}
	return userSession, nil
}

func (d defaultUserSessionModel) GenUserSessionTableName(userId int64) string {
	return "user_session_" + fmt.Sprintf("%02d", userId%(d.shards))
}

func NewUserSessionModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserSessionModel {
	return defaultUserSessionModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
