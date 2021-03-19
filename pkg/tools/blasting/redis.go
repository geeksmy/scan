package blasting

import (
	"github.com/go-redis/redis/v7"
)

func NewConnRedis(addr, user, pass string) bool {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: user,
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
