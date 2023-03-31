package hub

import (
	"fmt"
	"github.com/catfishlty/webhooks-hub/exp"
	"github.com/catfishlty/webhooks-hub/internal/check"
	"github.com/catfishlty/webhooks-hub/internal/common"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/catfishlty/webhooks-hub/internal/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (hub *Hub) AddUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.UserCreateRequest
		err := c.BindJSON(&request)
		exp.HandleBindJSON(err)
		exp.HandleRequestInvalid(request.Validate())
		err = hub.db.CreateUser(request.Username, request.Password)
		exp.HandleDB(err, "create user error")
		c.Status(http.StatusOK)
	}
}
func (hub *Hub) DeleteUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		exp.HandleRequestInvalid(check.ValidateId(id))
		err := hub.db.DeleteUser(id)
		exp.HandleDB(err, "delete user error")
		c.Status(http.StatusOK)
	}
}

func (hub *Hub) GetUserListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		exp.HandleRequestInvalid(check.ValidatePage(c.Param("page")))
		page, _ := strconv.Atoi(c.Param("page"))
		users, total, err := hub.db.GetUserList(page, common.PageSize)
		exp.HandleDB(err, "get user list error")
		c.JSON(http.StatusOK, common.ListResponse(users, total))
	}
}

func (hub *Hub) GetUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		exp.HandleRequestInvalid(check.ValidateId(id))
		user, err := hub.db.GetUser(id)
		exp.HandleDB(err, "get user error")
		c.JSON(http.StatusOK, user)
	}
}

func (hub *Hub) UpdateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		exp.HandleRequestInvalid(check.ValidateId(id))
		var request types.UserUpdateRequest
		err := c.BindJSON(&request)
		exp.HandleBindJSON(err)
		exp.HandleRequestInvalid(request.Validate())
		err = hub.db.UpdateUser(id, types.User{
			Username: request.Username,
		})
		exp.HandleDB(err, "update user error")
		c.Status(http.StatusOK)
	}
}

func (hub *Hub) UpdateUserResetPasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		exp.HandleRequestInvalid(check.ValidateId(id))
		var request types.UserResetPasswordRequest
		err := c.BindJSON(&request)
		exp.HandleBindJSON(err)
		exp.HandleRequestInvalid(request.Validate())
		err = hub.db.UpdatePassword(id, request.Password, request.NewPassword)
		exp.HandleDB(err, "reset password error")
		c.Status(http.StatusOK)
	}
}

func (hub *Hub) AddRuleHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request types.Rule
		err := c.BindJSON(&request)
		exp.HandleBindJSON(err)
		exp.HandleRequestInvalid(request.Validate())
		request.ID = utils.UUID()
		request.SendId = utils.UUID()
		request.Send.ID = request.SendId
		request.ReceiveId = utils.UUID()
		request.Receive.ID = request.ReceiveId
		id, err := hub.db.AddRule(request)
		exp.HandleDB(err, "add rule error")
		c.JSON(http.StatusOK, common.SingleResponse("id", id))
	}
}

func (hub *Hub) GetRuleHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		exp.HandleRequestInvalid(check.ValidateId(id))
		rule, err := hub.db.GetRule(id)
		exp.HandleDB(err, fmt.Sprintf("rule id %s not found", id))
		c.JSON(http.StatusOK, rule)
	}
}

func (hub *Hub) GetRuleListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		exp.HandleRequestInvalid(check.ValidatePage(c.Param("page")))
		page, _ := strconv.Atoi(c.Param("page"))
		rules, total, err := hub.db.GetRuleList(page)
		exp.HandleDB(err, "get rule list error")
		c.JSON(http.StatusOK, common.ListResponse(rules, total))
	}
}

func (hub *Hub) DeleteRuleHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := hub.db.RemoveRule(id)
		exp.HandleDB(err, "delete rule error")
		c.Status(http.StatusOK)
	}
}
