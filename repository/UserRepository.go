package repository

import (
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"time"
)

type User struct {
	Id       int64 `gorm:"primaryKey"`
	Email    string
	Password string
	CreateAt int64
	UpdateAt int64
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
	user.UpdateAt = user.CreateAt
	res := db.Create(user)
	if res.RowsAffected == 0 {
		return result.CreateUserErr
	}
	return result.Ok
}

func DeleteUserByEmail(email string) result.Result {
	db := tool.GetDatabase()
	var user User
	db.Where("email = ?",email).First(&user)
	res := db.Delete(&user)
	if res.RowsAffected == 0 {
		return result.UserNotFoundErr
	}
	res = db.Where(&Profile{
		UserId: user.Id,
	}).Delete(Profile{})
	if res.RowsAffected == 0 {
		return result.UserNotFoundErr
	}
	return result.Ok
}
