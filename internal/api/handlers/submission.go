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

// SubmitAssignment godoc
// @Summary Submit code for judging
// @Description Upload a .zip file containing source code for a specific problem.
// @Tags Submissions
// @Accept  multipart/form-data
// @Produce  json
// @Security Bearer
// @Param   problem_id formData string true "Problem ID (e.g., p1001)"
// @Param   file formData file true "Source code as a .zip file"
// @Success 200 {object} object{message=string, operatorId=string, userId=integer} "提交成功，開始評測"
// @Failure 400 {object} object{error=string} "請求錯誤"
// @Failure 500 {object} object{error=string} "伺服器內部錯誤"
// @Router /submissions [post]
func SubmitAssignment(c *gin.Context) {
	problemID := c.PostForm("problem_id")
	if problemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "提交失敗：缺少題目 ID (problem_id)"})
		return
	}

	var problem models.Problem
	if err := database.DB.Select("id").Where("id = ?", problemID).First(&problem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "提交失敗：找不到指定的題目"})
		return
	}

	val, exists := c.Get("user_id")
	uID, ok := val.(uint)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授權的操作"})
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
		OperatorID: operatorID,
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

	if _, err := os.Stat(filepath.Join(workspace, "CMakeLists.txt")); os.IsNotExist(err) {
		updateSubmissionStatus(operatorID, "SE")
		return
	}

	configCmd := exec.Command(
		"docker", "run", "--rm",
		"--network", "none",
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
		"--network", "none",
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

	fmt.Printf("[評測結束] ID: %s | 結果: %s | 耗時: %.3fms | 記憶體: %d KB\n",
		operatorID, result.Status, result.PeakTime*1000, result.PeakMemory)
}

func updateSubmissionStatus(operatorID string, status string) {
	database.DB.Model(&models.Submission{}).
		Where("operator_id = ?", operatorID).
		Update("status", status)
}

// GetSubmissionStatus godoc
// @Summary Get submission status
// @Description Retrieves the current status and result of a specific submission.
// @Tags Submissions
// @Produce  json
// @Security Bearer
// @Param   operatorId path string true "Operator ID of the submission"
// @Success 200 {object} object{operatorId=string, status=string, created_at=string, run_time=integer, run_memory=integer}
// @Failure 404 {object} object{error=string} "找不到該筆評測紀錄"
// @Router /submissions/{operatorId} [get]
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
		"run_time":   submission.RunTime,
		"run_memory": submission.RunMemory,
	})
}

// GetSubmissionLog godoc
// @Summary Get submission log file
// @Description Downloads a specific log file (configure, compile, or output) for a submission.
// @Tags Submissions
// @Produce  plain
// @Security Bearer
// @Param   operatorId path string true "Operator ID of the submission"
// @Param   type path string true "Log type" Enums(configure, compile, output)
// @Success 200 {file} file "Log file content"
// @Failure 400 {object} object{error=string} "無效的日誌類型"
// @Failure 404 {object} object{error=string} "找不到指定的日誌檔案"
// @Router /submissions/{operatorId}/logs/{type} [get]
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的日誌類型"})
		return
	}

	logPath := filepath.Join("storage", "workspaces", operatorID, fileName)

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到指定的日誌檔案"})
		return
	}

	c.File(logPath)
}

func getSubmissionsByUserID(userID string) ([]models.Submission, error) {
	var submissions []models.Submission
	err := database.DB.Preload("Problem").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&submissions).Error
	return submissions, err
}

// GetSubmissions godoc
// @Summary Get personal submission history
// @Description Retrieves a list of all submissions made by the currently authenticated user.
// @Tags Submissions
// @Produce  json
// @Security Bearer
// @Success 200 {array} models.Submission
// @Failure 401 {object} object{error=string} "未授權的操作"
// @Failure 500 {object} object{error=string} "無法取得提交紀錄"
// @Router /submissions [get]
func GetSubmissions(c *gin.Context) {
	val, exists := c.Get("user_id")
	currentUID, ok := val.(uint)
	if !exists || !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授權的操作"})
		return
	}

	submissions, err := getSubmissionsByUserID(fmt.Sprint(currentUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得提交紀錄"})
		return
	}
	c.JSON(http.StatusOK, submissions)
}

// GetUserSubmissions godoc
// @Summary Get a specific user's submission history
// @Description Retrieves a list of all submissions made by a specific user.
// @Tags Submissions
// @Produce  json
// @Param   user_id path integer true "User ID"
// @Success 200 {array} models.Submission
// @Router /users/{user_id}/submissions [get]
func GetUserSubmissions(c *gin.Context) {
	targetUserID := c.Param("user_id")

	if _, err := strconv.Atoi(targetUserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "無效的使用者 ID"})
		return
	}

	submissions, err := getSubmissionsByUserID(targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "無法取得提交紀錄"})
		return
	}

	c.JSON(http.StatusOK, submissions)
}

// GetSubmissionSource godoc
// @Summary Download submission source code
// @Description Downloads the original .zip file for a submission. Only accessible by the owner or an admin.
// @Tags Submissions
// @Produce  application/zip
// @Security Bearer
// @Param   operatorId path string true "Operator ID of the submission"
// @Success 200 {file} file "The submission's source code as a .zip file"
// @Router /submissions/{operatorId}/source [get]
func GetSubmissionSource(c *gin.Context) {
	operatorID := c.Param("operatorId")
	val, _ := c.Get("user_id")
	currentRole, _ := c.Get("role")

	currentUID, ok := val.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授權的操作"})
		return
	}

	var submission models.Submission
	if err := database.DB.Where("operator_id = ?", operatorID).First(&submission).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "找不到提交紀錄"})
		return
	}

	if currentRole != "Admin" && submission.UserID != currentUID {
		c.JSON(http.StatusForbidden, gin.H{"error": "你沒有權限查看此原始碼"})
		return
	}

	filePath := filepath.Join("storage", "submissions", operatorID+".zip")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "原始檔案不存在"})
		return
	}

	downloadName := fmt.Sprintf("submission_%s.zip", operatorID)
	c.FileAttachment(filePath, downloadName)
}
