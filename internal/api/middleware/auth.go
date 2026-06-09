package middleware

import (
	"errors"
	"net/http"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	jwtPkg "regs-backend/pkg/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func parseBearerToken(authHeader string) (string, bool) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	return token, token != ""
}

func authenticateToken(tokenString string) (*jwtPkg.Claims, error) {
	claims, err := jwtPkg.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	var blacklisted models.JwtBlacklist
	result := database.DB.Where("token = ?", tokenString).Limit(1).Find(&blacklisted)
	if result.RowsAffected > 0 {
		return nil, errors.New("token blacklisted")
	}

	return claims, nil
}

// OptionalAuthMiddleware validates Authorization when present (Guest routes).
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		tokenString, ok := parseBearerToken(authHeader)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的認證格式，請使用 Bearer <token>"})
			c.Abort()
			return
		}

		claims, err := authenticateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 Token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			if requiredRole == "Guest" {
				c.Next()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供認證標頭"})
			c.Abort()
			return
		}

		tokenString, ok := parseBearerToken(authHeader)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的認證格式，請使用 Bearer <token>"})
			c.Abort()
			return
		}

		claims, err := authenticateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 Token"})
			c.Abort()
			return
		}

		if claims.Role != "Admin" && requiredRole != "Guest" && claims.Role != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "權限不足"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
