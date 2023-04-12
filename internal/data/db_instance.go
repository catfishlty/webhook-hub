package data

import (
	"errors"
	"github.com/catfishlty/webhook-hub/internal/common"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func GetDatabase(dbType, dsn string) (gorm.Dialector, error) {
	switch dbType {
	case common.DBTypeSqlite:
		return sqlite.Open(dsn), nil
	case common.DBTypeMySQL:
		return mysql.Open(dsn), nil
	case common.DBTypePostgres:
		return postgres.Open(dsn), nil
	case common.DBTypeSQLServer:
		return sqlserver.Open(dsn), nil
	}
	return nil, errors.New("unsupported database type")
}
