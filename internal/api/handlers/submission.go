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
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無法讀取上傳的檔案"})
		return
	}

	operatorID := uuid.New().String()

	zipDir := filepath.Join("storage", "submissions")
	os.MkdirAll(zipDir, os.ModePerm)
	zipPath := filepath.Join(zipDir, operatorID+".zip")

	workspaceDir := filepath.Join("storage", "workspaces", operatorID)
	os.MkdirAll(workspaceDir, os.ModePerm)

	if err := c.SaveUploadedFile(file, zipPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "儲存檔案失敗"})
		return
	}

	if err := utils.Unzip(zipPath, workspaceDir); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解壓縮失敗"})
		return
	}
	go processSubmission(operatorID, workspaceDir)

	c.JSON(http.StatusAccepted, gin.H{
		"message":    "任務已受理，正在背景評測中",
		"operatorId": operatorID,
	})
}

func processSubmission(operatorID string, workspace string) {
	fmt.Printf("\n[背景任務啟動] OperatorID: %s\n", operatorID)

	status := judge.CompileProject(operatorID, workspace)

	if status == "Ready" {
		status = judge.RunAndJudge(operatorID, workspace)
	}

	fmt.Printf("[任務結束] 最終評測結果: %s\n", status)
}
