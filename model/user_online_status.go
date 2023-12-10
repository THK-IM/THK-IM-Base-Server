package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	UserOnlineRecord struct {
		UserId     int64  `gorm:"user_id"`
		OnlineTime int64  `gorm:"online_time"`
		ConnId     int64  `gorm:"conn_id"`
		Platform   string `gorm:"platform"`
	}

	UserOnlineRecordModel interface {
		GetUsersOnlineRecords(userIds []int64) ([]*UserOnlineRecord, error)
		UpdateUserOnlineRecord(userId, onlineTime, connId int64, platform string) error
	}

	defaultUserOnlineRecordModel struct {
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
		shards        int64
	}
)

func (d defaultUserOnlineRecordModel) GetUsersOnlineRecords(userIds []int64) ([]*UserOnlineRecord, error) {
	usersOnlineStatus := make([]*UserOnlineRecord, 0)
	sql := "select * from user_online_status_00 where user_id in ?"
	err := d.db.Raw(sql, userIds).Scan(&usersOnlineStatus).Error
	return usersOnlineStatus, err
}

func (d defaultUserOnlineRecordModel) UpdateUserOnlineRecord(userId, onlineTime, connId int64, platform string) (err error) {
	tx := d.db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	// 通过user_id和platform找到连接
	sqlStr := "select * from " + d.genUserOnlineRecordTable(userId) + " where user_id = ? and platform = ?"
	onlineStatus := &UserOnlineRecord{}
	err = tx.Raw(sqlStr, userId, platform).Scan(onlineStatus).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if onlineStatus.UserId <= 0 {
		// 插入
		sqlStr = "insert into " + d.genUserOnlineRecordTable(userId) +
			" (user_id, online_time, conn_id, platform) values (?, ?, ?, ?)"
		return tx.Exec(sqlStr, userId, onlineTime, connId, platform).Error
	} else {
		// 连接id不相等时更新
		sqlStr = "update " + d.genUserOnlineRecordTable(userId) +
			" set online_time = ?, conn_id = ? where user_id = ? and conn_id = ?"
		return tx.Exec(sqlStr, onlineTime, connId, userId, onlineStatus.ConnId).Error
	}
}

func (d defaultUserOnlineRecordModel) genUserOnlineRecordTable(userId int64) string {
	return fmt.Sprintf("user_online_record_%d", userId%(d.shards))
}

func NewUserOnlineRecordModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserOnlineRecordModel {
	return defaultUserOnlineRecordModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
