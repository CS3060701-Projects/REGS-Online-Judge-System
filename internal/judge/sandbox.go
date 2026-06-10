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
	"regs-backend/internal/problem"
	"regs-backend/pkg/utils"
	"strings"
	"syscall"
	"time"
)

// RunAndJudge 是評測的主入口。
// 若題目存在 settings.yaml，使用逐 preset 評測流程；否則 fallback 至舊 ctest 流程。
func RunAndJudge(operatorID string, workspace string, prob models.Problem) models.JudgeResult {
	problemRoot := prob.TestcasePath
	if problemRoot == "" {
		problemRoot = filepath.Join("testdata", prob.ID)
	}

	cmakePath := filepath.Join(problemRoot, "CMakeLists.txt")
	if _, err := os.Stat(cmakePath); err != nil {
		fmt.Printf("[%s] 題目資料夾缺少 CMakeLists.txt: %s\n", operatorID, cmakePath)
		return models.JudgeResult{Status: "SE"}
	}

	// 嘗試載入 settings.yaml
	settings, err := problem.LoadSettings(problemRoot)
	if err != nil || len(settings.Presets) == 0 {
		fmt.Printf("[%s] 未找到 settings.yaml 或 presets 為空，使用舊版 ctest 流程\n", operatorID)
		return RunAndJudgeCTest(operatorID, workspace, prob)
	}

	fmt.Printf("[%s] 載入 settings.yaml 成功，共 %d 個 preset\n", operatorID, len(settings.Presets))
	return RunAndJudgeWithSettings(operatorID, workspace, prob, settings)
}

// RunAndJudgeWithSettings 根據 settings.yaml 逐 preset 進行評測
func RunAndJudgeWithSettings(operatorID string, workspace string, prob models.Problem, settings *problem.Settings) models.JudgeResult {
	problemRoot := prob.TestcasePath
	if problemRoot == "" {
		problemRoot = filepath.Join("testdata", prob.ID)
	}

	absProblemRoot, err := filepath.Abs(problemRoot)
	if err != nil {
		fmt.Printf("[%s] 無法取得題目根目錄絕對路徑: %v\n", operatorID, err)
		return models.JudgeResult{Status: "SE"}
	}

	totalScore := settings.TotalScore()
	earnedScore := 0
	var presetResults []models.PresetResult
	overallStatus := "AC"
	var maxTime float64

	outputLogPath := filepath.Join(workspace, "output.log")
	os.Remove(outputLogPath)
	logFile, err := os.OpenFile(outputLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("[%s] 無法建立 output.log: %v\n", operatorID, err)
		return models.JudgeResult{Status: "SE"}
	}
	defer logFile.Close()

	for i, preset := range settings.Presets {
		fmt.Printf("[%s] === Preset %d (score=%d) ===\n", operatorID, i+1, preset.Score)
		logFile.WriteString(fmt.Sprintf("\n=== Preset %d (score=%d) ===\n", i+1, preset.Score))

		result := runSinglePreset(operatorID, workspace, absProblemRoot, prob, settings, preset, i, logFile)
		presetResults = append(presetResults, result)

		if result.Status == "AC" {
			earnedScore += result.Earned
		}

		if result.PeakTime > maxTime {
			maxTime = result.PeakTime
		}

		// 更新整體狀態：如果有任一 preset 不是 AC，整體狀態為非 AC
		if result.Status != "AC" && overallStatus == "AC" {
			overallStatus = "WA" // 預設改為 WA，若有更嚴重的狀態再覆蓋
		}

		fmt.Printf("[%s] Preset %d 結果: %s (得分: %d/%d)\n", operatorID, i+1, result.Status, result.Earned, result.Score)
	}

	// 如果全部 AC，整體就是 AC；如果部分通過，用 Partial AC 或 WA 表示
	if earnedScore == totalScore {
		overallStatus = "AC"
	} else if earnedScore > 0 {
		overallStatus = "WA" // 部分通過也標記 WA
	}

	fmt.Printf("[%s] 評測完成: 狀態=%s, 得分=%d/%d\n", operatorID, overallStatus, earnedScore, totalScore)

	// 將各 preset 的 configure.log 與 compile.log 合併至 workspace 根目錄
	// 確保 GET /submissions/{id}/logs/configure 與 /compile 在 settings.yaml 流程下也能正常回傳
	mergePresetLogs(workspace, len(settings.Presets))

	return models.JudgeResult{
		Status:        overallStatus,
		PeakTime:      maxTime,
		PeakMemory:    0,
		TotalScore:    totalScore,
		EarnedScore:   earnedScore,
		PresetResults: presetResults,
	}
}

