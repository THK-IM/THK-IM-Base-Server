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
	Object struct {
		Id         int64  `gorm:"id" json:"id"`
		SId        int64  `gorm:"s_id" json:"s_id"`
		Engine     string `gorm:"engine" json:"engine"`
		Key        string `gorm:"key" json:"key"`
		CreateTime int64  `gorm:"create_time" json:"create_time"`
	}

	ObjectModel interface {
		AddSessions(ids []int64, sId int64) error
		Insert(sId int64, engine, key string) (int64, error)
		FindObject(id int64) (*Object, error)
		FindObjectByUId(id, uId int64, usTableName string) (*Object, error)
	}

	defaultObjectModel struct {
		db            *gorm.DB
		shards        int64
		logger        *logrus.Entry
		snowflakeNode *snowflake.Node
	}
)

func (d defaultObjectModel) AddSessions(ids []int64, sId int64) error {
	objects := make([]*Object, 0)
	for _, id := range ids {
		object, err := d.FindObject(id)
		if err != nil {
			return err
		}
		objects = append(objects, object)
	}

	now := time.Now().UnixMilli()
	for _, object := range objects {
		object.SId = sId
		object.CreateTime = now
		err := d.db.Table(d.genObjectTableName(object.Id)).Clauses(clause.OnConflict{DoNothing: true}).Create(object).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (d defaultObjectModel) Insert(sId int64, engine, key string) (int64, error) {
	id := d.snowflakeNode.Generate().Int64()
	o := &Object{
		Id:         id,
		SId:        sId,
		Engine:     engine,
		Key:        key,
		CreateTime: time.Now().UnixMilli(),
	}
	tableName := d.genObjectTableName(id)
	return id, d.db.Table(tableName).Clauses(clause.OnConflict{DoNothing: true}).Create(o).Error
}

func (d defaultObjectModel) FindObject(id int64) (*Object, error) {
	tableName := d.genObjectTableName(id)
	sql := fmt.Sprintf("select * from %s where id = ? limit 0, 1", tableName)
	object := &Object{}
	err := d.db.Raw(sql, id).Scan(object).Error
	return object, err
}

func (d defaultObjectModel) FindObjectByUId(id, uId int64, usTableName string) (*Object, error) {
	tableName := d.genObjectTableName(id)
	sql := fmt.Sprintf("select t0.* from %s as t0 "+
		"inner join %s as t1 on t1.session_id = t0.s_id and t1.deleted = 0 "+
		"where t0.id = ? and t1.user_id = ? limit 0, 1", tableName, usTableName)
	object := &Object{}
	err := d.db.Raw(sql, id, uId).Scan(object).Error
	return object, err
}

func (d defaultObjectModel) genObjectTableName(id int64) string {
	return "object_" + fmt.Sprintf("%02d", id%(d.shards))
}

func NewObjectModel(db *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node, shards int64) ObjectModel {
	return defaultObjectModel{db: db, logger: logger, snowflakeNode: snowflakeNode, shards: shards}
}
