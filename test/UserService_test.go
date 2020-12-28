package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/service"
	"testing"
)

var email = "e4yb596h@meantinc.com"
var id int64

func TestVerifyEmail(t *testing.T) {
	verifyCode, res := service.VerifyEmail(email)
	if res.NotOK(){
		t.Error(res.Message)
	}
	fmt.Printf("verifyCode:%v\n",verifyCode)
}

func TestUserRegister(t *testing.T) {
	_ = repository.DeleteUserByEmail(email)

	vrfCode, res := service.VerifyEmail(email)
	if res.NotOK() {
		t.Error(res.Message)
	}
	testUser := dto.UserRegisterDto{
		Email:            email,
		VerificationCode: vrfCode,
		Password:         "123456",
		RePassword:       "123456",
		Name:             "Bob",
		Gender:           1,
		Description:      "I'm Bob.",
		AvatarUrl:        "http://example.com/avatar.jpg",
	}
	res = service.UserRegister(&testUser)
	if res.NotOK() {
		t.Error(res.Message)
	}
}

func TestUserLogin(t *testing.T) {
	testUser := dto.UserLoginDto{
		Email:    email,
		Password: "123456",
	}
	res := service.UserLogin(&testUser)
	if res.NotOK() {
		t.Error(res.Message)
	}
	dat := res.Data.(gin.H)
	id = dat["id"].(int64)
}

func TestGetUserProfileByUid(t *testing.T)  {
	var resProfile *dto.UserProfileDto
	resProfile,res := service.GetUserProfileByUid(id)
	if res.NotOK() {
		t.Error(res.Message)
	}
	fmt.Printf("user profile:%v\n",resProfile)
}