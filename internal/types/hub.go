package types

import (
	"gorm.io/datatypes"
)

type ReceiveRequest struct {
	ID        string            `json:"id"`
	Method    string            `json:"method"`
	Variables datatypes.JSONMap `json:"variables"`
	Validate  datatypes.JSONMap `json:"validate"`
}

type SendRequest struct {
	ID     string            `json:"id"`
	Url    string            `json:"url"`
	Method string            `json:"method"`
	Header datatypes.JSONMap `json:"header"`
	Query  datatypes.JSONMap `json:"query"`
	Body   string            `json:"body"`
}

type Rule struct {
	RuleItem
	ReceiveId string         `json:"-"`
	Receive   ReceiveRequest `json:"receive" gorm:"foreignkey:ReceiveId;references:ID"`
	SendId    string         `json:"-"`
	Send      SendRequest    `json:"send" gorm:"foreignkey:SendId;references:ID"`
}

type RuleItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	GroupId     string `json:"groupId"`
	IsAuth      bool   `json:"isAuth"`
}

type User struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"-,"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
}
