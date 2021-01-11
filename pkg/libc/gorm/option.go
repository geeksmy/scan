package gorm

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Option = func(db *gorm.DB) error

// set opt for the special db
func SetDBOption(db2 *gorm.DB, opts ...Option) error {
	for _, opt := range opts {
		if err := opt(db2); err != nil {
			return err
		}
	}

	return nil
}

func BlockGlobalUpdateOpt(v bool) Option {
	return func(db *gorm.DB) error {
		return db.BlockGlobalUpdate(v).Error
	}
}

func LogModOpt(v bool) Option {
	return func(db *gorm.DB) error {
		return db.LogMode(v).Error
	}
}

func SetMaxIdleConnsOpt(v int) Option {
	return func(db *gorm.DB) error {
		db.DB().SetMaxIdleConns(v)
		return nil
	}
}

func SetMaxOpenConnsOpt(v int) Option {
	return func(db *gorm.DB) error {
		db.DB().SetMaxOpenConns(v)
		return nil
	}
}

func SetConnMaxLifetimeOpt(t time.Duration) Option {
	return func(db *gorm.DB) error {
		db.DB().SetConnMaxLifetime(t)
		return nil
	}
}
