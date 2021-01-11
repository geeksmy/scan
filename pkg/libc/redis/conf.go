package redis

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// redis 的配置
/*
Redis:
  URI: redis://password@redis:6379/0
*/
type Conf struct {
	URI              string   `yaml:"uri"` // redis://password@host:6379/0
	PoolSize         int      `yaml:"pool_size"`
	Sentinel         bool     `yaml:"sentinel"`
	Password         string   `yaml:"password"`
	MasterName       string   `yaml:"master_name" default:"master"`
	SentinelAddrs    []string `yaml:"sentinel_addrs"`
	SentinelPassword string   `yaml:"sentinel_password"`
	DB               int      `yaml:"db"`
}

var C = Conf{URI: "redis://:password@redis:6379/0"}

func BindPflag(flagSet *pflag.FlagSet, keyPrefix string) {
	flagSet.String("redis_uri", "", "redis URI eg: redis://:password@host:port/db")
	_ = viper.BindPFlag(keyPrefix+".URI", flagSet.Lookup("redis_uri"))
}
