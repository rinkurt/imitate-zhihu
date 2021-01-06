package repository

import (
	"gorm.io/gorm"
	"imitate-zhihu/enum"
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
	Content     string
	CreatorId   int64
	AnswerCount int
	ViewCount   int
	CreateAt    int64
	UpdateAt    int64
}

type Hot struct {
	Id   int64 `gorm:"primaryKey"`
	Heat int
}

func SelectQuestions(uid int64, search string, cursor []int64, limit int, orderBy string) ([]Question, result.Result) {
	db := tool.GetDatabase()
	var questions []Question
	db = db.Where(&Question{CreatorId: uid})
	if search != "" {
		db = db.Where("title LIKE ? OR FIND_IN_SET(?,tag)", "%"+search+"%", search)
	}
	switch orderBy {
	case enum.ByHeat:
		if cursor[1] != -1 {
			db = db.Where("(view_count = ? AND id < ?) OR view_count < ?", cursor[0], cursor[1], cursor[0])
		}
		db = db.Order("view_count desc, id desc")
	case enum.ByTime:
		if cursor[1] != -1 {
			db = db.Where("(update_at = ? AND id < ?) OR update_at < ?", cursor[0], cursor[1], cursor[0])
		}
		db = db.Order("update_at desc, id desc")
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

func UpdateQuestion(question *Question) result.Result {
	db := tool.GetDatabase()
	question.UpdateAt = time.Now().Unix()
	db = db.Save(question)
	if db.RowsAffected == 0 {
		return result.UpdateQuestionErr
	}
	return result.Ok
}

func UpdateQuestionCounts(question *Question) result.Result {
	db := tool.GetDatabase()
	db = db.Model(question).Updates(map[string]interface{}{
		"answer_count":  gorm.Expr("answer_count + ?", question.AnswerCount),
		"comment_count": gorm.Expr("comment_count + ?", question.CommentCount),
		"view_count":    gorm.Expr("view_count + ?", question.ViewCount),
		"like_count":    gorm.Expr("like_count + ?", question.LikeCount),
	})
	if db.RowsAffected == 0 {
		return result.UpdateQuestionErr
	}
	return result.Ok
}

func DeleteQuestionById(id int64) result.Result {
	db := tool.GetDatabase()
	db = db.Delete(&Question{}, id)
	if db.RowsAffected == 0 {
		return result.DeleteQuestionErr
	}
	return result.Ok
}

func GetQuestionHeats() ([]Hot, result.Result) {
	db := tool.GetDatabase()
	var hots []Hot
	db = db.Model(&Question{}).Select("id, (view_count + answer_count * 5) as heat").Find(&hots)
	if db.Error != nil {
		return nil, result.HandleServerErr(db.Error)
	}
	return hots, result.Ok
}