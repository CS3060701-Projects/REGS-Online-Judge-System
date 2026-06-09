package judge

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regs-backend/internal/models"
	"strings"
	"syscall"
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

	if strings.Contains(outputStr, "No tests were found") {
		return models.JudgeResult{
			Status:     "SE",
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

		if isRuntimeError(outputStr, runErr) {
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

func isRuntimeError(output string, err error) bool {
	lower := strings.ToLower(output)
	if strings.Contains(lower, "segmentation fault") ||
		strings.Contains(lower, "core dumped") ||
		strings.Contains(lower, "bus error") ||
		strings.Contains(lower, "access violation") ||
		strings.Contains(lower, "abort") ||
		strings.Contains(lower, "sigabrt") ||
		strings.Contains(lower, "sigsegv") {
		return true
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
			if status.Signaled() && status.Signal() != syscall.SIGKILL {
				return true
			}
		}
	}

	return false
}

func resolveCMakeSource(problemRoot string) string {
	if _, err := os.Stat(filepath.Join(problemRoot, "solution", "CMakeLists.txt")); err == nil {
		return "/problem/solution"
	}
	return "/problem"
}

func cmakeConfigureArgs(problemRoot string) []string {
	cmakeSource := resolveCMakeSource(problemRoot)
	args := []string{
		"-S", cmakeSource,
		"-B", "build",
		"-D", "SOURCE_DIR=/upload",
		"-D", "SOURCE_ROOT=/upload",
		"-G", "Ninja",
	}
	if cmakeSource == "/problem" {
		args = append(args, "-D", "PROBLEM_ROOT=/problem")
	}
	return args
}

func RunConfigure(absWorkspace, absProblemRoot string, configLogPath string) error {
	configLog, err := os.Create(configLogPath)
	if err != nil {
		return err
	}
	defer configLog.Close()

	problemRoot := filepath.Clean(filepath.FromSlash(absProblemRoot))
	args := append([]string{
		"run", "--rm",
		"--network", BuildNetwork,
		"-v", absWorkspace + ":/upload",
		"-v", absProblemRoot + ":/problem",
		"-v", absWorkspace + ":/app",
		"-w", "/app",
		models.JUDGER_IMAGE,
		"cmake",
	}, cmakeConfigureArgs(problemRoot)...)

	cmdConfig := exec.Command("docker", args...)
	cmdConfig.Stdout = configLog
	cmdConfig.Stderr = configLog
	return cmdConfig.Run()
}

func RunBuild(absWorkspace, absProblemRoot string, compileLogPath string) error {
	compileLog, err := os.Create(compileLogPath)
	if err != nil {
		return err
	}
	defer compileLog.Close()

	cmdBuild := exec.Command("docker", "run", "--rm",
		"--network", BuildNetwork,
		"-v", absWorkspace+":/upload",
		"-v", absProblemRoot+":/problem",
		"-v", absWorkspace+":/app",
		"-w", "/app",
		models.JUDGER_IMAGE,
		"cmake",
		"--build", "build",
		"--verbose",
	)

	cmdBuild.Stdout = compileLog
	cmdBuild.Stderr = compileLog
	return cmdBuild.Run()
}
