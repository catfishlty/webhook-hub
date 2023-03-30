package data

import (
	"errors"
	"github.com/catfishlty/webhooks-hub/internal/common"
	"github.com/catfishlty/webhooks-hub/internal/types"
	"github.com/catfishlty/webhooks-hub/internal/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DB struct {
	orm  *gorm.DB
	salt string
}

func NewDB(gorm *gorm.DB, salt string) *DB {
	return &DB{
		orm:  gorm,
		salt: salt,
	}
}

func (db *DB) Init() {
	if db.GetUserCount() == 0 {
		if err := db.CreateUser(common.DefaultUsername, common.DefaultPassword); err != nil {
			log.Fatal("failed to create default user", err)
		}
	}
}

func (db *DB) Migrate() {
	err := db.orm.AutoMigrate(&types.Rule{}, &types.ReceiveRequest{}, &types.SendRequest{}, &types.User{})
	if err != nil {
		panic(err)
	}
}

func (db *DB) CheckUser(username, password string) (*types.User, error) {
	var user types.User
	result := db.orm.Where("username = ? and password = ?", username, utils.EncodePassword(password, db.salt)).First(&user)
	return &user, result.Error
}

func (db *DB) CreateUser(username, password string) error {
	result := db.orm.Create(&types.User{
		Id:        utils.UUID(),
		Username:  username,
		Password:  utils.EncodePassword(password, db.salt),
		AccessKey: utils.NewRandom().String(16),
		SecretKey: utils.NewRandom().String(32),
	})
	return result.Error
}

func (db *DB) GetUser(id string) (types.User, error) {
	var user types.User
	result := db.orm.Where("id = ?", id).First(&user)
	return user, result.Error
}

func (db *DB) GetUserList(page, pageSize int) ([]types.User, int64, error) {
	var users []types.User
	result := db.orm.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)
	return users, db.GetUserCount(), result.Error
}

func (db *DB) GetUserCount() int64 {
	var count int64
	db.orm.Model(&types.User{}).Count(&count)
	return count
}

func (db *DB) DeleteUser(id string) error {
	if db.GetUserCount() <= 1 {
		return errors.New("can't delete user, please ensure there's at least one user")
	}
	result := db.orm.Where("id = ?", id).Delete(&types.User{})
	return result.Error
}

func (db *DB) UpdateUser(id string, request types.User) error {
	if request.Password != "" {
		request.Password = ""
	}
	result := db.orm.Model(&types.User{}).Where("id = ?", id).Updates(request)
	return result.Error
}

func (db *DB) UpdatePassword(id, password, newPassword string) error {
	result := db.orm.Model(&types.User{}).
		Where("id = ?", id).
		Where("password = ?", utils.EncodePassword(password, db.salt)).
		Updates(types.User{
			Password: utils.EncodePassword(newPassword, db.salt),
		})
	return result.Error
}

func (db *DB) UpdatePasswordAdmin(id, password string) error {
	result := db.orm.Model(&types.User{}).
		Where("id = ?", id).
		Updates(types.User{
			Password: utils.EncodePassword(password, db.salt),
		})
	return result.Error
}

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

func (db *DB) RemoveRule(id string) error {
	result := db.orm.Select("Send", "Receive").Where("id = ?", id).Delete(&types.Rule{})
	return result.Error
}

func (db *DB) AddRule(request types.Rule) (string, error) {
	result := db.orm.Create(&request)
	return request.ID, result.Error
}
