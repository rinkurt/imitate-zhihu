package main

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/controller"
	"imitate-zhihu/middleware"
	"imitate-zhihu/tool"
)

func init() {
	//file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//if err == nil {
	//	logrus.SetOutput(file)
	//} else {
	//	logrus.Info("Failed to log to file, using default stderr")
	//}
	//logrus.SetLevel(logrus.WarnLevel)
}

func main() {
	config := tool.GetConfig()

	tool.InitDatabase()
	tool.InitLogger()

	gin.SetMode(config.Mode)
	engine := gin.Default()

	if config.LogFile != "" {
		engine.Use(middleware.LoggerToFile)
	}

	controller.RouteQuestionController(engine)
	controller.RouteUserController(engine)

	err := engine.Run(":" + config.Port)
	if err != nil {
		panic(err.Error())
	}
}
