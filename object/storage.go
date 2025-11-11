package object

import (
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
)

type Storage interface {
	UploadObject(key string, path string) (*string, error)
	GetUploadParams(key string) (string, string, map[string]string, error)
	GetDownloadUrl(key string, second int64) (*string, error)
	DeleteObjectsByKeys(keys []string) error
	KeyExists(key string) (bool, error)
	KeyFileSize(key string) (int64, error)
}

func NewStorage(logger *logrus.Entry, conf *conf.ObjectStorage) Storage {
	if conf.Engine == "oss" {
		return NewOssStorage(logger, conf)
	} else if conf.Engine == "minio" {
		return NewMinioStorage(logger, conf)
	} else if conf.Engine == "cloudflare" {
		return NewCloudFlareStorage(logger, conf)
	} else if conf.Engine == "s3" {
		return NewS3Storage(logger, conf)
	} else {
		panic("ObjectStorage Engine not supported")
	}
}
