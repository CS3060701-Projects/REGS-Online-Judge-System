package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	"strconv"

	"archive/zip"
	"io"

	"github.com/gin-gonic/gin"
)

// GetProblems godoc
// @Summary Get a list of problems
// @Description Retrieves a paginated list of visible problems.
// @Tags Problems
// @Produce  json
// @Param   page query int false "Page number" default(1)
// @Param   limit query int false "Items per page" default(10)
// @Success 200 {object} object{total=integer, page=integer, limit=integer, data=[]models.Problem}
// @Failure 500 {object} object{error=string} "查詢題目列表失敗"
// @Router /problems [get]
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

// GetProblem godoc
// @Summary Get a single problem's details
// @Description Retrieves the full description of a single visible problem.
// @Tags Problems
// @Produce  json
// @Param   id path string true "Problem ID"
// @Success 200 {object} object{description=string}
// @Failure 404 {object} object{error=string} "找不到該題目，或題目尚未公開"
// @Router /problems/{id} [get]
func GetProblem(c *gin.Context) {
	problemID := c.Param("id")

	var problem models.Problem
	result := database.DB.Select("description").
		Where("id = ? AND is_visible = ?", problemID, true).
		First(&problem)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到該題目，或題目尚未公開"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"description": problem.Description,
	})
}

// DownloadTestCases godoc
// @Summary Download problem test cases
// @Description (Admin only) Downloads all test cases for a problem as a single .zip file.
// @Tags Admin
// @Produce  application/zip
// @Security Bearer
// @Param   id path string true "Problem ID"
// @Router /problems/{id}/testcases [get]
func DownloadTestCases(c *gin.Context) {
	problemID := c.Param("id")

	var problem models.Problem
	if err := database.DB.Select("id").Where("id = ?", problemID).First(&problem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到該題目"})
		return
	}

	testDataDir := filepath.Join("testdata", problemID)

	info, err := os.Stat(testDataDir)
	if os.IsNotExist(err) || !info.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"error": "該題目尚未上傳任何測資檔案"})
		return
	}

	downloadName := fmt.Sprintf("problem_%s_testcases.zip", problemID)
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Header("Content-Type", "application/zip")

	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	err = filepath.Walk(testDataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, _ := filepath.Rel(testDataDir, path)
		f, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		fileContent, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fileContent.Close()

		_, err = io.Copy(f, fileContent)
		return err
	})

	if err != nil {
		fmt.Printf("[Error] 打包測資失敗: %v\n", err)
	}
}
