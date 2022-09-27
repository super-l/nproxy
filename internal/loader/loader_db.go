package loader

import (
	"errors"
	"github.com/super-l/nproxy/internal"
	"github.com/super-l/nproxy/services/model"
)

func InitDbLoader() error {
	var err error
	// 1:初始化数据库相关
	err = migrateDB()

	return err
}

// 1:初始化数据库相关
func migrateDB() error {
	db := internal.GetDbInstance()
	if db == nil {
		return errors.New("db instance is empty!")
	}
	_ = db.AutoMigrate(&model.Proxy{})
	_ = db.AutoMigrate(&model.ProxyApi{})
	return nil
}
