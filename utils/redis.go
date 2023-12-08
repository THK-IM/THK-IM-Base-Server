package utils

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

const evalBatchGet = "local result ={} " +
	"for i = 1, #(KEYS) do " +
	"result[i] = redis.call('get', KEYS[i]) " +
	"end " +
	"return result"

func BatchGet(redisClient *redis.Client, keys []string) ([]interface{}, error) {
	resp, errScript := redisClient.Eval(context.Background(), evalBatchGet, keys).Result()
	if errScript != nil {
		return nil, errScript
	}
	sliceResp, ok := resp.([]interface{})
	if ok {
		return sliceResp, nil
	} else {
		return nil, os.ErrInvalid
	}
}

func BatchGetString(redisClient *redis.Client, keys []string) ([]*string, error) {
	resp, errScript := redisClient.Eval(context.Background(), evalBatchGet, keys).Result()
	if errScript != nil {
		return nil, errScript
	}
	sliceResp, ok := resp.([]interface{})
	if ok {
		stringResp := make([]*string, 0)
		for _, p := range sliceResp {
			if p != nil {
				stringResp = append(stringResp, p.(*string))
			} else {
				stringResp = append(stringResp, nil)
			}
		}
		return stringResp, nil
	} else {
		return nil, os.ErrInvalid
	}
}
