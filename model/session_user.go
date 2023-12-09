package model

import (
	"bytes"
	"fmt"
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

const (
	SessionMember     = 1 // 普通成员 可以查询session信息，session会话历史消息
	SessionAdmin      = 2 // 管理员 模式下可以修改session基本信息，禁言单个用户, 添加/删除普通成员
	SessionSuperAdmin = 3 // 超级管理员 可以全员禁言, 添加/删除管理员
	SessionOwner      = 4 // 拥有者 可以添加超级管理员, 删除管理员，删除session
)

type (
	SessionUser struct {
		SessionId  int64 `gorm:"session_id" json:"session_id"`
		UserId     int64 `gorm:"user_id" json:"user_id"`
		Type       int   `gorm:"type" json:"type"`
		Role       int   `gorm:"role" json:"role"`
		Mute       int   `gorm:"mute" json:"mute"`
		Status     int   `gorm:"status" json:"status"`
		CreateTime int64 `gorm:"create_time" json:"create_time"`
		UpdateTime int64 `gorm:"update_time" json:"update_time"`
		Deleted    int8  `gorm:"deleted" json:"deleted"`
	}

	SessionUserModel interface {
		FindSessionUsersByMTime(sessionId, mTime int64, role *int, count int) ([]*SessionUser, error)
		FindAllSessionUsers(sessionId int64) ([]*SessionUser, error)
		FindSessionUsers(sessionId int64, userIds []int64) ([]*SessionUser, error)
		FindSessionUser(sessionId, userId int64) (*SessionUser, error)
		FindSessionUserCount(sessionId int64) (int, error)
		FindUIdsInSessionWithoutStatus(sessionId int64, status int, uIds []int64) []int64
		FindUIdsInSessionContainStatus(sessionId int64, status int, uIds []int64) []int64
		AddUser(session *Session, entityIds []int64, userIds []int64, role []int, maxCount int) ([]*UserSession, error)
		DelUser(session *Session, userIds []int64) (err error)
		UpdateUser(sessionId int64, userIds []int64, role, status *int, mute *string) (err error)
	}

	defaultSessionUserModel struct {
		shards        int64
		logger        *logrus.Entry
		db            *gorm.DB
		snowflakeNode *snowflake.Node
	}
)

func (d defaultSessionUserModel) FindSessionUsersByMTime(sessionId, mTime int64, role *int, count int) ([]*SessionUser, error) {
	sessionUser := make([]*SessionUser, 0)
	tableName := d.genSessionUserTableName(sessionId)
	var err error
	if role == nil {
		sqlStr := fmt.Sprintf("select * from %s where session_id = ? and deleted = 0 and update_time <= ? order by update_time desc limit 0, ?", tableName)
		err = d.db.Raw(sqlStr, sessionId, mTime, count).Scan(&sessionUser).Error
	} else {
		sqlStr := fmt.Sprintf("select * from %s where session_id = ? and deleted = 0 and role = ? and update_time <= ? order by update_time desc limit 0, ?", tableName)
		err = d.db.Raw(sqlStr, sessionId, *role, mTime, count).Scan(&sessionUser).Error
	}
	return sessionUser, err
}

func (d defaultSessionUserModel) FindAllSessionUsers(sessionId int64) ([]*SessionUser, error) {
	sessionUser := make([]*SessionUser, 0)
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select * from %s where session_id = ?", tableName)
	err := d.db.Raw(sqlStr, sessionId).Scan(&sessionUser).Error
	return sessionUser, err
}

func (d defaultSessionUserModel) FindSessionUsers(sessionId int64, userIds []int64) ([]*SessionUser, error) {
	sessionUser := make([]*SessionUser, 0)
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select * from %s where session_id = ? and user_id in ? and deleted = 0", tableName)
	err := d.db.Raw(sqlStr, sessionId, userIds).Scan(&sessionUser).Error
	return sessionUser, err
}

func (d defaultSessionUserModel) FindSessionUser(sessionId, userId int64) (*SessionUser, error) {
	sessionUser := &SessionUser{}
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select * from %s where session_id = ?  and user_id = ? and deleted = 0", tableName)
	err := d.db.Raw(sqlStr, sessionId, userId).Scan(sessionUser).Error
	return sessionUser, err
}

func (d defaultSessionUserModel) FindSessionUserCount(sessionId int64) (int, error) {
	count := 0
	tableName := d.genSessionUserTableName(sessionId)
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ?  and deleted = 0", tableName)
	err := d.db.Raw(sqlStr, sessionId).Scan(&count).Error
	return count, err
}

func (d defaultSessionUserModel) FindUIdsInSessionWithoutStatus(sessionId int64, status int, userIds []int64) []int64 {
	sessionUsers := make([]*SessionUser, 0)
	uIdsCondition := ""
	if len(userIds) > 0 {
		uIdsCondition = " and user_id in ? "
	}
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ? %s and status & ? = 0 and deleted = 0",
		d.genSessionUserTableName(sessionId), uIdsCondition)
	if len(userIds) > 0 {
		tx := d.db.Raw(sqlStr, sessionId, userIds, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	} else {
		tx := d.db.Raw(sqlStr, sessionId, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	}
	uIds := make([]int64, 0)
	for _, su := range sessionUsers {
		uIds = append(uIds, su.UserId)
	}
	return uIds
}

func (d defaultSessionUserModel) FindUIdsInSessionContainStatus(sessionId int64, status int, userIds []int64) []int64 {
	sessionUsers := make([]*SessionUser, 0)
	uIdsCondition := ""
	if len(userIds) > 0 {
		uIdsCondition = " and user_id in ? "
	}
	sqlStr := fmt.Sprintf("select user_id from %s where session_id = ? %s and status & ? > 0 and deleted = 0",
		d.genSessionUserTableName(sessionId), uIdsCondition)
	if len(userIds) > 0 {
		tx := d.db.Raw(sqlStr, sessionId, userIds, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	} else {
		tx := d.db.Raw(sqlStr, sessionId, status).Scan(&sessionUsers)
		if tx.Error != nil {
			d.logger.Error(tx.Error)
		}
	}
	uIds := make([]int64, 0)
	for _, su := range sessionUsers {
		uIds = append(uIds, su.UserId)
	}
	return uIds
}

func (d defaultSessionUserModel) AddUser(session *Session, entityIds []int64, userIds []int64, role []int, maxCount int) (userSessions []*UserSession, err error) {
	tx := d.db.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()
	count := 0
	tableName := d.genSessionUserTableName(session.Id)
	sqlStr := fmt.Sprintf("select count(0) from %s where session_id = ? and user_id not in ? and deleted = 0", tableName)
	if err = tx.Raw(sqlStr, session.Id, userIds).Scan(&count).Error; err != nil {
		return nil, err
	}

	if count > maxCount-len(userIds) {
		return nil, errorx.ErrGroupMemberCountBeyond
	}

	t := time.Now().UnixMilli()
	sql1 := "insert into " + d.genSessionUserTableName(session.Id) +
		" (session_id, user_id, role, type, create_time, update_time) values (?, ?, ?, ?, ?, ?) " +
		"on duplicate key update role = ?, deleted = ?, update_time = ? "

	userMute := 0
	if session.Mute == 1 {
		userMute = 1
	}
	userSessions = make([]*UserSession, 0)
	for index, id := range userIds {
		if err = tx.Exec(sql1, session.Id, id, role[index], session.Type, t, t, role[index], 0, t).Error; err != nil {
			return nil, err
		}

		sql2 := "insert into " + d.genUserSessionTableName(id) + " " +
			"(session_id, user_id, type, entity_id, role, name, remark, mute, ext_data, parent_id, " +
			"create_time, update_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)  " +
			"on duplicate key update  top = ?, role = ?, name = ?, remark = ?, mute = ?, " +
			"deleted = ?, ext_data = ?, parent_id = ?, update_time = ? "
		if err = tx.Exec(sql2, session.Id, id, session.Type, entityIds[index], role[index], session.Name,
			session.Remark, userMute, session.ExtData, 0, t, t, 0, role[index], session.Name,
			session.Remark, userMute, 0, session.ExtData, 0, t).Error; err != nil {
			return nil, err
		}
		userSession := &UserSession{
			SessionId:  session.Id,
			UserId:     userIds[index],
			Type:       session.Type,
			EntityId:   entityIds[index],
			Name:       session.Name,
			Remark:     session.Remark,
			Top:        0,
			Role:       role[index],
			Mute:       userMute,
			ExtData:    session.ExtData,
			ParentId:   0,
			Status:     0,
			CreateTime: t,
			UpdateTime: t,
			Deleted:    0,
		}
		userSessions = append(userSessions, userSession)
	}
	return
}

func (d defaultSessionUserModel) DelUser(session *Session, userIds []int64) (err error) {
	tx := d.db.Begin()
	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			err = tx.Commit().Error
		}
	}()
	t := time.Now().UnixMilli()
	sql1 := "update " + d.genSessionUserTableName(session.Id) +
		" set deleted = ?, update_time = ? where session_id = ? and user_id = ?"
	for _, id := range userIds {
		if err = tx.Exec(sql1, 1, t, session.Id, id).Error; err != nil {
			return err
		}

		sql2 := "update " + d.genUserSessionTableName(id) +
			" set deleted = ?, update_time = ? where session_id = ? and user_id = ?"
		if err = tx.Exec(sql2, 1, t, session.Id, id).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d defaultSessionUserModel) UpdateUser(sessionId int64, userIds []int64, role, status *int, mute *string) (err error) {
	if role == nil && status == nil && mute == nil {
		return nil
	}
	t := time.Now().UnixMilli()
	sqlBuffer := bytes.Buffer{}
	sqlBuffer.WriteString(fmt.Sprintf("update %s set ", d.genSessionUserTableName(sessionId)))
	if role != nil {
		sqlBuffer.WriteString(fmt.Sprintf(" role = %d, ", *role))
	}
	if status != nil {
		sqlBuffer.WriteString(fmt.Sprintf(" status = %d, ", *status))
	}
	if mute != nil {
		sqlBuffer.WriteString(fmt.Sprintf(" mute = %s, ", *mute))
	}
	sqlBuffer.WriteString(fmt.Sprintf(" update_time = %d where session_id = ? and user_id in ? ", t))
	return d.db.Exec(sqlBuffer.String(), sessionId, userIds).Error
}

func (d defaultSessionUserModel) genUserSessionTableName(userId int64) string {
	return "user_session_" + fmt.Sprintf("%02d", userId%(d.shards))
}

func (d defaultSessionUserModel) genSessionUserTableName(sessionId int64) string {
	return "session_user_" + fmt.Sprintf("%02d", sessionId%(d.shards))
}

func NewSessionUserModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) SessionUserModel {
	return defaultSessionUserModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
