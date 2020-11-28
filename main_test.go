package main

import (
	"fmt"
	"imitate-zhihu/result"
	"testing"
)

func TestHandleError(t *testing.T) {

}

type Bar struct {
	a int
}

func Foo() *Bar {
	bar := Bar{a: 1}
	return &bar
}

type Res *result.Result

func TestA(t *testing.T) {
	result.Ok.WithData("aaaaa")
	fmt.Println(result.Ok.Data)
}
