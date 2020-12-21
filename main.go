package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"imitate-zhihu/cache"
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

	c := cron.New(cron.WithSeconds())
	c.AddFunc("@every 3h", cache.SyncCount)
	c.Start()

	gin.SetMode(tool.Cfg.Mode)
	engine := gin.Default()

	if tool.Cfg.LogFile != "" {
		engine.Use(middleware.LoggerToFile)
	}

	controller.RouteQuestionController(engine)
	controller.RouteUserController(engine)

	err := engine.Run(":" + tool.Cfg.Port)
	if err != nil {
		panic(err.Error())
	}

}
