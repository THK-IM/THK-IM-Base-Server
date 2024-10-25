package filter

import (
	"bytes"
	"encoding/binary"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

type RedisBloomFilter struct {
	bitmap  BitFilter
	key     string
	expired time.Duration
	seeds   []uint32
	m       uint32
	n       uint32
	p       float64
}

func (r RedisBloomFilter) Clear() error {
	return r.bitmap.Clear(r.key, r.expired)
}

func (r RedisBloomFilter) Delete() error {
	return r.bitmap.Delete(r.key)
}

func (r RedisBloomFilter) Params() (m uint32, k uint32, n uint32, p float64) {
	return r.m, uint32(len(r.seeds)), r.n, r.p
}

func (r RedisBloomFilter) FalseRate(n uint32) float64 {
	// 计算误判率的公式
	m := r.m
	k := uint32(len(r.seeds))
	p := math.Pow(1-math.Pow(1-float64(k)/float64(m), float64(n)), float64(k))
	return p
}

func (r RedisBloomFilter) interfaceToBytes(i interface{}) ([]byte, error) {
	var buf bytes.Buffer
	switch v := i.(type) {
	case int:
		if err := binary.Write(&buf, binary.LittleEndian, int64(v)); err != nil {
			return nil, err
		}
	case string:
		buf.WriteString(v)
	default:
		if err := binary.Write(&buf, binary.LittleEndian, v); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (r RedisBloomFilter) hashFun(seed uint32, value interface{}) (uint32, error) {
	hash := seed
	bs, err := r.interfaceToBytes(value)
	if err != nil {
		return 0, err
	}
	for i := 0; i < len(bs); i++ {
		hash = hash*33 + uint32(bs[i])
	}
	return hash, nil
}

func (r RedisBloomFilter) AddOne(item interface{}) error {
	poss := make([]uint32, 0)
	for _, seed := range r.seeds {
		if hash, err := r.hashFun(seed, item); err != nil {
			return err
		} else {
			poss = append(poss, hash%r.m)
		}
	}
	return r.bitmap.AddPos(r.key, r.expired, poss...)
}

func (r RedisBloomFilter) Add(items []interface{}) error {
	poss := make([]uint32, 0)
	for _, item := range items {
		for _, seed := range r.seeds {
			if hash, err := r.hashFun(seed, item); err != nil {
				return err
			} else {
				poss = append(poss, hash%r.m)
			}
		}
	}
	return r.bitmap.AddPos(r.key, r.expired, poss...)
}

func (r RedisBloomFilter) CheckOneExisted(item interface{}) (bool, error) {
	poss := make([]uint32, 0)
	for _, seed := range r.seeds {
		if hash, err := r.hashFun(seed, item); err != nil {
			return false, err
		} else {
			poss = append(poss, hash%r.m)
		}
	}
	contains, err := r.bitmap.Contains(r.key, poss...)
	if err != nil {
		return false, err
	}
	var existed = true
	for _, c := range contains {
		if c == false {
			existed = false
			break
		}
	}
	return existed, nil
}

func (r RedisBloomFilter) CheckExisted(items []interface{}) ([]bool, error) {
	exists := make([]bool, 0)
	poss := make([]uint32, 0)
	for _, item := range items {
		for _, seed := range r.seeds {
			if hash, err := r.hashFun(seed, item); err != nil {
				return nil, err
			} else {
				poss = append(poss, hash%r.m)
			}
		}
	}
	contains, err := r.bitmap.Contains(r.key, poss...)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(contains); i += len(r.seeds) {
		var existed = true
		for j := 0; j < len(r.seeds); j++ {
			if contains[i+j] == false {
				existed = false
				break
			}
		}
		exists = append(exists, existed)
	}
	return exists, nil
}

type RedisFilterFactory struct {
	client *redis.Client
	logger *logrus.Entry
}

func (r RedisFilterFactory) NewBitFilter(maxBit uint32) BitFilter {
	return NewRedisBitFilter(
		r.client,
		r.logger.WithFields(
			logrus.Fields{"BitFilter": maxBit},
		),
		maxBit,
	)
}

func (r RedisFilterFactory) calculateBitSize(n uint32, p float64) int {
	return int(-float64(n) * math.Log(p) / (math.Pow(math.Log(2), 2)))
}

func (r RedisFilterFactory) NewBloomFilter(params *BloomFilterParams) BloomFilter {
	if params.p <= 0 || params.p > 1 {
		panic("NewBloomFilterByCountAndRate param p is out of range")
	}
	bitSize := r.calculateBitSize(params.n, params.p)
	if bitSize > math.MaxUint32 {
		bitSize = math.MaxUint32
	}
	m := uint32(bitSize)
	bitmap := NewRedisBitFilter(
		r.client,
		r.logger.WithFields(
			logrus.Fields{"BloomFilter": params.key},
		),
		m,
	)
	return &RedisBloomFilter{
		key:    params.key,
		bitmap: bitmap,
		seeds:  params.seeds,
		m:      m,
		n:      params.n,
		p:      params.p,
	}
}

func NewFactory(client *redis.Client, logger *logrus.Entry) Factory {
	return &RedisFilterFactory{
		client: client,
		logger: logger,
	}
}
