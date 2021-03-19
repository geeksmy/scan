package blasting

import (
	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

func NewConnPgSql(driverName string, dataSourceName string) bool {
	engine, err := xorm.NewEngine(driverName, dataSourceName)
	if err != nil {
		return false
	}
	engine.SetLogLevel(4)
	if err = engine.Ping(); err != nil {
		return false
	}
	defer engine.Close()
	return true
}
