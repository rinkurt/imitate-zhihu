package service

import (
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
)

func GetQuestions(search string, page int, size int) ([]dto.QuestionShortDto, result.Result) {
	offset := (page - 1) * size
	questions, res := repository.SelectQuestions(search, offset, size)
	if res.NotOK() {
		return nil, res
	}
	var questionDtos []dto.QuestionShortDto
	for _, question := range questions {
		questionDto := dto.QuestionShortDto{}
		model.Copy(&questionDto, &question)
		//userDto, res := GetUserById(question.CreatorId)
		//if !res.IsOK() {
		//	userDto = dto.AnonymousUser()
		//}
		//questionDto.Creator = userDto
		questionDtos = append(questionDtos, questionDto)
	}
	return questionDtos, result.Ok
}

func GetQuestionById(id int) (*dto.QuestionDetailDto, result.Result) {
	question, res := repository.SelectQuestionById(id)
	if res.NotOK() {
		return nil, res
	}
	questionDto := &dto.QuestionDetailDto{}
	model.Copy(questionDto, question)
	user, res := GetUserById(question.CreatorId)
	if res.NotOK() {
		user = dto.AnonymousUser()
	}
	questionDto.Creator = user
	return questionDto, result.Ok
}


func NewQuestion(userId int, questionDto *dto.QuestionCreateDto) result.Result {
	question := &repository.Question{}
	model.Copy(question, questionDto)
	question.CreatorId = userId
	res := repository.CreateQuestion(question)
	if res.NotOK() {
		return res
	}

	questionDetailDto := &dto.QuestionDetailDto{}
	model.Copy(questionDetailDto, question)
	user, res := GetUserById(userId)
	if res.NotOK() {
		user = dto.AnonymousUser()
	}
	questionDetailDto.Creator = user
	return result.Ok.WithData(questionDetailDto)
}