package hub

import (
	"github.com/catfishlty/webhook-hub/internal/common"
	"github.com/catfishlty/webhook-hub/internal/data"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type Hub struct {
	router *gin.Engine
	db     *data.DB
	sender *Sender
}

func NewHub(orm *gorm.DB, secretKey, salt string) *Hub {
	hub := &Hub{
		db:     data.NewDB(orm, salt),
		sender: NewSender(),
	}
	hub.db.Migrate()
	hub.initRouter(secretKey)
	hub.db.Init()
	return hub
}

func (hub *Hub) initRouter(secretKey string) {
	authMiddleware := hub.getAuthMiddleware(secretKey, common.JwtRealm)
	r := gin.New()
	r.Any("webhooks/:id", hub.webhookErrorHandler(), hub.webhookHandler())
	authGroup := r.Group("auth")
	{
		authGroup.POST("login", hub.adminErrorHandler(), authMiddleware.LoginHandler)
		authGroup.GET("refresh", hub.adminErrorHandler(), authMiddleware.RefreshHandler)
		authGroup.GET("logout", hub.adminErrorHandler(), authMiddleware.LogoutHandler)
	}
	apiGroup := r.Group("hub", authMiddleware.MiddlewareFunc(), hub.adminErrorHandler())
	{
		userGroup := apiGroup.Group("user")
		{
			userGroup.POST("", hub.AddUserHandler())
			userGroup.DELETE(":id", hub.DeleteUserHandler())
			userGroup.GET(":id", hub.GetUserHandler())
			userGroup.PUT(":id", hub.UpdateUserHandler())
			userGroup.PUT(":id/reset", hub.UpdateUserResetPasswordHandler())
		}
		usersGroup := apiGroup.Group("users")
		{
			usersGroup.GET(":page", hub.GetUserListHandler())
		}
		ruleGroup := apiGroup.Group("rule")
		{
			ruleGroup.POST("", hub.AddRuleHandler())
			ruleGroup.PUT(":id", hub.UpdateRuleHandler())
			ruleGroup.DELETE(":id", hub.DeleteRuleHandler())
			ruleGroup.GET(":id", hub.GetRuleHandler())
		}
		rulesGroup := apiGroup.Group("rules")
		{
			rulesGroup.GET(":page", hub.GetRuleListHandler())
		}
	}
	hub.router = r
}

func (hub *Hub) Server(port string) *http.Server {
	return &http.Server{
		Addr:    ":" + port,
		Handler: hub.router,
	}
}
