package controller

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/cache"
	"imitate-zhihu/middleware"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"net/http"
)

func RouteVoteController(engine *gin.Engine) {
	group := engine.Group("/vote")
	group.GET("/answer/:aid", middleware.JWTAuthMiddleware, GetAnswerVoteStatus)
	group.PUT("/answer/:aid", middleware.JWTAuthMiddleware, SetAnswerVoteStatus)
}

func GetAnswerVoteStatus(c *gin.Context) {
	uid, err := middleware.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.ContextErr.WithError(err))
		return
	}
	sAid := c.Param("aid")
	aid, err := tool.StrToInt64(sAid)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	res := cache.GetAnswerVoteStatus(uid, aid)
	if res.NotOK() {
		res = res.WithData(0)
	}
	c.JSON(http.StatusOK, res)
}

func SetAnswerVoteStatus(c *gin.Context) {
	uid, err := middleware.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.ContextErr.WithError(err))
		return
	}
	sAid := c.Param("aid")
	aid, err := tool.StrToInt64(sAid)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithErrorStr("invalid aid"))
		return
	}
	ss := c.Query("s")
	s, err := tool.StrToInt(ss)
	if err != nil || s < 0 || s > 2 {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithErrorStr("invalid status"))
		return
	}
	res := cache.SetAnswerVoteStatus(uid, aid, s)
	c.JSON(http.StatusOK, res)
}