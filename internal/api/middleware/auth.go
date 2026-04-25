package middleware

import (
	"net/http"
	"strings"

	"regs-backend/pkg/jwt"

	jwtlib "github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if requiredRole == "Guest" {
			c.Next()
			return
		}

		// 取得 Authorization Header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供認證憑證"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		tokenString = strings.TrimSpace(tokenString)

		token, err := jwtlib.Parse(tokenString, func(t *jwtlib.Token) (interface{}, error) {
			return jwt.VerifyKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無效或已過期的 Token"})
			return
		}

		claims, ok := token.Claims.(jwtlib.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "無法解析 Token 聲明"})
			return
		}

		userRole, _ := claims["role"].(string)
		userID, _ := claims["user_id"].(float64) // JWT 數字預設解析為 float64

		// admin > user > guest
		isAuthorized := false
		switch requiredRole {
		case "User":
			if userRole == "User" || userRole == "Admin" {
				isAuthorized = true
			}
		case "Admin":
			if userRole == "Admin" {
				isAuthorized = true
			}
		}

		if !isAuthorized {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "權限不足，拒絕存取"})
			return
		}

		c.Set("user_id", uint(userID))
		c.Next()
	}
}
