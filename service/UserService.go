package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func UserLogin(loginDto *dto.UserLoginDto) result.Result {
	user, res := repository.SelectUserByEmail(loginDto.Email)
	if res.NotOK() {
		return res
	}
	// decrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))
	if err != nil {
		return result.PasswordNotCorrectErr.WithError(err)
	}
	token, err := tool.GenToken(strconv.FormatInt(user.Id, 10))
	if err != nil {
		return result.HandleServerErr(err)
	}

	return res.WithData(gin.H{
		"id": user.Id,
		"token": token,
	})
}


func UserRegister(registerDto *dto.UserRegisterDto) result.Result {
	_, res := repository.SelectUserByEmail(registerDto.Email)
	if res.IsOK() {
		return result.EmailAlreadyExistErr
	}

	user := &repository.User{}
	model.Copy(user, registerDto)
	// encrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return result.HandleServerErr(err)
	}
	user.Password = string(hash)
	res = repository.CreateUser(user)
	if res.NotOK() {
		return res
	}

	// Create Profile
	profile := &repository.Profile{}
	model.Copy(profile, registerDto)
	profile.UserId = user.Id
	res = repository.CreateProfile(profile)
	if res.NotOK() {
		return res
	}

	token, err := tool.GenToken(strconv.FormatInt(user.Id, 10))
	if err != nil {
		return result.HandleServerErr(err)
	}

	return res.WithData(gin.H{
		"id": user.Id,
		"token": token,
	})
}

func GetUserProfileByUid(userId int64) (*dto.UserProfileDto, result.Result) {
	profile, res := repository.SelectProfileByUserId(userId)
	if res.NotOK() {
		return nil, res
	}
	userDto := &dto.UserProfileDto{}
	model.Copy(userDto, profile)
	userDto.Id = profile.UserId
	return userDto, result.Ok
}



func VerifyEmail(email string) result.Result {
   	verifyCode := GenValidateCode(6)

	mailTo := email
	subject := "您的邮箱验证码是："+verifyCode
	body := "您的验证码为:"+verifyCode+"\r\n。如果您本人没有通过登录验证请求此验证码，请立即前往“ 我的帐户 ”页面更改密码。\r\n如果您需要支持，请联系zhihu帮助。\r\n感谢您帮助我们一同维护您的帐户安全。\r\n\r\n祝您生活愉快，\r\nzhihu团队"
	err := SendMail(mailTo, subject, body)
	if err != nil {
		return result.EmailSendErr.WithError(err)
	}
	return result.Ok.WithData(gin.H{
		"verification_code": verifyCode,
	})
}

//生成有效验证码
func GenValidateCode(width int) string {
	numeric := [10]byte{0,1,2,3,4,5,6,7,8,9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[ rand.Intn(r) ])
	}
	return sb.String()
}

//发送邮件
func SendMail(mailTo string, subject string, body string) error {
	mailConn := map[string]string{
		"user": "zh_account_verify@163.com",
		"pass": "PVQIRHXFMBMFLFEI",//客户端授权密码
		"host": "smtp.163.com",
		"port": "465",
	}
	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int
	m := gomail.NewMessage()
	m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo)    //发送给某个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文
	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
	err := d.DialAndSend(m)
	return err
}

