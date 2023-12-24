package errorx

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

type ErrorX struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func NewErrorXFromResp(resp *resty.Response) *ErrorX {
	err := &ErrorX{}
	if resp.StatusCode() == http.StatusBadRequest ||
		resp.StatusCode() == http.StatusInternalServerError {
		contentType := resp.Header().Get("Content-Type")
		contentType = strings.ToLower(contentType)
		bytesBody := resp.Body()
		if strings.Contains(contentType, "application/json") && len(bytesBody) > 0 {
			_ = json.Unmarshal(bytesBody, err)
		}
	}
	if err.Code == 0 {
		err.Code = resp.StatusCode()
		err.Msg = fmt.Sprintf("http status code %d", resp.StatusCode())
	}
	return err
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
