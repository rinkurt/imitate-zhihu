package repository

import (
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"time"
)

type Profile struct {
	Id          int64 `gorm:"primaryKey"`
	UserId      int64
	Name        string
	Gender      int
	Description string
	AvatarUrl   string
	CreateAt    int64
	UpdateAt    int64
}

func SelectProfileByUserId(userId int64) (*Profile, result.Result) {
	db := tool.GetDatabase()
	profile := Profile{}
	res := db.Where(&Profile{UserId: userId}).First(&profile)
	if res.RowsAffected == 0 {
		return nil, result.UserNotFoundErr
	}
	return &profile, result.Ok
}

func CreateProfile(profile *Profile) result.Result {
	db := tool.GetDatabase()
	profile.Id = 0
	profile.CreateAt = time.Now().Unix()
	profile.UpdateAt = profile.CreateAt
	res := db.Create(profile)
	if res.RowsAffected == 0 {
		return result.CreateUserErr
	}
	return result.Ok
}
