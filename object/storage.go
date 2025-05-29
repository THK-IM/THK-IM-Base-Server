package object

import (
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
)

type Storage interface {
	UploadObject(key string, path string) (*string, error)
	GetUploadParams(key string) (string, string, map[string]string, error)
	GetDownloadUrl(key string) (*string, error)
	DeleteObjectsByKeys(keys []string) error
	KeyExists(key string) (bool, error)
	KeyFileSize(key string) (int64, error)
}

func NewStorage(logger *logrus.Entry, conf *conf.ObjectStorage) Storage {
	if conf.Engine == "oss" {
		return NewMinioStorage(logger, conf)
	} else if conf.Engine == "minio" {
		return NewOssStorage(logger, conf)
	} else {
		panic("ObjectStorage Engine not supported")
	}
}
