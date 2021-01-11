package test_tool

import (
	"database/sql"
	"testing"

	"github.com/jinzhu/gorm"
)

func MockedGORMDBForTest(t *testing.T, sqlDB *sql.DB) *gorm.DB {
	gormDB, err := gorm.Open("postgres", sqlDB)
	if err != nil {
		t.Error(err)
	}

	return gormDB
}
