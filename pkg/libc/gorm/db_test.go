package gorm

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConnectWithDSN(t *testing.T) {
	viper.SetEnvPrefix("DATABASE")
	_ = viper.BindEnv("DSN")
	_ = viper.Unmarshal(&C)

	db, err := ConnectWithDSN(C.DSN)
	assert.NoError(t, err)
	assert.IsType(t, &gorm.DB{}, db)
}
