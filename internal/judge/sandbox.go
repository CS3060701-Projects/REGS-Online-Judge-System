package judge

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regs-backend/internal/models"
	"strings"
	"time"
)

func RunAndJudge(operatorID string, workspace string, problem models.Problem) models.JudgeResult {
	absWorkspace, _ := filepath.Abs(workspace)
	absTestData, _ := filepath.Abs(filepath.Join("test_data", problem.ID))

	// 取得所有 .in 檔案
	testCases, _ := filepath.Glob(filepath.Join(absTestData, "*.in"))
	if len(testCases) == 0 {
		return models.JudgeResult{Status: "SE"}
	}

	var peakTime float64
	var peakMemory int64
	outputLogPath := filepath.Join(workspace, "output.log")
	os.Remove(outputLogPath) // 執行前清空舊日誌

	// 建立或開啟 output.log 用於記錄所有輸出
	logFile, _ := os.OpenFile(outputLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer logFile.Close()

	for _, inputPath := range testCases {
		testName := strings.TrimSuffix(filepath.Base(inputPath), ".in")
		expectedPath := strings.TrimSuffix(inputPath, ".in") + ".out"

		// 設定超時限制 (由 Problem 模型提供)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(problem.TimeLimit)*time.Millisecond)
		defer cancel()

		cmdRun := exec.CommandContext(ctx, "docker", "run", "--rm", "-i",
			"--network", "none",
			"--cpus", "1.0",
			"--memory", fmt.Sprintf("%dm", problem.MemoryLimit),
			"-v", fmt.Sprintf("%s:/app", absWorkspace),
			"-w", "/app",
			models.JUDGER_IMAGE,
			"/usr/bin/time", "-f", "METRIC:%e:%M", "./build/main",
		)

		inFile, _ := os.Open(inputPath)
		var stdoutBuf, stderrBuf bytes.Buffer
		cmdRun.Stdin = inFile
		cmdRun.Stdout = &stdoutBuf
		cmdRun.Stderr = &stderrBuf

		runErr := cmdRun.Run()
		inFile.Close()

		// --- 1. 解析效能數據 ---
		var currentTime float64
		var currentMemory int64
		for _, line := range strings.Split(stderrBuf.String(), "\n") {
			if strings.Contains(line, "METRIC:") {
				fmt.Sscanf(strings.TrimSpace(line), "METRIC:%f:%d", &currentTime, &currentMemory)
			}
		}
		if currentTime > peakTime {
			peakTime = currentTime
		}
		if currentMemory > peakMemory {
			peakMemory = currentMemory
		}

		// --- 2. 寫入日誌檔案 ---
		logFile.WriteString(fmt.Sprintf("=== Test %s (%.3fs, %d KB) ===\n", testName, currentTime, currentMemory))
		logFile.Write(stdoutBuf.Bytes())
		logFile.WriteString("\n")

		// --- 3. 判定結果 ---
		if ctx.Err() == context.DeadlineExceeded {
			return models.JudgeResult{Status: "TLE", PeakTime: peakTime, PeakMemory: peakMemory}
		}
		if runErr != nil {
			return models.JudgeResult{Status: "RE", PeakTime: peakTime, PeakMemory: peakMemory}
		}

		expectedOut, _ := os.ReadFile(expectedPath)
		if strings.TrimSpace(stdoutBuf.String()) != strings.TrimSpace(string(expectedOut)) {
			return models.JudgeResult{Status: "WA", PeakTime: peakTime, PeakMemory: peakMemory}
		}
	}

	return models.JudgeResult{Status: "AC", PeakTime: peakTime, PeakMemory: peakMemory}
}

func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "")
	return strings.TrimSpace(s)
}

func CompileProject(operatorID string, workspace string) string {
	absWorkspace, err := filepath.Abs(workspace)
	if err != nil {
		fmt.Printf("[%s] 取得絕對路徑失敗: %v\n", operatorID, err)
		return "SE"
	}
	absWorkspace = filepath.ToSlash(absWorkspace)

	buildPath := filepath.Join(workspace, "build")
	if err := os.RemoveAll(buildPath); err != nil {
		fmt.Printf("[%s] 清理舊編譯目錄失敗: %v\n", operatorID, err)
	}

	fmt.Printf("[%s] 開始編譯流程，Workspace: %s\n", operatorID, absWorkspace)

	configLogPath := filepath.Join(workspace, "configure.log")
	configLog, _ := os.Create(configLogPath)
	defer configLog.Close()

	// cmake -G Ninja -B build
	cmdConfig := exec.Command("docker", "run", "--rm",
		"-v", absWorkspace+":/workspace",
		"-w", "/workspace",
		models.JUDGER_IMAGE,
		"cmake", "-G", "Ninja", "-B", "build",
	)

	cmdConfig.Stdout = configLog
	cmdConfig.Stderr = configLog

	fmt.Printf("[%s] [1/2] 正在執行 CMake Configure...\n", operatorID)
	if err := cmdConfig.Run(); err != nil {
		fmt.Printf("[%s] Configure 失敗，請檢查日誌: %s\n", operatorID, configLogPath)
		return "CE"
	}

	compileLogPath := filepath.Join(workspace, "compile.log")
	compileLog, _ := os.Create(compileLogPath)
	defer compileLog.Close()

	// cmake --build build
	cmdBuild := exec.Command("docker", "run", "--rm",
		"-v", absWorkspace+":/workspace",
		"-w", "/workspace",
		models.JUDGER_IMAGE,
		"cmake", "--build", "build",
	)

	cmdBuild.Stdout = compileLog
	cmdBuild.Stderr = compileLog

	fmt.Printf("[%s] [2/2] 正在執行 CMake Build...\n", operatorID)
	if err := cmdBuild.Run(); err != nil {
		fmt.Printf("[%s] Build 失敗，請檢查日誌: %s\n", operatorID, compileLogPath)
		return "CE"
	}

	fmt.Printf("[%s] 編譯成功！\n", operatorID)
	return "Ready"
}
