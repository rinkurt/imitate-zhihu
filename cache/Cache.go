package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"imitate-zhihu/enum"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"strings"
	"time"
)

const ExpireTime = time.Hour * 6

func KeyVrfCode(email string) string {
	return "VrfCode:" + email
}

func KeyUser(uid int64) string {
	return "User:" + tool.Int64ToStr(uid)
}

func KeyQuestion(qid int64) string {
	return "Question:" + tool.Int64ToStr(qid)
}

// Write-cache for counts
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
func IncrCount(typ int, id int64, countType string, count int) result.Result {
	_, err := tool.Rdb.HIncrBy(context.Background(), KeyWrite(typ, id), countType, int64(count)).Result()
	if err != nil {
		tool.Logger.Error(err)
		return result.RedisErr.WithError(err)
	}
	return result.Ok
}

func ReadQuestionCounts(question *repository.Question) result.Result {
	counts, err := tool.Rdb.HGetAll(context.Background(), KeyWrite(enum.Question, question.Id)).Result()
	if err != nil {
		tool.Logger.Error(err)
		return result.RedisErr.WithError(err)
	}
	for k, v := range counts {
		c, err := tool.StrToInt(v)
		if err != nil {
			tool.Logger.Error("Cache value not integer!")
			continue
		}
		switch k {
		case enum.ViewCount:
			question.ViewCount += c
		case enum.AnswerCount:
			question.AnswerCount += c
		case enum.UpvoteCount:
			question.LikeCount += c
		case enum.CommentCount:
			question.CommentCount += c
		}
	}
	return result.Ok
}

func ReadAnswerCounts(answer *repository.Answer) result.Result {
	counts, err := tool.Rdb.HGetAll(context.Background(), KeyWrite(enum.Answer, answer.Id)).Result()
	if err != nil {
		tool.Logger.Error(err)
		return result.RedisErr.WithError(err)
	}
	for k, v := range counts {
		c, err := tool.StrToInt(v)
		if err != nil {
			tool.Logger.Error("Cache value not integer!")
			continue
		}
		switch k {
		case enum.ViewCount:
			answer.ViewCount += c
		case enum.UpvoteCount:
			answer.UpvoteCount += c
		case enum.DownvoteCount:
			answer.DownvoteCount += c
		case enum.CommentCount:
			answer.CommentCount += c
		}
	}
	return result.Ok
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

	SyncHotQuestions()
}



