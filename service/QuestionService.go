package service

import (
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
)

func GetQuestions(search string, page int, size int) result.Result {
	offset := (page - 1) * size
	questions, res := repository.SelectQuestions(search, offset, size)
	if !res.IsOK() {
		return res
	}
	var questionDtos []dto.QuestionDetailDto
	for _, question := range questions {
		questionDto := dto.QuestionDetailDto{}
		model.Copy(&questionDto, &question)
		userDto, res := GetUserById(question.CreatorId)
		if !res.IsOK() {
			userDto = dto.AnonymousUser()
		}
		questionDto.Creator = userDto
		questionDtos = append(questionDtos, questionDto)
	}
	return result.Ok.WithData(questionDtos)
}

func GetQuestionById(id int) result.Result {
	question := repository.SelectQuestionById(id)
	questionDto := dto.QuestionDetailDto{}
	model.Copy(&questionDto, &question)
	user, res := GetUserById(question.CreatorId)
	if !res.IsOK() {
		user = dto.AnonymousUser()
	}
	questionDto.Creator = user
	return result.Ok.WithData(questionDto)
}


func NewQuestion(userId int, questionDto *dto.QuestionCreateDto) result.Result {
	question := repository.Question{}
	model.Copy(&question, questionDto)
	question.CreatorId = userId
	res := repository.CreateQuestion(&question)
	if !res.IsOK() {
		return res
	}

	questionDetailDto := dto.QuestionDetailDto{}
	model.Copy(&questionDetailDto, &question)
	user, res := GetUserById(userId)
	if !res.IsOK() {
		user = dto.AnonymousUser()
	}
	questionDetailDto.Creator = user
	return result.Ok.WithData(questionDetailDto)
}