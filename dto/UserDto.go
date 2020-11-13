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

type UserDto struct {
	Id        int
	Name      string
	Email     string
	Token     string
	Bio       string
	AvatarUrl string `json:"avatar_url"`
}
