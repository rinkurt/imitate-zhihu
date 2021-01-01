package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
)

func KeyReadVote(uid int64) string {
	return "ReadVote:" + tool.Int64ToStr(uid)
}

func KeyWriteVote(uid int64) string {
	return "WriteVote:" + tool.Int64ToStr(uid)
}

func SyncAnswerVoteFromDB(uid int64) {
	votes := repository.SelectVotesByUid(uid)
	for _, vote := range votes {
		var state string
		if vote.IsUpvote {
			state = "1"
		} else {
			state = "2"
		}
		err := tool.Rdb.HSet(context.Background(), KeyReadVote(uid),
			tool.Int64ToStr(vote.AnswerId), state+":"+tool.Int64ToStr(vote.UpdateAt)).Err()
		if err != nil {
			tool.Logger.Error(err)
		}
	}
}

func GetAnswerVoteStatus(uid int64, aid int64) result.Result {
	vote, err := tool.Rdb.HGetAll(context.Background(), KeyReadVote(uid)).Result()
	if err != nil && err != redis.Nil {
		tool.Logger.Error(err)
		return result.RedisErr.WithError(err)
	}

	if len(vote) == 0 {
		tool.Rdb.HSet(context.Background(), KeyReadVote(uid), "", "")
		SyncAnswerVoteFromDB(uid)

		vote, err = tool.Rdb.HGetAll(context.Background(), KeyReadVote(uid)).Result()
		if err != nil && err != redis.Nil {
			tool.Logger.Error(err)
			return result.RedisErr.WithError(err)
		}

		tool.Rdb.Expire(context.Background(), KeyReadVote(uid), ExpireTime)
	}

	write, err := tool.Rdb.HGetAll(context.Background(), KeyWriteVote(uid)).Result()
	if err != nil && err != redis.Nil {
		tool.Logger.Error(err)
		return result.RedisErr.WithError(err)
	}

	for k, v := range write {
		vote[k] = v
	}

	val := vote[tool.Int64ToStr(aid)]
	status, _, _ := tool.ParseVoteVal(val)

	return result.Ok.WithData(status)
}
