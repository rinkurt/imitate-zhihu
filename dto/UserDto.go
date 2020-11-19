package dto

type UserLoginDto struct {
	Email    string `binding:"email"`
	Password string `binding:"max=20,min=6"`
}

type UserRegisterDto struct {
	Name       string `binding:"required"`
	Email      string `binding:"email"`
	Password   string `binding:"max=20,min=6"`
	RePassword string `json:"re_password" binding:"eqfield=Password"`
	Bio        string
	AvatarUrl  string `json:"avatar_url"`
}

type UserDetailDto struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url"`
}

type UserDto struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatar_url"`
}

var anonymousUser *UserDto = nil

func AnonymousUser() *UserDto {
	if anonymousUser == nil {
		anonymousUser = &UserDto{
			Id:   0,
			Name: "Anonymous",
		}
	}
	return anonymousUser
}
