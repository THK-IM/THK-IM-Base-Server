package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thk-im/thk-im-base-server/dto"
	"strconv"
	"strings"
)

const ClaimsKey = "Claims"

func Claims() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims := dto.ThkClaims{}
		traceID := context.Request.Header.Get(dto.TraceID)
		if strings.EqualFold(traceID, "") {
			traceID = uuid.NewString()
		}
		claims.PutValue(dto.TraceID, traceID)

		parentSpanID := context.Request.Header.Get(dto.SpanID)
		spanID := ""
		if parentSpanID == "" {
			parentSpanID = "0"
			spanID = "1"
		} else {
			i, err := strconv.Atoi(parentSpanID)
			if err == nil {
				spanID = fmt.Sprintf("%d", i+1)
			} else {
				spanID = "1"
			}
		}
		claims.PutValue(dto.ParentSpanID, parentSpanID)
		claims.PutValue(dto.SpanID, spanID)

		clientIP := context.Request.Header.Get(dto.OriginIP)
		if clientIP == "" {
			clientIP = context.ClientIP()
		}
		claims.PutValue(dto.OriginIP, clientIP)

		device := context.Request.Header.Get(dto.Device)
		claims.PutValue(dto.Device, device)

		timeZone := context.Request.Header.Get(dto.TimeZone)
		claims.PutValue(dto.TimeZone, timeZone)

		platform := context.Request.Header.Get(dto.Platform)
		claims.PutValue(dto.Platform, platform)

		version := context.Request.Header.Get(dto.Version)
		claims.PutValue(dto.Version, version)

		language := context.Request.Header.Get(dto.Language)
		claims.PutValue(dto.Language, language)

		token := context.Request.Header.Get(dto.JwtToken)
		token = strings.ReplaceAll(token, "Bearer ", "")
		token = strings.ReplaceAll(token, " ", "")
		claims.PutValue(dto.JwtToken, token)

		context.Set(ClaimsKey, claims)
		context.Next()
	}
}
