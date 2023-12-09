package conf

import (
	"flag"
	"fmt"
	"os"
)

var consulEndpoint = flag.String("config-consul-endpoint", "", "consul address")
var consulKey = flag.String("config-consul-key", "", "consul key")
var configFile = flag.String("config-file", "etc/server.yaml", "the config file")

func getConsul() (endpoint, key string) {
	if *consulEndpoint != "" && *consulKey != "" {
		return *consulEndpoint, *consulKey
	} else {
		return os.Getenv("config-consul-endpoint"), os.Getenv("config-consul-key")
	}
}

func LoadConfig() *Config {
	var (
		config Config
		err    error
	)
	cAddress, cKey := getConsul()
	if cAddress != "" && cKey != "" {
		config, err = LoadFromConsul(cAddress, cKey)
		if err != nil {
			panic(fmt.Sprintf("config read error: %v", err))
		}
	} else {
		config, err = Load(*configFile)
		if err != nil {
			panic(fmt.Sprintf("config read error: %v", err))
		}
	}
	return &config
}
