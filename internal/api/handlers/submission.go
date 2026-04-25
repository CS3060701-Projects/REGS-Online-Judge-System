package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"regs-backend/internal/database"
	"regs-backend/internal/models"
	"regs-backend/pkg/utils"
)

func SubmitAssignment(c *gin.Context) {
	problemID := c.DefaultPostForm("problem_id", "p1001")
	userIDStr := c.DefaultPostForm("user_id", "0") // 預設使用 2 以符合你的測試

	uID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id 格式錯誤"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請上傳檔案"})
		return
	}

	if !strings.HasSuffix(file.Filename, ".zip") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只接受 .zip 格式"})
		return
	}

	operatorID := uuid.New().String()
	workspace := filepath.Join("storage", "workspaces", operatorID)

	if err := os.MkdirAll(workspace, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立工作空間"})
		return
	}

	zipPath := filepath.Join(workspace, "submission.zip")
	if err := c.SaveUploadedFile(file, zipPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "檔案儲存失敗"})
		return
	}

	if err := utils.Unzip(zipPath, workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解壓縮失敗，請檢查檔案格式"})
		return
	}

	submission := models.Submission{
		OperatorID: operatorID,
		ProblemID:  problemID,
		UserID:     uint(uID),
		Status:     "Pending",
	}

	if err := database.DB.Create(&submission).Error; err != nil {
		fmt.Printf("[DB Error] 建立 Submission 失敗: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "資料庫紀錄建立失敗",
			"details": err.Error(),
		})
		return
	}

	JobQueue <- JudgeJob{
		OperatorID: operatorID,
		Workspace:  workspace,
		ProblemID:  problemID,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "提交成功，開始評測",
		"operatorId": operatorID,
		"problemId":  problemID,
		"userId":     uint(uID),
	})
}

func processSubmission(operatorID, workspace, problemID string) {
	absWorkspace, _ := filepath.Abs(workspace)
	absTestData, _ := filepath.Abs(filepath.Join("test_data", problemID))

	updateSubmissionStatus(operatorID, "Compiling")

	if _, err := os.Stat(filepath.Join(workspace, "CMakeLists.txt")); os.IsNotExist(err) {
		updateSubmissionStatus(operatorID, "SE")
		return
	}

	// 生成 configure.log 與判定 SE
	configCmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", absWorkspace),
		"-w", "/app",
		"yhlib/cs3060701",
		"cmake", "-G", "Ninja", "-B", "build",
	)
	configOut, err := configCmd.CombinedOutput()
	os.WriteFile(filepath.Join(workspace, "configure.log"), configOut, 0644) // 寫入日誌

	if err != nil {
		updateSubmissionStatus(operatorID, "SE")
		return
	}

	// 生成 compile.log 與判定 CE
	buildCmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", absWorkspace),
		"-w", "/app",
		"yhlib/cs3060701",
		"cmake", "--build", "build", "--verbose",
	)
	buildOut, err := buildCmd.CombinedOutput()
	os.WriteFile(filepath.Join(workspace, "compile.log"), buildOut, 0644) // 寫入日誌

	if err != nil {
		updateSubmissionStatus(operatorID, "CE")
		return
	}

	// 生成 output.log，判定 AC/WA/TLE/RE
	updateSubmissionStatus(operatorID, "Judging")
	testFiles, _ := os.ReadDir(absTestData)
	allPassed := true

	outLogFile, _ := os.OpenFile(filepath.Join(workspace, "output.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer outLogFile.Close()

	for _, file := range testFiles {
		if strings.HasSuffix(file.Name(), ".in") {
			testName := strings.TrimSuffix(file.Name(), ".in")
			inFile := filepath.Join(absTestData, file.Name())
			outFile := filepath.Join(absTestData, testName+".out")

			if _, err := os.Stat(outFile); os.IsNotExist(err) {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

			runCmd := exec.CommandContext(ctx, "docker", "run", "--rm", "-i",
				"--network", "none", // 禁止網路存取
				"-v", fmt.Sprintf("%s:/app", absWorkspace),
				"-w", "/app",
				"yhlib/cs3060701",
				"./build/main",
			)

			inData, _ := os.ReadFile(inFile)
			runCmd.Stdin = bytes.NewReader(inData)

			var outBuffer bytes.Buffer
			runCmd.Stdout = &outBuffer

			err := runCmd.Run()

			// 生成output.log
			outLogFile.WriteString(fmt.Sprintf("=== Test %s ===\n", testName))
			outLogFile.Write(outBuffer.Bytes())
			outLogFile.WriteString("\n")

			if ctx.Err() == context.DeadlineExceeded {
				updateSubmissionStatus(operatorID, "TLE")
				allPassed = false
				cancel()
				break
			} else if err != nil {
				updateSubmissionStatus(operatorID, "RE")
				allPassed = false
				cancel()
				break
			}

			expectedOut, _ := os.ReadFile(outFile)
			if strings.TrimSpace(string(expectedOut)) != strings.TrimSpace(outBuffer.String()) {
				updateSubmissionStatus(operatorID, "WA")
				allPassed = false
				cancel()
				break
			}
			cancel()
		}
	}

	if allPassed {
		updateSubmissionStatus(operatorID, "AC")
	}
}

func updateSubmissionStatus(operatorID string, status string) {
	err := database.DB.Model(&models.Submission{}).
		Where("operator_id = ?", operatorID).
		Update("status", status).Error

	if err != nil {
		fmt.Printf("狀態更新失敗 [%s -> %s]: %v\n", operatorID, status, err)
	}
}

func GetSubmissionStatus(c *gin.Context) {
	operatorID := c.Param("operatorId")

	var submission models.Submission
	if err := database.DB.Where("operator_id = ?", operatorID).First(&submission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到該筆評測紀錄"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"operatorId": submission.OperatorID,
		"status":     submission.Status,
		"created_at": submission.CreatedAt,
	})
}
