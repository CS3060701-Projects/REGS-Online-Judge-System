package handlers

import (
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
	"regs-backend/internal/judge"
	"regs-backend/internal/models"
	"regs-backend/pkg/utils"
)

func SubmitAssignment(c *gin.Context) {
	problemID := c.DefaultPostForm("problem_id", "p1001")

	var problem models.Problem
	if err := database.DB.Select("id").Where("id = ?", problemID).First(&problem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "提交失敗：找不到指定的題目"})
		return
	}

	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授權的操作"})
		return
	}

	var uID uint
	switch v := val.(type) {
	case float64:
		uID = uint(v)
	case uint:
		uID = v
	case int:
		uID = uint(v)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法解析使用者 ID"})
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
	archiveDir := filepath.Join("storage", "submissions")
	workspace := filepath.Join("storage", "workspaces", operatorID)

	if err := os.MkdirAll(archiveDir, os.ModePerm); err != nil || os.MkdirAll(workspace, os.ModePerm) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "伺服器資料夾建立失敗"})
		return
	}

	zipPath := filepath.Join(archiveDir, operatorID+".zip")
	if err := c.SaveUploadedFile(file, zipPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "檔案儲存失敗"})
		return
	}

	if err := utils.Unzip(zipPath, workspace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "解壓縮失敗，請檢查檔案格式"})
		return
	}

	submission := models.Submission{
		OperatorID: operatorID, // 確保你的 Model 欄位名稱正確
		ProblemID:  problemID,
		UserID:     uID,
		Status:     "Pending",
	}

	if err := database.DB.Create(&submission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "資料庫紀錄建立失敗", "details": err.Error()})
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
		"userId":     uID,
	})
}

func processSubmission(operatorID, workspace, problemID string) {
	var problem models.Problem
	if err := database.DB.First(&problem, "id = ?", problemID).Error; err != nil {
		updateSubmissionStatus(operatorID, "SE")
		return
	}

	absWorkspace, _ := filepath.Abs(workspace)
	updateSubmissionStatus(operatorID, "Compiling")

	configCmd := exec.Command("docker", "run", "--rm",
		"-v", absWorkspace+":/app", "-w", "/app",
		models.JUDGER_IMAGE, "cmake", "-G", "Ninja", "-B", "build",
	)
	configOut, _ := configCmd.CombinedOutput()
	os.WriteFile(filepath.Join(workspace, "configure.log"), configOut, 0644)
	if !configCmd.ProcessState.Success() {
		updateSubmissionStatus(operatorID, "SE")
		return
	}

	buildCmd := exec.Command("docker", "run", "--rm",
		"-v", absWorkspace+":/app", "-w", "/app",
		models.JUDGER_IMAGE, "cmake", "--build", "build",
	)
	buildOut, _ := buildCmd.CombinedOutput()
	os.WriteFile(filepath.Join(workspace, "compile.log"), buildOut, 0644)
	if !buildCmd.ProcessState.Success() {
		updateSubmissionStatus(operatorID, "CE")
		return
	}

	updateSubmissionStatus(operatorID, "Judging")
	result := judge.RunAndJudge(operatorID, workspace, problem)

	updateSubmissionStatus(operatorID, result.Status)

	database.DB.Model(&models.Submission{}).Where("operator_id = ?", operatorID).
		Updates(map[string]interface{}{
			"run_time":   int(result.PeakTime * 1000),
			"run_memory": result.PeakMemory,
		})

	fmt.Printf("[評測結束] ID: %s | 結果: %s | 峰值耗時: %.3fms | 記憶體: %d KB\n",
		operatorID, result.Status, result.PeakTime*1000, result.PeakMemory)
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
		"run_time":   submission.RunTime,   // ms
		"run_memory": submission.RunMemory, // KB
	})
}

func GetSubmissionLog(c *gin.Context) {
	operatorID := c.Param("operatorId")
	logType := c.Param("type")

	var fileName string
	switch logType {
	case "configure":
		fileName = "configure.log"
	case "compile":
		fileName = "compile.log"
	case "output":
		fileName = "output.log"
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的日誌類型，僅限 configure, compile, output"})
		return
	}

	logPath := filepath.Join("storage", "workspaces", operatorID, fileName)

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("找不到指定的日誌檔案: %s", fileName)})
		return
	}

	c.File(logPath)
}

func GetSubmissions(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未經授權的存取"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var submissions []models.Submission
	var total int64

	query := database.DB.Model(&models.Submission{}).Where("user_id = ?", userID)

	query.Count(&total)

	result := query.Preload("Problem").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&submissions)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查詢提交紀錄失敗"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total": total,
		"page":  page,
		"limit": limit,
		"data":  submissions,
	})
}

func GetUserSubmissions(c *gin.Context) {
	targetUserIDStr := c.Param("user_id")

	currentUserID, _ := c.Get("user_id")
	currentRole, _ := c.Get("role")

	targetUserID, err := strconv.ParseUint(targetUserIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的使用者 ID"})
		return
	}

	if currentRole != "Admin" && uint(targetUserID) != currentUserID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "你沒有權限查看他人的提交紀錄"})
		return
	}

	var submissions []models.Submission
	if err := database.DB.Preload("Problem").
		Where("user_id = ?", targetUserID).
		Order("created_at desc"). // 由新到舊排序
		Find(&submissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得提交紀錄"})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

func GetSubmissionSource(c *gin.Context) {
	operatorID := c.Param("operatorId")
	val, _ := c.Get("user_id")
	currentRole, _ := c.Get("role")

	var currentUserID uint
	if v, ok := val.(uint); ok {
		currentUserID = v
	} else if f, ok := val.(float64); ok {
		currentUserID = uint(f)
	}

	var submission models.Submission
	if err := database.DB.Where("operator_id = ?", operatorID).First(&submission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到該筆提交紀錄"})
		return
	}

	if currentRole != "Admin" && submission.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "你沒有權限查看此原始碼"})
		return
	}

	filePath := filepath.Join("storage", "submissions", operatorID+".zip")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "原始檔案不存在或已被清理"})
		return
	}

	downloadName := fmt.Sprintf("submission_%s.zip", operatorID)
	c.FileAttachment(filePath, downloadName)
}
