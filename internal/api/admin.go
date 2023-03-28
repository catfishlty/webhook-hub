package api

import (
	"fmt"
	"github.com/catfishlty/webhooks-hub/internal/common"
	"github.com/catfishlty/webhooks-hub/internal/data"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/catfishlty/webhooks-hub/internal/utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func AddRuleHandler(db *gorm.DB) func(c *gin.Context) {
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
		id, err := data.AddRule(db, request)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, common.SingleResponse("id", id))
	}
}

func GetRuleHandler(db *gorm.DB) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		rule, err := data.GetRule(db, id)
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusNotFound,
				Msg:  fmt.Sprintf("rule id %s not found", id),
			})
		}
		c.JSON(http.StatusOK, rule)
	}
}

func GetRuleListHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.Param("page"))
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusBadRequest,
				Msg:  "page must be a number",
			})
		}
		rules, total, err := data.GetRuleList(db, page)
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusInternalServerError,
				Msg:  "get rule list error",
			})
		}
		c.JSON(http.StatusOK, common.ListResponse(rules, total))
	}
}

func DeleteRuleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := data.RemoveRule(db, id)
		if err != nil {
			panic(&types.CommonError{
				Code: http.StatusInternalServerError,
				Msg:  "delete rule error",
			})
		}
		c.Status(http.StatusOK)
	}
}
