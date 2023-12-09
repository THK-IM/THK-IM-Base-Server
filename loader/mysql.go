package loader

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

type DBLogger struct {
	logger *logrus.Entry
}

func (l *DBLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *DBLogger) Info(ctx context.Context, arg0 string, args ...interface{}) {
	l.logger.WithContext(ctx).Info(arg0, args)
}

func (l *DBLogger) Warn(ctx context.Context, arg0 string, args ...interface{}) {
	l.logger.WithContext(ctx).Warn(arg0, args)
}

func (l *DBLogger) Error(ctx context.Context, arg0 string, args ...interface{}) {
	l.logger.WithContext(ctx).Error(arg0, args)
}

func (l *DBLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rowsAffected := fc()
	end := time.Now()
	duration := end.UnixMilli() - begin.UnixMilli()
	l.logger.WithContext(ctx).Tracef("sql: %s, affect rows: %d, duration: %d", sql, rowsAffected, duration)
}

func LoadMysql(entry *logrus.Entry, source *conf.MysqlSource) *gorm.DB {
	if source == nil {
		return nil
	}
	dbLogger := &DBLogger{
		logger: entry,
	}
	db, errOpen := gorm.Open(mysql.Open(fmt.Sprintf("%s%s", source.Endpoint, source.Uri)), &gorm.Config{
		Logger: dbLogger,
	})
	if errOpen != nil {
		panic(errOpen)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxIdleConns(source.MaxIdleConn)
	sqlDb.SetMaxOpenConns(source.MaxOpenConn)
	sqlDb.SetConnMaxIdleTime(time.Duration(source.ConnMaxIdleTime) * time.Second)
	sqlDb.SetConnMaxLifetime(time.Duration(source.ConnMaxLifeTime) * time.Second)

	return db
}