// mergePresetLogs 將各 preset_N/ 子目錄的 configure.log 與 compile.log
// 合併後寫入 workspace 根目錄，確保 log API 在 settings.yaml 流程下能正常查詢。
// 同時也建立 config.log（configure.log 的別名）以符合規格書要求。
func mergePresetLogs(workspace string, presetCount int) {
	type logEntry struct {
		filename string
		dest     string
		alias    string // 若不為空則同時複製一份別名
	}
	entries := []logEntry{
		{filename: "configure.log", dest: "configure.log", alias: "config.log"},
		{filename: "compile.log", dest: "compile.log"},
	}

	for _, entry := range entries {
		destPath := filepath.Join(workspace, entry.dest)
		f, err := os.Create(destPath)
		if err != nil {
			fmt.Printf("[mergePresetLogs] 無法建立 %s: %v\n", entry.dest, err)
			continue
		}

		for i := 0; i < presetCount; i++ {
			srcPath := filepath.Join(workspace, fmt.Sprintf("preset_%d", i), entry.filename)
			data, err := os.ReadFile(srcPath)
			if err != nil {
				// 該 preset 可能在 configure 階段就失敗，log 可能不存在，略過
				continue
			}
			f.WriteString(fmt.Sprintf("=== Preset %d ===\n", i+1))
			f.Write(data)
			f.WriteString("\n")
		}
		f.Close()

		// 建立別名（config.log → configure.log）
		if entry.alias != "" {
			aliasPath := filepath.Join(workspace, entry.alias)
			if data, err := os.ReadFile(destPath); err == nil {
				_ = os.WriteFile(aliasPath, data, 0644)
			}
		}
	}
}

