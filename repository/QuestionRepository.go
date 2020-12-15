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

type QuestionShortModel struct {
	Id          int64 `gorm:"primaryKey"`
	Title       string
	CreatorId   int64
	AnswerCount int
	ViewCount   int
	CreateAt    int64
	UpdateAt    int64
}

func SelectQuestions(search string, cursor int64, cid int64, limit int, order int) ([]QuestionShortModel, result.Result) {
	db := tool.GetDatabase()
	var questions []QuestionShortModel
	if search != "" {
		db = db.Where("title LIKE ? OR FIND_IN_SET(?,tag)", "%"+search+"%", search)
	}
	switch order {
	case tool.OrderByHeat:
		if cursor != 0 || cid != 0 {
			db = db.Where("(view_count = ? AND id > ?) OR view_count < ?", cursor, cid, cursor)
		}
		db = db.Order("view_count desc")
	case tool.OrderByTime:
		if cursor != 0 || cid != 0 {
			db = db.Where("(update_at = ? AND id > ?) OR update_at < ?", cursor, cid, cursor)
		}
		db = db.Order("update_at desc")
	}
	res := db.Model(&Question{}).Limit(limit).Find(&questions)
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
