package service

import (
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/cache"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
)

func GetHotQuestions(cursor int, size int) ([]dto.HotQuestionDto, result.Result) {
	hots, res := cache.GetHotQuestions(cursor, size)
	if res.NotOK() {
		return nil, res
	}
	if len(hots) == 0 {
		return nil, result.QuestionNotFoundErr
	}
	var ques []dto.HotQuestionDto
	for _, hot := range hots {
		id, err := tool.StrToInt64(hot.Member.(string))
		if err != nil {
			continue
		}
		question, res := repository.SelectQuestionById(id)
		if res.NotOK() {
			continue
		}
		hotQues := dto.HotQuestionDto{}
		hotQues.Answer = &dto.AnswerShortDto{}
		hotQues.Answer.Creator = &dto.UserProfileDto{}
		model.Copy(&hotQues, question)
		hotQues.Heat = int(hot.Score)

		answer, _ := repository.GetBestAnswerByQues(id)
		if answer != nil {
			model.Copy(hotQues.Answer, answer)

			user, _ := repository.SelectProfileByUserId(answer.CreatorId)
			if user == nil || user.Id == 0 {
				hotQues.Answer.Creator = dto.AnonymousUser()
			} else {
				model.Copy(hotQues.Answer.Creator, user)
				hotQues.Answer.Creator.Id = user.UserId
			}
		}

		ques = append(ques, hotQues)
	}
	return ques, result.Ok
}