// runSinglePreset 執行單一 preset 的完整流程：prepare → configure → build → run → compare
func runSinglePreset(operatorID string, workspace string, absProblemRoot string, prob models.Problem, settings *problem.Settings, preset problem.Preset, index int, logFile *os.File) models.PresetResult {
	result := models.PresetResult{
		Index: index,
		Score: preset.Score,
	}

	// 1. 建立此 preset 的臨時工作目錄
	presetDir := filepath.Join(workspace, fmt.Sprintf("preset_%d", index))
	os.RemoveAll(presetDir) // 清理舊的
	if err := os.MkdirAll(presetDir, os.ModePerm); err != nil {
		fmt.Printf("[%s][preset %d] 無法建立 preset 工作目錄: %v\n", operatorID, index, err)
		result.Status = "SE"
		return result
	}

	// 2. 複製原始上傳檔案到 preset 目錄
	if err := copyUploadFiles(workspace, presetDir); err != nil {
		fmt.Printf("[%s][preset %d] 複製上傳檔案失敗: %v\n", operatorID, index, err)
		result.Status = "SE"
		return result
	}

	// 3. 執行 replace 規則：用題目目錄中的檔案覆蓋到 preset 工作目錄
	problemRoot := prob.TestcasePath
	if problemRoot == "" {
		problemRoot = filepath.Join("testdata", prob.ID)
	}

	for _, r := range preset.Replace {
		srcPath := r.Source
		if !filepath.IsAbs(srcPath) {
			srcPath = filepath.Join(problemRoot, srcPath)
		}
		dstPath := filepath.Join(presetDir, r.Target)

		if _, err := os.Stat(srcPath); err != nil {
			fmt.Printf("[%s][preset %d] replace 來源不存在: %s\n", operatorID, index, srcPath)
			result.Status = "SE"
			return result
		}

		if err := utils.CopyFile(srcPath, dstPath); err != nil {
			fmt.Printf("[%s][preset %d] replace 複製失敗 %s -> %s: %v\n", operatorID, index, srcPath, dstPath, err)
			result.Status = "SE"
			return result
		}
		fmt.Printf("[%s][preset %d] Replace: %s -> %s\n", operatorID, index, r.Source, r.Target)
	}

	// 4. Configure
	absPresetDir, err := filepath.Abs(presetDir)
	if err != nil {
		result.Status = "SE"
		return result
	}
	absPresetDir = filepath.ToSlash(absPresetDir)
	absProblemRootSlash := filepath.ToSlash(absProblemRoot)

	configLogPath := filepath.Join(presetDir, "configure.log")
	if err := RunConfigure(absPresetDir, absProblemRootSlash, configLogPath); err != nil {
		fmt.Printf("[%s][preset %d] Configure 失敗: %v\n", operatorID, index, err)
		logFile.WriteString(fmt.Sprintf("Preset %d: Configure FAILED\n", index+1))
		result.Status = "SE"
		return result
	}

	// 5. Build
	compileLogPath := filepath.Join(presetDir, "compile.log")
	if err := RunBuild(absPresetDir, absProblemRootSlash, compileLogPath); err != nil {
		fmt.Printf("[%s][preset %d] 編譯失敗: %v\n", operatorID, index, err)
		logFile.WriteString(fmt.Sprintf("Preset %d: Compile FAILED\n", index+1))
		result.Status = "CE"
		return result
	}

	// 6. 執行程式並捕獲輸出
	// 使用 settings.yaml 的 limits.totalTime（毫秒）
	timeoutMs := settings.Limits.TotalTime
	if timeoutMs <= 0 {
		timeoutMs = prob.TimeLimit
	}
	timeout := time.Duration(timeoutMs*10) * time.Millisecond
	if timeout < 10*time.Second {
		timeout = 10 * time.Second
	}

	// 使用 settings.yaml 的 limits.memory（KB → MB）
	memLimit := problem.MemoryLimitMB(settings.Limits.Memory)
	if settings.Limits.Memory <= 0 {
		memLimit = prob.MemoryLimit
		if memLimit <= 0 {
			memLimit = 256
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 使用 ctest 執行此 preset 的測試（build 目錄中應只有一個 test case）
	cmdRun := exec.CommandContext(ctx, "docker", "run", "--rm",
		"--network", "none",
		"--cpus", "1.0",
		"--memory", fmt.Sprintf("%dm", memLimit),
		"-v", absPresetDir+":/upload",
		"-v", absProblemRootSlash+":/problem:ro",
		"-v", absPresetDir+":/app",
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
	result.PeakTime = elapsed

	outputStr := stdoutBuf.String()
	logFile.WriteString(fmt.Sprintf("Preset %d output:\n%s\n", index+1, outputStr))

	// 7. 檢查超時
	if ctx.Err() == context.DeadlineExceeded {
		result.Status = "TLE"
		return result
	}

	// 8. 檢查 Runtime Error
	if runErr != nil {
		if strings.Contains(outputStr, "Timeout") ||
			strings.Contains(outputStr, "TIMEOUT") ||
			strings.Contains(outputStr, "Test timeout") {
			result.Status = "TLE"
			return result
		}

		if isRuntimeError(outputStr, runErr) {
			result.Status = "RE"
			return result
		}
	}

	// 9. 比對輸出與 expected
	problemRootForExpected := prob.TestcasePath
	if problemRootForExpected == "" {
		problemRootForExpected = filepath.Join("testdata", prob.ID)
	}

	expectedContent, err := problem.LoadExpectedContent(preset.Expected, problemRootForExpected)
	if err != nil {
		fmt.Printf("[%s][preset %d] 載入預期輸出失敗: %v\n", operatorID, index, err)
		result.Status = "SE"
		return result
	}

	// 從 ctest 輸出中提取程式的 stdout (system-out 部分)
	// ctest -V 的輸出格式中，程式輸出會直接包含在內
	actualOutput := extractProgramOutput(outputStr)

	if compareOutput(actualOutput, expectedContent) {
		result.Status = "AC"
		result.Earned = preset.Score
	} else {
		fmt.Printf("[%s][preset %d] 輸出不匹配\n", operatorID, index)
		fmt.Printf("[%s][preset %d] 預期:\n%s\n", operatorID, index, expectedContent)
		fmt.Printf("[%s][preset %d] 實際:\n%s\n", operatorID, index, actualOutput)
		result.Status = "WA"

		// 如果 ctest 本身成功 (exit code 0) 但輸出不匹配，仍然是 WA
		// 如果 ctest 失敗 (exit code != 0)，可能是 RE
		if runErr != nil {
			result.Status = "WA"
		}
	}

	return result
}

// copyUploadFiles 將上傳的原始檔案複製到 preset 工作目錄
// 排除 preset_* 目錄、build 目錄和 log 檔案
func copyUploadFiles(srcWorkspace, dstDir string) error {
	return filepath.WalkDir(srcWorkspace, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcWorkspace, path)
		if err != nil {
			return err
		}

		// 跳過 preset_* 目錄、build 目錄、log 檔案
		if d.IsDir() {
			name := d.Name()
			if strings.HasPrefix(name, "preset_") || name == "build" {
				return filepath.SkipDir
			}
		}

		// 跳過 log 檔案
		if !d.IsDir() {
			name := d.Name()
			if strings.HasSuffix(name, ".log") {
				return nil
			}
		}

		target := filepath.Join(dstDir, rel)
		if d.IsDir() {
			return os.MkdirAll(target, os.ModePerm)
		}

		return utils.CopyFile(path, target)
	})
}

// extractProgramOutput 從 ctest -V 的輸出中提取第一個測試的實際輸出
// ctest -V 格式中，每行程式輸出都帶有測試編號前綴，例如 "1: Valid: xxx"
// 此函式只提取第一個測試的輸出，並去除 "N: " 前綴
func extractProgramOutput(ctestOutput string) string {
	lines := strings.Split(ctestOutput, "\n")
	var programLines []string
	inOutput := false
	firstTestNum := ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// 偵測程式輸出開始（在 "Test timeout computed" 之後）
		if strings.Contains(trimmed, "Test timeout computed to be:") {
			if firstTestNum == "" {
				inOutput = true
			}
			continue
		}

		// 偵測程式輸出結束（ctest 的結果行或下一個測試開始）
		if inOutput && (strings.Contains(trimmed, "Test #") ||
			strings.Contains(trimmed, "% tests passed") ||
			strings.Contains(trimmed, "Total Test time") ||
			strings.HasPrefix(trimmed, "Start ") ||
			strings.HasPrefix(trimmed, "UpdateCTestConfiguration") ||
			strings.HasPrefix(trimmed, "Test project")) {
			inOutput = false
			continue
		}

		if inOutput && trimmed != "" {
			// 去除 ctest -V 的行號前綴 "N: "（例如 "1: Valid: xxx" → "Valid: xxx"）
			cleaned := stripCtestPrefix(trimmed)

			// 確定第一個測試的編號
			if firstTestNum == "" {
				prefix := extractCtestPrefix(trimmed)
				if prefix != "" {
					firstTestNum = prefix
				}
			}

			// 只保留第一個測試的輸出
			if firstTestNum != "" {
				prefix := extractCtestPrefix(trimmed)
				if prefix != "" && prefix != firstTestNum {
					// 這是其他測試的輸出，跳過
					continue
				}
			}

			programLines = append(programLines, cleaned)
		}
	}

	return strings.Join(programLines, "\n")
}

