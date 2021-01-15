package services

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var runCtx context.Context

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	runCtx = context.Background()
}
