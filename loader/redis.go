package loader

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Base-Server/conf"
	"github.com/redis/go-redis/v9"
	"time"
)

func LoadRedis(source *conf.RedisSource) *redis.Client {
	if source == nil {
		return nil
	}
	opt, err := redis.ParseURL(fmt.Sprintf("%s%s", source.Endpoint, source.Uri))
	if err != nil {
		panic(err)
	}
	opt.ConnMaxLifetime = time.Duration(source.ConnMaxLifeTime) * time.Second
	opt.ConnMaxIdleTime = time.Duration(source.ConnMaxIdleTime) * time.Second
	opt.ReadTimeout = time.Duration(source.ConnTimeout) * time.Second
	opt.WriteTimeout = time.Duration(source.ConnTimeout) * time.Second
	opt.PoolTimeout = time.Duration(source.ConnTimeout) * time.Second
	opt.MaxIdleConns = source.MaxIdleConn
	opt.PoolSize = source.MaxOpenConn

	rdb := redis.NewClient(opt)
	return rdb
}
