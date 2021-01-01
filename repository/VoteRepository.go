package repository

import (
	"imitate-zhihu/tool"
)

type AnswerVote struct {
	Id       int64 `gorm:"primaryKey"`
	IsUpvote bool
	AnswerId int64
	UserId   int64
	UpdateAt int64
}


func SelectVotesByUid(uid int64) []AnswerVote {
	db := tool.GetDatabase()
	var votes []AnswerVote
	db = db.Where(&AnswerVote{UserId: uid}).Find(&votes)
	return votes
}
