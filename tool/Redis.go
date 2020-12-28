package tool

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var Rdb *redis.Client

func InitRedis() {
	if Rdb != nil {
		return
	}
	Rdb = redis.NewClient(&redis.Options{
		Addr:     Cfg.RedisAddr,
		Password: Cfg.RedisPassword, // no password set
		DB:       0,  // use default DB
	})
	_, err := Rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Redis Connection Failed")
		panic(err)
	}
}

