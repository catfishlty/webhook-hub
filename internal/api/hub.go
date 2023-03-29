package api

import (
	"github.com/catfishlty/webhooks-hub/internal/common"
	"github.com/catfishlty/webhooks-hub/internal/data"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Hub struct {
	GIN    *gin.Engine
	DB     *gorm.DB
	Sender *Sender
}

func NewHub(db *gorm.DB, secretKey string) *Hub {
	sender := NewSender()
	migrate(db)
	return &Hub{
		GIN:    newRouter(db, secretKey, sender),
		DB:     db,
		Sender: sender,
	}
}

func (hub *Hub) Init() {
	var userCount int64
	hub.DB.Model(&types.User{}).Count(&userCount)
	if userCount == 0 {
		if err := data.CreateUser(hub.DB, common.DefaultUsername, common.DefaultPassword); err != nil {
			log.Fatal("failed to create default user", err)
		}
	}
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&types.Rule{}, &types.ReceiveRequest{}, &types.SendRequest{}, &types.User{})
	if err != nil {
		panic(err)
	}
}

func newRouter(db *gorm.DB, secretKey string, sender *Sender) *gin.Engine {
	authMidware := getAuthMiddleware(db, "webhook-hub", secretKey)
	r := gin.New()
	r.Any("webhooks/:id", webhookErrorHandler(), webhookHandler(db, sender))
	authGroup := r.Group("auth")
	{
		authGroup.POST("login", adminErrorHandler(), authMidware.LoginHandler)
		authGroup.GET("refresh", adminErrorHandler(), authMidware.RefreshHandler)
		authGroup.GET("logout", adminErrorHandler(), authMidware.LogoutHandler)
	}
	apiGroup := r.Group("api", authMidware.MiddlewareFunc(), adminErrorHandler())
	{
		ruleGroup := apiGroup.Group("rule")
		{
			ruleGroup.POST("", AddRuleHandler(db))
			ruleGroup.DELETE(":id", DeleteRuleHandler(db))
			ruleGroup.GET(":id", GetRuleHandler(db))
		}
		rulesGroup := apiGroup.Group("rules")
		{
			rulesGroup.GET(":page", GetRuleListHandler(db))
		}
	}
	return r
}

func (hub *Hub) Server(port int) *http.Server {
	return &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: hub.GIN,
	}
}
