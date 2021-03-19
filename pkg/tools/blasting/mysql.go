package blasting

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"xorm.io/xorm"
)

func NewXormConnMysql(driverName string, dataSourceName string) bool {
	engine, err := xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		return false
	}
	engine.SetLogLevel(4)
	if err = engine.Ping(); err != nil {
		return false
	}
	engine.Close()
	return true
}

func NewGormConnMysql(driverName string, dataSourceName string) bool {
	engine, err := gorm.Open(driverName, dataSourceName)
	if err != nil {
		return false
	}
	// engine.SetLogger()
	engine.SetLogger(gorm.Logger{})
	if err = engine.DB().Ping(); err != nil {
		return false
	}

	engine.Close()
	return true
}
