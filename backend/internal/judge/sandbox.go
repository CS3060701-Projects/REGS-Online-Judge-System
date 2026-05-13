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
	problemRoot := problem.TestcasePath
	if problemRoot == "" {
		problemRoot = filepath.Join("testdata", problem.ID)
	}

	cmakePath := filepath.Join(problemRoot, "CMakeLists.txt")
	if _, err := os.Stat(cmakePath); err != nil {
		fmt.Printf("[%s] 題目資料夾缺少 CMakeLists.txt: %s\n", operatorID, cmakePath)
		return models.JudgeResult{Status: "SE"}
	}

	return RunAndJudgeCTest(operatorID, workspace, problem)
}

func RunAndJudgeCTest(operatorID string, workspace string, problem models.Problem) models.JudgeResult {
	absWorkspace, err := filepath.Abs(workspace)
	if err != nil {
		fmt.Printf("[%s] 無法取得 workspace 絕對路徑: %v\n", operatorID, err)
		return models.JudgeResult{Status: "SE"}
	}
	absWorkspace = filepath.ToSlash(absWorkspace)

	problemRoot := problem.TestcasePath
	if problemRoot == "" {
		problemRoot = filepath.Join("testdata", problem.ID)
	}

	absProblemRoot, err := filepath.Abs(problemRoot)
	if err != nil {
		fmt.Printf("[%s] 無法取得題目根目錄絕對路徑: %v\n", operatorID, err)
		return models.JudgeResult{Status: "SE"}
	}
	absProblemRoot = filepath.ToSlash(absProblemRoot)

	outputLogPath := filepath.Join(workspace, "output.log")
	os.Remove(outputLogPath)

	logFile, err := os.OpenFile(outputLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("[%s] 無法建立 output.log: %v\n", operatorID, err)
		return models.JudgeResult{Status: "SE"}
	}
	defer logFile.Close()

	timeout := time.Duration(problem.TimeLimit*10) * time.Millisecond
	if timeout < 10*time.Second {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmdRun := exec.CommandContext(ctx, "docker", "run", "--rm",
		"--network", "none",
		"--cpus", "1.0",
		"--memory", fmt.Sprintf("%dm", problem.MemoryLimit),
		"-v", absWorkspace+":/upload",
		"-v", absProblemRoot+":/problem:ro",
		"-v", absWorkspace+":/app",
		"-w", "/app",
		models.JUDGER_IMAGE,
		"ctest",
		"--test-dir", "build",
		"--output-on-failure",
		"-V",
	)

	var stdoutBuf bytes.Buffer
	cmdRun.Stdout = &stdoutBuf
	cmdRun.Stderr = &stdoutBuf

	start := time.Now()
	runErr := cmdRun.Run()
	elapsed := time.Since(start).Seconds()

	outputStr := stdoutBuf.String()

	fmt.Printf("[%s] CTest 執行結果:\n%s\n", operatorID, outputStr)

	logFile.WriteString("=== CTest Execution ===\n")
	logFile.Write(stdoutBuf.Bytes())
	logFile.WriteString("\n")

	if ctx.Err() == context.DeadlineExceeded {
		return models.JudgeResult{
			Status:     "TLE",
			PeakTime:   elapsed,
			PeakMemory: 0,
		}
	}

	if runErr != nil {
		if strings.Contains(outputStr, "Timeout") ||
			strings.Contains(outputStr, "TIMEOUT") ||
			strings.Contains(outputStr, "Test timeout") {
			return models.JudgeResult{
				Status:     "TLE",
				PeakTime:   elapsed,
				PeakMemory: 0,
			}
		}

		if strings.Contains(outputStr, "Segmentation fault") ||
			strings.Contains(outputStr, "segmentation fault") ||
			strings.Contains(outputStr, "core dumped") ||
			strings.Contains(outputStr, "Bus error") ||
			strings.Contains(outputStr, "Access violation") {
			return models.JudgeResult{
				Status:     "RE",
				PeakTime:   elapsed,
				PeakMemory: 0,
			}
		}

		return models.JudgeResult{
			Status:     "WA",
			PeakTime:   elapsed,
			PeakMemory: 0,
		}
	}

	return models.JudgeResult{
		Status:     "AC",
		PeakTime:   elapsed,
		PeakMemory: 0,
	}
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
		"--network", "none",
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
		"--network", "none",
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
