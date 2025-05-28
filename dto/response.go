package dto

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/thk-im/thk-im-base-server/errorx"
	"github.com/thk-im/thk-im-base-server/i18n"
	"net/http"
)

var Localize i18n.Localize

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (eR *ErrorResponse) Localize(language string) {
	if Localize != nil {
		eR.Message = Localize.Get(eR.Message, language)
	}
}

func ResponseStatusCode(ctx *gin.Context, httpStatusCode int, err error) {
	if err != nil {
		var e *errorx.ErrorX
		if errors.As(err, &e) {
			rsp := &ErrorResponse{
				Code:    e.Code,
				Message: e.Message,
			}
			ctx.JSON(httpStatusCode, rsp)
		} else {
			rsp := &ErrorResponse{
				Code:    httpStatusCode,
				Message: err.Error(),
			}
			ctx.JSON(httpStatusCode, rsp)
		}
	} else {
		ctx.Status(httpStatusCode)
	}

}

func ResponseForbidden(ctx *gin.Context) {
	rsp := &ErrorResponse{
		Code:    http.StatusForbidden,
		Message: "StatusForbidden",
	}
	claims := ctx.MustGet(ClaimsKey).(ThkClaims)
	rsp.Localize(claims.GetLanguage())
	ctx.JSON(http.StatusForbidden, rsp)
}

func ResponseUnauthorized(ctx *gin.Context) {
	rsp := &ErrorResponse{
		Code:    http.StatusUnauthorized,
		Message: "StatusUnauthorized",
	}
	claims := ctx.MustGet(ClaimsKey).(ThkClaims)
	rsp.Localize(claims.GetLanguage())
	ctx.JSON(http.StatusUnauthorized, rsp)
}

func ResponseBadRequest(ctx *gin.Context) {
	rsp := &ErrorResponse{
		Code:    http.StatusBadRequest,
		Message: "BadRequest",
	}
	claims := ctx.MustGet(ClaimsKey).(ThkClaims)
	rsp.Localize(claims.GetLanguage())
	ctx.JSON(http.StatusBadRequest, rsp)
}

func ResponseInternalServerError(ctx *gin.Context, err error) {
	var e *errorx.ErrorX
	if errors.As(err, &e) {
		if e.Code <= 5000000 {
			rsp := &ErrorResponse{
				Code:    e.Code,
				Message: e.Message,
			}
			claims := ctx.MustGet(ClaimsKey).(ThkClaims)
			rsp.Localize(claims.GetLanguage())
			ctx.JSON(http.StatusBadRequest, rsp)
		} else {
			rsp := &ErrorResponse{
				Code:    e.Code,
				Message: e.Message,
			}
			claims := ctx.MustGet(ClaimsKey).(ThkClaims)
			rsp.Localize(claims.GetLanguage())
			ctx.JSON(http.StatusInternalServerError, rsp)
		}
	} else {
		e = errorx.ErrInternalServerError
		rsp := &ErrorResponse{
			Code:    e.Code,
			Message: e.Message,
		}
		claims := ctx.MustGet(ClaimsKey).(ThkClaims)
		rsp.Localize(claims.GetLanguage())
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
