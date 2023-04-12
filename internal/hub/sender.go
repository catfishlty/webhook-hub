package hub

import (
	"errors"
	"github.com/catfishlty/webhook-hub/internal/types"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Sender struct {
	resty *resty.Client
}

func (s *Sender) Send(send types.RestySendRequest) (*resty.Response, error) {
	r := s.resty.R().
		SetHeaders(send.Header).
		SetQueryParams(send.Query)
	if send.IsForm {
		r.SetFormData(send.Form)
	} else {
		r.SetBody(send.Body)
	}
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
