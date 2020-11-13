package controller

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/dto"
	"imitate-zhihu/service"
	"imitate-zhihu/tool"
	"net/http"
)

func RouteUserController(engine *gin.Engine) {
	group := engine.Group("/user")
	group.POST("/login", UserLogin)
	group.POST("/register", UserRegister)
}


func UserLogin(c *gin.Context) {
	userDto := dto.UserLoginDto{}
	err := c.BindJSON(&userDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, tool.RequestFormatErr.Show())
		return
	}
	res := service.UserLogin(&userDto)
	if res.IsOK() {
		c.JSON(http.StatusOK, res.Show())
	} else {
		c.JSON(http.StatusBadRequest, res.Show())
	}
}


func UserRegister(c *gin.Context) {
	registerDto := dto.UserRegisterDto{}
	err := c.BindJSON(&registerDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, tool.RequestFormatErr.WithData(err.Error()).Show())
		return
	}
	res := service.UserRegister(&registerDto)
	if res.IsOK() {
		c.JSON(http.StatusOK, res.Show())
	} else {
		c.JSON(http.StatusBadRequest, res.Show())
	}
}