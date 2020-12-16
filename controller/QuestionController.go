package controller

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/dto"
	"imitate-zhihu/middleware"
	"imitate-zhihu/result"
	"imitate-zhihu/service"
	"imitate-zhihu/tool"
	"net/http"
	"strconv"
	"strings"
)

func RouteQuestionController(engine *gin.Engine) {
	group := engine.Group("/question")
	group.GET("", GetQuestions)
	group.GET("/:question_id", GetQuestionById)
	group.POST("", middleware.JWTAuthMiddleware, NewQuestion)
	group.PUT("/:question_id", middleware.JWTAuthMiddleware, UpdateQuestionById)
}

func GetQuestions(c *gin.Context) {
	cursor := c.Query("cursor")
	split := strings.Split(cursor, ",")
	var cur int64 = 0
	var cid int64 = 0
	if len(split) == 2 {
		cur, _ = tool.StringToInt64(split[0])
		cid, _ = tool.StringToInt64(split[1])
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		size = 10
	}
	search := c.Query("search")
	orderBy := c.Query("orderby")
	var order int
	switch orderBy {
	case "time":
		order = tool.OrderByTime
	case "heat":
		order = tool.OrderByHeat
	default:
		order = tool.OrderByTime
	}
	q, res := service.GetQuestions(search, cur, cid, size, order)
	nextCursor := ""
	if len(q) > 0 {
		tail := q[len(q)-1]
		switch order {
		case tool.OrderByTime:
			nextCursor = tool.Int64ToString(tail.UpdateAt) + "," + tool.Int64ToString(tail.Id)
		case tool.OrderByHeat:
			nextCursor = strconv.Itoa(tail.ViewCount) + "," + tool.Int64ToString(tail.Id)
		}
	}
	c.JSON(http.StatusOK, res.WithData(gin.H{
		"next_cursor": nextCursor,
		"questions": q,
	}))
}


func GetQuestionById(c *gin.Context) {
	qid := c.Param("question_id")
	id, err := tool.StringToInt64(qid)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	q, res := service.GetQuestionById(id)
	c.JSON(http.StatusOK, res.WithData(q))
}


func NewQuestion(c *gin.Context) {
	userId, err := middleware.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.TokenErr.WithError(err))
		return
	}
	questionDto := dto.QuestionCreateDto{}
	err = c.ShouldBindJSON(&questionDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	res := service.NewQuestion(userId, &questionDto)
	c.JSON(http.StatusOK, res)
}

func UpdateQuestionById(c *gin.Context) {
	uid, err := middleware.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.TokenErr.WithError(err))
		return
	}
	sQid := c.Param("question_id")
	qid, err := tool.StringToInt64(sQid)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	questionDto := dto.QuestionCreateDto{}
	err = c.ShouldBindJSON(&questionDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	res := service.UpdateQuestionById(uid, qid, questionDto)
	c.JSON(http.StatusOK, res)
}