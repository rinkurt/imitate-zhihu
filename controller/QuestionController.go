package controller

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/dto"
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
	group.POST("", JWTAuthMiddleware, NewQuestion)
}

func GetQuestions(c *gin.Context) {
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page == 0 {
		page = 1
	}
	size, err := strconv.Atoi(c.Query("size"))
	if err != nil || size < 10 {
		size = 10
	}
	search := c.Query("search")
	q, res := service.GetQuestions(search, page, size)
	c.JSON(http.StatusOK, res.WithData(q))
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
	sUserId, exists := c.Get("user_id")
	userId, err := tool.StringToInt64(sUserId.(string))
	if err != nil || !exists {
		c.JSON(http.StatusInternalServerError,
			result.ContextErr.WithErrorStr("get user_id failed"))
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