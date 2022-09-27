package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type sqliteDb struct {
	Db   *gorm.DB
	Path string
}

func NewSqliteDb(path string) (*sqliteDb, error) {
	gormDB, err := gorm.Open(sqlite.Open(path))
	//gormDB.Logger = logger.Default.LogMode(logger.Silent)
	if err != nil {
		return nil, err
	}
	return &sqliteDb{gormDB, path}, nil
}

func (sqlite *sqliteDb) Close() error {
	sqlDB, dbErr := sqlite.Db.DB()
	if dbErr != nil {
		return dbErr
	}
	return sqlDB.Close()
}
