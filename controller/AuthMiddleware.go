package controller

import (
	"github.com/gin-gonic/gin"
	"imitate-zhihu/result"
	"imitate-zhihu/tool"
	"net/http"
	"strings"
)

func JWTAuthMiddleware(c *gin.Context) {
	// Token放在Header的Authorization中，并使用Bearer开头
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, result.EmptyAuth)
		c.Abort()
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusUnauthorized, result.AuthFormatErr)
		c.Abort()
		return
	}
	mc, err := tool.ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, result.TokenErr.WithError(err))
		c.Abort()
		return
	}

	// 将当前请求的username信息保存到请求的上下文c上
	c.Set("user_id", mc.UserId)
	// 后续的处理函数可以用过c.Get("user_id")来获取当前请求的用户信息
	c.Next()

}
