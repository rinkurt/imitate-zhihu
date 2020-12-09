package main

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/controller"
	"imitate-zhihu/middleware"
	"imitate-zhihu/tool"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
)


func main() {
	if pid := syscall.Getpid(); pid != 1 {
		ioutil.WriteFile("server.pid", []byte(strconv.Itoa(pid)), 0777)
		defer os.Remove("server.pid")
	}

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
