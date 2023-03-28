package api

import (
	"errors"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/catfishlty/webhooks-hub/internal/utils"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Sender struct {
	resty *resty.Client
}

func (s *Sender) Send(send types.SendRequest) (*resty.Response, error) {
	r := s.resty.R().SetHeaders(utils.JsonToStringMap(send.Header)).SetQueryParams(utils.JsonToStringMap(send.Query)).SetBody(send.Body)
	switch send.Method {
	case http.MethodGet:
		return r.Get(send.Url)
	case http.MethodPost:
		return r.Post(send.Url)
	case http.MethodPut:
		return r.Put(send.Url)
	case http.MethodDelete:
		return r.Delete(send.Url)
	default:
		return nil, errors.New("unknown method " + send.Method)
	}
}

func NewSender() *Sender {
	return &Sender{
		resty: resty.New(),
	}
}
