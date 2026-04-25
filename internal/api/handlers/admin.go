package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regs-backend/internal/database"
	"regs-backend/internal/models"
	"regs-backend/pkg/utils"

	"github.com/gin-gonic/gin"
)

func CreateProblem(c *gin.Context) {
	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的題目格式"})
		return
	}

	if problem.TestcasePath == "" {
		problem.TestcasePath = filepath.Join("test_data", problem.ID)
	}

	var existing models.Problem
	err := database.DB.Unscoped().Where("id = ?", problem.ID).First(&existing).Error

	if err == nil {
		if existing.DeletedAt.Valid {
			err := database.DB.Unscoped().Model(&existing).Updates(map[string]interface{}{
				"deleted_at":    nil,
				"title":         problem.Title,
				"description":   problem.Description,
				"time_limit":    problem.TimeLimit,
				"memory_limit":  problem.MemoryLimit,
				"testcase_path": problem.TestcasePath,
				"is_visible":    problem.IsVisible,
			}).Error

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "復活舊題目失敗"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "已成功覆蓋並重新啟用被刪除的題目", "problem": existing})
			return
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "題目 ID 已存在且在使用中"})
			return
		}
	}

	if err := database.DB.Create(&problem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立題目"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "題目建立成功", "problem": problem})
}

func UploadTestData(c *gin.Context) {
	problemID := c.Param("id")

	var problem models.Problem
	if err := database.DB.Where("id = ?", problemID).First(&problem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到指定題目"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請上傳測資檔案 (.zip)"})
		return
	}
	testDataDir := problem.TestcasePath
	if testDataDir == "" {
		testDataDir = filepath.Join("test_data", problemID)
	}

	os.RemoveAll(testDataDir)
	if err := os.MkdirAll(testDataDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立測資目錄"})
		return
	}

	zipPath := filepath.Join(testDataDir, "temp_testdata.zip")
	if err := c.SaveUploadedFile(file, zipPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "檔案儲存失敗"})
		return
	}

	if err := utils.Unzip(zipPath, testDataDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解壓縮失敗，請確認檔案格式"})
		return
	}

	os.Remove(zipPath)

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("題目 %s 測資上傳並更新完成", problemID)})
}

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
		testDataDir = filepath.Join("test_data", problemID)
	}

	if err := os.RemoveAll(testDataDir); err != nil {
		fmt.Printf("[Warning] 測資目錄刪除失敗: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("題目 %s 已成功刪除，測資檔案已清除", problemID),
	})
}
