package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/rpc"
	"github.com/thk-im/thk-im-base-server/server"
)

const (
	tokenKey = "Token"
	uidKey   = "Uid"
)

func UserTokenAuth(appCtx *server.Context) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get(tokenKey)
		if token == "" {
			appCtx.Logger().Warn("token nil error")
			dto.ResponseUnauthorized(context)
			context.Abort()
			return
		}
		req := rpc.GetUserIdByTokenReq{Token: token}
		res, err := appCtx.RpcUserApi().GetUserIdByToken(req)
		if err != nil {
			appCtx.Logger().Warn("token error")
			dto.ResponseUnauthorized(context)
			context.Abort()
		} else {
			context.Set(uidKey, res.UserId)
			context.Next()
		}
	}
}
