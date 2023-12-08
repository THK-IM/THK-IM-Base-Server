package loader

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Base-Server/conf"
	rotate "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

type LogFileWriter struct {
	Writer io.Writer
}

func (l LogFileWriter) Write(p []byte) (n int, err error) {
	_, _ = fmt.Println(string(p))
	return l.Writer.Write(p)
}

func LoadLogger(serverName string, logg *conf.Logger) *logrus.Entry {
	path := fmt.Sprintf("%s/%s", logg.Dir, serverName)
	logFileWriter, err := rotate.New(
		fmt.Sprintf("%s.%s.log", path, "%Y%m%d%H%M"),
		rotate.WithLinkName(path),
		rotate.WithMaxAge(time.Duration(logg.RetainAge)*time.Hour),
		rotate.WithRotationTime(time.Duration(logg.RotationAge)*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	level, errL := logrus.ParseLevel(logg.Level)
	if errL != nil {
		level = logrus.TraceLevel
	}
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetLevel(level)
	logger.SetOutput(LogFileWriter{Writer: logFileWriter})
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "@timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "caller",
		},
	})
	return logger.WithField("index_name", logg.IndexName)
}
