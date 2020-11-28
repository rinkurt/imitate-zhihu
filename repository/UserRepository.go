package repository

import (
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"time"
)

type User struct {
	Id          int `gorm:"primaryKey"`
	Name        string
	Email       string
	Password    string
	Token       string
	CreateAt    int64
	Description string
	AvatarUrl   string `json:"avatar_url"`
}

//func SetUserToken(user *User, token string) result.Result {
//	db := tool.GetDatabase()
//	res := db.Model(user).Update("token", token)
//	if res.RowsAffected == 0 {
//		return result.SetTokenErr
//	}
//	user.Token = token
//	return result.Ok
//}

func SelectUserById(id int) (*User, result.Result) {
	db := tool.GetDatabase()
	user := User{}
	res := db.First(&user, id)
	if res.RowsAffected == 0 {
		return nil, result.UserNotFoundErr
	}
	return &user, result.Ok
}

func SelectUserByEmail(email string) (*User, result.Result) {
	db := tool.GetDatabase()
	user := User{}
	res := db.Where(&User{Email: email}).First(&user)
	if res.RowsAffected == 0 {
		return nil, result.UserNotFoundErr
	}
	return &user, result.Ok
}

func CreateUser(user *User) result.Result {
	db := tool.GetDatabase()
	user.Id = 0
	user.CreateAt = time.Now().Unix()
	res := db.Create(user)
	if res.RowsAffected == 0 {
		return result.CreateUserErr
	}
	return result.Ok
}
