package main

import (
	"context"
	"fmt"
	"github.com/thk-im/thk-im-base-server/conf"
	"github.com/thk-im/thk-im-base-server/server"
	"github.com/thk-im/thk-im-base-server/utils"
	"os"
	"time"
)

func main() {
	dir, errWd := os.Getwd()
	if errWd != nil {
		panic(errWd)
	}
	fmt.Println(dir)

	config, err := conf.Load(dir + "/etc/server.yaml")
	if err != nil {
		panic(err)
	}

	srvContext := &server.Context{}
	srvContext.Init(&config)

	keys := make([]string, 0)
	for i := 0; i < 10000; i++ {
		keys = append(keys, utils.GetRandomString(12))
	}
	resp, errGet := utils.BatchGetString(srvContext.RedisCache(), keys)
	if errGet != nil {
		fmt.Println(errGet)
	} else {
		fmt.Println("len:", len(resp))
	}

	errSet := srvContext.RedisCache().Set(context.Background(), "22", 11, time.Hour).Err()
	if errSet != nil {
		fmt.Println(errSet)
	}

	ok, errDel := utils.DelKeyByValue(srvContext.RedisCache(), "22", 11)
	if errDel != nil {
		fmt.Println(errDel)
	} else {
		fmt.Println("ok: ", ok)
	}
}
