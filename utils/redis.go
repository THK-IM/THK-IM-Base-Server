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

const evalDelKeyByValue = `if redis.call("GET", KEYS[1]) == ARGV[1] then ` +
	`return redis.call("DEL", KEYS[1]) ` +
	`else ` +
	`return 0 ` +
	`end`

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

func DelKeyByValue(redisClient *redis.Client, key string, value interface{}) (bool, error) {
	resp, err := redisClient.Eval(context.Background(), evalDelKeyByValue, []string{key}, []interface{}{value}).Result()
	if err != nil {
		return false, err
	}
	if resp == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
