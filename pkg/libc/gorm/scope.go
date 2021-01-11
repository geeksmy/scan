package gorm

import (
	"github.com/jinzhu/gorm"
)

// PaginationScope 分页 Scope,
// @param lastID: 前一页的最后一行记录的 ID
// @param perPage: 每页行数, default 100
func PaginationScope(lastID, perPage int) Scope {
	if perPage <= 0 {
		perPage = 100
	}

	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id>?", lastID).Limit(perPage)
	}
}
