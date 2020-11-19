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
