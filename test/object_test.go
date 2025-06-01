package test

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/object"
	"testing"
)

func TestMinioObject(t *testing.T) {
	logger := logrus.New()
	loggerEntry := logrus.NewEntry(logger)
	storageConf := &conf.ObjectStorage{
		Endpoint: "http://minio.thkim.com",
		Bucket:   "thk",
		AK:       "F94GlSuFJ5RNrDZ2PjNx",
		SK:       "L9lGsURsZzTZHeGAOyg5NkR8yRNYcZhrjrqrWqaY",
		Region:   "us-east-1",
	}
	key := "session-1855133845351829506/1855133462856471628/1855139947246260224-fFIPsn6D_thumb.png"
	minioStorage := object.NewMinioStorage(loggerEntry, storageConf)
	existed, err := minioStorage.KeyExists(key)
	if err != nil {
		fmt.Println("err: ", err.Error())
		t.Failed()
	}
	if !existed {
		t.Failed()
	} else {
		errDel := minioStorage.DeleteObjectsByKeys([]string{key})
		if errDel != nil {
			t.Failed()
		}
	}

}

func TestOssObject(t *testing.T) {
	logger := logrus.New()
	loggerEntry := logrus.NewEntry(logger)
	storageConf := &conf.ObjectStorage{
		Endpoint: "https://oss-cn-beijing.aliyuncs.com",
		Bucket:   "trtsdd",
		AK:       "fdsdfsf",
		SK:       "dfsfsfs",
		Region:   "cn-beijing",
	}
	key := "zwtest/image1.png"
	storage := object.NewOssStorage(loggerEntry, storageConf)

	url, method, formData, errSign := storage.GetUploadParams(key)
	if errSign != nil {
		fmt.Println(errSign)
		t.Fail()
	}
	fmt.Printf("curl -i -X %s ", method)
	for k, v := range formData {
		fmt.Printf("-F %s=%s ", k, v)
	}
	fmt.Printf("-F file=@./etc/image1.png")
	fmt.Printf(" %s\n", url)

	existed, err := storage.KeyExists(key)
	if err != nil {
		fmt.Println("err: ", err.Error())
		t.Failed()
	}
	if !existed {
		t.Failed()
	} else {
		downloadUrl, errGet := storage.GetDownloadUrl(key)
		if errGet != nil {
			fmt.Println("err: ", errGet.Error())
			t.Failed()
		} else {
			fmt.Println(downloadUrl)
		}
		errDel := storage.DeleteObjectsByKeys([]string{key})
		if errDel != nil {
			fmt.Println("err: ", errDel.Error())
			t.Failed()
		}
	}

}
