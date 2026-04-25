package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"regs-backend/internal/database"
	"regs-backend/internal/models"
	"regs-backend/pkg/utils"
)

func SubmitAssignment(c *gin.Context) {
	// 1. 從 form-data 讀取 problem_id 與 user_id
	problemID := c.DefaultPostForm("problem_id", "p1001")
	userIDStr := c.DefaultPostForm("user_id", "2") // 預設使用 2 以符合你的測試

	// 將 user_id (字串) 轉換為 uint (無號整數)
	uID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id 格式錯誤"})
		return
	}

	// 2. 接收上傳的 ZIP 檔案
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "請上傳檔案"})
		return
	}

	// 3. 檢查副檔名 (簡單驗證)
	if !strings.HasSuffix(file.Filename, ".zip") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "只接受 .zip 格式"})
		return
	}

	// 4. 生成唯一的 OperatorID (UUID)
	operatorID := uuid.New().String()
	workspace := filepath.Join("storage", "workspaces", operatorID)

	// 5. 建立該次評測的獨立工作空間
	if err := os.MkdirAll(workspace, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法建立工作空間"})
		return
	}

	// 6. 將檔案儲存到本地
	zipPath := filepath.Join(workspace, "submission.zip")
	if err := c.SaveUploadedFile(file, zipPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "檔案儲存失敗"})
		return
	}

	// 7. 解壓縮 ZIP
	if err := utils.Unzip(zipPath, workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解壓縮失敗，請檢查檔案格式"})
		return
	}

	// 8. 建立資料庫紀錄 (將 uID 正確寫入)
	submission := models.Submission{
		OperatorID: operatorID,
		ProblemID:  problemID,
		UserID:     uint(uID),
		Status:     "Pending",
	}

	// 若寫入失敗，回傳詳細的 err.Error() 方便除錯
	if err := database.DB.Create(&submission).Error; err != nil {
		fmt.Printf("[DB Error] 建立 Submission 失敗: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "資料庫紀錄建立失敗",
			"details": err.Error(),
		})
		return
	}

	// 9. 啟動背景評測任務
	go processSubmission(operatorID, workspace, problemID)

	// 10. 回傳成功訊息給前端/Postman
	c.JSON(http.StatusOK, gin.H{
		"message":    "提交成功，開始評測",
		"operatorId": operatorID,
		"problemId":  problemID,
		"userId":     uint(uID),
	})
}

func processSubmission(operatorID, workspace, problemID string) {
	fmt.Printf("\n[背景任務啟動] OperatorID: %s, 題目: %s\n", operatorID, problemID)

	// 1. 取得工作空間與測資的絕對路徑 (Docker volume 掛載必須使用絕對路徑)
	absWorkspace, _ := filepath.Abs(workspace)
	absTestData, _ := filepath.Abs(filepath.Join("test_data", problemID))

	// 2. 更新狀態為「編譯中」
	updateSubmissionStatus(operatorID, "Compiling")

	// ==========================================
	// 階段一：檢查 CMakeLists.txt 是否存在
	// ==========================================
	cmakePath := filepath.Join(workspace, "CMakeLists.txt")
	if _, err := os.Stat(cmakePath); os.IsNotExist(err) {
		fmt.Printf("[評測中斷] 找不到 CMakeLists.txt\n")
		updateSubmissionStatus(operatorID, "CE") // Compile Error
		return
	}

	// ==========================================
	// 階段二：透過 Docker 執行 CMake 與 Make
	// ==========================================
	// 這裡假設你使用的編譯環境映像檔為 gcc:latest 或 cmake 相關的 image
	// 指令: cd /app && cmake . && make
	compileCmd := exec.Command("docker", "run", "--rm",
		"-v", fmt.Sprintf("%s:/app", absWorkspace),
		"-w", "/app",
		"yhlib/cs3060701", // ⚠️ 請根據你的 Docker 映像檔名稱進行修改
		"sh", "-c", "cmake . && make",
	)

	compileOut, err := compileCmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[編譯失敗] %v\n輸出訊息:\n%s\n", err, string(compileOut))
		updateSubmissionStatus(operatorID, "CE") // Compile Error
		return
	}
	fmt.Println("[編譯成功] 準備進入評測階段")

	// 3. 更新狀態為「評測中」
	updateSubmissionStatus(operatorID, "Judging")

	// ==========================================
	// 階段三：執行測資評測 (讀取 .in，比對 .out)
	// ==========================================
	// 讀取該題目的測資資料夾
	testFiles, err := os.ReadDir(absTestData)
	if err != nil {
		fmt.Printf("[系統錯誤] 無法讀取測資目錄: %s\n", absTestData)
		updateSubmissionStatus(operatorID, "SE") // System Error
		return
	}

	allPassed := true

	for _, file := range testFiles {
		if strings.HasSuffix(file.Name(), ".in") {
			testName := strings.TrimSuffix(file.Name(), ".in")
			inFile := filepath.Join(absTestData, file.Name())
			outFile := filepath.Join(absTestData, testName+".out")

			if _, err := os.Stat(outFile); os.IsNotExist(err) {
				continue
			}

			runCmd := exec.Command("docker", "run", "--rm", "-i",
				"--network", "none",
				"--memory", "256m",
				"--cpus", "1.0",
				"--pids-limit", "50",
				"-v", fmt.Sprintf("%s:/app", absWorkspace),
				"-w", "/app",
				"yhlib/cs3060701",
				"./main",
			)

			// 讀取測資輸入
			inData, _ := os.ReadFile(inFile)
			runCmd.Stdin = bytes.NewReader(inData)

			var outBuffer bytes.Buffer
			runCmd.Stdout = &outBuffer

			// 執行程式
			err := runCmd.Run()
			if err != nil {
				fmt.Printf("測資 %s 執行錯誤 (RE): %v\n", testName, err)
				updateSubmissionStatus(operatorID, "RE") // Runtime Error
				allPassed = false
				break
			}

			// 讀取標準答案並進行比對 (去除頭尾空白與換行)
			expectedOut, _ := os.ReadFile(outFile)
			expectedStr := strings.TrimSpace(string(expectedOut))
			actualStr := strings.TrimSpace(outBuffer.String())

			if expectedStr != actualStr {
				fmt.Printf("測資 %s 答案錯誤 (WA)\n預期: %s\n實際: %s\n", testName, expectedStr, actualStr)
				updateSubmissionStatus(operatorID, "WA") // Wrong Answer
				allPassed = false
				break
			} else {
				fmt.Printf("測資 %s 通過!\n", testName)
			}
		}
	}

	if allPassed {
		fmt.Printf("OperatorID: %s 全數測資通過 (AC)!\n", operatorID)
		updateSubmissionStatus(operatorID, "AC") // Accepted
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
