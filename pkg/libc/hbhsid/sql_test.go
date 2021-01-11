package hbhsid

import (
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gotest.tools/assert"
)

func connectDB(t *testing.T) *gorm.DB {
	dsn, ok := os.LookupEnv("DATABASE_DSN")
	if !ok {
		t.Fatalf("no database dsn")
	}

	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db failed %s", err)
	}

	db.SingularTable(true)
	db.LogMode(true)

	return db
}

func TestID_Scan(t *testing.T) {
	db := connectDB(t)
	if err := db.AutoMigrate(BaseHBHSIDModel{}).Error; err != nil {
		t.Fatalf("auto migrate model failed %s", err)
	}

	m := NewBaseHBHSIDModel()
	t.Logf("create id: %d", m.ID.orig)

	if err := db.Create(&m).Error; err != nil {
		t.Fatalf("create row error %s", err)
	}

	m2 := NewBaseHBHSIDModel()
	// 使用上一次创建成功的 pk+100000, 防止测试的时候冲突
	m2.ID = New(m.ID.Origin() + 100000)
	if err := db.Create(&m2).Error; err != nil {
		t.Fatalf("create row with pk error %s", err)
	}
	t.Logf("%+v", m2)

	m3 := BaseHBHSIDModel{}
	if err := db.Find(&m3, "id=?", m2.ID).Error; err != nil {
		t.Fatalf("get row by pk failed: %s", err)
	}

	assert.Equal(t, m2.ID.Origin(), m3.ID.Origin())

	// db.Exec(`TRUNCATE TABLE base_hbhs_id_model RESTART IDENTITY ;`)
}
