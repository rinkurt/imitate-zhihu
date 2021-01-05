package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
)

const KeyHotQuestions = "HotQuestions"

func SyncHotQuestions() {
	tool.Rdb.Del(context.Background(), KeyHotQuestions)

	heats, _ := repository.GetQuestionHeats()
	for _, heat := range heats {
		tool.Rdb.ZAdd(context.Background(), KeyHotQuestions, &redis.Z{
			Member: heat.Id,
			Score: float64(heat.Heat),
		})
	}

	ansHeats, _ := repository.GetAnswerHeatsGroupByQuestion()
	for _, heat := range ansHeats {
		tool.Rdb.ZIncrBy(context.Background(), KeyHotQuestions,
			float64(heat.Heat), tool.Int64ToStr(heat.QuestionId))
	}

}


func GetHotQuestions(cursor int, size int) ([]redis.Z, result.Result) {
	q, err := tool.Rdb.ZRevRangeWithScores(context.Background(), KeyHotQuestions,
		int64(cursor), int64(cursor+size-1)).Result()
	if err != nil && err != redis.Nil {
		return nil, result.RedisErr.WithError(err)
	}
	return q, result.Ok
}