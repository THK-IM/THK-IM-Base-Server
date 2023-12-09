package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/server"
	"net"
	"strings"
)

func WhiteIpAuth(appCtx *server.Context) gin.HandlerFunc {
	ipWhiteList := appCtx.Config().IpWhiteList
	ips := strings.Split(ipWhiteList, ",")
	return func(context *gin.Context) {
		ip := context.ClientIP()
		appCtx.Logger().Infof("RemoteAddr: %s", ip)
		if isIpValid(ip, ips) {
			dto.ResponseForbidden(context)
			context.Abort()
		} else {
			context.Next()
		}
	}
}

func isIpValid(clientIp string, whiteIpList []string) bool {
	ip := net.ParseIP(clientIp)
	for _, whiteIp := range whiteIpList {
		_, ipNet, err := net.ParseCIDR(whiteIp)
		if err != nil {
			return false
		}
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}
