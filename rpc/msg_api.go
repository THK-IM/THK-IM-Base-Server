package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/THK-IM/THK-IM-Base-Server/conf"
	"github.com/THK-IM/THK-IM-Base-Server/dto"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

const (
	msgApiPostUserOnlineStatusUrl = "/system/user/online"
)

type (
	MsgApi interface {
		PostUserOnlineStatus(req dto.PostUserOnlineReq) error
	}

	defaultMsgApi struct {
		endpoint string
		logger   *logrus.Entry
		client   *resty.Client
	}
)

func (d defaultMsgApi) PostUserOnlineStatus(req dto.PostUserOnlineReq) error {
	dataBytes, err := json.Marshal(req)
	if err != nil {
		d.logger.Error(err)
		return err
	}
	url := fmt.Sprintf("%s%s", d.endpoint, msgApiPostUserOnlineStatusUrl)
	res, errRequest := d.client.R().
		SetHeader("Content-Type", contentType).
		SetBody(dataBytes).
		Post(url)
	if errRequest != nil {
		return errRequest
	}
	if res.StatusCode() != http.StatusOK {
		e := errors.New(string(res.Body()))
		d.logger.Error(e)
		return e
	} else {
		return nil
	}
}

func NewMsgApi(sdk conf.Sdk, logger *logrus.Entry) MsgApi {
	return defaultMsgApi{
		endpoint: sdk.Endpoint,
		logger:   logger.WithField("rpc", sdk.Name),
		client: resty.New().
			SetTransport(&http.Transport{
				MaxIdleConns:    10,
				MaxConnsPerHost: 10,
				IdleConnTimeout: 30 * time.Second,
			}).
			SetTimeout(5 * time.Second).
			SetRetryCount(3).
			SetRetryWaitTime(15 * time.Second).
			SetRetryMaxWaitTime(5 * time.Second),
	}
}
