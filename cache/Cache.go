package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"imitate-zhihu/enum"
	"imitate-zhihu/repository"
	"imitate-zhihu/tool"
	"strings"
	"time"
)

const ExpireTime = time.Hour * 3

func KeyVrfCode(email string) string {
	return "VrfCode:" + email
}

func KeyUser(uid int64) string {
	return "User:" + tool.Int64ToStr(uid)
}

// Write-cache for counts
func KeyWrite(typ int, id int64) string {
	return "Write:" + tool.IntToStr(typ) + ":" + tool.Int64ToStr(id)
}

// Read-cache for vote user lists
func KeyReadVote(typ int, id int64) string {
	return "ReadVote:" + tool.IntToStr(typ) + ":" + tool.Int64ToStr(id)
}

// Write-cache for vote user lists
func KeyWriteVote(typ int, id int64) string {
	return "WriteVote:" + tool.IntToStr(typ) + ":" + tool.Int64ToStr(id)
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

func SyncCount() {
	fmt.Println("SyncCount")
	keys, _ := tool.Rdb.Keys(context.Background(), "Write:*").Result()
	for _, key := range keys {
		split := strings.Split(key, ":")
		if len(split) != 3 {
			continue
		}
		typ, err1 := tool.StrToInt(split[1])
		id, err2 := tool.StrToInt64(split[2])
		if err1 != nil || err2 != nil {
			continue
		}
		counts, err := tool.Rdb.HGetAll(context.Background(), key).Result()
		if err != nil {
			tool.Logger.Error(err)
			continue
		}

		switch typ {
		case enum.Question:
			question := &repository.Question{Id: id}
			question.AnswerCount = tool.StrToDefaultInt(counts[enum.AnswerCount])
			question.CommentCount = tool.StrToDefaultInt(counts[enum.CommentCount])
			question.ViewCount = tool.StrToDefaultInt(counts[enum.ViewCount])
			question.LikeCount = tool.StrToDefaultInt(counts[enum.UpvoteCount])
			res := repository.UpdateQuestionCounts(question)
			if res.NotOK() {
				tool.Logger.Error(res.Error())
				continue
			}
		case enum.Answer:
			answer := &repository.Answer{Id: id}
			answer.ViewCount = tool.StrToDefaultInt(counts[enum.ViewCount])
			answer.UpvoteCount = tool.StrToDefaultInt(counts[enum.UpvoteCount])
			answer.DownvoteCount = tool.StrToDefaultInt(counts[enum.DownvoteCount])
			answer.CommentCount = tool.StrToDefaultInt(counts[enum.CommentCount])
			res := repository.UpdateAnswerCounts(answer)
			if res.NotOK() {
				tool.Logger.Error(res.Error())
				continue
			}
		}
		tool.Rdb.Del(context.Background(), key)
	}
}



