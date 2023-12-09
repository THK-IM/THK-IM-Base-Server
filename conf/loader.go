package conf

import (
	"flag"
	"fmt"
	"os"
)

var configConsulEndpoint = flag.String("config-consul-endpoint", "", "config consul address")
var configConsulKey = flag.String("config-consul-key", "", "config consul key")
var configPath = flag.String("config-path", "", "config file path")

func init() {
	flag.Parse()
}

func getConfigConsul() (endpoint, key string) {
	if *configConsulEndpoint != "" && *configConsulKey != "" {
		return *configConsulEndpoint, *configConsulKey
	} else {
		return os.Getenv("config-consul-endpoint"), os.Getenv("config-consul-key")
	}
}

func getConfigPath() string {
	if *configPath != "" {
		return *configPath
	} else {
		return os.Getenv("config-path")
	}
}

func LoadConfig(path string) *Config {
	var (
		config Config
		err    error
	)
	consulAddress, consulKey := getConfigConsul()
	if consulAddress != "" && consulKey != "" {
		config, err = LoadFromConsul(consulAddress, consulKey)
		if err != nil {
			panic(fmt.Sprintf("config read error: %v", err))
		}
		return &config
	}

	confPath := getConfigPath()
	if confPath != "" {
		config, err = Load(confPath)
		if err != nil {
			panic(fmt.Sprintf("config read error: %v", err))
		}
		return &config
	}

	if path != "" {
		config, err = Load(path)
		if err != nil {
			panic(fmt.Sprintf("config read error: %v", err))
		}
		return &config
	}
	panic("load config error")
}
