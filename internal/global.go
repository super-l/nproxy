package internal

import (
	"github.com/super-l/nproxy/internal/config"
	"github.com/super-l/nproxy/pkg/db/mysql"
	"github.com/super-l/nproxy/pkg/db/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// 获取数据库链接
func GetDbInstance() *gorm.DB {
	if db == nil {
		if config.GetConfig().Db.Type == "mysql" {
			dbConfig := mysql.Config{
				Host:     config.GetConfig().Db.Mysql.Host,
				Port:     config.GetConfig().Db.Mysql.Port,
				Username: config.GetConfig().Db.Mysql.Username,
				Password: config.GetConfig().Db.Mysql.Password,
				Database: config.GetConfig().Db.Mysql.Database,
				Prefix:   "nsproxy_",
			}
			dbData, err1 := mysql.NewMysqlDb(dbConfig)
			if err1 == nil {
				db = dbData.Db
			} else {
				SLogger.GetStdoutLogger().Error(err1.Error())
			}
		} else if config.GetConfig().Db.Type == "sqlite" {
			dbData, err2 := sqlite.NewSqliteDb(config.GetConfig().Db.Sqlite.DbPath)
			if err2 == nil {
				db = dbData.Db
			} else {
				SLogger.GetStdoutLogger().Error(err2.Error())
			}
			db = dbData.Db
		} else {
			SLogger.GetStdoutLogger().Error("dbtype is error")
		}
	}
	return db
}
