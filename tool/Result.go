package tool

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

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
	return res.code == 1
}

func (res Result) WithData(data interface{}) Result {
	return Result{res.code, res.message, data}
}


func NewResult(code int, message string) Result {
	return Result{code: code, message: message}
}

func NewResultWithData(code int, message string, data interface{}) Result {
	return Result{code, message, data}
}



// Result definitions
var (
	Ok                    = Result{code: 1, message: "OK"}
	RequestFormatErr      = Result{code: 400, message: "Request Format Error"}
	ServerErr             = Result{code: 500, message: "Server Error"}
	UserNotFoundErr       = Result{code: 2001, message: "User Not Found"}
	PasswordNotCorrectErr = Result{code: 2002, message: "Password not correct"}
	EmailAlreadyExistErr  = Result{code: 2003, message: "Email already exists"}
	CreateUserErr         = Result{code: 2004, message: "Create user error"}
	SetTokenErr           = Result{code: 2005, message: "Set token error"}
	MapperErr             = Result{code: 3000, message: "Mapper error"}
)
