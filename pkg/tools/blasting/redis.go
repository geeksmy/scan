package blasting

import (
	"github.com/go-redis/redis/v7"
)

func NewConnRedis(addr, pass string) bool {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pass,
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return false
	}

	_ = client.Close()
	return true
}
