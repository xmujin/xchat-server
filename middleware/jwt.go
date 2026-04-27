package middleware

import (
	"strings"
	"xchat-server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if authHeader == "" {
			c.JSON(401, gin.H{
				"error": "unauthorized",
			})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])

		if err != nil {
			c.JSON(401, gin.H{
				"error": "invalid token",
			})
			c.Abort()
			return
		}

		// 把用户信息存进 context
		c.Set("username", claims.Username)

		c.Next()
	}
}
