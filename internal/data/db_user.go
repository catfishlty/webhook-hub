package data

import (
	"errors"
	"github.com/catfishlty/webhook-hub/internal/types"
	"github.com/catfishlty/webhook-hub/internal/utils"
)

func (db *DB) CheckUser(username, password string) (*types.User, error) {
	var user types.User
	result := db.orm.Where("username = ? and password = ?", username, utils.EncodePassword(password, db.salt)).First(&user)
	return &user, result.Error
}

func (db *DB) CreateUser(username, password string) error {
	result := db.orm.Create(&types.User{
		UID:       utils.UUID(),
		Username:  username,
		Password:  utils.EncodePassword(password, db.salt),
		AccessKey: utils.NewRandom().String(16),
		SecretKey: utils.NewRandom().String(32),
	})
	return result.Error
}

func (db *DB) GetUser(id string) (types.User, error) {
	var user types.User
	result := db.orm.Where("uid = ?", id).First(&user)
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
	result := db.orm.Where("uid = ?", id).Delete(&types.User{})
	return result.Error
}

func (db *DB) UpdateUser(id string, request types.User) error {
	if request.Password != "" {
		request.Password = ""
	}
	result := db.orm.Model(&types.User{}).Where("uid = ?", id).Updates(request)
	return result.Error
}

func (db *DB) UpdatePassword(id, password, newPassword string) error {
	result := db.orm.Model(&types.User{}).
		Where("uid = ?", id).
		Where("password = ?", utils.EncodePassword(password, db.salt)).
		Updates(types.User{
			Password: utils.EncodePassword(newPassword, db.salt),
		})
	return result.Error
}

func (db *DB) UpdatePasswordAdmin(id, password string) error {
	result := db.orm.Model(&types.User{}).
		Where("uid = ?", id).
		Updates(types.User{
			Password: utils.EncodePassword(password, db.salt),
		})
	return result.Error
}
