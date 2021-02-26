package dao

import (
	"time"

	"scan/config"

	// init mysql driver
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) {
	zap.L().Debug("connect db", zap.String("dsn", cfg.Database.DSN))

	var err error

	DB, err = gorm.Open("mysql", cfg.Database.DSN)
	if err != nil {
		zap.L().Panic("connect db failed", zap.Error(err))
	}

	if cfg.Debug {
		DB.LogMode(cfg.Debug)
	}

	DB.SingularTable(true)

	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	if cfg.Database.MaxIdleConns > 0 {
		DB.DB().SetMaxIdleConns(cfg.Database.MaxIdleConns)
	}

	// SetMaxOpenCons 设置数据库的最大连接数量。
	if cfg.Database.MaxOpenConns > 0 {
		DB.DB().SetMaxOpenConns(cfg.Database.MaxOpenConns)
	}

	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	if cfg.Database.ConnMaxLifetime != "" {
		maxLifetime, err := time.ParseDuration(cfg.Database.ConnMaxLifetime)
		if err != nil {
			zap.L().Panic("db ConnMaxLifetime parse failed", zap.Error(err))
		}

		DB.DB().SetConnMaxLifetime(maxLifetime)
	}

	if err := DB.DB().Ping(); err != nil {
		zap.L().Panic("ping db failed", zap.Error(err))
	}
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func AutoMigrateDB() {
	query := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci")
	if err := query.AutoMigrate(
	// TODO 如果有业务需要model，请先修改这里
	).Error; err != nil {
		zap.L().Panic("migrate db fail", zap.Error(err))
	}
}
