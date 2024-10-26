package test

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/pool"
	"testing"
	"time"
)

func TestRecommendPool(t *testing.T) {
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
	recommendPool := factory.NewRecommendPool("re", 10*1000)

	elements := make([]*pool.Element, 0)
	for i := int64(0); i < 10*1000; i++ {
		element := &pool.Element{
			Id:    i,
			Score: float64(i),
		}
		elements = append(elements, element)
	}
	err = recommendPool.Add(elements...)
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}

	count := int64(0)
	count, err = recommendPool.ElementCount()
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	fmt.Println(count)

	count, err = recommendPool.ElementCountByRange(0, 10)
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	fmt.Println(count)

	count, err = recommendPool.RemoveByScore(1000)
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	fmt.Println(count)

	count, err = recommendPool.ElementCount()
	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	fmt.Println(count)

	st := make([]*pool.RecommendStrategy, 0)
	st1 := &pool.RecommendStrategy{
		Type:             pool.StrategyScore,
		Count:            10,
		RepeatRetryCount: 2,
	}
	st = append(st, st1)
	st2 := &pool.RecommendStrategy{
		Type:             pool.StrategyRandom,
		Count:            10,
		RepeatRetryCount: 1,
	}
	st = append(st, st2)
	uId := int64(1)
	recommendElements, errRecommend := recommendPool.FetchElements(uId, st)
	if errRecommend != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
	recommendIds := make([]int64, 0)
	for _, element := range recommendElements {
		fmt.Println(element)
		recommendIds = append(recommendIds, element.Id)
	}
	fmt.Println(len(recommendElements))
	err = recommendPool.AddUserRecord(uId, time.Hour, recommendIds...)

	if err != nil {
		fmt.Println(err)
		t.Failed()
		return
	}
}
