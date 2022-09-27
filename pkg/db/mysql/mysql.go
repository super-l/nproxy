package mysql

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"
)

type mysqlDb struct {
	Db          *gorm.DB
	TablePrefix string
}

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Prefix   string
}

func NewMysqlDb(config Config) (*mysqlDb, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database)

	sqlDB, sqlErr := sql.Open("mysql", dsn)

	if sqlErr != nil {
		return nil, sqlErr
	}

	gormDB, gormErr := gorm.Open(
		mysql.New(mysql.Config{Conn: sqlDB}),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   config.Prefix, // 指定表前缀，修改默认表名
				SingularTable: true,          // 设置全局表名禁用复数
			},
			//Logger: logger.Default.LogMode(logger.Info),
		})
	if gormErr != nil {
		return nil, gormErr
	}
	gormDB.Logger = logger.Default.LogMode(logger.Silent)
	sqlDB, dbErr := gormDB.DB()
	if dbErr != nil {
		return nil, dbErr
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	return &mysqlDb{
		Db:          gormDB,
		TablePrefix: config.Prefix,
	}, nil
}

func (mysql *mysqlDb) Close() error {
	sqlDB, dbErr := mysql.Db.DB()
	if dbErr != nil {
		return dbErr
	}
	return sqlDB.Close()
}
