// sqlmock 实现 sql-mock 的接口
package hbhsid

import (
	"database/sql/driver"
)

type HBHSIDArg struct{}

func (HBHSIDArg) Match(value driver.Value) bool {
	_, ok := value.(int64)
	return ok
}
