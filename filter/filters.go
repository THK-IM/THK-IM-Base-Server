package filter

import (
	"math"
	"time"
)

var defaultSeeds = []uint32{31, 33, 37}

type (
	// BitFilter 最多存储uint32类型Max值4294967295个元素
	BitFilter interface {
		Init(key string, ex time.Duration) error
		Clear(key string, ex time.Duration) error
		AddPos(key string, ex time.Duration, pos ...uint32) error
		Contains(key string, pos ...uint32) (contains []bool, err error)
		AllPos(key string) ([]uint32, error)
		Count(key string) (uint32, error)
		Delete(key string) error
	}

	BloomFilter interface {
		Delete() error
		AddOne(item interface{}) error
		Add(items []interface{}) error
		// CheckOneExisted 返回true 表示可能存在，false表示一定不存在
		CheckOneExisted(items interface{}) (bool, error)
		CheckExisted(items []interface{}) ([]bool, error)
		// FalseRate 计算插入n个元素后的误差率
		FalseRate(n uint32) float64
		// Params 参数m: bitmap容量, k:hash函数数量, n:设计元素数量，p:设计误差率
		Params() (m uint32, k uint32, n uint32, p float64)
	}

	BloomFilterParams struct {
		key     string
		expired time.Duration
		n       uint32
		p       float64
		seeds   []uint32
	}
)

func DefaultBloomFilterParams(key string, expired time.Duration) *BloomFilterParams {
	return NewBloomFilterParams(key, expired, 10*1024, 0.1, defaultSeeds)
}

// NewBloomFilterParams
// key 索引key
// expired 过期时间
// n 最大放入元素数量 (0, math.MaxUint32) n为math.MaxUint32时 所占内存为512M（redis string 最大内存）
// p 容忍最大误差 (0, 1]
// hash函数种子数组
func NewBloomFilterParams(key string, expired time.Duration, n uint32, p float64, seeds []uint32) *BloomFilterParams {
	if n > math.MaxUint32 {
		panic("NewBloomFilterParams n must be <= math.MaxUint32")
	}
	if p <= 0 || p > 1 {
		panic("NewBloomFilterParams p must be between 0 and 1")
	}
	return &BloomFilterParams{
		key:     key,
		expired: expired,
		n:       n,
		p:       p,
		seeds:   seeds,
	}
}
