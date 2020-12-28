package test

import (
	"fmt"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/service"
	"testing"
)

func TestVerifyEmail(t *testing.T) {
	testEmail := "1812754991@qq.com"
	verifyCode, res := service.VerifyEmail(testEmail)
	if res.NotOK(){
		t.Error(res.Message)
	}
	fmt.Printf("verifyCode:%v\n",verifyCode)
}

func TestUserRegister(t *testing.T) {
	email := "1812754991@qq.com"
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
	res = repository.DeleteUserByEmail(email)
}

func TestUserLogin(t *testing.T) {
	testUser := dto.UserLoginDto{
		Email:    "1812754991@qq.com",
		Password: "123456",
	}
	res := service.UserLogin(&testUser)
	if res.NotOK() {
		t.Error(res.Message)
	}
}

func TestGetUserProfileByUid(t *testing.T)  {
	var resProfile *dto.UserProfileDto
	testID := int64(5)
	resProfile,res := service.GetUserProfileByUid(testID)
	if res.NotOK() {
		t.Error(res.Message)
	}
	fmt.Printf("user profile:%v\n",resProfile)
}