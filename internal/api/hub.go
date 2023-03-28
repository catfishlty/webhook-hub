package api

import (
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type Hub struct {
	GIN    *gin.Engine
	DB     *gorm.DB
	Sender *Sender
}

func NewHub(db *gorm.DB) *Hub {
	sender := NewSender()
	return &Hub{
		GIN:    newRouter(db, sender),
		DB:     db,
		Sender: sender,
	}
}

func (hub *Hub) Migrate() {
	err := hub.DB.AutoMigrate(&types.Rule{}, &types.ReceiveRequest{}, &types.SendRequest{})
	if err != nil {
		panic(err)
	}
}

func newRouter(db *gorm.DB, sender *Sender) *gin.Engine {
	r := gin.New()
	r.Any("/webhooks/:id", webhookErrorHandler(), webhookHandler(db, sender))
	apiGroup := r.Group("/api", adminErrorHandler())
	{
		ruleGroup := apiGroup.Group("/rule")
		{
			ruleGroup.POST("", AddRuleHandler(db))
			ruleGroup.DELETE(":id", DeleteRuleHandler(db))
			ruleGroup.GET(":id", GetRuleHandler(db))
		}
		rulesGroup := apiGroup.Group("/rules")
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
