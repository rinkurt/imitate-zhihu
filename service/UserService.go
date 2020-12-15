package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"strconv"
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

	// verify code
	val, err := tool.Rdb.Get(context.Background(), tool.KeyVrfCode(registerDto.Email)).Result()
	if err != nil && err != redis.Nil {
		return result.HandleServerErr(err)
	}
	if val != registerDto.VerificationCode {
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
	// find in cache
	userDto := &dto.UserProfileDto{}
	found := tool.CacheGet(tool.KeyUser(userId), userDto)
	if found {
		return userDto, result.Ok
	}

	profile, res := repository.SelectProfileByUserId(userId)
	if res.NotOK() {
		return nil, res
	}

	model.Copy(userDto, profile)
	userDto.Id = profile.UserId

	// save in cache
	tool.CacheSet(tool.KeyUser(userId), userDto)
	return userDto, result.Ok
}

func VerifyEmail(email string) (string, result.Result) {
   	vrfCode := tool.GenValidateCode(6)
	//tool.CodeCache[email]=vrfCode
	err := tool.Rdb.Set(context.Background(), tool.KeyVrfCode(email), vrfCode, time.Minute*30).Err()
	if err != nil {
		return "", result.HandleServerErr(err)
	}
	mailTo := email
	subject := "您的邮箱验证码是："+ vrfCode
	body := "您的验证码为："+ vrfCode + "，有效期 30 分钟。"
	err = tool.SendMail(mailTo, subject, body)
	if err != nil {
		return "", result.EmailSendErr.WithError(err)
	}
	return vrfCode, result.Ok
}







