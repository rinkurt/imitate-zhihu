package controller

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/service"
	"net/http"
	"strconv"
)

func RouteQuestionController(engine *gin.Engine) {
	group := engine.Group("/question")
	group.GET("", GetQuestions)
	group.GET("/:question_id", GetQuestionById)
	group.POST("", NewQuestion)
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
	res := service.GetQuestions(search, page, size)
	c.JSON(http.StatusOK, res.Show())
}


// TODO
func GetQuestionById(c *gin.Context) {
	qid := c.Param("question_id")
	id, err := strconv.Atoi(qid)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ShowBadRequest(err.Error()))
	}
	question := repository.SelectQuestionById(id)
	c.JSON(http.StatusOK, question)
}

// TODO
func NewQuestion(c *gin.Context) {
	questionDto := dto.QuestionDto{}
	err := c.BindJSON(&questionDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.ShowBadRequest(err.Error()))
	}
	c.JSON(http.StatusOK, questionDto)
}