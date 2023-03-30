package hub

import (
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (hub *Hub) webhookErrorHandler() gin.HandlerFunc {
	return hub.commonErrorHandler("Webhook")
}

func (hub *Hub) adminErrorHandler() gin.HandlerFunc {
	return hub.commonErrorHandler("Admin")
}

func (*Hub) commonErrorHandler(errorType string) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if customError := recovered.(*types.CommonError); customError != nil {
			msg := customError.Msg
			if msg == "" && customError.Err != nil {
				msg = customError.Err.Error()
			} else {
				msg = "Unknown error"
			}
			c.JSON(customError.Code, gin.H{
				"message": customError.Msg,
			})
			log.Errorf("[%s] code=%d, msg=%s", errorType, customError.Code, customError.Msg)
			return
		}
		if baseErr := recovered.(error); baseErr != nil {
			log.Errorf("[%s] code=%d, msg=%s", errorType, http.StatusInternalServerError, baseErr.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": baseErr.Error(),
			})
			return
		}
		log.Errorf("[%s] Unhandled error: %v", errorType, recovered)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unhandled error",
		})
	})
}
