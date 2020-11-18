package service

import (
	"fmt"
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
)

func GetQuestions(search string, page int, size int) result.Result {
	offset := (page - 1) * size
	var questions []repository.Question
	var res result.Result
	if search == "" {
		questions, res = repository.SelectQuestionsWithLimit(offset, size)
	} else {
		questions, res = repository.SelectQuestionsBySearchWithLimit(search, offset, size)
	}
	if !res.IsOK() {
		return res
	}
	var questionDtos []dto.QuestionDto
	for _, question := range questions {
		questionDto := dto.QuestionDto{}
		model.Copy(&questionDto, &question)
		userDto, res := GetUserById(question.CreatorId)
		if !res.IsOK() {
			// TODO: 匿名用户
			fmt.Println("User Not Found")
		}
		questionDto.Creator = userDto
		questionDtos = append(questionDtos, questionDto)
	}
	return result.Ok.WithData(questionDtos)
}
