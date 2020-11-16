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

func SelectAllQuestions() ([]Question, result.Result) {
	db := tool.GetDatabase()
	var users []Question
	res := db.Find(&users)
	if res.RowsAffected == 0 {
		return nil, result.QuestionNotFoundErr.HandleError(res.Error)
	}
	return users, result.Ok
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
