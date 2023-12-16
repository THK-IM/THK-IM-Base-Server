package rpc

//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"github.com/go-resty/resty/v2"
//	"github.com/sirupsen/logrus"
//	"github.com/thk-im/thk-im-base-server/conf"
//	"net/http"
//	"time"
//)
//
//const (
//	userApiGetUserIdByTokenUrl string = "/user"
//	userApiOnlineStatusUrl     string = "/user/online"
//	contentType                string = "application/json"
//)
//
//type (
//	PostUserOnlineReq struct {
//		UserId    int64  `json:"user_id"`
//		IsOnline  bool   `json:"is_online"`
//		Timestamp int64  `json:"timestamp"`
//		ConnId    int64  `json:"conn_id"`
//		Platform  string `json:"platform"`
//	}
//
//	GetUserIdByTokenReq struct {
//		Platform string `json:"platform"`
//		Token    string `json:"token"`
//	}
//
//	GetUserIdByTokenRes struct {
//		UserId int64 `json:"user_id"`
//	}
//
//	UserApi interface {
//		PostUserOnlineStatus(req PostUserOnlineReq) error
//		GetUserIdByToken(req GetUserIdByTokenReq) (*GetUserIdByTokenRes, error)
//	}
//
//	defaultUserApi struct {
//		endpoint string
//		logger   *logrus.Entry
//		client   *resty.Client
//	}
//)
//
//func (d defaultUserApi) PostUserOnlineStatus(req PostUserOnlineReq) error {
//	dataBytes, err := json.Marshal(req)
//	if err != nil {
//		d.logger.Error(err)
//		return err
//	}
//	url := fmt.Sprintf("%s%s", d.endpoint, userApiOnlineStatusUrl)
//	res, errRequest := d.client.R().
//		SetHeader("Content-Type", contentType).
//		SetBody(dataBytes).
//		Post(url)
//	if errRequest != nil {
//		return errRequest
//	}
//	if res.StatusCode() != http.StatusOK {
//		e := errors.New(string(res.Body()))
//		d.logger.Error(e)
//		return e
//	} else {
//		return nil
//	}
//}
//
//func (d defaultUserApi) GetUserIdByToken(req GetUserIdByTokenReq) (*GetUserIdByTokenRes, error) {
//	url := fmt.Sprintf("%s%s?token=%s&platform=%s", d.endpoint, userApiGetUserIdByTokenUrl, req.Token, req.Platform)
//	res, errRequest := d.client.R().
//		SetHeader("Content-Type", contentType).
//		Get(url)
//	if errRequest != nil {
//		return nil, errRequest
//	}
//	if res.StatusCode() != http.StatusOK {
//		e := errors.New(string(res.Body()))
//		d.logger.Error(e)
//		return nil, e
//	} else {
//		tokenRes := &GetUserIdByTokenRes{}
//		e := json.Unmarshal(res.Body(), tokenRes)
//		if e != nil {
//			d.logger.Error(e)
//			return nil, e
//		} else {
//			return tokenRes, nil
//		}
//	}
//}
//
//func NewUserApi(sdk conf.Sdk, logger *logrus.Entry) UserApi {
//	return defaultUserApi{
//		endpoint: sdk.Endpoint,
//		logger:   logger.WithField("rpc", sdk.Name),
//		client: resty.New().
//			SetTransport(&http.Transport{
//				MaxIdleConns:    10,
//				MaxConnsPerHost: 10,
//				IdleConnTimeout: 30 * time.Second,
//			}).
//			SetTimeout(5 * time.Second).
//			SetRetryCount(3).
//			SetRetryWaitTime(15 * time.Second).
//			SetRetryMaxWaitTime(5 * time.Second),
//	}
//}
