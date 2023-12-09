package model

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type (
	SessionObject struct {
		Id         int64 `gorm:"id" json:"id"`
		SId        int64 `gorm:"s_id" json:"s_id"`
		FromUserId int64 `gorm:"from_user_id" json:"from_user_id"`
		ClientId   int64 `gorm:"client_id" json:"client_id"`
		CreateTime int64 `gorm:"create_time" json:"create_time"`
	}

	SessionObjectModel interface {
		AddSessionObjects(sId int64, fromUIds, clientMsgIds []int64, newFromUId, newClientMsgId, newSId int64) ([]int64, error)
		Insert(id, sId, fromUId, clientId int64) (int64, error)
	}

	defaultSessionObjectModel struct {
		db            *gorm.DB
		shards        int64
		logger        *logrus.Entry
		snowflakeNode *snowflake.Node
	}
)

func (d defaultSessionObjectModel) AddSessionObjects(sId int64, fromUIds, clientMsgIds []int64, newFromUId, newClientMsgId, newSId int64) ([]int64, error) {
	db := d.db
	ids := make([]int64, 0)
	tableName := d.genSessionObjectTableName(sId)
	sql := fmt.Sprintf("select * from %s where s_id = ? and from_user_id in ? and client_id in ?", tableName)
	objects := make([]*SessionObject, 0)
	err := db.Raw(sql, sId, fromUIds, clientMsgIds).Scan(&objects).Error
	if err != nil {
		return nil, err
	}
	if len(objects) > 0 {
		now := time.Now().UnixMilli()
		for i, object := range objects {
			ids = append(ids, object.Id)
			object.Id = objects[i].Id
			object.SId = newSId
			object.FromUserId = newFromUId
			object.ClientId = newClientMsgId
			object.CreateTime = now
		}
		newTableName := d.genSessionObjectTableName(newSId)
		err = db.Table(newTableName).Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(objects, len(objects)).Error
	}
	return ids, err
}

func (d defaultSessionObjectModel) Insert(id, sId, fromUId, clientId int64) (int64, error) {
	o := &SessionObject{
		Id:         id,
		SId:        sId,
		FromUserId: fromUId,
		ClientId:   clientId,
		CreateTime: time.Now().UnixMilli(),
	}
	tableName := d.genSessionObjectTableName(sId)
	return id, d.db.Table(tableName).Clauses(clause.OnConflict{DoNothing: true}).Create(o).Error
}

func (d defaultSessionObjectModel) genSessionObjectTableName(sId int64) string {
	return "session_object_" + fmt.Sprintf("%02d", sId%(d.shards))
}

func NewSessionObjectModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) SessionObjectModel {
	return defaultSessionObjectModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
