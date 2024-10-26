package filter

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"
)

type (
	RedisBitFilter struct {
		client *redis.Client
		logger *logrus.Entry
		maxBit uint32
	}
)

func (r RedisBitFilter) Init(key string, ex time.Duration) error {
	_, err := r.client.SetNX(context.Background(), key, "", ex).Result()
	return err
}

func (r RedisBitFilter) AddPos(key string, ex time.Duration, pos ...uint32) error {
	args := make([]interface{}, 0)
	args = append(args, "OVERFLOW")
	args = append(args, "FAIL")
	for _, p := range pos {
		args = append(args, "set")
		args = append(args, "u1")
		args = append(args, fmt.Sprintf("#%d", p%r.maxBit))
		args = append(args, "1")
	}
	_, err := r.client.BitField(context.Background(), key, args...).Result()
	if err != nil {
		return err
	}
	_, err = r.client.Expire(context.Background(), key, ex).Result()
	return err
}

func (r RedisBitFilter) Clear(key string, ex time.Duration) error {
	_, err := r.client.Set(context.Background(), key, "", ex).Result()
	return err
}

func (r RedisBitFilter) Contains(key string, pos ...uint32) ([]bool, error) {
	args := make([]interface{}, 0)
	args = append(args, "OVERFLOW")
	args = append(args, "FAIL")
	for _, p := range pos {
		args = append(args, "get")
		args = append(args, "u1")
		args = append(args, fmt.Sprintf("#%d", p%r.maxBit))
	}
	res, err := r.client.BitField(context.Background(), key, args).Result()
	if err == nil {
		contains := make([]bool, 0)
		for _, v := range res {
			if v > 0 {
				contains = append(contains, true)
			} else {
				contains = append(contains, false)
			}
		}
		return contains, nil
	}
	return nil, err
}

func (r RedisBitFilter) AllPos(key string) ([]uint32, error) {
	res, err := r.client.Get(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}
	poss := make([]uint32, 0)
	bs := []byte(res)
	for i, b := range bs {
		for j := 7; j >= 0; j-- {
			if b&(1<<uint(j)) > 0 {
				pos := uint32(i*8) + 7 - uint32(j)
				poss = append(poss, pos)
			}
		}
	}
	return poss, nil
}

func (r RedisBitFilter) Count(key string) (uint32, error) {
	bitCount := &redis.BitCount{
		Start: 0,
		End:   -1,
	}
	res, err := r.client.BitCount(context.Background(), key, bitCount).Result()
	if err == nil {
		uRes := uint32(res)
		return uRes, nil
	}
	return 0, err
}

func (r RedisBitFilter) Delete(key string) error {
	_, err := r.client.Del(context.Background(), key).Result()
	return err
}

func NewRedisBitFilter(client *redis.Client, logger *logrus.Entry, maxBit uint32) BitFilter {
	redisBitmapFilter := &RedisBitFilter{client: client, logger: logger, maxBit: maxBit}
	return redisBitmapFilter
}
