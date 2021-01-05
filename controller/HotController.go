package controller

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/result"
	"imitate-zhihu/service"
	"imitate-zhihu/tool"
	"net/http"
)

func RouteHotController(engine *gin.Engine) {
	group := engine.Group("/hot")
	group.GET("", GetHotQuestions)
}

func GetHotQuestions(c *gin.Context) {
	cursor, err := tool.StrToInt(c.DefaultQuery("cursor", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithErrorStr("Cursor format error"))
		return
	}
	size, err := tool.StrToInt(c.DefaultQuery("size", "10"))
	if err != nil {
		size = 10
	}
	q, res := service.GetHotQuestions(cursor, size)
	if res.NotOK() {
		c.JSON(http.StatusOK, res)
		return
	}
	c.JSON(http.StatusOK, res.WithData(gin.H{
		"next_cursor": cursor + size,
		"questions": q,
	}))
}
