package hub

import (
	"fmt"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (hub *Hub) webhookHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		rule, err := hub.db.GetRule(id)
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusNotFound,
				Msg:  fmt.Sprintf("webhook id: %s not found", id),
			})
			return
		}
		resp, err := hub.sender.Send(rule.Send)
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
