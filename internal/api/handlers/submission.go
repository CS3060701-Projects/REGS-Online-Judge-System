package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"regs-backend/pkg/utils"

	"regs-backend/internal/judge"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SubmitCode(c *gin.Context) {
	// 1. 取得上傳的檔案 (假設表單的 key 叫做 "file")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無法讀取上傳的檔案"})
		return
	}

	// 2. 產生唯一的 operatorId
	operatorID := uuid.New().String()

	// 3. 建立儲存路徑
	// 存放原始 zip 檔的地方
	zipDir := filepath.Join("storage", "submissions")
	os.MkdirAll(zipDir, os.ModePerm)
	zipPath := filepath.Join(zipDir, operatorID+".zip")

	// 存放解壓縮後代碼的地方 (這就是要掛載到 Docker 的 Workspace)
	workspaceDir := filepath.Join("storage", "workspaces", operatorID)
	os.MkdirAll(workspaceDir, os.ModePerm)

	// 4. 儲存上傳的 Zip 檔
	if err := c.SaveUploadedFile(file, zipPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "儲存檔案失敗"})
		return
	}

	// 5. 解壓縮到 Workspace
	if err := utils.Unzip(zipPath, workspaceDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解壓縮失敗"})
		return
	}

	// 6. [重點] 在背景啟動判題邏輯 (Goroutine)
	go processSubmission(operatorID, workspaceDir)

	// 7. 立刻回傳 operatorId 給前端
	c.JSON(http.StatusAccepted, gin.H{
		"message":    "任務已受理，正在背景評測中",
		"operatorId": operatorID,
	})
}

// 這是我們之後要串接 Docker 的地方
func processSubmission(operatorID string, workspace string) {
	fmt.Printf("\n[背景任務啟動] OperatorID: %s\n", operatorID)

	// 1. 編譯
	status := judge.CompileProject(operatorID, workspace)

	// 2. 如果編譯成功，才執行判題
	if status == "Ready" {
		status = judge.RunAndJudge(operatorID, workspace)
	}

	fmt.Printf("[任務結束] 最終評測結果: %s\n", status)
}
