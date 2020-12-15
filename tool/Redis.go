package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var Rdb *redis.Client
const CacheExpireTime = time.Hour * 6

func init() {
	if Rdb == nil {
		config := GetConfig()
		Rdb = redis.NewClient(&redis.Options{
			Addr:     config.RedisAddr,
			Password: config.RedisPassword, // no password set
			DB:       0,  // use default DB
		})
		_, err := Rdb.Ping(context.Background()).Result()
		if err != nil {
			fmt.Println("Redis Connection Failed")
			panic(err)
		}
	}
}

func CacheGet(key string, val interface{}) error {
	str, err := Rdb.Get(context.Background(), key).Result()
	if err != nil {
		if err != redis.Nil {
			Logger.Error(err)
		}
		return err
	}
	err = json.Unmarshal([]byte(str), val)
	return err
}

func CacheSet(key string, val interface{}) {
	bytes, err := json.Marshal(val)
	if err != nil {
		Logger.Error(err)
		return
	}
	err = Rdb.Set(context.Background(), key, bytes, CacheExpireTime).Err()
	if err != nil {
		Logger.Error(err)
	}
}

func KeyVrfCode(email string) string {
	return "VrfCode:" + email
}

func KeyUser(uid int64) string {
	return "User:" + Int64ToString(uid)
}