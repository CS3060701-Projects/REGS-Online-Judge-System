package handlers

import (
	"math"
	"net/http"
	"regs-backend/internal/database"
	"regs-backend/internal/models"

	"github.com/gin-gonic/gin"
)

type ProblemStatsResponse struct {
	TotalSubmissions int64          `json:"total_submissions"`
	ACCount          int64          `json:"ac_count"`
	AcceptanceRate   float64        `json:"acceptance_rate"`
	StatusDist       map[string]int `json:"status_distribution"`
}

type UserStatsResponse struct {
	TotalSubmissions int64          `json:"total_submissions"`
	SolvedCount      int64          `json:"solved_count"`
	AcceptanceRate   float64        `json:"acceptance_rate"`
	StatusDist       map[string]int `json:"status_distribution"`
}

// GetProblemStats godoc
// @Summary Get statistics for a problem
// @Description Retrieves submission statistics for a specific problem, including total submissions, AC count, and status distribution.
// @Tags Statistics
// @Produce  json
// @Param   problem_id path string true "Problem ID"
// @Success 200 {object} ProblemStatsResponse
// @Router /stats/problems/{problem_id} [get]
func GetProblemStats(c *gin.Context) {
	problemID := c.Param("problem_id")

	var stats ProblemStatsResponse
	stats.StatusDist = make(map[string]int)

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

	for _, r := range results {
		stats.StatusDist[r.Status] = r.Count
		stats.TotalSubmissions += int64(r.Count)
		if r.Status == "AC" {
			stats.ACCount = int64(r.Count)
		}
	}

	if stats.TotalSubmissions > 0 {
		rate := (float64(stats.ACCount) / float64(stats.TotalSubmissions)) * 100
		stats.AcceptanceRate = math.Round(rate*100) / 100
	}

	c.JSON(http.StatusOK, stats)
}

// GetUserStats godoc
// @Summary Get statistics for a user
// @Description Retrieves submission statistics for a specific user, including total submissions, solved count, and status distribution.
// @Tags Statistics
// @Produce  json
// @Param   user_id path integer true "User ID"
// @Success 200 {object} UserStatsResponse
// @Router /stats/users/{user_id} [get]
func GetUserStats(c *gin.Context) {
	targetUserID := c.Param("user_id")

	var stats UserStatsResponse
	stats.StatusDist = make(map[string]int)

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

	for _, r := range results {
		stats.StatusDist[r.Status] = r.Count
		stats.TotalSubmissions += int64(r.Count)
	}

	database.DB.Model(&models.Submission{}).
		Where("user_id = ? AND status = ?", targetUserID, "AC").
		Distinct("problem_id").
		Count(&stats.SolvedCount)

	if stats.TotalSubmissions > 0 {
		rate := (float64(stats.StatusDist["AC"]) / float64(stats.TotalSubmissions)) * 100
		stats.AcceptanceRate = math.Round(rate*100) / 100
	}

	c.JSON(http.StatusOK, stats)
}
