package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	SingleSessionType     = 1
	GroupSessionType      = 2
	SuperGroupSessionType = 3
)

type (
	Session struct {
		Id         int64   `gorm:"id" json:"id"`
		Name       string  `gorm:"name" json:"name"`
		Remark     string  `gorm:"remark" json:"remark"`
		Type       int     `gorm:"type" json:"type"`
		Mute       int8    `gorm:"mute" json:"mute"`
		ExtData    *string `json:"ext_data" json:"ext_data"`
		CreateTime int64   `gorm:"create_time" json:"create_time"`
		UpdateTime int64   `gorm:"update_time" json:"update_time"`
		Deleted    int8    `gorm:"deleted" json:"deleted"`
	}

	SessionModel interface {
		UpdateSession(sessionId int64, name, remark *string, mute *int, extData *string) error
		FindSession(sessionId int64) (*Session, error)
		CreateEmptySession(sessionType int, extData *string, name string, remark string) (*Session, error)
	}

	defaultSessionModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultSessionModel) UpdateSession(sessionId int64, name, remark *string, mute *int, extData *string) error {
	if name == nil && remark == nil && mute == nil {
		return nil
	}
	updateMap := make(map[string]interface{})
	if name != nil {
		updateMap["name"] = *name
	}
	if remark != nil {
		updateMap["remark"] = *remark
	}
	if mute != nil {
		updateMap["mute"] = *mute
	}
	if extData != nil {
		updateMap["ext_data"] = *extData
	}
	updateMap["update_time"] = time.Now().UnixMilli()
	return d.db.Table(d.genSessionTableName(sessionId)).Where("id = ?", sessionId).Updates(updateMap).Error
}

func (d defaultSessionModel) FindSession(sessionId int64) (*Session, error) {
	sqlStr := "select * from " + d.genSessionTableName(sessionId) + " where id = ? and deleted = 0"
	session := &Session{}
	err := d.db.Raw(sqlStr, sessionId).Scan(session).Error
	return session, err
}

func (d defaultSessionModel) CreateEmptySession(sessionType int, extData *string, name string, remark string) (*Session, error) {
	sessionId := int64(d.snowflakeNode.Generate())
	currTime := time.Now().UnixMilli()
	session := Session{
		Id:         sessionId,
		Type:       sessionType,
		Name:       name,
		Remark:     remark,
		ExtData:    extData,
		CreateTime: currTime,
		UpdateTime: currTime,
	}
	err := d.db.Table(d.genSessionTableName(sessionId)).Create(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (d defaultSessionModel) genSessionTableName(sessionId int64) string {
	return "session_" + fmt.Sprintf("%02d", sessionId%(d.shards))
}

func NewSessionModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) SessionModel {
	return defaultSessionModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
