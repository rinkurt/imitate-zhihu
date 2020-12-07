package repository

import (
	"gorm.io/gorm"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"time"
)

type Question struct {
	Id           int64 `gorm:"primaryKey"`
	Title        string
	Content      string
	CreatorId    int64
	Tag          string
	AnswerCount  int
	CommentCount int
	ViewCount    int
	LikeCount    int
	CreateAt     int64
	UpdateAt     int64
}


func SelectQuestions(search string, offset int, limit int, order string) ([]Question, result.Result) {
	db := tool.GetDatabase()
	var questions []Question
	if search != "" {
		db = db.Where("title LIKE ?", "%" + search + "%").
			Or("FIND_IN_SET(?,tag)", search)
	}
	switch order {
	case "heat":
		db = db.Order("view_count desc")
		// TODO: 更加丰富的热度判断标准
	case "time":
		db = db.Order("update_at desc")
	}
	res := db.Limit(limit).Offset(offset).Find(&questions)
	if res.RowsAffected == 0 {
		return questions, result.QuestionNotFoundErr
	}
	return questions, result.Ok
}

func SelectQuestionById(id int64) (*Question, result.Result) {
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

func AddQuestionViewCount(id int64, count int) result.Result {
	db := tool.GetDatabase()
	res := db.Model(&Question{Id: id}).Update("view_count",
		gorm.Expr("view_count + ?", count))
	if res.RowsAffected == 0 {
		return result.UpdateViewCountErr
	}
	return result.Ok
}