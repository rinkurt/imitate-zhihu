package controller

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/repository"
	"imitate-zhihu/dto"
	"strconv"
)

func RouteQuestionController(engine *gin.Engine) {
	group := engine.Group("/question")
	group.GET("/:question_id", GetQuestionById)
	group.POST("", NewQuestion)
}

func GetQuestionById(c *gin.Context) {
	qid := c.Param("question_id")
	id, err := strconv.Atoi(qid)
	if err != nil {
		c.String(400, "Wrong Parameter")
	}
	question := repository.SelectQuestionById(id)
	c.JSON(200, question)
}

func NewQuestion(c *gin.Context) {
	questionDto := dto.QuestionDto{}
	err := c.BindJSON(&questionDto)
	if err != nil {
		c.String(400, "BindJson Error")
	}
	c.JSON(200, questionDto)
}