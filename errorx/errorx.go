package errorx

import "fmt"

type ErrorX struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewErrorX(code int, msg string) *ErrorX {
	return &ErrorX{Code: code, Msg: msg}
}

func New(msg string) *ErrorX {
	return &ErrorX{Code: 0, Msg: msg}
}

func (e *ErrorX) Error() string {
	return fmt.Sprintf("[code: %d, msg: %s]", e.Code, e.Msg)
}
