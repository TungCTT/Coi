package middleware

import (
	"net/http"
	"strings"

	jwtpkg "coi/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware là Gin middleware xác thực JWT.
//
// Middleware trong Gin là func(c *gin.Context) — khác với net/http là func(next http.Handler).
// Gin middleware dùng c.Next() để chuyển sang handler tiếp theo,
// hoặc c.Abort() để dừng chain và trả response ngay lập tức.
//
// Cách dùng: router.GET("/protected", AuthMiddleware(), handler)
// hoặc:       group := router.Group("/", AuthMiddleware())
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Yêu cầu xác thực — thiếu Authorization header",
			})
			c.Abort()
			return
		}

		_, tokenString, found := strings.Cut(authHeader, "Bearer ")
		if !found || tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Định dạng Authorization header không hợp lệ, yêu cầu: Bearer <token>",
			})
			c.Abort()
			return
		}

		claims, err := jwtpkg.ValidateAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token không hợp lệ hoặc đã hết hạn",
			})
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
