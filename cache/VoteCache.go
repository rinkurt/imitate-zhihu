package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"imitate-zhihu/enum"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"strings"
	"time"
)

func KeyReadVote(uid int64) string {
	return "ReadVote:" + tool.Int64ToStr(uid)
}

func KeyWriteVote(uid int64) string {
	return "WriteVote:" + tool.Int64ToStr(uid)
}

func KeyVoteZSet(uid int64) string {
	return "VoteZSet:" + tool.Int64ToStr(uid)
}

func SyncAnswerVoteFromDB(uid int64) {
	votes := repository.SelectVotesByUid(uid)
	// 删除旧缓存
	err := tool.Rdb.Del(context.Background(), KeyReadVote(uid)).Err()
	if err != nil {
		tool.Logger.Error(err)
	}
	tool.Rdb.HSet(context.Background(), KeyReadVote(uid), "", "")
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
	tool.Rdb.Expire(context.Background(), KeyReadVote(uid), ExpireTime)
}

func GetAnswerVoteStatus(uid int64, aid int64) result.Result {
	exists, err := tool.Rdb.Exists(context.Background(), KeyReadVote(uid)).Result()
	if err != nil && err != redis.Nil {
		return result.RedisErr.WithError(err)
	}

	// key 在缓存中未找到，从数据库读取
	if exists == 0 {
		SyncAnswerVoteFromDB(uid)
	}

	val, err := tool.Rdb.HGet(context.Background(), KeyReadVote(uid), tool.Int64ToStr(aid)).Result()
	if err != nil && err != redis.Nil {
		return result.RedisErr.WithError(err)
	}

	wVal, err := tool.Rdb.HGet(context.Background(), KeyWriteVote(uid), tool.Int64ToStr(aid)).Result()
	if err != nil && err != redis.Nil {
		return result.RedisErr.WithError(err)
	}

	if wVal != "" {
		val = wVal
	}

	status, _, _ := tool.ParseVoteVal(val)
	return result.Ok.WithData(status)
}

func SetAnswerVoteStatus(uid int64, aid int64, state int) result.Result {
	res := GetAnswerVoteStatus(uid, aid)
	if res.NotOK() {
		return res
	}
	oldState := res.Data.(int)

	err := tool.Rdb.HSet(context.Background(), KeyWriteVote(uid), aid,
		tool.IntToStr(state)+":"+tool.Int64ToStr(time.Now().Unix())).Err()
	if err != nil {
		tool.Logger.Error(err)
		return result.RedisErr.WithError(err)
	}

	// 先从旧状态->0
	switch oldState {
	case 1:
		res = IncrCount(enum.Answer, aid, enum.UpvoteCount, -1)
	case 2:
		res = IncrCount(enum.Answer, aid, enum.DownvoteCount, -1)
	}

	if res.NotOK() {
		return res
	}

	// 再从0->新状态
	switch state {
	case 1:
		res = IncrCount(enum.Answer, aid, enum.UpvoteCount, 1)
	case 2:
		res = IncrCount(enum.Answer, aid, enum.DownvoteCount, 1)
	}

	if res.NotOK() {
		return res
	}

	return result.Ok
}

func SyncAnswerVote() {
	fmt.Println("SyncAnswerVote")
	keys, err := tool.Rdb.Keys(context.Background(), "WriteVote:*").Result()
	if err != nil && err != redis.Nil {
		tool.Logger.Error(err)
		return
	}
	if len(keys) == 0 {
		return
	}

	for _, key := range keys {
		val, err := tool.Rdb.HGetAll(context.Background(), key).Result()
		if err != nil {
			continue
		}

		split := strings.Split(key, ":")
		if len(split) < 2 {
			continue
		}
		uid, err := tool.StrToInt64(split[1])
		if err != nil {
			continue
		}

		var delAids []int64
		var insVotes []repository.AnswerVote

		for k, v := range val {
			aid, err := tool.StrToInt64(k)
			if err != nil {
				continue
			}
			s, t, err := tool.ParseVoteVal(v)
			if err != nil {
				continue
			}

			delAids = append(delAids, aid)

			if s != 0 {
				vote := repository.AnswerVote{
					AnswerId: aid,
					UserId: uid,
					UpdateAt: t,
				}
				if s == 1 {
					vote.IsUpvote = true
				}
				insVotes = append(insVotes, vote)
			}
		}

		repository.DeleteVotes(uid, delAids)
		repository.CreateVotes(insVotes)

		SyncAnswerVoteFromDB(uid)
		err = tool.Rdb.Del(context.Background(), key).Err()
		if err != nil {
			tool.Logger.Error(err)
		}
	}
}

func GetVotedAnswerIds(uid int64, cursor int, size int) ([]int64, result.Result) {
	if cursor == 0 {
		// 更新列表
		err := tool.Rdb.Del(context.Background(), KeyVoteZSet(uid)).Err()
		if err != nil && err != redis.Nil {
			return []int64{}, result.RedisErr.WithError(err)
		}

		exists, _ := tool.Rdb.Exists(context.Background(), KeyReadVote(uid)).Result()
		if exists == 0 {
			SyncAnswerVoteFromDB(uid)
		}

		rVal, _ := tool.Rdb.HGetAll(context.Background(), KeyReadVote(uid)).Result()
		wVal, _ := tool.Rdb.HGetAll(context.Background(), KeyWriteVote(uid)).Result()
		for k, v := range wVal {
			rVal[k] = v
		}

		for k, v := range rVal {
			state, updateAt, err := tool.ParseVoteVal(v)
			if err != nil {
				continue
			}
			if state == 1 {
				tool.Rdb.ZAdd(context.Background(), KeyVoteZSet(uid), &redis.Z{
					Score: float64(updateAt),
					Member: k,
				})
			}
		}
	}

	ids, err := tool.Rdb.ZRevRange(context.Background(), KeyVoteZSet(uid), int64(cursor), int64(cursor+size-1)).Result()
	if err != nil && err != redis.Nil {
		return []int64{}, result.RedisErr.WithError(err)
	}

	ret := make([]int64, len(ids))
	for k, v := range ids {
		id, err := tool.StrToInt64(v)
		if err != nil {
			id = 0
		}
		ret[k] = id
	}
	return ret, result.Ok
}