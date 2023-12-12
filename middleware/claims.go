package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/server"
	"github.com/thk-im/thk-im-base-server/utils"
)

const ClaimsKey = "claims"

func Claims(appCtx *server.Context) gin.HandlerFunc {
	return func(context *gin.Context) {
		traceId := context.Request.Header.Get(dto.TraceId)
		lang := context.Request.Header.Get("Accept-Language")
		if traceId == "" {
			traceId = utils.GetRandomString(10)
		}
		claims := dto.MapClaims{}
		claims.PutValue(dto.TraceId, traceId)
		claims.PutValue(dto.Language, lang)
		context.Set(ClaimsKey, claims)
	}
}
