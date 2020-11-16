package service

import (
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/jeevatkm/go-model.v1"
	"imitate-zhihu/dto"
	"imitate-zhihu/repository"
	"imitate-zhihu/result"
)

func UserLogin(loginDto *dto.UserLoginDto) result.Result {
	user, res := repository.SelectUserByEmail(loginDto.Email)
	if !res.IsOK() {
		return res
	}
	// decrypt
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))
	if err != nil {
		return result.PasswordNotCorrectErr.HandleError(err)
	}
	res = repository.SetUserToken(&user, uuid.NewV4().String())
	if !res.IsOK() {
		return res
	}
	userDto := dto.UserDto{}
	model.Copy(&userDto, &user)
	return res.WithData(&userDto)
}


func UserRegister(registerDto *dto.UserRegisterDto) result.Result {
	_, res := repository.SelectUserByEmail(registerDto.Email)
	if res.IsOK() {
		return result.EmailAlreadyExistErr
	}
	user := repository.User{}
	model.Copy(&user, registerDto)
	// encrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return result.OtherErr.HandleError(err)
	}
	user.Password = string(hash)
	res = repository.CreateUser(&user)
	if !res.IsOK() {
		return res
	}
	res = repository.SetUserToken(&user, uuid.NewV4().String())
	if !res.IsOK() {
		return res
	}
	userDto := dto.UserDto{}
	model.Copy(&userDto, &user)
	return res.WithData(userDto)
}