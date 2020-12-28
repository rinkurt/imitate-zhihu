package test

import "imitate-zhihu/tool"

func init() {
	tool.InitConfig("../config")
	tool.InitLogger()
	tool.InitDatabase("zhihu_test")
	tool.InitRedis()
}
