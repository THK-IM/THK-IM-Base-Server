package errorx

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

type ErrorX struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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
		err.Message = fmt.Sprintf("http status code %d", resp.StatusCode())
	}
	return err
}

func NewErrorX(code int, message string) *ErrorX {
	return &ErrorX{Code: code, Message: message}
}

func New(message string) *ErrorX {
	return &ErrorX{Code: 0, Message: message}
}

func (e *ErrorX) Error() string {
	return fmt.Sprintf("[code: %d, message: %s]", e.Code, e.Message)
}
