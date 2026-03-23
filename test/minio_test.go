package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	endpoint    = "192.168.0.102:9000"
	minioAK     = "clawtravel"
	minioSk     = "clawtravel_pwd"
	minioBucket = "clawtravel"
)

func TestMinioGenerateSignPostParams(t *testing.T) {
	s3Client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAK, minioSk, ""),
		Secure: false,
		Region: "us-east-1",
	})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	policy := minio.NewPostPolicy()
	err = policy.SetBucket(minioBucket)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = policy.SetKey("user/image1.png")
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	// Expires in 10 days.
	err = policy.SetExpires(time.Now().UTC().AddDate(0, 0, 10))
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	// Returns form data for POST form request.
	url, formData, errSign := s3Client.PresignedPostPolicy(context.Background(), policy)
	if errSign != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Printf("curl ")
	for k, v := range formData {
		fmt.Printf("-F %s=%s ", k, v)
	}
	fmt.Printf("-F file=@./etc/image1.png ")
	fmt.Printf("%s\n", url)
}
