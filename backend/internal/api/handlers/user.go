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

// Register godoc
// @Summary Register a new user
// @Description Creates a new user account. The first user registered will be an admin.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   user body AuthRequest true "User Registration Info"
// @Success 201 {object} object{message=string, user_id=integer, role=string} "註冊成功"
// @Failure 400 {object} object{error=string} "格式錯誤"
// @Failure 409 {object} object{error=string} "使用者名稱已被註冊"
// @Failure 500 {object} object{error=string} "伺服器內部錯誤"
// @Router /users/register [post]
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

	newUser := models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         "User",
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "建立使用者失敗"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "註冊成功", "user_id": newUser.ID, "role": newUser.Role})
}

// Login godoc
// @Summary Log in a user
// @Description Authenticates a user and returns a JWT token.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param   credentials body AuthRequest true "User Login Credentials"
// @Success 200 {object} object{message=string, token=string, role=string} "登入成功"
// @Failure 400 {object} object{error=string} "格式錯誤"
// @Failure 401 {object} object{error=string} "帳號或密碼錯誤"
// @Failure 500 {object} object{error=string} "Token 生成失敗"
// @Router /users/login [post]
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

// Logout godoc
// @Summary Log out the current user
// @Description Invalidates the current user's JWT by adding it to a blacklist.
// @Tags Users
// @Security Bearer
// @Success 200 {object} object{message=string} "登出成功"
// @Failure 400 {object} object{error=string} "無效的請求標頭"
// @Failure 401 {object} object{error=string} "無效的 Token"
// @Failure 500 {object} object{error=string} "登出操作失敗"
// @Router /users/logout [post]
func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) < 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的請求標頭"})
		return
	}
	tokenString := authHeader[7:]

	claims, err := jwtPkg.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "無效的 Token"})
		return
	}

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

// GetMe godoc
// @Summary Get current user's profile
// @Description Retrieves the profile information of the currently authenticated user.
// @Tags Users
// @Produce  json
// @Security Bearer
// @Success 200 {object} models.User
// @Failure 401 {object} object{error=string} "未授權的操作"
// @Failure 404 {object} object{error=string} "找不到使用者資料"
// @Router /users/me [get]
func GetMe(c *gin.Context) {
	val, exists := c.Get("user_id")
	userID, ok := val.(uint)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授權的操作"})
		return
	}

	var user models.User
	err := database.DB.Select("id", "username", "role", "created_at").
		First(&user, userID).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到使用者資料", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
