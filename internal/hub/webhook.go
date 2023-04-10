package hub

import (
	"fmt"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/catfishlty/webhooks-hub/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
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
		variables, data, err := utils.GetVariables(rule.Receive, c)
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusBadRequest,
				Err:  err,
			})
		}
		err = utils.ValidateVariables(variables)
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusBadRequest,
				Err:  err,
			})
		}
		resp, err := hub.sendRequest(rule, variables, data)
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

func (hub *Hub) sendRequest(rule types.Rule, variables map[string]types.VariableItem, data []byte) (*resty.Response, error) {
	send := rule.Send.ToResty()
	utils.ReplaceVariables(&send, variables)
	return hub.sender.Send(send)
}
