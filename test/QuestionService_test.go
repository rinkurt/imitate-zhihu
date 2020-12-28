package test

import (
	"fmt"
	"imitate-zhihu/dto"
	"imitate-zhihu/service"
	"testing"
)

func TestNewQuestion(t *testing.T)  {
	testQuestion  := dto.QuestionCreateDto{
		Title:   "如何看待中国量子计算原型机九章问世？",
		Content: "这毫无疑问是中国量子计算方面的卓越进展。",
		Tag:     "中国、中国科学技术大学、量子计算与量子信息",
	}
	testID := int64(5)
	res := service.NewQuestion(testID,&testQuestion)
	if res.NotOK() {
		t.Error(res.Message)
	}
}

func TestGetQuestionById(t *testing.T)  {
	testID := int64(1)
	var resQuestionDetail *dto.QuestionDetailDto
	resQuestionDetail,res := service.GetQuestionById(testID)
	if res.NotOK() {
		t.Error(res.Message)
	}
	fmt.Println(resQuestionDetail)
}
