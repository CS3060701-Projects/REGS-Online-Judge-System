package judge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RunAndJudge(operatorID string, workspace string) string {
	absWorkspace, _ := filepath.Abs(workspace)
	absWorkspace = filepath.ToSlash(absWorkspace)

	outputLogPath := filepath.Join(workspace, "output.log")
	outputLog, _ := os.Create(outputLogPath)
	defer outputLog.Close()

	// 執行指令：在容器內跑 build/main (或是你在 CMake 裡 add_executable 的名字)
	// 這裡加上 --network none 滿足安全性要求
	cmdRun := exec.Command("docker", "run", "--rm",
		"-v", absWorkspace+":/workspace",
		"-w", "/workspace",
		"--network", "none",
		"yhlib/cs3060701",
		"./build/main",
	)

	cmdRun.Stdout = outputLog
	cmdRun.Stderr = outputLog

	fmt.Println("[3/3] 正在執行程式並收集輸出...")
	if err := cmdRun.Run(); err != nil {
		fmt.Println("執行失敗 (RE):", err)
		return "RE" // Runtime Error
	}

	// 讀取剛剛輸出的結果進行比對
	userOutput, _ := os.ReadFile(outputLogPath)
	actualOutput := strings.TrimSpace(string(userOutput))

	// 這裡先寫死一個標準答案，之後可以從資料庫讀取題目測資
	expectedOutput := "Hello, REGS System!"

	if actualOutput == expectedOutput {
		return "AC"
	} else {
		return "WA"
	}
}

// CompileProject 執行兩階段的 CMake 編譯，並回傳狀態碼 (SE, CE, 或 Ready)
func CompileProject(operatorID string, workspace string) string {
	// Docker 的 Volume 掛載 (-v) 必須使用「絕對路徑」
	absWorkspace, err := filepath.Abs(workspace)
	if err != nil {
		fmt.Println("取得絕對路徑失敗:", err)
		return "System Error"
	}
	absWorkspace = filepath.ToSlash(absWorkspace)

	fmt.Println("開始編譯流程，Workspace:", absWorkspace)

	// ==========================================
	// 階段一：預處理與配置 (Configure)
	// ==========================================
	configLogPath := filepath.Join(workspace, "configure.log")
	configLog, _ := os.Create(configLogPath)
	defer configLog.Close()

	// 組合 Docker 指令
	cmdConfig := exec.Command("docker", "run", "--rm",
		"-v", absWorkspace+":/workspace", // 掛載目錄
		"-w", "/workspace", // 設定容器內的工作目錄
		"yhlib/cs3060701",                     // 題目指定的映像檔
		"cmake", "-G", "Ninja", "-B", "build", // 執行的指令
	)

	// 將標準輸出與錯誤輸出都導向到 configure.log 檔案中
	cmdConfig.Stdout = configLog
	cmdConfig.Stderr = configLog

	fmt.Println("[1/2] 正在執行 CMake Configure...")
	if err := cmdConfig.Run(); err != nil {
		fmt.Println("Configure 失敗:", err)
		return "SE" // 初始化錯誤 (Setup Error)
	}

	// ==========================================
	// 階段二：編譯 (Build)
	// ==========================================
	compileLogPath := filepath.Join(workspace, "compile.log")
	compileLog, _ := os.Create(compileLogPath)
	defer compileLog.Close()

	cmdBuild := exec.Command("docker", "run", "--rm",
		"-v", absWorkspace+":/workspace",
		"-w", "/workspace",
		"yhlib/cs3060701",
		"cmake", "--build", "build", "--verbose",
	)

	cmdBuild.Stdout = compileLog
	cmdBuild.Stderr = compileLog

	fmt.Println("[2/2] 正在執行 CMake Build...")
	if err := cmdBuild.Run(); err != nil {
		fmt.Println("Build 失敗:", err)
		return "CE" // 編譯錯誤 (Compilation Error)
	}

	fmt.Println("編譯成功！可執行檔已準備就緒。")
	return "Ready"
}
