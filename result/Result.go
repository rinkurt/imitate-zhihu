package result

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strconv"
)

// Usage:
// Use WithData() to bring normal errors,
// In repository layer, DO NOT bring models,
// bring DTOs in service layer instead.
// Use HandleError() to print stack when bringing error as data.

type Result struct {
	// Hide from outside
	code    int
	message string
	data    interface{}
}

func (res Result) Error() string {
	return strconv.Itoa(res.code) + ":" + res.message
}

// Get message in map (gin.H), generally for JSON.
func (res Result) Show() gin.H {
	return gin.H{
		"code":    res.code,
		"message": res.message,
		"data":    res.data,
	}
}

func (res Result) IsOK() bool {
	return res.code == 0
}

func (res Result) WithData(data interface{}) Result {
	return Result{res.code, res.message, data}
}

// Print stack details while bringing error as data.
func (res Result) HandleError(err error) Result {
	ers := errors.WithStack(err)
	fmt.Printf("%+v\n", ers)
	return Result{res.code, res.message, err.Error()}
}

func NewResult(code int, message string) Result {
	return Result{code: code, message: message}
}

func NewResultWithData(code int, message string, data interface{}) Result {
	return Result{code, message, data}
}

// For wrong request format, such as bind JSON error
func ShowBadRequest(data interface{}) gin.H {
	return RequestFormatErr.WithData(data).Show()
}

// For Authorization Fail
func ShowAuthErr(data interface{}) gin.H {
	return Result{code: 1001, message: "Authorization Error", data: data}.Show()
}

// For Errors in Controller Layer
func ShowControllerErr(data interface{}) gin.H {
	return Result{code: 1002, message: "Controller Error", data: data}.Show()
}

// Result definitions
var (
	Ok                    = Result{code: 0, message: "OK"}
	RequestFormatErr      = Result{code: 1, message: "Request Format Error"}
	OtherErr              = Result{code: 2001, message: "Other Error"}
	UserNotFoundErr       = Result{code: 2002, message: "User Not Found"}
	PasswordNotCorrectErr = Result{code: 2003, message: "Password not correct"}
	EmailAlreadyExistErr  = Result{code: 2004, message: "Email already exists"}
	CreateUserErr         = Result{code: 2005, message: "Create user error"}
	SetTokenErr           = Result{code: 2006, message: "Set token error"}
	QuestionNotFoundErr   = Result{code: 2007, message: "Question Not Found"}
	CreateQuestionErr     = Result{code: 2008, message: "Create question error"}
)
