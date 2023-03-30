package hub

import (
	"fmt"
	"github.com/catfishlty/webhooks-hub/internal/common"
	"github.com/catfishlty/webhooks-hub/internal/data"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/catfishlty/webhooks-hub/internal/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (hub *Hub) AddUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request types.UserCreateRequest
		err := c.BindJSON(&request)
		if err != nil {
			log.Error("json parse error", err)
			panic(err)
		}
		err = hub.db.CreateUser(request.Username, request.Password)
		if err != nil {
			log.Errorf("create user error: %v", err)
			panic(err)
		}
		c.Status(http.StatusOK)
	}
}
func (hub *Hub) DeleteUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := hub.db.DeleteUser(id)
		if err != nil {
			log.Errorf("delete user error: %v", err)
			panic(err)
		}
		c.Status(http.StatusOK)
	}
}

func (hub *Hub) GetUserListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Param("page"))
		if err != nil {
			log.Errorf("get user list error: %v", err)
			panic(err)
		}
		users, total, err := hub.db.GetUserList(page, common.PageSize)
		if err != nil {
			log.Errorf("get user list error: %v", err)
			panic(err)
		}
		c.JSON(http.StatusOK, common.ListResponse(users, total))
	}
}

func (hub *Hub) GetUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		user, err := hub.db.GetUser(id)
		if err != nil {
			log.Errorf("get user error: %v", err)
			panic(err)
		}
		c.JSON(http.StatusOK, user)
	}
}

func (hub *Hub) UpdateUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var request types.UserUpdateRequest
		err := c.BindJSON(&request)
		if err != nil {
			log.Error("json parse error", err)
			panic(err)
		}
		err = hub.db.UpdateUser(id, types.User{
			Username: request.Username,
		})
		if err != nil {
			log.Errorf("update user error: %v", err)
			panic(err)
		}
		c.Status(http.StatusOK)
	}
}

func (hub *Hub) UpdateUserResetPasswordHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var request types.UserResetPasswordRequest
		err := c.BindJSON(&request)
		if err != nil {
			log.Error("json parse error", err)
			panic(err)
		}
		err = hub.db.UpdatePassword(id, request.Password, request.NewPassword)
		if err != nil {
			log.Errorf("reset password error: %v", err)
			panic(err)
		}
		c.Status(http.StatusOK)
	}
}

func (hub *Hub) AddRuleHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		var request types.Rule
		err := c.BindJSON(&request)
		if err != nil {
			log.Error("json parse error", err)
			panic(err)
		}
		request.ID = utils.UUID()
		request.SendId = utils.UUID()
		request.Send.ID = request.SendId
		request.ReceiveId = utils.UUID()
		request.Receive.ID = request.ReceiveId
		id, err := data.AddRule(hub.db, request)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, common.SingleResponse("id", id))
	}
}

func (hub *Hub) GetRuleHandler() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		rule, err := data.GetRule(hub.db, id)
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusNotFound,
				Msg:  fmt.Sprintf("rule id %s not found", id),
			})
		}
		c.JSON(http.StatusOK, rule)
	}
}

func (hub *Hub) GetRuleListHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Param("page"))
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusBadRequest,
				Msg:  "page must be a number",
			})
		}
		rules, total, err := data.GetRuleList(hub.db, page)
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusInternalServerError,
				Msg:  "get rule list error",
			})
		}
		c.JSON(http.StatusOK, common.ListResponse(rules, total))
	}
}

func (hub *Hub) DeleteRuleHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := data.RemoveRule(hub.db, id)
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusInternalServerError,
				Msg:  "delete rule error",
			})
		}
		c.Status(http.StatusOK)
	}
}
