package test

import (
	"fmt"
	"imitate-zhihu/dto"
	"imitate-zhihu/service"
	"testing"
)

var qid int64

func TestNewQuestion(t *testing.T)  {
	TestUserRegister(t)
	TestUserLogin(t)

	testQuestion  := dto.QuestionCreateDto{
		Title:   "如何看待中国量子计算原型机九章问世？",
		Content: "这毫无疑问是中国量子计算方面的卓越进展。",
		Tag:     "中国、中国科学技术大学、量子计算与量子信息",
	}
	testID := Uid
	res := service.NewQuestion(testID,&testQuestion)
	if res.NotOK() {
		t.Error(res.Message)
	}
	dat := res.Data.(*dto.QuestionDetailDto)
	fmt.Println(dat)
	qid = dat.Id
	fmt.Println(dat.Id)
}

func TestGetQuestions(t *testing.T)  {
	testCursor := []int64{-1,-1}
	testSearch := "中国"
	testOrder := "heat"
	questions,res := service.GetQuestions(Uid,testSearch,testCursor,5,testOrder)
	if res.NotOK() {
		t.Error(res.Message)
	}
	for i,v := range questions {
		fmt.Printf("question%d:%v\n",i+1,v)
	}
}

func TestGetQuestionById(t *testing.T)  {
	testID := qid
	var resQuestionDetail *dto.QuestionDetailDto
	resQuestionDetail,res := service.GetQuestionById(testID)
	if res.NotOK() {
		t.Error(res.Message)
	}
	fmt.Println(resQuestionDetail)
}

func TestGetQuestionTitle(t *testing.T)  {
	testID := qid
	var resQuestionTitle *dto.QuestionTitleDto
	resQuestionTitle,res := service.GetQuestionTitle(testID)
	if res.NotOK() {
		t.Error(res.Message)
	}
	fmt.Println(resQuestionTitle)
}

func TestUpdateQuestionById(t *testing.T)  {
	testQuestion  := dto.QuestionCreateDto{
		Title:   "有什么烧脑的电影推荐？",
		Content: "有什么烧脑的电影推荐？",
		Tag:     "电影、电影推荐",
	}
	res := service.UpdateQuestionById(Uid,qid,testQuestion)
	if res.NotOK() {
		t.Error(res.Message)
	}
	dat := res.Data.(*dto.QuestionDetailDto)
	fmt.Println(dat)
}

func TestDeleteQuestionById(t *testing.T)  {
	res := service.DeleteQuestionById(Uid,qid)
	if res.NotOK() {
		t.Error(res.Message)
	}
}



