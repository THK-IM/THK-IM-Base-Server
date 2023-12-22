package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thk-im/thk-im-base-server/dto"
	"strconv"
)

const ClaimsKey = "THK-Claims"

func Claims() gin.HandlerFunc {
	return func(context *gin.Context) {
		claims := dto.ThkClaims{}
		traceID := context.Request.Header.Get(dto.TraceID)
		if traceID == "" {
			traceID = uuid.New().String()
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

		clientIP := context.Request.Header.Get(dto.ClientOriginIP)
		if clientIP == "" {
			clientIP = context.ClientIP()
		}
		claims.PutValue(dto.ClientOriginIP, clientIP)

		clientPlatform := context.Request.Header.Get(dto.ClientPlatform)
		claims.PutValue(dto.ClientPlatform, clientPlatform)

		clientVersion := context.Request.Header.Get(dto.ClientVersion)
		claims.PutValue(dto.ClientVersion, clientVersion)

		lang := context.Request.Header.Get(dto.Language)
		claims.PutValue(dto.Language, lang)

		context.Set(ClaimsKey, claims)
	}
}
