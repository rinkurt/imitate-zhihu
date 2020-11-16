package repository

import (
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"time"
)

type User struct {
	Id        int `gorm:"primaryKey"`
	Name      string
	Email     string
	Password  string
	Token     string
	GmtCreate int64
	Bio       string
	AvatarUrl string `json:"avatar_url"`
}

func SetUserToken(user *User, token string) result.Result {
	db := tool.GetDatabase()
	res := db.Model(user).Update("token", token)
	if res.RowsAffected == 0 {
		return result.SetTokenErr.WithData(res.Error.Error())
	}
	user.Token = token
	return result.Ok
}

func SelectUserByEmail(email string) (User, result.Result) {
	db := tool.GetDatabase()
	user := User{Email: email}
	res := db.Where(&user).First(&user)
	if res.RowsAffected == 0 {
		return user, result.UserNotFoundErr.WithData(res.Error.Error())
	}
	return user, result.Ok
}

func CreateUser(user *User) result.Result {
	db := tool.GetDatabase()
	user.Id = 0
	user.GmtCreate = time.Now().Unix()
	res := db.Create(user)
	if res.RowsAffected == 0 {
		return result.CreateUserErr.WithData(res.Error.Error())
	}
	return result.Ok
}
