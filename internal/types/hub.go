package types

import (
	"github.com/catfishlty/webhooks-hub/internal/check"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/datatypes"
)

type ReceiveRequest struct {
	ID        string            `json:"id"`
	Method    string            `json:"method"`
	Variables datatypes.JSONMap `json:"variables"`
	Validator datatypes.JSONMap `json:"Validator"`
}

func (r ReceiveRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Method, validation.Required, check.IsHttpMethod),
	)
}

func (r ReceiveRequest) ValidateFunc() func(interface{}) error {
	return func(interface{}) error {
		return r.Validate()
	}
}

type SendRequest struct {
	ID     string            `json:"id"`
	Url    string            `json:"url"`
	Method string            `json:"method"`
	Header datatypes.JSONMap `json:"header"`
	Query  datatypes.JSONMap `json:"query"`
	Body   string            `json:"body"`
}

func (r SendRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Url, validation.Required, is.RequestURL),
		validation.Field(&r.Method, validation.Required, check.IsHttpMethod),
		validation.Field(&r.Header, validation.Required),
	)
}

func (r SendRequest) ValidateFunc() func(interface{}) error {
	return func(interface{}) error {
		return r.Validate()
	}
}

type Rule struct {
	RuleItem
	ReceiveId string         `json:"-"`
	Receive   ReceiveRequest `json:"receive" gorm:"foreignkey:ReceiveId;references:ID"`
	SendId    string         `json:"-"`
	Send      SendRequest    `json:"send" gorm:"foreignkey:SendId;references:ID"`
}

func (r Rule) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RuleItem, validation.Required, validation.By(r.RuleItem.ValidateFunc())),
		validation.Field(&r.Receive, validation.Required, validation.By(r.Receive.ValidateFunc())),
		validation.Field(&r.Send, validation.Required, validation.By(r.Send.ValidateFunc())),
	)
}

type RuleItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupId     string `json:"groupId"`
	IsAuth      bool   `json:"isAuth"`
}

func (r RuleItem) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required, is.UUIDv4),
		validation.Field(&r.Name, validation.Required, validation.Min(1)),
		validation.Field(&r.IsAuth, validation.Required, validation.In(true, false)),
	)
}

func (r RuleItem) ValidateFunc() func(interface{}) error {
	return func(interface{}) error {
		return r.Validate()
	}
}

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"-,"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}
