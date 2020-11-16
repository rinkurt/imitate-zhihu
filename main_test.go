package main

import (
	"fmt"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"testing"
)

func TestHandleError(t *testing.T) {
	res := result.Ok.HandleError(result.Ok)
	fmt.Println(res.Show())
}

func TestA(t *testing.T) {
	questions, res := repository.SelectAllQuestions()
	fmt.Println(res)
	fmt.Println(questions)
}
