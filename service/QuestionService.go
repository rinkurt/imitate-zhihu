package service

import (
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
)

func GetQuestions(search string, cursor int64, cid int64, size int, order int) ([]dto.QuestionShortDto, result.Result) {
	//offset := (page - 1) * size
	questions, res := repository.SelectQuestions(search, cursor, cid, size, order)
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

	quesDetail := &dto.QuestionDetailDto{}
	model.Copy(quesDetail, question)
	user, _ := GetUserProfileByUid(userId)
	quesDetail.Creator = user
	return result.Ok.WithData(quesDetail)
}

func UpdateQuestionById(uid int64, qid int64, quesDto dto.QuestionCreateDto) result.Result {
	ques, res := repository.SelectQuestionById(qid)
	if res.NotOK() {
		return res
	}
	if ques.CreatorId != uid {
		return result.UnauthorizedOpr
	}

	model.Copy(ques, quesDto)
	res = repository.UpdateQuestion(ques)
	if res.NotOK() {
		return res
	}

	quesDetail := &dto.QuestionDetailDto{}
	model.Copy(quesDetail, ques)
	user, _ := GetUserProfileByUid(uid)
	quesDetail.Creator = user
	return result.Ok.WithData(quesDetail)
}