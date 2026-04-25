package judge

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const TIME_LIMIT = 2 // second

func RunAndJudge(operatorID string, workspace string, problemID string) string {
	absWorkspace, _ := filepath.Abs(workspace)
	absWorkspace = filepath.ToSlash(absWorkspace)

	testDataDir := filepath.Join("test_data", problemID)
	testCases, err := filepath.Glob(filepath.Join(testDataDir, "*.in"))

	if err != nil || len(testCases) == 0 {
		fmt.Printf("[%s] 錯誤: 找不到題目 %s 的測資\n", operatorID, problemID)
		return "SE"
	}

	for _, inputPath := range testCases {
		expectedPath := strings.TrimSuffix(inputPath, ".in") + ".out"
		outputLogPath := filepath.Join(workspace, "output.log")
		os.Remove(outputLogPath)

		outputLog, _ := os.Create(outputLogPath)

		ctx, cancel := context.WithTimeout(context.Background(), TIME_LIMIT*time.Second)
		defer cancel()

		cmdRun := exec.CommandContext(ctx, "docker", "run", "--rm", "-i",
			"-v", absWorkspace+":/workspace",
			"-w", "/workspace",
			"--network", "none",
			"yhlib/cs3060701",
			"./build/main",
		)

		// 串接測資輸入
		inFile, _ := os.Open(inputPath)
		cmdRun.Stdin = inFile
		cmdRun.Stdout = outputLog
		cmdRun.Stderr = outputLog

		runErr := cmdRun.Run()
		inFile.Close()
		outputLog.Close()

		// 判定結果
		if ctx.Err() == context.DeadlineExceeded {
			return "TLE"
		}
		if runErr != nil {
			return "RE"
		}

		// 字串比對 (強化版)
		userOut, _ := os.ReadFile(outputLogPath)
		expectedOut, _ := os.ReadFile(expectedPath)

		if cleanString(string(userOut)) != cleanString(string(expectedOut)) {
			// 如果 WA，印出 Debug 資訊
			fmt.Printf("DEBUG [%s]: User[%s] != Expected[%s]\n",
				operatorID, cleanString(string(userOut)), cleanString(string(expectedOut)))
			return "WA"
		}
	}

	return "AC"
}

// 輔助函式：清理字串
func cleanString(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "")
	return strings.TrimSpace(s)
}

func CompileProject(operatorID string, workspace string) string {
	// 1. 取得絕對路徑並轉換為 Docker 友善的正斜線格式
	absWorkspace, err := filepath.Abs(workspace)
	if err != nil {
		fmt.Printf("[%s] 取得絕對路徑失敗: %v\n", operatorID, err)
		return "SE"
	}
	absWorkspace = filepath.ToSlash(absWorkspace)

	// 2. [關鍵修復] 強制清理舊的 build 資料夾
	// 這能確保 Docker 跑的是你「現在」上傳的程式，而不是映像檔內的舊殘留
	buildPath := filepath.Join(workspace, "build")
	if err := os.RemoveAll(buildPath); err != nil {
		fmt.Printf("[%s] 清理舊編譯目錄失敗: %v\n", operatorID, err)
		// 如果是被 OneDrive 鎖定導致無法刪除，這裡會報錯，建議搬離 OneDrive
	}

	fmt.Printf("[%s] 開始編譯流程，Workspace: %s\n", operatorID, absWorkspace)

	// ==========================================
	// 階段一：CMake Configure (產生 Ninja 建置檔)
	// ==========================================
	configLogPath := filepath.Join(workspace, "configure.log")
	configLog, _ := os.Create(configLogPath)
	defer configLog.Close()

	// 這裡執行 cmake -G Ninja -B build
	cmdConfig := exec.Command("docker", "run", "--rm",
		"-v", absWorkspace+":/workspace",
		"-w", "/workspace",
		"yhlib/cs3060701",
		"cmake", "-G", "Ninja", "-B", "build",
	)

	cmdConfig.Stdout = configLog
	cmdConfig.Stderr = configLog

	fmt.Printf("[%s] [1/2] 正在執行 CMake Configure...\n", operatorID)
	if err := cmdConfig.Run(); err != nil {
		fmt.Printf("[%s] Configure 失敗，請檢查日誌: %s\n", operatorID, configLogPath)
		return "CE"
	}

	// ==========================================
	// 階段二：CMake Build (正式編譯)
	// ==========================================
	compileLogPath := filepath.Join(workspace, "compile.log")
	compileLog, _ := os.Create(compileLogPath)
	defer compileLog.Close()

	// 這裡執行 cmake --build build
	cmdBuild := exec.Command("docker", "run", "--rm",
		"-v", absWorkspace+":/workspace",
		"-w", "/workspace",
		"yhlib/cs3060701",
		"cmake", "--build", "build",
	)

	cmdBuild.Stdout = compileLog
	cmdBuild.Stderr = compileLog

	fmt.Printf("[%s] [2/2] 正在執行 CMake Build...\n", operatorID)
	if err := cmdBuild.Run(); err != nil {
		fmt.Printf("[%s] Build 失敗 (可能是語法錯誤)，請檢查日誌: %s\n", operatorID, compileLogPath)
		return "CE"
	}

	fmt.Printf("[%s] 編譯成功！\n", operatorID)
	return "Ready"
}