// stripCtestPrefix 去除 ctest -V 的 "N: " 前綴
// 例如 "1: Valid: xxx" → "Valid: xxx"
// 如果沒有前綴則原樣返回
func stripCtestPrefix(line string) string {
	// ctest 前綴格式: "數字: "
	for i, ch := range line {
		if ch == ':' && i > 0 {
			// 檢查冒號前面是否全部是數字
			allDigits := true
			for _, d := range line[:i] {
				if d < '0' || d > '9' {
					allDigits = false
					break
				}
			}
			if allDigits && i+1 < len(line) && line[i+1] == ' ' {
				return line[i+2:]
			}
			break
		}
		if ch < '0' || ch > '9' {
			break
		}
	}
	return line
}

// extractCtestPrefix 提取 ctest -V 的測試編號前綴
// 例如 "1: Valid: xxx" → "1"，無前綴返回 ""
func extractCtestPrefix(line string) string {
	for i, ch := range line {
		if ch == ':' && i > 0 {
			allDigits := true
			for _, d := range line[:i] {
				if d < '0' || d > '9' {
					allDigits = false
					break
				}
			}
			if allDigits && i+1 < len(line) && line[i+1] == ' ' {
				return line[:i]
			}
			break
		}
		if ch < '0' || ch > '9' {
			break
		}
	}
	return ""
}

// compareOutput 逐行比對兩個輸出（忽略尾部空白和 \r）
func compareOutput(actual, expected string) bool {
	// 正規化：統一換行符號，去除尾部空白
	actual = strings.ReplaceAll(actual, "\r\n", "\n")
	expected = strings.ReplaceAll(expected, "\r\n", "\n")
	actual = strings.TrimRight(actual, "\n ")
	expected = strings.TrimRight(expected, "\n ")

	actualLines := strings.Split(actual, "\n")
	expectedLines := strings.Split(expected, "\n")

	if len(actualLines) != len(expectedLines) {
		return false
	}

	for i := range actualLines {
		a := strings.TrimRight(actualLines[i], " \t\r")
		e := strings.TrimRight(expectedLines[i], " \t\r")
		if a != e {
			return false
		}
	}

	return true
}

// ====== 以下為舊版 ctest 流程（向後相容）======

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
	
	// 根據期末專案規格：Exit Code 非 0 視為 RE
	// CTest 若捕捉到非零退出，會印出 "Exit code X" 或 "Exception: ..."
	if strings.Contains(lower, "segmentation fault") ||
		strings.Contains(lower, "core dumped") ||
		strings.Contains(lower, "bus error") ||
		strings.Contains(lower, "access violation") ||
		strings.Contains(lower, "abort") ||
		strings.Contains(lower, "sigabrt") ||
		strings.Contains(lower, "sigsegv") ||
		strings.Contains(lower, "exception:") ||
		(strings.Contains(lower, "exit code") && !strings.Contains(lower, "exit code 0")) {
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
