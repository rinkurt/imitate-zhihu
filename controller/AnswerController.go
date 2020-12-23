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
)

func RouteAnswerController(engine *gin.Engine) {
	group := engine.Group("/answer")
	{
		group.GET("", GetAnswers)
		group.GET("/:answer_id", GetAnswerById)
		group.POST("", middleware.JWTAuthMiddleware, NewAnswer)
		group.PUT("/:answer_id", middleware.JWTAuthMiddleware, UpdateAnswerById)
		group.DELETE("/:answer_id", middleware.JWTAuthMiddleware, DeleteAnswerById)
	}
}

func NewAnswer(c *gin.Context) {
	userId, err := middleware.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.ContextErr.WithError(err))
		return
	}
	answerDto := dto.AnswerCreateDto{}
	answerDto.QuestionId, err = tool.StrToInt64(c.Query("qid"))
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	err = c.ShouldBind(&answerDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	res := service.NewAnswer(userId, &answerDto)
	c.JSON(http.StatusOK, res)
}

func GetAnswerById(c *gin.Context) {
	id := c.Param("answer_id")
	answerId, err := tool.StrToInt64(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	answerDetail, res := service.GetAnswerById(answerId)
	c.JSON(http.StatusOK, res.WithData(answerDetail))
}

func UpdateAnswerById(c *gin.Context) {
	userId, err := middleware.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.ContextErr.WithError(err))
		return
	}
	id := c.Param("answer_id")
	answerId, err := tool.StrToInt64(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	answer := dto.AnswerCreateDto{}
	err = c.ShouldBind(&answer)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	answerDetail, res := service.UpdateAnswerById(userId, answerId, &answer)
	c.JSON(http.StatusOK, res.WithData(answerDetail))
}

func DeleteAnswerById(c *gin.Context) {
	userId, err := middleware.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.UnauthorizedOpr)
		return
	}
	answerId := c.Param("answer_id")
	id, err := tool.StrToInt64(answerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	res := service.DeleteAnswerById(userId, id)
	c.JSON(http.StatusOK, res)
}

func GetAnswers(c *gin.Context) {
	qid := c.Query("qid")
	questionId, err := tool.StrToInt64(qid)
	if err != nil { //没有转换成功，说明请求失败
		c.JSON(http.StatusBadRequest, result.BadRequest.WithErrorStr("Missing question id"))
		return
	}
	cu := c.DefaultQuery("cursor", "-1,-1") //游标，缺省情况下取(-1,-1)，即从排好序的列表的第一个开始取size个记录

	cursor, err := tool.ParseCursor(cu)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithErrorStr("Cursor format error"))
		return
	}

	s := c.DefaultQuery("size", "5") //每页记录数,缺省情况下取5
	size, err := tool.StrToInt(s)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}
	orderBy := c.DefaultQuery("orderby", enum.ByTime) //获取排序方式，默认为时间戳降序
	answers, res := service.GetAnswers(questionId, cursor, size, orderBy)

	nextCursor := ""
	if len(answers) > 0 {
		tail := answers[len(answers) - 1]
		switch orderBy {
		case enum.ByTime:
			nextCursor = tool.Int64ToStr(tail.UpdateAt) + "," + tool.Int64ToStr(tail.Id)
		case enum.ByHeat:
			nextCursor = tool.IntToStr(tail.ViewCount) + "," + tool.Int64ToStr(tail.Id)
		case enum.ByUpvote:
			nextCursor = tool.IntToStr(tail.UpvoteCount) + "," + tool.Int64ToStr(tail.Id)
		}
	}

	c.JSON(http.StatusOK, res.WithData(gin.H{
		"next_cursor": nextCursor,
		"answers":     answers,
	}))

}
