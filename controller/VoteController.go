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
	c.JSON(http.StatusOK, res)
}