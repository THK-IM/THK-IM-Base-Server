package object

import (
	"context"
	"fmt"
	"github.com/h2non/filetype"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"os"
	"strings"
	"time"
)

type MinioStorage struct {
	logger *logrus.Entry
	conf   *conf.ObjectStorage
	client *minio.Client
}

func (m MinioStorage) UploadObject(key string, path string) (*string, error) {
	buf, errBuf := os.ReadFile(path)
	if errBuf != nil {
		return nil, errBuf
	}
	kind, errKind := filetype.Match(buf)
	if errKind != nil {
		return nil, errKind
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileInfo, errInfo := file.Stat()
	if errInfo != nil {
		return nil, errInfo
	}
	options := minio.PutObjectOptions{
		ContentType: kind.MIME.Value,
	}
	info, errPut := m.client.PutObject(context.Background(), m.conf.Bucket, key, file, fileInfo.Size(), options)
	if errPut != nil {
		return nil, errPut
	} else {
		return &info.Key, nil
	}
}

func (m MinioStorage) GetUploadParams(key string) (string, string, map[string]string, error) {
	policy := minio.NewPostPolicy()
	err := policy.SetBucket(m.conf.Bucket)
	if err != nil {
		return "", "", nil, err
	}
	err = policy.SetKey(key)
	if err != nil {
		return "", "", nil, err
	}
	// Expires in 10 days.
	err = policy.SetExpires(time.Now().UTC().Add(time.Minute * 5))
	if err != nil {
		return "", "", nil, err
	}
	// Returns form data for POST form request.
	uploadUrl, formData, errSign := m.client.PresignedPostPolicy(context.Background(), policy)
	if errSign != nil {
		return "", "", nil, errSign
	}
	params := make(map[string]string, 0)
	for k, v := range formData {
		params[k] = v
	}
	// params["success_action_status"] = "200"
	return uploadUrl.String(), "POST", params, nil
}

func (m MinioStorage) GetDownloadUrl(key string) (*string, error) {
	preSignedURL, err := m.client.PresignedGetObject(context.Background(), m.conf.Bucket, key, time.Minute*100, nil)
	if err != nil {
		return nil, nil
	}
	absolutPath := preSignedURL.String()
	return &absolutPath, nil
}

func NewMinioStorage(logger *logrus.Entry, conf *conf.ObjectStorage) Storage {
	secure := false
	endpoint := conf.Endpoint
	if strings.HasPrefix(conf.Endpoint, "https") {
		secure = true
	}
	endpoint = strings.Replace(endpoint, "http://", "", -1)
	endpoint = strings.Replace(endpoint, "https://", "", -1)

	fmt.Println(fmt.Sprintf("endpoint: %s", endpoint))
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(conf.AK, conf.SK, ""),
		Secure: secure,
		Region: conf.Region,
	})
	if err != nil {
		panic(err)
	}
	return &MinioStorage{logger: logger, conf: conf, client: client}
}
