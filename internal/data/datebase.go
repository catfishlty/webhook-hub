package data

import (
	"github.com/catfishlty/webhooks-hub/internal/common"
	"github.com/catfishlty/webhooks-hub/internal/types"
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
	err := db.orm.AutoMigrate(&types.Rule{}, &types.User{})
	if err != nil {
		panic(err)
	}
}
