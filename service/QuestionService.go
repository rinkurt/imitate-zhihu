package service

import (
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
)

func GetQuestions(search string, page int, size int, order string) ([]dto.QuestionShortDto, result.Result) {
	offset := (page - 1) * size
	questions, res := repository.SelectQuestions(search, offset, size, order)
	if res.NotOK() {
		return nil, res
	}
	var questionDtos []dto.QuestionShortDto
	for _, question := range questions {
		questionDto := dto.QuestionShortDto{}
		model.Copy(&questionDto, &question)
		questionDtos = append(questionDtos, questionDto)
	}
	return questionDtos, result.Ok
}

func GetQuestionById(id int64) (*dto.QuestionDetailDto, result.Result) {
	question, res := repository.SelectQuestionById(id)
	if res.NotOK() {
		return nil, res
	}
	res = repository.AddQuestionViewCount(id, 1)
	if res.NotOK() {
		tool.Logger.Error("Failed in Adding View Count")
	}
	questionDto := &dto.QuestionDetailDto{}
	model.Copy(questionDto, question)
	user, res := GetUserProfileByUid(question.CreatorId)
	if res.NotOK() {
		user = dto.AnonymousUser()
	}
	questionDto.Creator = user
	return questionDto, result.Ok
}


func NewQuestion(userId int64, questionDto *dto.QuestionCreateDto) result.Result {
	question := &repository.Question{}
	model.Copy(question, questionDto)
	question.CreatorId = userId
	res := repository.CreateQuestion(question)
	if res.NotOK() {
		return res
	}

	questionDetailDto := &dto.QuestionDetailDto{}
	model.Copy(questionDetailDto, question)
	user, res := GetUserProfileByUid(userId)
	if res.NotOK() {
		user = dto.AnonymousUser()
	}
	questionDetailDto.Creator = user
	return result.Ok.WithData(questionDetailDto)
}