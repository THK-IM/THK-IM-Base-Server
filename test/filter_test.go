package test

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/filter"
	"math/rand"
	"testing"
	"time"
)

func TestFilter(t *testing.T) {
	opt, err := redis.ParseURL("redis://:dev123456@redis.yujianmeet.cn:16379")
	if err != nil {
		t.Failed()
		fmt.Println(err)
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

	fl := filter.NewRedisBitmap(rdb, loggerEntry, 1024)
	poss := make([]uint32, 0)
	for i := 0; i < 100; i++ {
		if i%2 == 0 {
			poss = append(poss, rand.Uint32()%100)
		}
	}
	for i := 0; i < len(poss); i++ {
		for j := i; j < len(poss); j++ {
			if poss[i] > poss[j] {
				poss[i], poss[j] = poss[j], poss[i]
			}
		}
	}
	fmt.Println(time.Now().UnixMilli(), poss, len(poss))
	key := "filter_0"
	err = fl.Init(key, time.Hour)
	err = fl.Clear(key, time.Hour)
	err = fl.AddPos(key, time.Hour, poss...)
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	fmt.Println(time.Now().UnixMilli(), "AddPos success")

	count, cErr := fl.Count(key)
	if cErr != nil {
		fmt.Println(cErr)
		t.Failed()
		return
	}
	fmt.Println(time.Now().UnixMilli(), count)

	allPoss, pErr := fl.AllPos(key)
	if pErr != nil {
		fmt.Println(pErr)
		t.Failed()
		return
	}
	fmt.Println(time.Now().UnixMilli(), "allPoss", allPoss, len(allPoss))

	for i := 10; i < 20; i++ {
		poss = append(poss, uint32(i))
	}
	contains, ctErr := fl.Contains(key, poss...)
	if ctErr != nil {
		fmt.Println(ctErr)
		t.Failed()
		return
	}
	fmt.Println("contains", contains, len(contains))

	//err = fl.Clear(key, time.Hour)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
}
