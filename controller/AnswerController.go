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

	group.GET("", GetAnswers)
	group.GET("/:answer_id", GetAnswerById)
	group.POST("", middleware.JWTAuthMiddleware, NewAnswer)
	group.PUT("/:answer_id", middleware.JWTAuthMiddleware, UpdateAnswerById)
	group.DELETE("/:answer_id", middleware.JWTAuthMiddleware, DeleteAnswerById)
}

var emptyAnswer = &dto.AnswerDetailDto{Creator: &dto.UserProfileDto{}}

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
	if res.NotOK() {
		res = res.WithData(emptyAnswer)
	}
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
	if res.NotOK() {
		answerDetail = emptyAnswer
	}
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
	if res.NotOK() {
		answerDetail = emptyAnswer
	}
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
	voteBy := c.Query("voteby")
	voteUid, err := tool.StrToInt64(voteBy)
	if err != nil {
		voteUid = 0
	}

	s := c.DefaultQuery("size", "10")
	size, err := tool.StrToInt(s)
	if err != nil {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithError(err))
		return
	}

	// 进入按点赞查询，该部分需要从缓存查询，单独处理
	if voteUid != 0 {
		// 此处 cursor 格式不同，单独处理
		qCursor := c.Query("cursor")
		cursor, err := tool.StrToInt(qCursor)
		if err != nil {
			cursor = 0
		}

		answers, res := service.GetAnswersByVoteUser(voteUid, cursor, size)
		if answers == nil {
			answers = []dto.AnswerDetailDto{}
		}

		c.JSON(http.StatusOK, res.WithData(gin.H{
			"next_cursor": tool.IntToStr(cursor + size),
			"answers":     answers,
		}))
		// 停止后续处理
		return
	}

	qid := c.Query("qid")
	questionId, err := tool.StrToInt64(qid)
	if err != nil {
		questionId = 0
	}

	uid := c.Query("uid")
	userId, err := tool.StrToInt64(uid)
	if err != nil {
		userId = 0
	}

	cursor, err := tool.ParseCursor(c.DefaultQuery("cursor", "-1,-1"))
	if err != nil || len(cursor) < 2 {
		c.JSON(http.StatusBadRequest, result.BadRequest.WithErrorStr("Cursor format error"))
		return
	}

	orderBy := c.DefaultQuery("orderby", enum.ByUpvote) //获取排序方式，默认为点赞数
	answers, res := service.GetAnswers(questionId, userId, cursor, size, orderBy)
	if answers == nil {
		answers = []dto.AnswerDetailDto{}
	}

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
