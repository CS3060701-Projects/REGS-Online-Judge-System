package handlers

import (
	"net/http"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	jwtPkg "regs-backend/pkg/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤，需要 username 與 password"})
		return
	}

	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "使用者名稱已被註冊"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密碼加密失敗"})
		return
	}

	role := "User"
	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)
	if totalUsers == 0 {
		role = "Admin"
	}

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

func Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤，需要 username 與 password"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "帳號或密碼錯誤"})
		return
	}

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

func Logout(c *gin.Context) {
	// 1. 從 Header 取得 Token (格式通常是 "Bearer <token>")
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求標頭"})
		return
	}
	tokenString := authHeader[7:]

	// 2. 解析 Token 取得 Claims (主要是為了拿 exp 過期時間)
	claims, err := jwtPkg.ParseToken(tokenString) // 確保有引用你的 jwt 套件
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 Token"})
		return
	}

	// 3. 將 Token 存入黑名單
	blacklist := models.JwtBlacklist{
		Token:     tokenString,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	if err := database.DB.Create(&blacklist).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登出操作失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "登出成功"})
}

func GetMe(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授權的操作"})
		return
	}

	userID := val.(uint)

	var user models.User
	err := database.DB.Select("id", "username", "role", "created_at").
		First(&user, userID).Error

	if err != nil {
		// 如果還是出現 record not found，請檢查資料庫內是否真的有 ID 為 userID 的資料
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到使用者資料", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
