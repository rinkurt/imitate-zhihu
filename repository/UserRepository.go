package repository

import (
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

func SetUserToken(user *User, token string) tool.Result {
	db := tool.GetDatabase()
	res := db.Model(user).Update("token", token)
	if res.RowsAffected == 0 {
		return tool.SetTokenErr.WithData(res.Error.Error())
	}
	return tool.Ok
}

func SelectUserByEmail(email string) (User, tool.Result) {
	db := tool.GetDatabase()
	user := User{Email: email}
	res := db.Where(&user).First(&user)
	if res.RowsAffected == 0 {
		return user, tool.UserNotFoundErr
	}
	return user, tool.Ok
}

func CreateUser(user *User) tool.Result {
	db := tool.GetDatabase()
	user.Id = 0
	user.GmtCreate = time.Now().Unix()
	res := db.Create(user)
	if res.RowsAffected == 0 {
		return tool.CreateUserErr.WithData(res.Error.Error())
	}
	return tool.Ok.WithData(user)
}
