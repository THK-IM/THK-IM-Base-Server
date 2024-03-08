package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/server"
	"net"
	"strings"
)

func WhiteIpAuth(appCtx *server.Context) gin.HandlerFunc {
	ipWhiteList := appCtx.Config().IpWhiteList
	ips := strings.Split(strings.ReplaceAll(ipWhiteList, " ", ""), ",")
	return func(context *gin.Context) {
		ip := context.ClientIP()
		claims := context.MustGet(ClaimsKey).(dto.ThkClaims)
		if !isIpValid(ip, ips) {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Errorf("RemoteAddr forbidden: %s %v", ip, ips)
			dto.ResponseForbidden(context)
			context.Abort()
		} else {
			appCtx.Logger().WithFields(logrus.Fields(claims)).Infof("RemoteAddr: %s", ip)
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
