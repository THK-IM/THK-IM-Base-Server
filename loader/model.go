package loader

import (
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/model"
	"gorm.io/gorm"
	"os"
)

func LoadModels(modeConfigs []conf.Model, database *gorm.DB, logger *logrus.Entry, snowflakeNode *snowflake.Node) map[string]interface{} {
	modelMap := make(map[string]interface{}, 0)
	for _, ms := range modeConfigs {
		var m interface{}
		if ms.Name == "session" {
			m = model.NewSessionModel(database, logger, snowflakeNode, ms.Shards)
		} else if ms.Name == "session_message" {
			m = model.NewSessionMessageModel(database, logger, snowflakeNode, ms.Shards)
		} else if ms.Name == "session_user" {
			m = model.NewSessionUserModel(database, logger, snowflakeNode, ms.Shards)
		} else if ms.Name == "user_session" {
			m = model.NewUserSessionModel(database, logger, snowflakeNode, ms.Shards)
		} else if ms.Name == "user_message" {
			m = model.NewUserMessageModel(database, logger, snowflakeNode, ms.Shards)
		} else if ms.Name == "object" {
			m = model.NewObjectModel(database, logger, snowflakeNode, ms.Shards)
		} else if ms.Name == "session_object" {
			m = model.NewSessionObjectModel(database, logger, snowflakeNode, ms.Shards)
		} else if ms.Name == "user_online_status" {
			m = model.NewUserOnlineStatusModel(database, logger, snowflakeNode, ms.Shards)
		}
		modelMap[ms.Name] = m
	}
	return modelMap
}

func LoadTables(modeConfigs []conf.Model, database *gorm.DB) error {
	for _, ms := range modeConfigs {
		path := fmt.Sprintf("./sql/%s.sql", ms.Name)
		buffer, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		for i := int64(0); i < ms.Shards; i++ {
			sql := fmt.Sprintf(string(buffer), fmt.Sprintf("%02d", i))
			err = database.Exec(sql).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}
