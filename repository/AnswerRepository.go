package repository

import (
	"imitate-zhihu/dto"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"strconv"
	"strings"
	"time"
)

type Answer struct {
	Id           	int64 `gorm:"primaryKey"`
	Content      	string
	ViewCount    	int
	UpvoteCount  	int
	DownvoteCount  	int
	CommentCount 	int
	CreateAt     	int64
	UpdateAt     	int64
	CreatorId    	int64
	QuestionId    	int64
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

func SelectAnswerById(answerId int64) (*Answer,result.Result)  {
	db := tool.GetDatabase()
	var ans Answer
	err := db.First(&ans,answerId).Error
	if err != nil {
		return nil,result.AnswerNotFoundErr
	}
	return &ans, result.Ok
}

func UpdateAnswer(answerId int64, answer *dto.AnswerCreateDto) (*Answer,result.Result) {
	db := tool.GetDatabase()
	var ans Answer
	err := db.First(&ans,answerId).Error
	if err != nil {//查找失败
		return nil, result.UpdateAnswerErr
	}
	err = db.Model(&ans).Updates(map[string]interface{}{
		"content":answer.Content,
		"update_at":time.Now().Unix(),
	}).Error
	if err != nil {//更新失败
		return nil,result.UpdateAnswerErr
	}
	return &ans,result.Ok
}

func DeleteAnswerById(answerId int64) result.Result  {
	var answer Answer
	db := tool.GetDatabase()
	err := db.First(&answer,answerId).Error
	if err != nil {
		return result.AnswerNotFoundErr
	}
	err = db.Delete(&answer).Error
	if err != nil {
		return result.DeleteAnswerErr
	}
	return result.Ok
}

func SelectAnswers(questionId int64, cursor []int64, size int64, orderby string) ([]Answer, result.Result, string)  {
	var answers []Answer
	var nextCursor string
	db := tool.GetDatabase()
	str := make([]string,2)
	switch orderby {
	case "time":
		if cursor[1] == -1 {//说明直接查询前size个数据,并返回最后一个answer的时间戳和id
			err := db.Where("question_id = ?",questionId).Order("update_at desc").Limit(int(size)).Find(&answers).Error
			if err != nil {
				return nil, result.AnswerNotFoundErr,"-1,-1"
			}
		}else {//如果不是默认的cursor，则从cursor所指向记录的下一个记录开始查起
			err := db.Where("question_id = ? AND update_at <= ?",questionId,cursor[0]).Order("update_at desc").Not("id = ?",cursor[1]).Limit(int(size)).Find(&answers).Error
			if err != nil {
				return nil, result.AnswerNotFoundErr, "-1,-1"
			}
		}
		if len(answers) == 0 {
			return nil, result.AnswerNotFoundErr,"-1,-1"
		}
		str[0] = strconv.FormatInt(answers[len(answers)-1].UpdateAt,10)
		str[1] = strconv.FormatInt(answers[len(answers)-1].Id,10)
		nextCursor = strings.Join(str,",")//找到最后一个的时间戳+id作为下一个游标

	case "heat":
		if cursor[1] == -1 {
			err := db.Where("question_id = ?",questionId).Order("view_count desc").Limit(int(size)).Find(&answers).Error
			if err != nil {
				return nil, result.AnswerNotFoundErr,"-1,-1"
			}
		}else {
			err := db.Where("question_id = ? AND update_at <= ?",questionId,cursor[0]).Order("view_count desc").Not("id = ?",cursor[1]).Limit(int(size)).Find(&answers).Error
			if err != nil {
				return nil, result.AnswerNotFoundErr, "-1,-1"
			}
		}
		if len(answers) == 0 {
			return nil, result.AnswerNotFoundErr,"-1,-1"
		}
		str[0] = strconv.FormatInt(answers[len(answers)-1].UpdateAt,10)
		str[1] = strconv.FormatInt(answers[len(answers)-1].Id,10)
		nextCursor = strings.Join(str,",")

	case "upvote":
		if cursor[1] == -1 {
			err := db.Where("question_id = ?",questionId).Order("upvote_count desc").Limit(int(size)).Find(&answers).Error
			if err != nil {
				return nil, result.AnswerNotFoundErr,"-1,-1"
			}
		}else {
			err := db.Where("question_id = ? AND update_at <= ?",questionId,cursor[0]).Order("upvote_count desc").Not("id = ?",cursor[1]).Limit(int(size)).Find(&answers).Error
			if err != nil {
				return nil, result.AnswerNotFoundErr, "-1,-1"
			}
		}
		if len(answers) == 0 {
			return nil, result.AnswerNotFoundErr,"-1,-1"
		}
		str[0] = strconv.FormatInt(answers[len(answers)-1].UpdateAt,10)
		str[1] = strconv.FormatInt(answers[len(answers)-1].Id,10)
		nextCursor = strings.Join(str,",")
	}
	return answers,result.Ok,nextCursor
}