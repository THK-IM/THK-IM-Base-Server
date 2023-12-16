package dto

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/errorx"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ResponseForbidden(ctx *gin.Context) {
	rsp := &ErrorResponse{
		Code:    http.StatusForbidden,
		Message: "StatusForbidden",
	}
	ctx.JSON(http.StatusForbidden, rsp)
}

func ResponseUnauthorized(ctx *gin.Context) {
	rsp := &ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: "StatusUnauthorized",
	}
	ctx.JSON(http.StatusUnauthorized, rsp)
}

func ResponseBadRequest(ctx *gin.Context) {
	rsp := &ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "BadRequest",
	}
	ctx.JSON(http.StatusBadRequest, rsp)
}

func ResponseInternalServerError(ctx *gin.Context, err error) {
	var e *errorx.ErrorX
	if errors.As(err, &e) {
		if e.Code <= 5000000 {
			rsp := &ErrorResponse{
				Code:    e.Code,
				Message: e.Msg,
			}
			ctx.JSON(http.StatusBadRequest, rsp)
		} else {
			rsp := &ErrorResponse{
				Code:    e.Code,
				Message: e.Msg,
			}
			ctx.JSON(http.StatusInternalServerError, rsp)
		}
	}
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		ctx.Status(http.StatusOK)
	} else {
		ctx.JSON(http.StatusOK, data)
	}
}

func Redirect302(ctx *gin.Context, url string) {
	ctx.Redirect(302, url)
}
