package service

import (
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
	//res = repository.SetUserToken(&user, uuid.NewV4().String())
	//if !res.IsOK() {
	//	return res
	//}
	token, err := tool.GenToken(strconv.Itoa(user.Id))
	if err != nil {
		return result.HandleServerErr(err)
	}
	user.Token = token

	userDto := dto.UserDetailDto{}
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
		return result.HandleServerErr(err)
	}
	user.Password = string(hash)
	res = repository.CreateUser(&user)
	if res.NotOK() {
		return res
	}
	//res = repository.SetUserToken(&user, uuid.NewV4().String())
	//if !res.IsOK() {
	//	return res
	//}
	token, err := tool.GenToken(strconv.Itoa(user.Id))
	if err != nil {
		return result.HandleServerErr(err)
	}
	user.Token = token

	userDto := dto.UserDetailDto{}
	model.Copy(&userDto, &user)
	return res.WithData(userDto)
}


// TODO: Cache
func GetUserById(id int) (*dto.UserDto, result.Result) {
	user, res := repository.SelectUserById(id)
	if res.NotOK() {
		return nil, res
	}
	userDto := dto.UserDto{}
	model.Copy(&userDto, &user)
	return &userDto, result.Ok
}