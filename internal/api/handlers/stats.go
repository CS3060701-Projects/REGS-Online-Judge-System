package handlers

import (
	"net/http"
	"regs-backend/internal/database"
	"regs-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func GetProblemStats(c *gin.Context) {
	problemID := c.Param("problem_id")

	var stats struct {
		TotalSubmissions int64          `json:"total_submissions"`
		ACCount          int64          `json:"ac_count"`
		StatusDist       map[string]int `json:"status_distribution"`
	}

	database.DB.Model(&models.Submission{}).Where("problem_id = ?", problemID).Count(&stats.TotalSubmissions)
	database.DB.Model(&models.Submission{}).Where("problem_id = ? AND status = ?", problemID, "AC").Count(&stats.ACCount)

	type Result struct {
		Status string
		Count  int
	}
	var results []Result
	database.DB.Model(&models.Submission{}).
		Select("status, count(*) as count").
		Where("problem_id = ?", problemID).
		Group("status").
		Scan(&results)

	stats.StatusDist = make(map[string]int)
	for _, r := range results {
		stats.StatusDist[r.Status] = r.Count
	}

	c.JSON(http.StatusOK, stats)
}

func GetUserStats(c *gin.Context) {
	targetUserID := c.Param("user_id")

	var stats struct {
		TotalSubmissions int64          `json:"total_submissions"`
		SolvedCount      int64          `json:"solved_count"` // 唯一 AC 的題目數量
		StatusDist       map[string]int `json:"status_distribution"`
	}

	database.DB.Model(&models.Submission{}).Where("user_id = ?", targetUserID).Count(&stats.TotalSubmissions)

	database.DB.Model(&models.Submission{}).
		Where("user_id = ? AND status = ?", targetUserID, "AC").
		Distinct("problem_id").
		Count(&stats.SolvedCount)

	type Result struct {
		Status string
		Count  int
	}
	var results []Result
	database.DB.Model(&models.Submission{}).
		Select("status, count(*) as count").
		Where("user_id = ?", targetUserID).
		Group("status").
		Scan(&results)

	stats.StatusDist = make(map[string]int)
	for _, r := range results {
		stats.StatusDist[r.Status] = r.Count
	}

	c.JSON(http.StatusOK, stats)
}
