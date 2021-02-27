package test

import (
	"fmt"
	"imitate-zhihu/dto"
	"imitate-zhihu/service"
	"testing"
)

var aid int64

func TestNewAnswer(t *testing.T)  {
	TestUserRegister(t)
	TestUserLogin(t)
	TestNewQuestion(t)

	testAnswer := dto.AnswerCreateDto{
		QuestionId: qid,
		Content:    "这是一条新回答",
	}
	testId := Uid
	res := service.NewAnswer(testId,&testAnswer)
	if res.NotOK() {
		t.Error(res.Message)
	}
	dat := res.Data.(*dto.AnswerDetailDto)
	fmt.Println(dat)
	aid = dat.Id
	fmt.Println(dat.Id)
}

func TestGetAnswers(t *testing.T)  {
	testQid := qid
	testUid := Uid
	testCursor := []int64{-1,-1}
	testOrder := "heat"
	testSize := 5
	answers,res := service.GetAnswers(testQid,testUid,testCursor,testSize,testOrder)
	if res.NotOK() {
		t.Error(res.Message)
	}
	for i,v := range answers {
		fmt.Printf("answer%d:%v\n",i+1,v)
	}
}

func TestGetAnswerById(t *testing.T)  {
	testId := aid
	testOrder := "upvote"
	var resAnswerDetail *dto.AnswerDetailDto
	resAnswerDetail,res := service.GetAnswerById(testId,testOrder)
	if res.NotOK() {
		t.Error(res.Message)
	}
	fmt.Println(resAnswerDetail)
}

func TestUpdateAnswerById(t *testing.T)  {
	testAnswer := dto.AnswerCreateDto{
		QuestionId: qid,
		Content:    "这是一条更新后的新回答",
	}
	testId := Uid
	testAid := qid
	resAnswerDetail,res := service.UpdateAnswerById(testId,testAid,&testAnswer)
	if res.NotOK() {
		t.Error(res.Message)
	}
	fmt.Println(resAnswerDetail)
}

func TestGetAnswersByVoteUser(t *testing.T)  {
	testId := Uid
	testCursor := 0
	testSize := 5
	resAnswerDetail,res := service.GetAnswersByVoteUser(testId,testCursor,testSize)
	if res.NotOK() {
		t.Error(res.Message)
	}
	for i,v := range resAnswerDetail {
		fmt.Printf("answer%d:%v\n",i+1,v)
	}
}

func TestDeleteAnswerById(t *testing.T)  {
	testId := Uid
	testAid := aid
	res := service.DeleteAnswerById(testId,testAid)
	if res.NotOK() {
		t.Error(res.Message)
	}
}