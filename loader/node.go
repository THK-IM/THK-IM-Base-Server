package loader

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/thk-im/thk-im-base-server/conf"
	"os"
	"strconv"
	"time"
)

var luaSetNodeId = "local value = redis.call('get', KEYS[1]) " +
	"if (value == false) then" +
	"   redis.call('setex', KEYS[1], ARGV[1], ARGV[2])" +
	"   return 0 " +
	"elseif (value == ARGV[2])  then" +
	"   redis.call('setex', KEYS[1], ARGV[1], ARGV[2])" +
	"   return 0 " +
	"else " +
	"   return 1 " +
	"end"

func LoadNodeId(config *conf.Config, client *redis.Client) (workerId int64, startTime int64) {
	nodeConfig := config.Node
	serverName := config.Name
	ctx := context.Background()
	startTime = time.Now().UnixMicro()
	nodeId := os.Getenv("NODE_ID")
	if nodeId != "" {
		if id, err := strconv.Atoi(nodeId); err == nil {
			workerId = int64(id)
			return
		}
	}
	workerId = 1
	expireTime := nodeConfig.PollingInterval
	for true {
		keys := []string{fmt.Sprintf("app/%s/%d", serverName, workerId)}
		result, err := client.Eval(ctx, luaSetNodeId, keys, expireTime+5, startTime).Result()
		if err != nil {
			panic(err)
		}
		if op, ok := result.(int64); ok {
			if op == 0 {
				break
			} else {
				workerId++
				if workerId >= nodeConfig.MaxCount {
					panic(errors.New("worker max number beyond"))
				}
			}
		} else {
			panic(errors.New("redis eval err"))
		}
	}
	keepNodeAlive(config, client, workerId, startTime)
	return
}

func keepNodeAlive(config *conf.Config, client *redis.Client, workerId, startTime int64) {
	interval := config.Node.PollingInterval
	serverName := config.Name
	ctx := context.Background()
	ticker := time.NewTicker(time.Second * time.Duration(interval))
	go func(n int64, value int64) {
		for range ticker.C {
			key := fmt.Sprintf("app/%s/%d", serverName, n)
			exp := time.Second * time.Duration(interval)
			_, err := client.SetEx(ctx, key, value, exp).Result()
			if err != nil {
				panic(err)
			}
		}
	}(workerId, startTime)
}
