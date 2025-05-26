package test

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/object"
	"testing"
)

func TestObject(t *testing.T) {
	logger := logrus.New()
	loggerEntry := logrus.NewEntry(logger)
	storageConf := &conf.ObjectStorage{
		Endpoint: "http://minio.thkim.com",
		Bucket:   "thk",
		AK:       "F94GlSuFJ5RNrDZ2PjNx",
		SK:       "L9lGsURsZzTZHeGAOyg5NkR8yRNYcZhrjrqrWqaY",
		Region:   "us-east-1",
	}
	key := "session-1855133845351829506/1855133462856471628/1855139947246260224-fFIPsn6D.png"
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
