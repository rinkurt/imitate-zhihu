package main

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/controller"
	"imitate-zhihu/tool"
)


func main() {
	config := tool.GetConfig()

	tool.InitDatabase()

	gin.SetMode(config.Mode)
	engine := gin.Default()

	controller.RouteQuestionController(engine)
	controller.RouteUserController(engine)

	err := engine.Run(":" + config.Port)
	if err != nil {
		panic(err.Error())
	}
}
