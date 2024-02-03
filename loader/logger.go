package loader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	elastic "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	rotate "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"io"
	"strings"
	"time"
)

type ElasticHook struct {
	client    *elastic.Client    // es的客户端
	host      string             // es 的 host
	index     string             // 获取索引的名字
	levels    []logrus.Level     // 日志的级别 info，error
	ctx       context.Context    // 上下文
	ctxCancel context.CancelFunc // 上下文cancel的函数，
	fireFunc  fireFunc           // 需要实现这个
}

type fireFunc func(entry *logrus.Entry, hook *ElasticHook) error

// NewElasticHook 新建一个es hook对象
func NewElasticHook(client *elastic.Client, level logrus.Level, index string) (*ElasticHook, error) {
	return newHookFuncAndFireFunc(client, level, index, syncFireFunc)
}

// 新建一个hook
func newHookFuncAndFireFunc(client *elastic.Client, level logrus.Level, index string, fireFunc fireFunc) (*ElasticHook, error) {
	var levels []logrus.Level
	for _, l := range []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	} {
		if l <= level {
			levels = append(levels, l)
		}
	}

	ctx, cancel := context.WithCancel(context.TODO())

	return &ElasticHook{
		client:    client,
		index:     index,
		levels:    levels,
		ctx:       ctx,
		ctxCancel: cancel,
		fireFunc:  fireFunc,
	}, nil
}

func createEsLog(entry *logrus.Entry) map[string]interface{} {
	level := entry.Level.String()
	m := make(map[string]interface{})
	for k, v := range entry.Data {
		m[k] = v
	}
	m["message"] = entry.Message
	m["caller"] = fmt.Sprintf("%s/%s:%d", entry.Caller.File, entry.Caller.Function, entry.Caller.Line)
	m["level"] = level
	m["@timestamp"] = entry.Time.UTC().Format(time.RFC3339)
	return m
}

func syncFireFunc(entry *logrus.Entry, hook *ElasticHook) error {
	data, errCreate := json.Marshal(createEsLog(entry))
	if errCreate != nil {
		return errCreate
	}

	index := hook.index
	if entry.Data["index_name"] != nil {
		if indexName, ok := entry.Data["index_name"].(string); ok {
			index = indexName
		}
	}

	req := esapi.IndexRequest{
		Index:   index,
		Body:    bytes.NewReader(data),
		Refresh: "true",
	}

	res, err := req.Do(hook.ctx, hook.client)
	if err != nil {
		return err
	}
	var r map[string]interface{}
	errDecode := json.NewDecoder(res.Body).Decode(&r)
	return errDecode
}

func (hook *ElasticHook) Fire(entry *logrus.Entry) error {
	return hook.fireFunc(entry, hook)
}

func (hook *ElasticHook) Levels() []logrus.Level {
	return hook.levels
}

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

	index := logg.IndexName
	if index == "" {
		index = serverName
	}

	if logg.ElasticEndpoint != "" {
		address := strings.Split(logg.ElasticEndpoint, ",")
		esClient, errClient := elastic.NewClient(elastic.Config{
			Addresses: address,
		})
		if errClient != nil {
			panic(errClient)
		}
		hook, errHook := NewElasticHook(esClient, level, index)
		if errHook != nil {
			panic(errHook)
		}
		logger.AddHook(hook)
	}

	return logger.WithField("index_name", index)
}
