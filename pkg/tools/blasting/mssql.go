package blasting

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

func NewGormConnMssql(driverName string, dataSourceName string) bool {
	engine, err := gorm.Open(driverName, dataSourceName)
	if err != nil {
		return false
	}
	if err = engine.DB().Ping(); err != nil {
		return false
	}
	defer engine.Close()
	return true
}
