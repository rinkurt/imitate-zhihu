package main

import (
	"fmt"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"testing"
)

func TestHandleError(t *testing.T) {
	res := result.Ok.HandleError(result.Ok)
	fmt.Println(res.Show())
}

type Bar struct {
	a int
}

func Foo() *Bar {
	bar := Bar{a: 1}
	return &bar
}

func TestA(t *testing.T) {
	token, _ := tool.GenToken("1")
	fmt.Println(token)
}
