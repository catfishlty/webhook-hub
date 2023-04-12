package types

import (
	"github.com/catfishlty/webhook-hub/internal/check"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ReceiveRequest struct {
	Method    string            `json:"method"`
	Variables datatypes.JSONMap `json:"variables"`
}

type VariableItem struct {
	Key      string `json:"key"`
	Assign   string `json:"assign"`
	Validate string `json:"validate"`
	Value    string `json:"-"`
	Type     string `json:"-"`
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

type SendBase struct {
	Url    string `json:"url"`
	Method string `json:"method"`
	IsForm bool   `json:"isForm"`
}

type SendRequest struct {
	SendBase
	Header datatypes.JSONMap `json:"header"`
	Query  datatypes.JSONMap `json:"query"`
	Form   datatypes.JSONMap `json:"form"`
	Body   datatypes.JSONMap `json:"body"`
}

func (r SendRequest) ToResty() RestySendRequest {
	return RestySendRequest{
		SendBase: r.SendBase,
		Header:   JsonToStringMap(r.Header),
		Query:    JsonToStringMap(r.Query),
		Form:     JsonToStringMap(r.Form),
		Body:     r.Body,
	}
}

func JsonToStringMap(data datatypes.JSONMap) map[string]string {
	m := make(map[string]string)
	for k, v := range data {
		m[k] = v.(string)
	}
	return m
}

type RestySendRequest struct {
	SendBase
	Header map[string]string `json:"header"`
	Query  map[string]string `json:"query"`
	Form   map[string]string `json:"form"`
	Body   map[string]any    `json:"body"`
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
	Receive ReceiveRequest `json:"receive" gorm:"embedded;embeddedPrefix:receive_"`
	Send    SendRequest    `json:"send" gorm:"embedded;embeddedPrefix:send_"`
}

func (r Rule) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.RuleItem, validation.Required, validation.By(r.RuleItem.ValidateFunc())),
		validation.Field(&r.Receive, validation.Required, validation.By(r.Receive.ValidateFunc())),
		validation.Field(&r.Send, validation.Required, validation.By(r.Send.ValidateFunc())),
	)
}

type RuleItem struct {
	gorm.Model
	UID         string `json:"uid" gorm:"unique,not null"`
	Name        string `json:"name" gorm:"unique,not null"`
	Description string `json:"description"`
	GroupId     string `json:"groupId"`
	IsAuth      bool   `json:"isAuth"`
	IsForward   bool   `json:"isForward"`
}

func (r RuleItem) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.UID, is.UUIDv4),
		validation.Field(&r.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.IsAuth, validation.Required, validation.In(true, false)),
	)
}

func (r RuleItem) ValidateFunc() func(any) error {
	return func(any) error {
		return r.Validate()
	}
}

type User struct {
	gorm.Model
	UID       string `json:"uid" gorm:"unique,not null"`
	Username  string `json:"username" gorm:"unique,not null"`
	Password  string `json:"-,"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}
