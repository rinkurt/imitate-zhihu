package repository

import (
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"time"
)

type Question struct {
	Id           int `gorm:"primaryKey"`
	Title        string
	Content      string
	CreatorId    int
	Tag          string
	AnswerCount  int
	CommentCount int
	ViewCount    int
	LikeCount    int
	CreateAt     int64
	UpdateAt     int64
}


func SelectQuestions(search string, offset int, limit int) ([]Question, result.Result) {
	db := tool.GetDatabase()
	var questions []Question
	if search != "" {
		db = db.Where("title LIKE ?", "%" + search + "%").
			Or("FIND_IN_SET(?,tag)", search)
	}
	res := db.Limit(limit).Offset(offset).Find(&questions)
	if res.RowsAffected == 0 {
		return questions, result.QuestionNotFoundErr
	}
	return questions, result.Ok
}

func SelectQuestionById(id int) (*Question, result.Result) {
	db := tool.GetDatabase()
	question := Question{}
	res := db.First(&question, id)
	if res.RowsAffected == 0 {
		return nil, result.QuestionNotFoundErr
	}
	return &question, result.Ok
}

func CreateQuestion(question *Question) result.Result {
	db := tool.GetDatabase()
	question.CreateAt = time.Now().Unix()
	question.UpdateAt = question.CreateAt
	res := db.Create(question)
	if res.RowsAffected == 0 {
		return result.CreateQuestionErr
	}
	return result.Ok
}
