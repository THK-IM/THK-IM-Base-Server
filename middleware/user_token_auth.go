package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/dto"
	"github.com/thk-im/thk-im-base-server/rpc"
	"github.com/thk-im/thk-im-base-server/server"
)

const (
	TokenKey = "Token"
	UidKey   = "Uid"
)

func UserTokenAuth(appCtx *server.Context) gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.Header.Get(TokenKey)
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
			context.Set(UidKey, res.UserId)
			context.Next()
		}
	}
}
