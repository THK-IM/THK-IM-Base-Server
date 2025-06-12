package object

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/errorx"
	"mime"
	"os"
	"path/filepath"
	"time"
)

type (
	CloudFlareStorage struct {
		logger *logrus.Entry
		conf   *conf.ObjectStorage
		client *s3.Client
	}
)

func (s CloudFlareStorage) UploadObject(key string, path string) (*string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	mimeType := mime.TypeByExtension(filepath.Ext(path))
	_, errPut := s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(s.conf.Bucket),
		Key:         aws.String(key),
		Body:        file,
		ContentType: aws.String(mimeType),
	})
	if errPut != nil {
		return nil, errPut
	} else {
		url := fmt.Sprintf("%s/%s/%s", s.conf.Endpoint, s.conf.Bucket, key)
		return &url, nil
	}
}

func (s CloudFlareStorage) GetUploadParams(key string) (string, string, map[string]string, error) {
	//preSignClient := s3.NewPresignClient(s.client)
	//output, err := preSignClient.PresignPostObject(context.Background(), &s3.PutObjectInput{
	//	Bucket: aws.String(s.conf.Bucket),
	//	Key:    aws.String(key),
	//}, func(options *s3.PresignPostOptions) {
	//	options.Expires = 10 * time.Minute
	//})
	//if err != nil {
	//	return "", "", map[string]string{}, err
	//} else {
	//	if output == nil {
	//		return "", "", map[string]string{}, errorx.ErrInternalServerError
	//	} else {
	//		return output.URL, "POST", output.Values, nil
	//	}
	//}

	mimeType := mime.TypeByExtension(filepath.Ext(key))
	preSignClient := s3.NewPresignClient(s.client)
	output, err := preSignClient.PresignPutObject(context.Background(), &s3.PutObjectInput{
		Bucket:      aws.String(s.conf.Bucket),
		Key:         aws.String(key),
		ContentType: aws.String(mimeType),
	}, func(options *s3.PresignOptions) {
		options.Expires = 10 * time.Minute
	})
	if err != nil {
		return "", "", map[string]string{}, err
	} else {
		if output == nil {
			return "", "", map[string]string{}, errorx.ErrInternalServerError
		} else {
			return output.URL, "PUT", map[string]string{}, nil
		}
	}
}

func (s CloudFlareStorage) GetDownloadUrl(key string) (*string, error) {
	// 创建预签名客户端
	preSignClient := s3.NewPresignClient(s.client)
	output, err := preSignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.conf.Bucket),
		Key:    aws.String(key),
	}, func(options *s3.PresignOptions) {
		options.Expires = 10 * time.Minute
	})
	url := ""
	if err == nil && output != nil {
		url = output.URL
	}
	return &url, err
}

func (s CloudFlareStorage) DeleteObjectsByKeys(keys []string) error {
	objects := make([]types.ObjectIdentifier, 0)
	for _, key := range keys {
		object := types.ObjectIdentifier{Key: aws.String(key)}
		objects = append(objects, object)
	}
	_, err := s.client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
		Bucket: aws.String(s.conf.Bucket),
		Delete: &types.Delete{
			Objects: objects,
		},
	})
	return err
}

func (s CloudFlareStorage) KeyExists(key string) (bool, error) {
	output, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(s.conf.Bucket),
		Key:    aws.String(key),
	})
	return output != nil, err
}

func (s CloudFlareStorage) KeyFileSize(key string) (int64, error) {
	output, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(s.conf.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, err
	}
	if output != nil && output.ContentLength != nil {
		return *(output.ContentLength), nil
	}
	return 0, nil
}

func NewCloudFlareStorage(logger *logrus.Entry, conf *conf.ObjectStorage) Storage {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(conf.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.AK, conf.SK, "")),
		config.WithBaseEndpoint(conf.Endpoint),
	)
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg)
	return &CloudFlareStorage{logger: logger, conf: conf, client: client}
}
