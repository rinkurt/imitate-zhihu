package repository

import (
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
)

type Question struct {
	Id           int
	Title        string
	Description  string
	CreatorId    int
	Tag          string
	CommentCount int
	ViewCount    int
	LikeCount    int
	GmtCreate    int64
	GmtModified  int64
}

func SelectQuestionsWithLimit(offset int, limit int) ([]Question, result.Result) {
	db := tool.GetDatabase()
	var questions []Question
	res := db.Limit(limit).Offset(offset).Find(&questions)
	if res.RowsAffected == 0 {
		return questions, result.QuestionNotFoundErr
	}
	return questions, result.Ok
}

func SelectQuestionsBySearchWithLimit(search string,
	offset int, limit int) ([]Question, result.Result) {
	db := tool.GetDatabase()
	var questions []Question
	res := db.Where("title LIKE ?", search).Or("tag LIKE ?", search).
		Limit(limit).Offset(offset).Find(&questions)
	if res.RowsAffected == 0 {
		return questions, result.QuestionNotFoundErr
	}
	return questions, result.Ok
}

func SelectQuestionById(id int) Question {
	db := tool.GetDatabase()
	question := Question{}
	db.First(&question, id)
	return question
}

func CreateQuestion(question Question) {
	db := tool.GetDatabase()
	db.Create(&question)
}
