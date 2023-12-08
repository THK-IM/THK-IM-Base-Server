package locker

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	lockCommand = `if redis.call("GET", KEYS[1]) == false then ` +
		`return redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2]) ` +
		`else ` +
		`return "FAILED" ` +
		`end`
	delCommand = `if redis.call("GET", KEYS[1]) == ARGV[1] then ` +
		`return redis.call("DEL", KEYS[1]) ` +
		`else ` +
		`return 0 ` +
		`end`
)

type redisLock struct {
	client    *redis.Client
	logger    *logrus.Entry
	key       string
	value     string
	waitMs    int
	timeoutMs int
}

func (r redisLock) Lock() (bool, error) {
	ctx := context.Background()
	wait := 0
	for {
		resp, err := r.client.Eval(ctx, lockCommand, []string{r.key}, []string{r.value, strconv.Itoa(r.timeoutMs)}).Result()
		if err != nil {
			return false, err
		}
		if resp != nil {
			reply, ok := resp.(string)
			if ok && reply == "OK" {
				return true, nil
			}
		}
		// retry
		waitMs := 50
		time.Sleep(time.Millisecond * time.Duration(waitMs))
		wait += waitMs
		if wait >= r.waitMs {
			break
		}
	}
	return false, nil
}

func (r redisLock) IsLocked() (bool, error) {
	ctx := context.Background()
	v, err := r.client.Get(ctx, r.key).Result()
	return strings.EqualFold(v, r.value), err
}

func (r redisLock) Release() (bool, error) {
	ctx := context.Background()
	v, err := r.client.Eval(ctx, delCommand, []string{r.key}, []string{r.value}).Int()
	return v == 1, err
}

type RedisLockerFactory struct {
	client *redis.Client
	logger *logrus.Entry
}

func (f RedisLockerFactory) NewLocker(key string, waitMs int, timeoutMs int) Locker {
	return &redisLock{
		client:    f.client,
		logger:    f.logger,
		key:       key,
		value:     uuid.New().String(),
		waitMs:    waitMs,
		timeoutMs: timeoutMs,
	}
}

func NewRedisLockerFactory(client *redis.Client, logger *logrus.Entry) Factory {
	return &RedisLockerFactory{
		client: client,
		logger: logger,
	}
}
