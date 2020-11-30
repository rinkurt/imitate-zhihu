package dto

type UserLoginDto struct {
	Email    string `binding:"email"`
	Password string `binding:"max=20,min=6"`
}

type UserRegisterDto struct {
	Email       string `binding:"email"`
	Password    string `binding:"max=20,min=6"`
	RePassword  string `json:"re_password" binding:"eqfield=Password"`
	Name        string `binding:"required"`
	Gender      int
	Description string
	AvatarUrl   string `json:"avatar_url"`
}

type UserProfileDto struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Gender      int    `json:"gender"`
	Description string `json:"description"`
	AvatarUrl   string `json:"avatar_url"`
}

var anonymousUser *UserProfileDto = nil

func AnonymousUser() *UserProfileDto {
	if anonymousUser == nil {
		anonymousUser = &UserProfileDto{
			Id:   0,
			Name: "Anonymous",
		}
	}
	return anonymousUser
}
