package conf

import (
	"errors"
	"flag"
	"os"
)

var configConsulEndpoint = flag.String("config-consul-endpoint", "", "config consul address")
var configConsulKey = flag.String("config-consul-key", "", "config consul key")
var configPath = flag.String("config-path", "", "config file path")

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

func LoadConfig(path string, config interface{}) error {
	flag.Parse()
	consulAddress, consulKey := getConfigConsul()
	if consulAddress != "" && consulKey != "" {
		return LoadFromConsul(consulAddress, consulKey, config)
	}

	confPath := getConfigPath()
	if confPath != "" {
		return Load(confPath, config)
	}

	if path != "" {
		return Load(path, config)
	}
	return errors.New("config not init")
}
