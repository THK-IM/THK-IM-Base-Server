package test

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/filter"
	"math"
	"testing"
	"time"
)

func TestBloomFilter(t *testing.T) {
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

	factory := filter.NewRedisFactory(rdb, loggerEntry)
	//params := filter.DefaultBloomFilterParams("bl", time.Hour)
	params := filter.NewBloomFilterParams("bl", time.Hour, math.MaxUint32/2, 0.1, []uint32{31, 35, 37})
	bloomFilter := factory.NewBloomFilter(params)
	counts := []uint32{10, 100, 1000, 10000, 100000, 1000000}
	m, k, n, p := bloomFilter.Params()
	fmt.Println(fmt.Sprintf("m %d, k: %d, n %d,  p: %f", m, k, n, p))
	for _, c := range counts {
		rate := bloomFilter.FalseRate(c)
		fmt.Println(fmt.Sprintf("c %d, rate: %f", c, rate))
	}

	elements := make([]interface{}, 0)
	for i := 0; i < 10; i++ {
		elements = append(elements, i)
	}
	err = bloomFilter.Add(elements)
	if err != nil {
		t.Failed()
		fmt.Println(err)
	}

	elements = make([]interface{}, 0)
	elements = append(elements, 1)
	elements = append(elements, 3)
	elements = append(elements, 2)
	elements = append(elements, 11)
	elements = append(elements, 9)
	elements = append(elements, 11)
	elements = append(elements, 124)
	existed, errEx := bloomFilter.CheckExisted(elements)
	if errEx != nil {
		t.Failed()
		fmt.Println(errEx)
	}
	fmt.Println(existed)
}
