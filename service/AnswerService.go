package service

import (
	"context"
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/cache"
	"imitate-zhihu/dto"
	"imitate-zhihu/enum"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
)

func NewAnswer(userId int64, answerCreateDto *dto.AnswerCreateDto) result.Result {
	answer := &repository.Answer{}
	model.Copy(answer, answerCreateDto)
	answer.CreatorId = userId
	res := repository.CreateAnswer(answer)
	if res.NotOK() {
		return res
	}
	answerDetail := &dto.AnswerDetailDto{}
	model.Copy(answerDetail, answer)
	user, _ := GetUserProfileByUid(userId)
	answerDetail.Creator = user
	return result.Ok.WithData(answerDetail)
}

func GetAnswerById(id int64) (*dto.AnswerDetailDto, result.Result) {
	answerDetail := &dto.AnswerDetailDto{}
	answer, res := repository.SelectAnswerById(id)
	if res.NotOK() {
		return nil, res
	}

	cache.IncrViewCount(enum.Answer, id, enum.ViewCount, 1)

	model.Copy(answerDetail, answer)

	// Read counts from cache
	counts, err := tool.Rdb.HGetAll(context.Background(), cache.KeyWrite(enum.Answer, id)).Result()
	if err != nil {
		tool.Logger.Error(err)
	}
	for k, v := range counts {
		c, err := tool.StrToInt(v)
		if err != nil {
			tool.Logger.Error("Cache value not integer!")
			break
		}
		switch k {
		case enum.ViewCount:
			answerDetail.ViewCount += c
		case enum.UpvoteCount:
			answerDetail.UpvoteCount += c
		case enum.CommentCount:
			answerDetail.CommentCount += c
		}
	}

	user, res := GetUserProfileByUid(answer.CreatorId)
	if res.NotOK() {
		user = dto.AnonymousUser()
	}
	answerDetail.Creator = user
	return answerDetail, result.Ok
}

func UpdateAnswerById(userId int64, answerId int64, answerDto *dto.AnswerCreateDto) (*dto.AnswerDetailDto, result.Result) {
	answer, res := repository.SelectAnswerById(answerId)
	if res.NotOK() {
		return nil, result.AnswerNotFoundErr
	}
	if answer.CreatorId != userId {
		return nil, result.UnauthorizedOpr
	}
	ans, res := repository.UpdateAnswer(answerId, answerDto)
	if res.NotOK() {
		return nil, res
	}
	answerDetail := &dto.AnswerDetailDto{}
	model.Copy(answerDetail, ans)
	user, res := GetUserProfileByUid(userId)
	if res.NotOK() {
		user = dto.AnonymousUser()
	}
	answerDetail.Creator = user
	return answerDetail, result.Ok
}

func DeleteAnswerById(userId int64, answerId int64) result.Result {
	answer, res := repository.SelectAnswerById(answerId)
	if res.NotOK() {
		return result.AnswerNotFoundErr
	}
	if userId != answer.CreatorId {
		return result.UnauthorizedOpr
	}
	res = repository.DeleteAnswerById(answerId)
	return res
}

func GetAnswers(questionId int64, cursor []int64, size int, orderby string) ([]dto.AnswerDetailDto, result.Result) {
	ans, res := repository.SelectAnswers(questionId, cursor, size, orderby)
	if res.NotOK() {
		return nil, result.AnswerNotFoundErr
	}
	answers := make([]dto.AnswerDetailDto, len(ans))
	for i := 0; i < len(ans); i++ {
		model.Copy(&answers[i], &ans[i])
		userProfile, res := GetUserProfileByUid(ans[i].CreatorId)
		if res.NotOK() {
			userProfile = dto.AnonymousUser()
		}
		answers[i].Creator = userProfile
	}
	return answers, result.Ok
}
