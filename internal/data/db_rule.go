package data

import (
	"errors"
	"fmt"
	"github.com/catfishlty/webhook-hub/internal/common"
	"github.com/catfishlty/webhook-hub/internal/types"
)

func (db *DB) GetRule(id string) (types.Rule, error) {
	var rule types.Rule
	result := db.orm.Where("uid = ?", id).First(&rule)
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

func (db *DB) ExistRule(id string) error {
	var count int64
	db.orm.Model(&types.Rule{}).Where("uid = ?", id).Count(&count)
	if count > 0 {
		return nil
	}
	return errors.New(fmt.Sprintf("rule uid: %s not found", id))
}

func (db *DB) AddRule(request types.Rule) (string, error) {
	result := db.orm.Create(&request)
	return request.UID, result.Error
}

func (db *DB) UpdateRule(request types.Rule) error {
	result := db.orm.Model(&types.Rule{}).Where("uid = ?", request.UID).Updates(request)
	return result.Error
}

func (db *DB) RemoveRule(id string) error {
	result := db.orm.Select("Send", "Receive").Where("uid = ?", id).Delete(&types.Rule{})
	return result.Error
}
