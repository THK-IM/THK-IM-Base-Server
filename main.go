package main

import (
	"fmt"
	"github.com/THK-IM/THK-IM-Base-Server/conf"
	"github.com/THK-IM/THK-IM-Base-Server/server"
	"github.com/THK-IM/THK-IM-Base-Server/utils"
	"os"
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
}
