package controller

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/dto"
	"imitate-zhihu/enum"
	"imitate-zhihu/middleware"
	"imitate-zhihu/result"
	"imitate-zhihu/service"
	"imitate-zhihu/tool"
	"net/http"
	"strconv"
)

func RouteQuestionController(engine *gin.Engine) {
	group := engine.Group("/question")
	group.GET("", GetQuestions)
	group.GET("/:question_id", GetQuestionById)
	group.POST("", middleware.JWTAuthMiddleware, NewQuestion)
	group.PUT("/:question_id", middleware.JWTAuthMiddleware, UpdateQuestionById)
	group.DELETE("/:question_id", middleware.JWTAuthMiddleware, DeleteQuestionById)
}

func GetQuestions(c *gin.Context) {
	cursor, err := tool.ParseCursor(c.DefaultQuery("cursor", "-1,-1"))
	if err != nil || len(cursor) < 2 {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithErrorStr("Cursor format error"))
		return
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil {
		size = 10
	}
	search := c.Query("search")
	orderBy := c.DefaultQuery("orderby", enum.ByTime)
	q, res := service.GetQuestions(search, cursor, size, orderBy)
	nextCursor := ""
	if len(q) > 0 {
		tail := q[len(q)-1]
		switch orderBy {
		case enum.ByTime:
			nextCursor = tool.Int64ToStr(tail.UpdateAt) + "," + tool.Int64ToStr(tail.Id)
		case enum.ByHeat:
			nextCursor = strconv.Itoa(tail.ViewCount) + "," + tool.Int64ToStr(tail.Id)
		}
	}
	c.JSON(http.StatusOK, res.WithData(gin.H{
		"next_cursor": nextCursor,
		"questions": q,
	}))
}


func GetQuestionById(c *gin.Context) {
	qid := c.Param("question_id")
	id, err := tool.StrToInt64(qid)
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
		c.JSON(http.StatusUnauthorized, result.ContextErr.WithError(err))
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
		c.JSON(http.StatusUnauthorized, result.ContextErr.WithError(err))
		return
	}
	sQid := c.Param("question_id")
	qid, err := tool.StrToInt64(sQid)
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

func DeleteQuestionById(c *gin.Context) {
	uid, err := middleware.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.ContextErr.WithError(err))
		return
	}
	sQid := c.Param("question_id")
	qid, err := tool.StrToInt64(sQid)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	res := service.DeleteQuestionById(uid, qid)
	c.JSON(http.StatusOK, res)
}