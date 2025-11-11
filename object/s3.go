package object

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"mime"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type (
	S3Storage struct {
		logger *logrus.Entry
		conf   *conf.ObjectStorage
		client *s3.Client
	}

	s3PostPolicy struct {
		Expiration string        `json:"expiration"`
		Conditions []interface{} `json:"conditions"`
	}
)

// formatDates generates formatted date strings for use in AWS policies and signatures
func (s S3Storage) formatDates(timeStamp time.Time) (shortDate, amzDate string) {
	shortDate = timeStamp.Format("20060102")
	amzDate = timeStamp.Format("20060102T150405Z")
	return
}

// createPolicy generates a base64 encoded JSON policy document with given conditions
func (s S3Storage) createPolicy(bucket, keyPrefix, credential, amzDate string, expiration time.Time) string {
	conditions := []interface{}{
		map[string]string{"bucket": bucket},
		map[string]string{"x-amz-algorithm": "AWS4-HMAC-SHA256"},
		map[string]string{"x-amz-credential": credential},
		map[string]string{"x-amz-date": amzDate},
		[]interface{}{"content-length-range", 0, 1024 * 1024 * 1024}, // 1GB
		[]interface{}{"starts-with", "$key", keyPrefix},
		[]interface{}{"starts-with", "$Content-Type", ""}, // e.g. image/
	}

	policy := s3PostPolicy{
		Expiration: expiration.UTC().Format(time.RFC3339),
		Conditions: conditions,
	}

	policyBytes, _ := json.Marshal(policy)
	return base64.StdEncoding.EncodeToString(policyBytes)
}

// generateSignature creates an AWS v4 signature for the provided policy
func (s S3Storage) generateSignature(policy string, sk, region, date string) (string, error) {
	signingKey := s.deriveSigningKey(sk, date, region, "s3")
	signature := s.hmacSHA256(signingKey, policy)
	return hex.EncodeToString(signature), nil
}

// deriveSigningKey derives the signing key used for AWS signature v4
func (s S3Storage) deriveSigningKey(secret, date, region, service string) []byte {
	kDate := s.hmacSHA256([]byte("AWS4"+secret), date)
	kRegion := s.hmacSHA256(kDate, region)
	kService := s.hmacSHA256(kRegion, service)
	return s.hmacSHA256(kService, "aws4_request")
}

// hmacSHA256 performs HMAC-SHA256 hashing algorithm with given key and data
func (s S3Storage) hmacSHA256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

//// printPostFields displays the constructed fields needed for S3 POST requests
//func printPostFields(bucketName, policy string, creds aws.Credentials, signature, credential, amzDate string) {
//	fmt.Printf("URL: https://%s.s3.%s.amazonaws.com.cn\n", bucketName, Region)
//	fmt.Println("Fields:")
//	//fmt.Printf("Key: %s\n", key)
//	fmt.Printf("AWSAccessKeyId: %s\n", creds.AccessKeyID)
//	fmt.Printf("Policy: %s\n", policy)
//	fmt.Printf("x-amz-signature: %s\n", signature)
//	fmt.Printf("x-amz-credential: %s\n", credential)
//	fmt.Printf("x-amz-algorithm: AWS4-HMAC-SHA256\n")
//	fmt.Printf("x-amz-date: %s\n", amzDate)
//}

//// generateCURLCommand generates a CURL command for uploading a file to S3
//func generateCURLCommand(bucketName, key, contentType, policy, credential, amzDate, signature string) string {
//	filename := filepath.Base(key)
//	curlTemplate := `curl -X POST \
//  -F "key=%s" \
//  -F "Content-Type=%s" \
//  -F "X-Amz-Credential=%s" \
//  -F "X-Amz-Algorithm=AWS4-HMAC-SHA256" \
//  -F "X-Amz-Date=%s" \
//  -F "Policy=%s" \
//  -F "X-Amz-Signature=%s" \
//  -F "file=@%s" \
//  https://%s.s3.%s.amazonaws.com.cn/`
//
//	return fmt.Sprintf(curlTemplate, key, contentType, credential, amzDate, policy, signature, filename, bucketName, Region)
//}

func (s S3Storage) UploadObject(key string, path string) (*string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
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

func (s S3Storage) GetUploadParams(key string) (string, string, map[string]string, error) {
	// Define the object key
	keyPrefix := ""
	// Set time values
	timeStamp := time.Now().UTC()
	shortDate, amzDate := s.formatDates(timeStamp)

	// Create the AWS credential string
	credential := fmt.Sprintf("%s/%s/%s/s3/aws4_request", s.conf.AK, shortDate, s.conf.Region)

	// Create the policy
	expiration := timeStamp.Add(15 * time.Minute)
	policy := s.createPolicy(s.conf.Bucket, keyPrefix, credential, amzDate, expiration)

	// Sign the policy
	signature, errSign := s.generateSignature(policy, s.conf.SK, s.conf.Region, shortDate)
	if errSign != nil {
		return "", "", nil, errSign
	}
	params := make(map[string]string)
	params["key"] = key
	params["bucket"] = s.conf.Bucket
	params["X-Amz-Credential"] = credential
	params["X-Amz-Algorithm"] = "AWS4-HMAC-SHA256"
	params["X-Amz-Date"] = shortDate
	params["Policy"] = policy
	params["X-Amz-Signature"] = signature
	return s.conf.Endpoint, "POST", params, nil

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
}

func (s S3Storage) GetDownloadUrl(key string, second int64) (*string, error) {
	// 创建预签名客户端
	preSignClient := s3.NewPresignClient(s.client)
	output, err := preSignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(s.conf.Bucket),
		Key:    aws.String(key),
	}, func(options *s3.PresignOptions) {
		options.Expires = time.Duration(second) * time.Second
	})
	url := ""
	if err == nil && output != nil {
		if s.conf.Cdn != "" {
			re := regexp.MustCompile(`https?://(?:www\.)?([^/]+)`)
			url = re.ReplaceAllString(output.URL, s.conf.Cdn)
		} else {
			url = output.URL
		}
	}
	return &url, err
}

func (s S3Storage) DeleteObjectsByKeys(keys []string) error {
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

func (s S3Storage) KeyExists(key string) (bool, error) {
	output, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(s.conf.Bucket),
		Key:    aws.String(key),
	})
	return output != nil, err
}

func (s S3Storage) KeyFileSize(key string) (int64, error) {
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

func NewS3Storage(logger *logrus.Entry, conf *conf.ObjectStorage) Storage {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(conf.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(conf.AK, conf.SK, "")),
		config.WithBaseEndpoint(conf.Endpoint),
	)
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Storage{logger: logger, conf: conf, client: client}
}
