package test

import (
	"imitate-zhihu/tool"
	"testing"
)

func TestMain(m *testing.M) {
	tool.InitConfig("../config")
	tool.InitLogger()
	tool.InitDatabase("zhihu_test")
	tool.InitRedis()
	m.Run()
}
