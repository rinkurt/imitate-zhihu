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
	AnswerCount  int
	CommentCount int
	ViewCount    int
	LikeCount    int
	GmtCreate    int64
	GmtModified  int64
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
