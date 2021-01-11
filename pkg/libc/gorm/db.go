package gorm

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var defaultDB *gorm.DB

// DB 为了 UnitTest 可以 mock
var DB func() *gorm.DB = globalDB

func globalDB() *gorm.DB {
	return defaultDB
}

// set opt for global db
func SetOption(opts ...Option) error {
	return SetDBOption(defaultDB, opts...)
}

// 通过  dsn, Options 创建 db
func ConnectWithDSN(dsn string, opts ...Option) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.DB().Ping(); err != nil {
		return nil, err
	}

	for _, opt := range opts {
		if err := opt(db); err != nil {
			return nil, err
		}
	}

	return db, nil
}

// 通过 Conf 创建全局 db
func ConnectGlobalDB(conf Conf) error {
	opts := conf.ToOptions()
	db, err := ConnectWithDSN(conf.DSN, opts...)

	if err != nil {
		return err
	}

	defaultDB = db

	return nil
}

// 使用 conf.C 创建全局 db 的便捷方式
func Connect() error {
	return ConnectGlobalDB(C)
}
