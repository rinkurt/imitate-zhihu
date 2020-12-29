package repository

import (
	"imitate-zhihu/result"
)

func UpvoteAnswer(answerId int64, userId int64) result.Result {
	return result.Ok
}