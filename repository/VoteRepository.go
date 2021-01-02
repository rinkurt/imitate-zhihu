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


func CreateVotes(votes []AnswerVote) {
	db := tool.GetDatabase()
	db.Create(votes)
}


func DeleteVotes(uid int64, aids []int64) {
	db := tool.GetDatabase()
	db = db.Where(map[string]interface{}{
		"user_id": uid,
		"answer_id": aids,
	})
	db = db.Delete(&AnswerVote{})
}
