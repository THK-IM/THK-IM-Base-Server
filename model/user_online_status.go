package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	UserOnlineStatus struct {
		UserId     int64  `gorm:"user_id"`
		OnlineTime int64  `gorm:"online_time"`
		ConnId     int64  `gorm:"conn_id"`
		Platform   string `gorm:"platform"`
	}

	UserOnlineStatusModel interface {
		GetUsersOnlineStatus(userIds []int64) ([]*UserOnlineStatus, error)
		UpdateUserOnlineStatus(userId, onlineTime, connId int64, platform string) error
		GetOnlineUserIds(userIds []int64, onlineTime int64) ([]int64, error)
	}

	defaultUserOnlineStatusModel struct {
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
		shards        int64
	}
)

func (d defaultUserOnlineStatusModel) GetUsersOnlineStatus(userIds []int64) ([]*UserOnlineStatus, error) {
	usersOnlineStatus := make([]*UserOnlineStatus, 0)
	sql := "select * from user_online_status_00 where user_id in ?"
	err := d.db.Raw(sql, userIds).Scan(&usersOnlineStatus).Error
	return usersOnlineStatus, err
}

func (d defaultUserOnlineStatusModel) UpdateUserOnlineStatus(userId, onlineTime, connId int64, platform string) (err error) {
	tx := d.db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	// 通过user_id和platform找到连接
	sqlStr := "select * from " + d.genUserOnlineStatusTable(userId) + " where user_id = ? and platform = ?"
	onlineStatus := &UserOnlineStatus{}
	err = tx.Raw(sqlStr, userId, platform).Scan(onlineStatus).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if onlineStatus.UserId <= 0 {
		// 插入
		sqlStr = "insert into " + d.genUserOnlineStatusTable(userId) +
			" (user_id, online_time, conn_id, platform) values (?, ?, ?, ?)"
		return tx.Exec(sqlStr, userId, onlineTime, connId, platform).Error
	} else {
		// 连接id不相等时更新
		sqlStr = "update " + d.genUserOnlineStatusTable(userId) +
			" set online_time = ?, conn_id = ? where user_id = ? and conn_id = ?"
		return tx.Exec(sqlStr, onlineTime, connId, userId, onlineStatus.ConnId).Error
	}
}

func (d defaultUserOnlineStatusModel) GetOnlineUserIds(userIds []int64, onlineTime int64) ([]int64, error) {
	sqlStr := "select user_id from user_online_status_00 where user_id in ? and online_time >= ?"
	onlineUserIds := make([]int64, 0)
	tx := d.db.Raw(sqlStr, userIds, onlineTime).Scan(&onlineUserIds)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return onlineUserIds, nil
}

func (d defaultUserOnlineStatusModel) genUserOnlineStatusTable(userId int64) string {
	return "user_online_status_" + fmt.Sprintf("%02d", userId%(d.shards))
}

func NewUserOnlineStatusModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) UserOnlineStatusModel {
	return defaultUserOnlineStatusModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
