package redis

import (
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestConnector_Connect(t *testing.T) {
	viper.SetEnvPrefix("REDIS")
	_ = viper.BindEnv("URI")
	_ = viper.Unmarshal(&C)
	t.Logf("redis uri: %s", C.URI)

	connector, err := NewConnector(C)
	assert.NoError(t, err)
	if err != nil {
		return
	}

	cli, err := connector.Connect()
	assert.NoError(t, err)
	assert.IsType(t, &redis.Client{}, cli)
}
