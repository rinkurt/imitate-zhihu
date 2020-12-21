package cache

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"imitate-zhihu/tool"
	"time"
)

const ExpireTime = time.Hour * 3

func KeyVrfCode(email string) string {
	return "VrfCode:" + email
}

func KeyUser(uid int64) string {
	return "User:" + tool.Int64ToStr(uid)
}

func KeyWrite(typ int, id int64) string {
	return "Write:" + tool.IntToStr(typ) + ":" + tool.Int64ToStr(id)
}

func Get(key string, val interface{}) bool {
	// reset expire time
	tool.Rdb.Expire(context.Background(), key, ExpireTime)
	str, err := tool.Rdb.Get(context.Background(), key).Result()
	if err != nil {
		if err != redis.Nil {
			tool.Logger.Error(err)
		}
		return false
	}
	err = json.Unmarshal([]byte(str), val)
	return err == nil
}

func Set(key string, val interface{}) {
	bytes, err := json.Marshal(val)
	if err != nil {
		tool.Logger.Error(err)
		return
	}
	err = tool.Rdb.Set(context.Background(), key, bytes, ExpireTime).Err()
	if err != nil {
		tool.Logger.Error(err)
	}
}

// typ: Question, Comment...
// countType: ViewCount, UpvoteCount...
// types definition see package enum
func IncrViewCount(typ int, id int64, countType string, count int) {
	_, err := tool.Rdb.HIncrBy(context.Background(), KeyWrite(typ, id), countType, int64(count)).Result()
	if err != nil {
		tool.Logger.Error(err)
	}
}



