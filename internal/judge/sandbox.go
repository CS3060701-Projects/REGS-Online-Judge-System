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

	userOutput, _ := os.ReadFile(outputLogPath)
	actualOutput := strings.TrimSpace(string(userOutput))

	expectedOutput := "Hello, REGS System!"

	if actualOutput == expectedOutput {
		return "AC"
	} else {
		return "WA"
	}
}

func CompileProject(operatorID string, workspace string) string {
	absWorkspace, err := filepath.Abs(workspace)
	if err != nil {
		fmt.Println("取得絕對路徑失敗:", err)
		return "System Error"
	}
	absWorkspace = filepath.ToSlash(absWorkspace)

	fmt.Println("開始編譯流程，Workspace:", absWorkspace)

	configLogPath := filepath.Join(workspace, "configure.log")
	configLog, _ := os.Create(configLogPath)
	defer configLog.Close()

	cmdConfig := exec.Command("docker", "run", "--rm",
		"-v", absWorkspace+":/workspace",
		"-w", "/workspace",
		"yhlib/cs3060701",
		"cmake", "-G", "Ninja", "-B", "build",
	)

	cmdConfig.Stdout = configLog
	cmdConfig.Stderr = configLog

	fmt.Println("[1/2] 正在執行 CMake Configure...")
	if err := cmdConfig.Run(); err != nil {
		fmt.Println("Configure 失敗:", err)
		return "SE"
	}

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
		return "CE"
	}

	fmt.Println("編譯成功！可執行檔已準備就緒。")
	return "Ready"
}
