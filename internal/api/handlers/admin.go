package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	problemPkg "regs-backend/internal/problem"
	"regs-backend/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// findSettingsYAML 在 root 目錄下（含子目錄）搜尋第一個 settings.yaml 的路徑
func findSettingsYAML(root string) string {
	var found string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || found != "" {
			return nil
		}
		if !info.IsDir() && info.Name() == "settings.yaml" {
			found = filepath.Dir(path)
			return filepath.SkipAll
		}
		return nil
	})
	return found
}

// CreateProblem godoc
// @Summary Create or update a problem (with test data)
// @Description (Admin only) Creates or restores a problem and uploads its test data in a single request using multipart/form-data.
// @Description time_limit and memory_limit are read automatically from settings.yaml inside the zip (limits.totalTime / limits.memory).
// @Tags Admin
// @Accept  multipart/form-data
// @Produce  json
// @Security Bearer
// @Param   id          formData string true  "Problem ID"
// @Param   title       formData string true  "Problem title"
// @Param   description formData string false "Problem description (Markdown)"
// @Param   is_visible  formData bool   false "Whether the problem is publicly visible (default true)"
// @Param   file        formData file   true  "Test cases as a .zip file (must contain settings.yaml)"
// @Success 200 {object} object{message=string, problem=models.Problem}
// @Failure 400 {object} object{error=string}
// @Failure 500 {object} object{error=string}
// @Router /problems [put]
func CreateProblem(c *gin.Context) {
	// --- Parse form fields ---
	problemID := c.PostForm("id")
	if problemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要欄位: id"})
		return
	}
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少必要欄位: title"})
		return
	}
	description := c.PostForm("description")

	isVisible := true
	if v := c.PostForm("is_visible"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			isVisible = b
		}
	}

	// --- Require zip file ---
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請上傳測資 zip 檔案 (file)"})
		return
	}

	testcasePath := filepath.Join("testdata", problemID)

	// --- Extract zip first so we can read settings.yaml ---
	os.RemoveAll(testcasePath)
	if err := os.MkdirAll(testcasePath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立測資目錄"})
		return
	}
	zipPath := filepath.Join(testcasePath, "temp_testdata.zip")
	if err := c.SaveUploadedFile(file, zipPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "zip 檔案儲存失敗"})
		return
	}
	if err := utils.Unzip(zipPath, testcasePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解壓縮失敗，請確認檔案格式"})
		return
	}
	os.Remove(zipPath)

	// --- Read limits from settings.yaml ---
	timeLimit := 1000
	memoryLimit := 256
	if settingsRoot := findSettingsYAML(testcasePath); settingsRoot != "" {
		if settings, err := problemPkg.LoadSettings(settingsRoot); err == nil {
			if settings.Limits.TotalTime > 0 {
				timeLimit = settings.Limits.TotalTime
			}
			if settings.Limits.Memory > 0 {
				memoryLimit = problemPkg.MemoryLimitMB(settings.Limits.Memory)
			}
		} else {
			fmt.Printf("[Warning] 無法解析 settings.yaml: %v\n", err)
		}
	} else {
		fmt.Printf("[Warning] 題目 %s 的 zip 中未找到 settings.yaml，使用預設限制\n", problemID)
	}

	problem := models.Problem{
		ID:           problemID,
		Title:        title,
		Description:  description,
		TimeLimit:    timeLimit,
		MemoryLimit:  memoryLimit,
		TestcasePath: testcasePath,
		IsVisible:    isVisible,
	}

	// --- Upsert into DB ---
	var existing models.Problem
	dbErr := database.DB.Unscoped().Where("id = ?", problemID).First(&existing).Error

	if dbErr == nil {
		// Record exists (active or soft-deleted) → update it
		updates := map[string]interface{}{
			"title":         problem.Title,
			"description":   problem.Description,
			"time_limit":    problem.TimeLimit,
			"memory_limit":  problem.MemoryLimit,
			"testcase_path": problem.TestcasePath,
			"is_visible":    problem.IsVisible,
		}
		if existing.DeletedAt.Valid {
			updates["deleted_at"] = nil
		}
		if err := database.DB.Unscoped().Model(&existing).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新題目失敗"})
			return
		}
		database.DB.Unscoped().First(&existing, "id = ?", problemID)
		problem = existing
	} else {
		// New record
		if err := database.DB.Create(&problem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立題目"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("題目 %s 建立/更新成功，測資已上傳（時限: %dms，記憶體: %dMB）", problemID, timeLimit, memoryLimit),
		"problem": problem,
	})
}

// DeleteProblem godoc
// @Summary Delete a problem
// @Description (Admin only) Soft-deletes a problem from the database and removes its associated test data files.
// @Tags Admin
// @Security Bearer
// @Param   id path string true "Problem ID"
// @Router /problems/{id} [delete]
func DeleteProblem(c *gin.Context) {
	problemID := c.Param("id")

	var problem models.Problem
	if err := database.DB.Where("id = ?", problemID).First(&problem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到指定題目"})
		return
	}

	if err := database.DB.Delete(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "刪除資料庫紀錄失敗"})
		return
	}
	testDataDir := problem.TestcasePath
	if testDataDir == "" {
		testDataDir = filepath.Join("testdata", problemID)
	}

	if err := os.RemoveAll(testDataDir); err != nil {
		fmt.Printf("[Warning] 測資目錄刪除失敗: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("題目 %s 已成功刪除，測資檔案已清除", problemID),
	})
}
