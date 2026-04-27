package middleware

import (
	"net/http"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	jwtPkg "regs-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
)

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

		if len(authHeader) < 7 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的認證格式"})
			c.Abort()
			return
		}
		tokenString := authHeader[7:]

		claims, err := jwtPkg.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 Token"})
			c.Abort()
			return
		}

		var blacklisted models.JwtBlacklist
		result := database.DB.Where("token = ?", tokenString).Limit(1).Find(&blacklisted)

		if result.RowsAffected > 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "此 Token 已登出，請重新登入"})
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
