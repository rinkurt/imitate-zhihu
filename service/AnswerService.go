package service

import (
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
)

func NewAnswer(userId int64, answerCreateDto *dto.AnswerCreateDto) result.Result {
	answer := &repository.Answer{}
	model.Copy(answer,answerCreateDto)
	answer.CreatorId = userId
	res := repository.CreateAnswer(answer)
	if res.NotOK() {
		return res
	}
	answerDetail := &dto.AnswerDetailDto{}
	model.Copy(answerDetail,answer)
	user,_ := GetUserProfileByUid(userId)
	answerDetail.Creator = user
	return result.Ok.WithData(answerDetail)
}

func GetAnswerById(answerId int64)  (*dto.AnswerDetailDto, result.Result){
	answerDetail := &dto.AnswerDetailDto{}
	answer, res := repository.SelectAnswerById(answerId)
	if res.NotOK() {
		return nil,res
	}
	model.Copy(answerDetail,answer)
	user,res := GetUserProfileByUid(answer.CreatorId)
	if res.NotOK() {
		user = dto.AnonymousUser()
	}
	answerDetail.Creator = user
	return answerDetail,result.Ok
}

func UpdateAnswerById(userId int64, answerId int64, answerDto *dto.AnswerCreateDto)(*dto.AnswerDetailDto,result.Result)  {
	answer,res := repository.SelectAnswerById(answerId)
	if res.NotOK() {
		return nil,result.AnswerNotFoundErr
	}
	if answer.CreatorId != userId {
		return nil,result.UnauthorizedOpr
	}
	ans,res := repository.UpdateAnswer(answerId,answerDto)
	if res.NotOK() {
		return nil,res
	}
	answerDetail := &dto.AnswerDetailDto{}
	model.Copy(answerDetail,ans)
	user, res := GetUserProfileByUid(userId)
	if res.NotOK(){
		user = dto.AnonymousUser()
	}
	answerDetail.Creator = user
	return answerDetail,result.Ok
}

func DeleteAnswerById(userId int64, answerId int64) result.Result  {
	answer,res := repository.SelectAnswerById(answerId)
	if res.NotOK() {
		return result.AnswerNotFoundErr
	}
	if userId != answer.CreatorId {
		return result.UnauthorizedOpr
	}
	res = repository.DeleteAnswerById(answerId)
	return res
}

func GetAnswers(questionId int64, cursor []int64, size int64, orderby string) (*[]dto.AnswerDetailDto, string, result.Result)  {
	ans, res, nextCursor := repository.SelectAnswers(questionId,cursor,size,orderby)
	if res.NotOK() {
		return nil,"-1,-1",result.AnswerNotFoundErr
	}
	answers := make([]dto.AnswerDetailDto,len(ans))
	for i := 0;i < len(ans);i++ {
		model.Copy(&answers[i],&ans[i])
		userProfile,res := GetUserProfileByUid(ans[i].CreatorId)
		if res.NotOK() {
			userProfile = dto.AnonymousUser()
		}
		answers[i].Creator = userProfile
	}
	return &answers,nextCursor,result.Ok
}