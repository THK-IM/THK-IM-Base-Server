package dto

import (
	"github.com/THK-IM/THK-IM-Server/pkg/errorx"
	"github.com/gin-gonic/gin"
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
	if e, ok := err.(*errorx.ErrorX); ok {
		if e.Code <= 500000 {
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
	} else {
		rsp := &ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}
		ctx.JSON(http.StatusInternalServerError, rsp)
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
