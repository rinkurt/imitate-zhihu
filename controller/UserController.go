package controller

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/dto"
	"imitate-zhihu/result"
	"imitate-zhihu/service"
	"net/http"
	"strconv"
)

func RouteUserController(engine *gin.Engine) {
	group := engine.Group("/user")
	group.POST("/login", UserLogin)
	group.POST("/register", UserRegister)
	group.GET("/profile/:user_id", GetUserProfile)
	group.GET("/verify", VerifyEmail)
}


func UserLogin(c *gin.Context) {
	userDto := dto.UserLoginDto{}
	resp := &dto.LoginResponseDto{}
	err := c.ShouldBindJSON(&userDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	res := service.UserLogin(&userDto)
	if res.NotOK() {
		res = res.WithData(resp)
	}
	c.JSON(http.StatusOK, res)
}


func UserRegister(c *gin.Context) {
	registerDto := dto.UserRegisterDto{}
	resp := &dto.LoginResponseDto{}
	err := c.ShouldBindJSON(&registerDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	res := service.UserRegister(&registerDto)
	if res.NotOK() {
		res = res.WithData(resp)
	}
	if res.IsServerErr() {
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	c.JSON(http.StatusOK, res)
}

func GetUserProfile(c *gin.Context) {
	sUserId := c.Param("user_id")
	userId, err := strconv.ParseInt(sUserId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	profile, res := service.GetUserProfileByUid(userId)
	if res.NotOK() {
		res = res.WithData(dto.AnonymousUser())
	}
	c.JSON(http.StatusOK, res.WithData(profile))
}

func VerifyEmail(c *gin.Context) {
	email := c.Query("email")
	_, res := service.VerifyEmail(email)
	c.JSON(http.StatusOK, res)
}