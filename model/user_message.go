package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MsgStatusAcked      = 1
	MsgStatusClientRead = 2 // 客户端已读
	MsgStatusServerRead = 4 // 服务端已读
	MsgStatusRead       = MsgStatusClientRead | MsgStatusServerRead
	MsgStatusReedit     = 8
)

type (
	UserMessage struct {
		MsgId      int64   `gorm:"msg_id" json:"msg_id"`
		ClientId   int64   `gorm:"client_id" json:"client_id"`
		UserId     int64   `gorm:"user_id" json:"user_id"`
		SessionId  int64   `gorm:"session_id" json:"session_id"`
		FromUserId int64   `gorm:"from_user_id" json:"from_user_id"`
		MsgType    int     `gorm:"msg_type" json:"msg_type"`
		MsgContent string  `gorm:"msg_content" json:"msg_content"`
		ReplyMsgId *int64  `gorm:"reply_msg_id" json:"reply_msg_id"`
		AtUsers    *string `gorm:"at_users" json:"at_users"`
		ExtData    *string `gorm:"ext_data" json:"ext_data"`
		Status     int     `gorm:"status" json:"status"`
		CreateTime int64   `gorm:"create_time" json:"create_time"`
		UpdateTime int64   `gorm:"update_time" json:"update_time"`
		Deleted    int8    `gorm:"deleted" json:"deleted"`
	}

	UserMessageModel interface {
		FindUserMessages(userId, sessionId int64, messageId []int64) ([]*UserMessage, error)
		FindUserMessage(userId, sessionId, messageId int64) (*UserMessage, error)
		InsertUserMessage(m *UserMessage) error
		AckUserMessages(userId int64, sessionId int64, messageIds []int64) error
		GetUserMessages(userId int64, ctime int64, offset, count int) ([]*UserMessage, error)
		DeleteMessages(userId int64, sessionId int64, messageIds []int64, from, to *int64) error
		DeleteMessagesBySessionId(userId int64, sessionId int64) error
		UpdateUserMessage(userId int64, sessionId int64, msgIds []int64, status int, content *string) error
	}

	defaultUserMessageModel struct {
		shards        int64
		db            *gorm.DB
		logger        *logrus.Entry
		snowflakeNode *snowflake.Node
	}
)

func (d defaultUserMessageModel) FindUserMessages(userId, sessionId int64, messageId []int64) ([]*UserMessage, error) {
	results := make([]*UserMessage, 0)
	strSql := "select * from " + d.genUserMessageTableName(userId) + " where user_id = ? and session_id = ? and msg_id in ?"
	err := d.db.Raw(strSql, userId, sessionId, messageId).Scan(results).Error
	return results, err
}

func (d defaultUserMessageModel) FindUserMessage(userId, sessionId, messageId int64) (*UserMessage, error) {
	result := &UserMessage{}
	strSql := "select * from " + d.genUserMessageTableName(userId) + " where user_id = ? and session_id = ? and msg_id = ?"
	err := d.db.Raw(strSql, userId, sessionId, messageId).Scan(result).Error
	return result, err
}

func (d defaultUserMessageModel) InsertUserMessage(m *UserMessage) error {
	return d.db.Table(d.genUserMessageTableName(m.UserId)).Clauses(clause.OnConflict{DoNothing: true}).Create(m).Error
}

func (d defaultUserMessageModel) AckUserMessages(userId int64, sessionId int64, messageIds []int64) error {
	sqlStr := fmt.Sprintf("update %s set status = (status | 1) where user_id = ?  and session_id = ? and msg_id in ? ",
		d.genUserMessageTableName(userId))
	err := d.db.Exec(sqlStr, userId, sessionId, messageIds).Error
	return err
}

func (d defaultUserMessageModel) GetUserMessages(userId int64, ctime int64, offset, count int) ([]*UserMessage, error) {
	result := make([]*UserMessage, 0)
	strSql := "select * from " + d.genUserMessageTableName(userId) + " where user_id = ? and deleted = 0 and (status = 0 or create_time > ?) order by create_time limit ? offset ?"

	tx := d.db.Raw(strSql, userId, ctime, count, offset).Scan(&result)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return result, nil
}

func (d defaultUserMessageModel) DeleteMessages(userId int64, sessionId int64, messageIds []int64, from, to *int64) error {
	if len(messageIds) > 0 {
		sqlStr := fmt.Sprintf("update %s set deleted = 1 where user_id = ? and session_id = ? and msg_id in ?",
			d.genUserMessageTableName(userId))
		err := d.db.Exec(sqlStr, userId, sessionId, messageIds).Error
		return err
	} else if from != nil && to != nil {
		sqlStr := fmt.Sprintf(
			"update %s set deleted = 1 where user_id = ? and session_id = ? and create_time >= ? and create_time <= ?",
			d.genUserMessageTableName(userId))
		err := d.db.Exec(sqlStr, userId, sessionId, from, to).Error
		return err
	} else {
		return nil
	}
}

func (d defaultUserMessageModel) DeleteMessagesBySessionId(userId int64, sessionId int64) error {
	sqlStr := fmt.Sprintf(
		"update %s set deleted = 1 where user_id = ? and session_id = ? ",
		d.genUserMessageTableName(userId))
	err := d.db.Exec(sqlStr, userId, sessionId).Error
	return err
}

func (d defaultUserMessageModel) UpdateUserMessage(userId int64, sessionId int64, msgIds []int64, status int, content *string) error {
	updateContent := ""
	if content != nil {
		updateContent = fmt.Sprintf(", msg_content = '%s' ", *content)
	}
	sqlStr := fmt.Sprintf(
		"update %s set status = status | ? %s where user_id = ? and session_id = ? and msg_id in ? ",
		d.genUserMessageTableName(userId), updateContent)
	err := d.db.Exec(sqlStr, status, userId, sessionId, msgIds).Error
	return err
}

func (d defaultUserMessageModel) genUserMessageTableName(userId int64) string {
	return "user_message_" + fmt.Sprintf("%02d", userId%(d.shards))
}

func NewUserMessageModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserMessageModel {
	return defaultUserMessageModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
