package test

import (
	"fmt"
	"imitate-zhihu/service"
	"testing"
)

func TestGetHotQuestions(t *testing.T)  {
	testCursor := 0
	testSize := 10
	resHotQuestionDto,res := service.GetHotQuestions(testCursor,testSize)
	if res.NotOK() {
		t.Error(res.Message)
	}
	for i,v := range resHotQuestionDto {
		fmt.Printf("HotQuestion%d:%v\n",i,v)
	}
}