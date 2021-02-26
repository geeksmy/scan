package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Strings []string

func (s *Strings) Scan(src interface{}) error {
	switch typ := src.(type) {
	default:
		return fmt.Errorf("%s not supported", typ)
	case []byte:
		return json.Unmarshal(src.([]byte), s)
	}
}

func (s Strings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type BaseModel struct {
	ID        int        `gorm:"primary_key"`
	CreatedAt time.Time  `sql:"type:timestamp(6);default:current_timestamp(6)"`
	UpdatedAt time.Time  `sql:"type:timestamp(6);default:current_timestamp(6)"`
	DeletedAt *time.Time `sql:"index"`
}

type BaseUUIDModel struct {
	ID        string     `gorm:"primary_key;type:varchar(36);not null;"`
	CreatedAt time.Time  `msgpack:"-"`
	UpdatedAt time.Time  `msgpack:"-"`
	DeletedAt *time.Time `sql:"index" msgpack:"-"`
}

func NewBaseUUIDModel() BaseUUIDModel {
	return BaseUUIDModel{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func IsDuplicateError(err error) bool {
	mysqlErr, ok := err.(*mysql.MySQLError)
	if ok {
		if mysqlErr.Number == 1062 {
			return true
		}
	}

	return false
}

// FilterDeleted 过滤软删除的数据
// require: Model 需要 embed `BaseUUIDModel`
// usage:   if excludeDelete {
//				db = FilterDeleted(dao.db)
//          }
//          db = db.Where(whereCondition)
//          db.Find(&rows)
func FilterDeleted(db *gorm.DB) *gorm.DB {
	return db.Where("deleted_at IS NULL")
}

// SoftDelete 软删除数据
// require: Model 必须 embed `BaseUUIDModel`
// 因为 `UpdateColumn` 是 execute 接口, 所以不能使用 `Scope`
// usage:   db = db.Where(obj.ID)
//          err := SoftDelete(dao.db).Error
func SoftDelete(db *gorm.DB) *gorm.DB {
	return db.UpdateColumn("deleted_at", time.Now())
}
