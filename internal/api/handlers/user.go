package handlers

import (
	"net/http"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	jwtPkg "regs-backend/pkg/jwt" // 替換成你的 jwt 套件路徑

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 定義接收請求的資料格式
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 註冊 API
func Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤，需要 username 與 password"})
		return
	}

	// 檢查使用者是否已經存在
	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "使用者名稱已被註冊"})
		return
	}

	// 使用 bcrypt 加密密碼 (Cost 預設為 10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密碼加密失敗"})
		return
	}

	// 第一個註冊的人預設給 Admin 權限，後續的人都是 User (方便你測試)
	role := "User"
	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)
	if totalUsers == 0 {
		role = "Admin"
	}

	// 建立使用者並存入資料庫
	newUser := models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         role,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "建立使用者失敗"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "註冊成功", "user_id": newUser.ID, "role": newUser.Role})
}

// 登入 API
func Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤，需要 username 與 password"})
		return
	}

	// 根據帳號去資料庫找人
	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤"})
		return
	}

	// 比對密碼是否正確
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤"})
		return
	}

	// 密碼正確，生成 JWT Token
	token, err := jwtPkg.GenerateToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token 生成失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登入成功",
		"token":   token,
		"role":    user.Role,
	})
}
