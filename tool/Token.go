package tool

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"time"
)

type MyClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

var MySecret = []byte("c2f7e3c0-0267-44bd-a57a-162198e07784")

// Expire after 15 days.
//var TokenExpireDuration = time.Hour * 24 * 15


func GenToken(userId string) (string, error) {
	c := MyClaims{
		userId, // 自定义字段
		jwt.StandardClaims{
			//ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "imitate-zhihu",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return token.SignedString(MySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return MySecret, nil
		})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
