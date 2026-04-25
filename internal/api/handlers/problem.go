package handlers

import (
	"net/http"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetProblems(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var problems []models.Problem
	var total int64

	query := database.DB.Model(&models.Problem{}).Where("is_visible = ?", true)
	query.Count(&total)

	result := query.Select("id", "title", "time_limit", "memory_limit", "created_at").
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Find(&problems)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢題目列表失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page":  page,
		"limit": limit,
		"data":  problems,
	})
}

func GetProblem(c *gin.Context) {
	problemID := c.Param("id")

	var problem models.Problem
	// 使用 Select 只撈取 description 欄位，節省記憶體與傳輸時間
	result := database.DB.Select("description").
		Where("id = ? AND is_visible = ?", problemID, true).
		First(&problem)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到該題目，或題目尚未公開"})
		return
	}

	// 只回傳 description 欄位
	c.JSON(http.StatusOK, gin.H{
		"description": problem.Description,
	})
}
