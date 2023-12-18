package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/server"
)

const ClaimsKey = "THK-Claims"

func Claims(appCtx *server.Context) gin.HandlerFunc {
	return func(context *gin.Context) {
		traceId := context.Request.Header.Get(dto.TraceId)
		lang := context.Request.Header.Get(dto.Language)
		if traceId == "" {
			traceId = uuid.New().String()
		}
		claims := dto.ThkClaims{}
		claims.PutValue(dto.TraceId, traceId)
		claims.PutValue(dto.Language, lang)
		context.Set(ClaimsKey, claims)
	}
}
