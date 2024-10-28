package test

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/pool"
	"testing"
	"time"
)

func TestMatchPool(t *testing.T) {
	opt, err := redis.ParseURL("redis://:dev123456@redis.yujianmeet.cn:16379")
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	opt.ConnMaxLifetime = 3 * time.Second
	opt.ConnMaxIdleTime = 3 * time.Second
	opt.ReadTimeout = 3 * time.Second
	opt.WriteTimeout = 3 * time.Second
	opt.PoolTimeout = 3 * time.Second
	opt.MaxIdleConns = 3
	opt.PoolSize = 3
	rdb := redis.NewClient(opt)
	loggerEntry := logrus.New().WithFields(logrus.Fields{})

	factory := pool.NewRedisPoolFactory(rdb, loggerEntry)
	matchPool := factory.NewMatchPool("ma")
	_, err = matchPool.Add("1", "2", "3")
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	count := int64(0)
	count, err = matchPool.Count()
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	fmt.Println(count)

	contain := false
	contain, err = matchPool.Contain("1")
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	fmt.Println(contain)

	_, err = matchPool.Remove("1", "2")
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}

	contain, err = matchPool.Contain("1")
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	fmt.Println(contain)

	_, err = matchPool.Add("8")
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}

	count, err = matchPool.Count()
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	fmt.Println(count)

	var id *string = nil
	id, err = matchPool.Match("7", 10, func(id string, candidateId string) (matchedId *string, putBlack bool, err error) {
		if candidateId == "8" {
			return &candidateId, false, nil
		} else {
			return nil, true, nil
		}
	})

	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	if id != nil {
		fmt.Println(*id)
	}

	count, err = matchPool.Count()
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	fmt.Println(count)
}
