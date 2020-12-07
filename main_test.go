package main

import (
	"fmt"
	"testing"
)

type Temp interface {
	FuncA()
	FuncB()
}

type Impl struct {

}

func (i *Impl) FuncA() {
	fmt.Println("Impl:FuncA")
}

func (i *Impl) FuncB() {
	fmt.Println("Impl:FuncB")
}

type Bar struct {
	*Impl
}

func (b *Bar) FuncA() {
	fmt.Println("Bar:FuncA")
}

func Func(t Temp) {
	t.FuncA()
	t.FuncB()
}

func TestA(t *testing.T) {
	b := Bar{}
	Func(&b)
}
