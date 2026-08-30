// Package cache handles cache interactions
package cache

import (
	"github.com/redis/go-redis/v9"
)

func NewCache() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	defer rdb.Close()
	return rdb
}
