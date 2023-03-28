package data

import (
	"github.com/catfishlty/webhooks-hub/internal/common"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"gorm.io/gorm"
)

func GetRule(db *gorm.DB, id string) (types.Rule, error) {
	var rule types.Rule
	result := db.Preload("Receive").Preload("Send").Where("id = ?", id).First(&rule)
	return rule, result.Error
}

func GetRuleList(db *gorm.DB, page int) ([]types.RuleItem, int64, error) {
	var rules []types.RuleItem
	result := db.Model(&types.Rule{}).Offset((page - 1) * common.PageSize).Limit(common.PageSize).Find(&rules)
	var count int64
	db.Model(&types.Rule{}).Count(&count)
	return rules, count, result.Error
}

func RemoveRule(db *gorm.DB, id string) error {
	result := db.Select("Send", "Receive").Where("id = ?", id).Delete(&types.Rule{})
	return result.Error
}

func AddRule(db *gorm.DB, request types.Rule) (string, error) {
	result := db.Create(&request)
	return request.ID, result.Error
}
