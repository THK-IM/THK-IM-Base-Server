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
	endpoint  = "192.168.1.4:9000"
	thkImAk   = "SnK8hH6GvrTtlTETooKU"
	thkImSk   = "EA3wI8dVMwC7yo5QU5hizRmwaSs7wdZSA0U7JBnC"
	thkBucket = "thk"
)

func TestMinioGenerateSignPostParams(t *testing.T) {
	s3Client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(thkImAk, thkImSk, ""),
		Secure: false,
		Region: "us-east-1",
	})
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

	policy := minio.NewPostPolicy()
	err = policy.SetBucket(thkBucket)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = policy.SetKey("user/sample-111s.mp3")
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
	fmt.Printf("-F file=@/etc/sample-15s.mp3 ")
	fmt.Printf("%s\n", url)
}
