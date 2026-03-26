package errorx

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var (
	ErrParamsError         = NewErrorX(4000000, "Params Error")
	ErrNotFound            = NewErrorX(4000001, "Not Found")
	ErrNotSupportReCommit  = NewErrorX(4000002, "Not support Recommit")
	ErrPermission          = NewErrorX(4000003, "Permission denied")
	ErrInternalServerError = NewErrorX(5000000, "Internal Server err")
	ErrServerBusy          = NewErrorX(5000001, "Server busy")
)

func translateReqParamsErr(err error, obj interface{}) string {
	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		val := reflect.ValueOf(obj).Elem()
		typ := val.Type()

		for _, e := range errs {
			field, _ := typ.FieldByName(e.Field())
			msg := field.Tag.Get("errMsg")
			if msg != "" {
				return msg
			}
		}
	}
	return err.Error()
}

func ErrReqParamsValidation(err error, obj interface{}) *ErrorX {
	msg := translateReqParamsErr(err, obj)
	return NewErrorX(4000000, msg)
}
