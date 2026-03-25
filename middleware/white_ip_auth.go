package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/thk-im/thk-im-base-server/dto"
)

func WhiteIpAuth(ipWhiteList string, logger *logrus.Entry) gin.HandlerFunc {
	ips := strings.Split(strings.ReplaceAll(ipWhiteList, " ", ""), ",")
	return func(context *gin.Context) {
		ip := context.ClientIP()
		claims := context.MustGet(dto.ClaimsKey).(dto.ThkClaims)
		if !isIpValid(ip, ips) {
			logger.WithFields(logrus.Fields(claims)).Errorf("RemoteAddr forbidden: %s %v", ip, ips)
			dto.ResponseForbidden(context)
			context.Abort()
		} else {
			logger.WithFields(logrus.Fields(claims)).Tracef("RemoteAddr: %s", ip)
			context.Next()
		}
	}
}

func isIpValid(clientIp string, whiteIpList []string) bool {
	ip := net.ParseIP(clientIp)
	if ip == nil {
		return false
	}

	for _, whiteIp := range whiteIpList {

		// 1️⃣ 先尝试 CIDR
		_, ipNet, err := net.ParseCIDR(whiteIp)
		if err == nil {
			if ipNet.Contains(ip) {
				return true
			}
			continue
		}

		// 2️⃣ 再尝试单 IP
		if ip.Equal(net.ParseIP(whiteIp)) {
			return true
		}
	}

	return false
}
