package gorm

import (
	"testing"

	"gotest.tools/assert"
)

func TestConf_ToOptions(t *testing.T) {
	conf := Conf{
		DSN:             "this_is_dsn",
		LogMode:         true,
		MaxIdleConns:    10,
		MaxOpenConns:    10,
		ConnMaxLifetime: 100,
	}

	options := conf.ToOptions()
	assert.Equal(t, 4, len(options))
}
