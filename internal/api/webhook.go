package api

import (
	"fmt"
	"github.com/catfishlty/webhooks-hub/internal/data"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func webhookHandler(db *gorm.DB, sender *Sender) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		rule, err := data.GetRule(db, id)
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusNotFound,
				Msg:  fmt.Sprintf("webhook id: %s not found", id),
			})
			return
		}
		resp, err := sender.Send(rule.Send)
		if err != nil {
			panic(&types.CommonError{
				Code: resp.StatusCode(),
				Err:  err,
			})
		}
		for k, values := range resp.Header() {
			for _, val := range values {
				c.Header(k, val)
			}
		}
		c.String(resp.StatusCode(), resp.String())
	}
}
