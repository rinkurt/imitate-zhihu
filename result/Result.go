package result

import (
	"github.com/pkg/errors"
	"imitate-zhihu/tool"
	"strconv"
)

// Usage:
// Use WithData() to bring normal errors,
// In repository layer, DO NOT bring models,
// bring DTOs in service layer instead.
// Use HandleError() to print stack when bringing error as Data.

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (res Result) Error() string {
	return strconv.Itoa(res.Code) + ": " + res.Message
}

func (res Result) IsOK() bool {
	return res.Code == 0
}

func (res Result) NotOK() bool {
	return res.Code != 0
}

func (res Result) IsServerErr() bool {
	return res.Code == 500
}

func (res Result) WithData(data interface{}) Result {
	res.Data = data
	return res
}

func (res Result) WithError(err error) Result {
	res.Message += ": " + err.Error()
	return res
}

func (res Result) WithErrorStr(str string) Result {
	res.Message += ": " + str
	return res
}

func HandleServerErr(err error) Result {
	ers := errors.WithStack(err)
	tool.Logger.Errorf("%+v\n", ers)
	return ServerErr.WithError(err)
}

// Result definitions
var (
	Ok         = Result{Code: 0, Message: "OK"}
	BadRequest = Result{Code: 400, Message: "Bad Request"}
	ServerErr  = Result{Code: 500, Message: "Server Error"}
	RedisErr  = Result{Code: 501, Message: "Redis Error"}
	GraphDeleteErr =Result{Code: 502, Message: "GraphDeleteErr"}

	// 100_: Authorization
	EmptyAuth       = Result{Code: 1001, Message: "Empty Authorization"}
	AuthFormatErr   = Result{Code: 1002, Message: "Authorization Format Error"}
	TokenErr        = Result{Code: 1003, Message: "Token Error"}
	ContextErr      = Result{Code: 1004, Message: "Context Error"}
	UnauthorizedOpr = Result{Code: 1005, Message: "Unauthorized operation"}

	// 20__: User
	// 200_: User Login
	UserNotFoundErr       = Result{Code: 2001, Message: "User Not Found"}
	PasswordNotCorrectErr = Result{Code: 2002, Message: "Password not correct"}
	// 201_: User Register
	EmailAlreadyExistErr  = Result{Code: 2011, Message: "Email already exists"}
	CreateUserErr         = Result{Code: 2012, Message: "DB Create user error"}
	EmailSendErr          = Result{Code: 2013, Message: "Email send error"}
	WrongVerificationCode = Result{Code: 2014, Message: "Wrong verification code, or code expired."}

	// 21__: Question
	// 210_: Get Question
	QuestionNotFoundErr = Result{Code: 2101, Message: "Question Not Found"}
	UpdateViewCountErr  = Result{Code: 2102, Message: "Failed in updating view count"}
	// 211_: Create Question
	CreateQuestionErr = Result{Code: 2111, Message: "DB create question error"}
	// 212_: Update Question
	UpdateQuestionErr = Result{Code: 2121, Message: "DB update question error"}
	// 213_: Delete Question
	DeleteQuestionErr = Result{Code: 2131, Message: "DB delete question error"}

	// 22__: Answer
	// 220_: Get Answer
	AnswerNotFoundErr = Result{Code: 2201, Message: "Answer Not Found"}
	// 221_: Create Answer
	CreateAnswerErr = Result{Code: 2211, Message: "DB create Answer error"}
	// 222_: Update Answer
	UpdateAnswerErr = Result{Code: 2221, Message: "DB update Answer error"}
	// 223_: Delete Answer
	DeleteAnswerErr = Result{Code: 2231, Message: "DB delete Answer error"}
)
