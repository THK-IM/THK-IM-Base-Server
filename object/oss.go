package object

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"hash"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type OssStorage struct {
	logger *logrus.Entry
	conf   *conf.ObjectStorage
	client *oss.Client
	bucket *oss.Bucket
	cred   oss.CredentialsProvider
}

func (o OssStorage) UploadObject(key string, path string) (*string, error) {
	bucket, err := o.client.Bucket(o.conf.Bucket)
	if err != nil {
		return nil, err
	}

	err = bucket.PutObjectFromFile(key, path)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/%s", o.conf.Endpoint, o.conf.Bucket, key)
	return &url, nil
}

func (o OssStorage) GetUploadParams(key string) (string, string, map[string]string, error) {
	cred := o.cred.GetCredentials()
	product := "oss"
	// 构建Post Policy
	utcTime := time.Now().UTC()
	date := utcTime.Format("20060102")
	expiration := utcTime.Add(10 * time.Minute)
	policyMap := map[string]any{
		"expiration": expiration.Format("2006-01-02T15:04:05.000Z"),
		"conditions": []any{
			map[string]string{"bucket": o.conf.Bucket},
			map[string]string{"x-oss-signature-version": "OSS4-HMAC-SHA256"},
			map[string]string{"x-oss-credential": fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request",
				cred.GetAccessKeyID(), date, o.conf.Region, product)}, // 凭证
			map[string]string{"x-oss-date": utcTime.Format("20060102T150405Z")},
			//[]any{"eq", "$success_action_status", "200"},
		},
	}

	// 将Post Policy序列化为JSON字符串
	policy, err := json.Marshal(policyMap)
	if err != nil {
		log.Fatalf("json.Marshal fail, err:%v", err)
	}

	// 将Post Policy编码为Base64字符串
	stringToSign := base64.StdEncoding.EncodeToString(policy)

	// 生成签名密钥
	hmacHash := func() hash.Hash { return sha256.New() }
	signingKey := "aliyun_v4" + cred.GetAccessKeySecret()
	h1 := hmac.New(hmacHash, []byte(signingKey))
	_, _ = io.WriteString(h1, date)
	h1Key := h1.Sum(nil)

	h2 := hmac.New(hmacHash, h1Key)
	_, _ = io.WriteString(h2, o.conf.Region)
	h2Key := h2.Sum(nil)

	h3 := hmac.New(hmacHash, h2Key)
	_, _ = io.WriteString(h3, product)
	h3Key := h3.Sum(nil)

	h4 := hmac.New(hmacHash, h3Key)
	_, _ = io.WriteString(h4, "aliyun_v4_request")
	h4Key := h4.Sum(nil)

	// 计算Post签名
	h := hmac.New(hmacHash, h4Key)
	_, _ = io.WriteString(h, stringToSign)
	signature := hex.EncodeToString(h.Sum(nil))

	params := make(map[string]string)
	params["key"] = key
	params["policy"] = stringToSign
	params["x-oss-signature-version"] = "OSS4-HMAC-SHA256"
	params["x-oss-credential"] = fmt.Sprintf("%v/%v/%v/%v/aliyun_v4_request", cred.GetAccessKeyID(), date, o.conf.Region, product)
	params["x-oss-date"] = utcTime.Format("20060102T150405Z")
	params["x-oss-signature"] = signature

	url := fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/", o.conf.Bucket, o.conf.Region)
	return url, "POST", params, nil
}

func (o OssStorage) GetDownloadUrl(key string) (*string, error) {
	signedURL, err := o.bucket.SignURL(key, oss.HTTPGet, 600)
	if err != nil {
		return nil, err
	} else {
		return &signedURL, nil
	}
}

func (o OssStorage) DeleteObjectsByKeys(keys []string) error {
	_, err := o.bucket.DeleteObjects(keys)
	return err
}

func (o OssStorage) KeyExists(key string) (bool, error) {
	isExist, err := o.bucket.IsObjectExist(key)
	return isExist, err
}

func (o OssStorage) KeyFileSize(key string) (int64, error) {
	header, err := o.bucket.GetObjectMeta(key)
	if err != nil {
		return 0, err
	}
	sizeStr := header.Get("Content-Length")
	size, errSize := strconv.ParseInt(sizeStr, 10, 64)
	if errSize != nil {
		return 0, errSize
	}
	return size, nil
}

func NewOssStorage(logger *logrus.Entry, conf *conf.ObjectStorage) Storage {
	errSetEnv := os.Setenv("OSS_ACCESS_KEY_ID", conf.AK)
	if errSetEnv != nil {
		panic(errSetEnv)
	}
	errSetEnv = os.Setenv("OSS_ACCESS_KEY_SECRET", conf.SK)
	if errSetEnv != nil {
		panic(errSetEnv)
	}
	provider, err := oss.NewEnvironmentVariableCredentialsProvider()
	if err != nil {
		panic(err)
	}
	clientOptions := []oss.ClientOption{oss.SetCredentialsProvider(&provider)}
	clientOptions = append(clientOptions, oss.Region(conf.Region))
	// 设置预签名版本
	clientOptions = append(clientOptions, oss.AuthVersion(oss.AuthV4))
	client, errClient := oss.New(conf.Endpoint, conf.AK, conf.SK, clientOptions...)
	if errClient != nil {
		panic(errClient)
	}
	bucket, errBucket := client.Bucket(conf.Bucket)
	if errBucket != nil {
		panic(errBucket)
	}
	return &OssStorage{logger: logger, conf: conf, client: client, bucket: bucket, cred: &provider}
}
