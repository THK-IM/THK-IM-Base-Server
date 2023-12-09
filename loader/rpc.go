package loader

import (
	"github.com/THK-IM/THK-IM-Server/pkg/conf"
	"github.com/THK-IM/THK-IM-Server/pkg/rpc"
	"github.com/sirupsen/logrus"
)

func LoadSdks(sdkConfigs []conf.Sdk, logger *logrus.Entry) map[string]interface{} {
	sdkMap := make(map[string]interface{}, 0)
	for _, c := range sdkConfigs {
		if c.Name == "user-api" {
			userApi := rpc.NewUserApi(c, logger)
			sdkMap[c.Name] = userApi
		} else if c.Name == "msg-api" {
			userApi := rpc.NewMsgApi(c, logger)
			sdkMap[c.Name] = userApi
		}
	}
	return sdkMap
}
