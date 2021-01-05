package repository

import (
	"gorm.io/gorm"
	"imitate-zhihu/dto"
	"imitate-zhihu/enum"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"time"
)

type Answer struct {
	Id            int64 `gorm:"primaryKey"`
	Content       string
	ViewCount     int
	UpvoteCount   int
	DownvoteCount int
	CommentCount  int
	CreateAt      int64
	UpdateAt      int64
	CreatorId     int64
	QuestionId    int64
}

type HotAnswer struct {
	QuestionId int64
	Heat       int
}

func CreateAnswer(answer *Answer) result.Result {
	db := tool.GetDatabase()
	answer.CreateAt = time.Now().Unix()
	answer.UpdateAt = answer.CreateAt
	res := db.Create(answer)
	if res.RowsAffected == 0 {
		return result.CreateAnswerErr
	}
	return result.Ok
}

func SelectAnswerById(answerId int64) (*Answer, result.Result) {
	db := tool.GetDatabase()
	var ans Answer
	err := db.First(&ans, answerId).Error
	if err != nil {
		return nil, result.AnswerNotFoundErr
	}
	return &ans, result.Ok
}

func UpdateAnswer(answerId int64, answer *dto.AnswerCreateDto) (*Answer, result.Result) {
	db := tool.GetDatabase()
	var ans Answer
	err := db.First(&ans, answerId).Error
	if err != nil { //查找失败
		return nil, result.UpdateAnswerErr
	}
	err = db.Model(&ans).Updates(map[string]interface{}{
		"content":   answer.Content,
		"update_at": time.Now().Unix(),
	}).Error
	if err != nil { //更新失败
		return nil, result.UpdateAnswerErr
	}
	return &ans, result.Ok
}

func DeleteAnswerById(answerId int64) result.Result {
	var answer Answer
	db := tool.GetDatabase()
	err := db.First(&answer, answerId).Error
	if err != nil {
		return result.AnswerNotFoundErr
	}
	err = db.Delete(&answer).Error
	if err != nil {
		return result.DeleteAnswerErr
	}
	return result.Ok
}

func SelectAnswers(questionId int64, cursor []int64, size int, orderBy string) ([]Answer, result.Result) {
	var answers []Answer
	db := tool.GetDatabase()
	db = db.Where("question_id = ?", questionId)
	switch orderBy {
	case enum.ByTime:
		if cursor[1] != -1 {
			db = db.Where("(update_at = ? AND id > ?) OR update_at < ?", cursor[0], cursor[1], cursor[0])
		}
		db = db.Order("update_at desc")
	case enum.ByHeat:
		if cursor[1] != -1 {
			db = db.Where("(view_count = ? AND id > ?) OR view_count < ?", cursor[0], cursor[1], cursor[0])
		}
		db = db.Order("view_count desc")
	case enum.ByUpvote:
		if cursor[1] != -1 {
			db = db.Where("(upvote_count = ? AND id > ?) OR upvote_count < ?", cursor[0], cursor[1], cursor[0])
		}
		db = db.Order("upvote_count desc")
	}
	db = db.Limit(size).Find(&answers)
	if db.RowsAffected == 0 {
		return nil, result.AnswerNotFoundErr
	}
	return answers, result.Ok
}

func UpdateAnswerCounts(answer *Answer) result.Result {
	db := tool.GetDatabase()
	db = db.Model(answer).Updates(map[string]interface{}{
		"view_count":     gorm.Expr("view_count + ?", answer.ViewCount),
		"upvote_count":   gorm.Expr("upvote_count + ?", answer.UpvoteCount),
		"downvote_count": gorm.Expr("downvote_count + ?", answer.DownvoteCount),
		"comment_count":  gorm.Expr("comment_count + ?", answer.CommentCount),
	})
	if db.RowsAffected == 0 {
		return result.UpdateAnswerErr
	}
	return result.Ok
}

func GetAnswerHeatsGroupByQuestion() ([]HotAnswer, result.Result) {
	db := tool.GetDatabase()
	var hots []HotAnswer
	db = db.Model(Answer{}).Select("question_id, upvote_count * 2 as heat").
		Group("question_id").Find(&hots)
	if db.Error != nil {
		return nil, result.HandleServerErr(db.Error)
	}
	return hots, result.Ok
}

func GetBestAnswerByQues(qid int64) (*Answer, result.Result) {
	db := tool.GetDatabase()
	answer := &Answer{}
	db.Where(&Answer{QuestionId: qid}).Order("upvote_count desc").Take(answer)
	if db.RowsAffected == 0 {
		return nil, result.AnswerNotFoundErr
	}
	return answer, result.Ok
}
