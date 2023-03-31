package data

import (
	"github.com/catfishlty/webhooks-hub/internal/common"
	"github.com/catfishlty/webhooks-hub/internal/types"
)

func (db *DB) GetRule(id string) (types.Rule, error) {
	var rule types.Rule
	result := db.orm.Preload("Receive").Preload("Send").Where("id = ?", id).First(&rule)
	return rule, result.Error
}

func (db *DB) GetRuleList(page int) ([]types.RuleItem, int64, error) {
	var rules []types.RuleItem
	result := db.orm.Model(&types.Rule{}).Offset((page - 1) * common.PageSize).Limit(common.PageSize).Find(&rules)
	return rules, db.GetRuleCount(), result.Error
}

func (db *DB) GetRuleCount() int64 {
	var count int64
	db.orm.Model(&types.Rule{}).Count(&count)
	return count
}

func (db *DB) AddRule(request types.Rule) (string, error) {
	result := db.orm.Create(&request)
	return request.ID, result.Error
}

func (db *DB) RemoveRule(id string) error {
	result := db.orm.Select("Send", "Receive").Where("id = ?", id).Delete(&types.Rule{})
	return result.Error
}
