package service

import (
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/cache"
	"imitate-zhihu/dto"
	"imitate-zhihu/enum"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
)

func NewAnswer(userId int64, answerCreateDto *dto.AnswerCreateDto) result.Result {
	answer := &repository.Answer{}
	model.Copy(answer, answerCreateDto)
	answer.CreatorId = userId
	res := repository.CreateAnswer(answer)
	if res.NotOK() {
		return res
	}

	// incr answer count
	res = cache.IncrCount(enum.Question, answerCreateDto.QuestionId, enum.AnswerCount, 1)
	if res.NotOK() {
		return res
	}

	answerDetail := &dto.AnswerDetailDto{}
	model.Copy(answerDetail, answer)
	user, _ := GetUserProfileByUid(userId)
	answerDetail.Creator = user
	return result.Ok.WithData(answerDetail)
}

func GetAnswerById(id int64, order string) (*dto.AnswerDetailDto, result.Result) {
	answerDetail := &dto.AnswerDetailDto{}
	answer, res := repository.SelectAnswerById(id)
	if res.NotOK() {
		return nil, res
	}

	res = cache.IncrCount(enum.Answer, id, enum.ViewCount, 1)
	if res.NotOK() {
		return nil, res
	}

	res = cache.ReadAnswerCounts(answer)
	if res.NotOK() {
		return nil, res
	}

	model.Copy(answerDetail, answer)

	user, res := GetUserProfileByUid(answer.CreatorId)
	if res.NotOK() {
		user = dto.AnonymousUser()
	}
	answerDetail.Creator = user
	answerDetail.NextId = repository.GetNextAnswerId(answer, order)
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
	if res.NotOK() {
		return res
	}

	// decr answer count
	res = cache.IncrCount(enum.Question, answer.QuestionId, enum.AnswerCount, -1)
	return res
}

func GetAnswers(questionId int64, userId int64, cursor []int64, size int, orderBy string) ([]dto.AnswerDetailDto, result.Result) {
	ans, res := repository.SelectAnswers(questionId, userId, cursor, size, orderBy)
	if res.NotOK() {
		return nil, result.AnswerNotFoundErr
	}
	answers := make([]dto.AnswerDetailDto, len(ans))

	profile := &dto.UserProfileDto{}
	if userId != 0 {
		profile, _ = GetUserProfileByUid(userId)
	}

	for i := 0; i < len(ans); i++ {
		res = cache.ReadAnswerCounts(&ans[i])
		if res.NotOK() {
			return nil, res
		}

		model.Copy(&answers[i], &ans[i])
		if userId != 0 {
			answers[i].Creator = profile
		} else {
			userProfile, _ := GetUserProfileByUid(ans[i].CreatorId)
			answers[i].Creator = userProfile
		}
	}
	return answers, result.Ok
}

func GetAnswersByVoteUser(uid int64, cursor int, size int) ([]dto.AnswerDetailDto, result.Result) {
	ids, res := cache.GetVotedAnswerIds(uid, cursor, size)
	if res.NotOK() {
		return nil, res
	}
	var ret []dto.AnswerDetailDto
	for _, id := range ids {
		if id == 0 {
			continue
		}
		answer, _ := repository.SelectAnswerById(id)
		if answer == nil {
			continue
		}
		ansDto := dto.AnswerDetailDto{}
		model.Copy(&ansDto, answer)

		profile, _ := GetUserProfileByUid(answer.CreatorId)
		ansDto.Creator = profile

		ret = append(ret, ansDto)
	}
	if len(ret) == 0 {
		return nil, result.AnswerNotFoundErr
	}
	return ret, result.Ok
}
