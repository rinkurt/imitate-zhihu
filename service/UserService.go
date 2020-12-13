package service

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"strconv"
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

	// verify code
	if tool.CodeCache[registerDto.Email] != registerDto.VerificationCode {
		return result.WrongVerificationCode
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

func VerifyEmail(email string) (string, result.Result) {
   	verifyCode := tool.GenValidateCode(6)
	tool.CodeCache[email]=verifyCode
	mailTo := email
	subject := "您的邮箱验证码是："+verifyCode
	body := "您的验证码为："+verifyCode
	err := tool.SendMail(mailTo, subject, body)
	if err != nil {
		return "", result.EmailSendErr.WithError(err)
	}
	return verifyCode, result.Ok
}







