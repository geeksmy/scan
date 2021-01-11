package hbhsid

import (
	"database/sql/driver"
	"fmt"
	"math"
)

func (id *ID) Scan(src interface{}) error {
	switch v := src.(type) {
	case int64:
		if v > math.MaxUint32 {
			return fmt.Errorf("sql.Scan error beacuse %d overflow", v)
		}
		*id = New(uint32(v))
		return nil
	default:
		return fmt.Errorf("sql.Scan error because %t not supported", src)
	}
}

func (id ID) Value() (driver.Value, error) {
	return int64(id.orig), nil
}
