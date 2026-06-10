package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"regs-backend/internal/database"
	"regs-backend/internal/models"
)

const (
	baseURL      = "http://localhost:8081/api"
	frameworkDir = "114framework"
)

type E2ETester struct {
	adminToken string
	userToken  string
	client     *http.Client
}

func main() {
	database.Connect()
	tester := &E2ETester{
		client: &http.Client{Timeout: 30 * time.Second},
	}

	fmt.Println("=== 開始 REGS 端到端 (E2E) 測試 ===")

	tester.authTests()
	tester.uploadProblems()
	tester.runSubmissions()
	tester.queryTests()

	fmt.Println("=== E2E 測試全部通過！ ===")
}

// ---------------------------
// 1. 認證測試
// ---------------------------
func (t *E2ETester) authTests() {
	fmt.Println("[測試] 註冊與登入...")

	adminEmail := fmt.Sprintf("admin_%d@test.com", time.Now().Unix())
	userEmail := fmt.Sprintf("user_%d@test.com", time.Now().Unix())

	// 註冊 Admin 與 User
	t.register(adminEmail, "admin", "Admin")
	t.adminToken = t.login(adminEmail, "admin")

	t.register(userEmail, "user", "User")
	t.userToken = t.login(userEmail, "user")

	fmt.Println("  -> 登入成功")
}

func (t *E2ETester) register(email, password, role string) {
	data := map[string]string{"username": email, "password": password, "role": role}
	b, _ := json.Marshal(data)
	resp, err := t.client.Post(baseURL+"/users/register", "application/json", bytes.NewBuffer(b))
	if err != nil {
		panic("註冊請求失敗: " + err.Error())
	}
	defer resp.Body.Close()

	if role == "Admin" {
		// 強制將測試帳號升級為 Admin (避免因資料庫非空導致無法獲得權限)
		database.DB.Model(&models.User{}).Where("username = ?", email).Update("role", "Admin")
	}
}

func (t *E2ETester) login(email, password string) string {
	data := map[string]string{"username": email, "password": password}
	b, _ := json.Marshal(data)
	resp, err := t.client.Post(baseURL+"/users/login", "application/json", bytes.NewBuffer(b))
	if err != nil {
		panic("登入請求失敗: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("登入失敗: %s, 輸出: %s", resp.Status, string(body)))
	}

	var result struct {
		Token string `json:"token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Token
}

// ---------------------------
// 2. 題目上傳測試
// ---------------------------
func (t *E2ETester) uploadProblems() {
	fmt.Println("[測試] 上傳 114framework 題目...")

	for i := 1; i <= 6; i++ {
		probID := fmt.Sprintf("114FinalQ%03d", i)
		probDir := filepath.Join(frameworkDir, probID)

		// Check if directory exists
		if stat, err := os.Stat(probDir); err != nil || !stat.IsDir() {
			fmt.Printf("  -> 跳過 %s (目錄不存在)\n", probID)
			continue
		}

		fmt.Printf("  -> 打包並上傳題目 %s ...\n", probID)
		zipData := createProblemZip(probDir)
		t.uploadProblemAPI(probID, "題目 "+probID, zipData)
	}
}

func createProblemZip(probDir string) []byte {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	filepath.Walk(probDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		// 忽略 solution 目錄，測資上傳不需要包含 solution
		if strings.Contains(path, "solution") || strings.Contains(path, "node_modules") || strings.Contains(path, ".git") {
			return nil
		}
		rel, _ := filepath.Rel(probDir, path)
		f, _ := w.Create(rel)
		content, _ := os.ReadFile(path)
		f.Write(content)
		return nil
	})
	w.Close()
	return buf.Bytes()
}

func (t *E2ETester) uploadProblemAPI(id, title string, zipData []byte) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("id", id)
	writer.WriteField("title", title)
	writer.WriteField("description", "這是一道自動上傳的題目")
	writer.WriteField("is_visible", "true")

	part, _ := writer.CreateFormFile("file", "testcases.zip")
	part.Write(zipData)
	writer.Close()

	req, _ := http.NewRequest("PUT", baseURL+"/problems", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+t.adminToken)

	resp, err := t.client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("上傳題目 %s 失敗: %v, %s", id, err, string(b)))
	}
	resp.Body.Close()
}

// ---------------------------
// 3. 測資評測測試 (AC/WA/CE/RE/TLE)
// ---------------------------
func (t *E2ETester) runSubmissions() {
	fmt.Println("[測試] 執行測資評測 (AC, WA, CE, RE, TLE)...")
	states := []string{"AC", "WA", "CE", "RE", "TLE"}

	for i := 1; i <= 6; i++ {
		probID := fmt.Sprintf("114FinalQ%03d", i)
		solDir := filepath.Join(frameworkDir, probID, "solution")

		if stat, err := os.Stat(solDir); err != nil || !stat.IsDir() {
			continue
		}

		for _, state := range states {
			fmt.Printf("  -> 測試 %s 的 %s 狀態...\n", probID, state)
			zipData := createSubmissionZip(solDir, state)
			opID := t.submitAssignment(probID, zipData)
			t.pollAndAssert(opID, state)
		}
	}
}

func createSubmissionZip(solDir, state string) []byte {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	injected := false

	filepath.Walk(solDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		// 不要打包既有的 zip
		if strings.HasSuffix(info.Name(), ".zip") {
			return nil
		}

		rel, _ := filepath.Rel(solDir, path)
		content, _ := os.ReadFile(path)

		// 找一個學生的檔案來注入錯誤碼 (排除系統會蓋掉的 test.h 與 entrypoint.cpp)
		if !injected && (strings.HasSuffix(info.Name(), ".h") || strings.HasSuffix(info.Name(), ".cpp")) {
			if info.Name() != "test.h" && info.Name() != "entrypoint.cpp" && info.Name() != "CMakeLists.txt" && info.Name() != "case.h" {
				content = injectCode(content, state)
				injected = true
			}
		}

		f, _ := w.Create(rel)
		f.Write(content)
		return nil
	})
	w.Close()
	return buf.Bytes()
}

func injectCode(orig []byte, state string) []byte {
	var injection string
	switch state {
	case "AC":
		return orig
	case "WA":
		injection = "\n#ifndef E2E_INJECTED_H\n#define E2E_INJECTED_H\n#include <iostream>\nnamespace {\n\tstruct _Wa {\n\t\t_Wa() { std::cout << \"Wrong Answer Garbage Output\\n\"; }\n\t} _wa;\n}\n#endif\n"
	case "CE":
		injection = "\n#ifndef E2E_INJECTED_H\n#define E2E_INJECTED_H\n#error \"Intended Compile Error\"\n#endif\n"
	case "RE":
		injection = "\n#ifndef E2E_INJECTED_H\n#define E2E_INJECTED_H\n#include <stdlib.h>\nnamespace {\n\tstruct _Re {\n\t\t_Re() { abort(); }\n\t} _re;\n}\n#endif\n"
	case "TLE":
		injection = "\n#ifndef E2E_INJECTED_H\n#define E2E_INJECTED_H\nnamespace {\n\tstruct _Tle {\n\t\t_Tle() { while(true) {} }\n\t} _tle;\n}\n#endif\n"
	}
	return append(orig, []byte(injection)...)
}

func (t *E2ETester) submitAssignment(probID string, zipData []byte) string {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	writer.WriteField("problem_id", probID)

	part, _ := writer.CreateFormFile("file", "submission.zip")
	part.Write(zipData)
	writer.Close()

	req, _ := http.NewRequest("POST", baseURL+"/submissions", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+t.userToken)

	resp, err := t.client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("上傳代碼失敗: %v, %s", err, string(b)))
	}
	defer resp.Body.Close()

	var result struct {
		OperatorID string `json:"operatorId"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.OperatorID
}

func (t *E2ETester) pollAndAssert(operatorID, expectedState string) {
	for i := 0; i < 300; i++ { // 輪詢最多 300 次 (300秒)
		time.Sleep(1 * time.Second)

		req, _ := http.NewRequest("GET", baseURL+"/submissions/"+operatorID, nil)
		req.Header.Set("Authorization", "Bearer "+t.userToken)

		resp, err := t.client.Do(req)
		if err != nil {
			continue
		}

		var result struct {
			Status string `json:"status"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if result.Status != "Pending" && result.Status != "Judging" && result.Status != "Compiling" {
			if result.Status != expectedState {
				panic(fmt.Sprintf("斷言失敗！預期: %s, 實際得到: %s (操作ID: %s)", expectedState, result.Status, operatorID))
			}
			fmt.Printf("    => 驗證成功 (%s)\n", result.Status)
			return
		}
	}
	panic("評測超時！狀態一直卡在 Pending/Judging")
}

// ---------------------------
// 4. API 查詢功能測試
// ---------------------------
func (t *E2ETester) queryTests() {
	fmt.Println("[測試] 查詢功能 API...")

	// GET /problems
	resp, _ := http.Get(baseURL + "/problems")
	if resp.StatusCode != 200 {
		panic("取得題目列表失敗")
	}
	resp.Body.Close()

	fmt.Println("  -> GET /problems 成功")

	// GET /submissions (個人紀錄)
	req, _ := http.NewRequest("GET", baseURL+"/submissions", nil)
	req.Header.Set("Authorization", "Bearer "+t.userToken)
	resp, _ = t.client.Do(req)
	if resp.StatusCode != 200 {
		panic("取得個人提交紀錄失敗")
	}
	resp.Body.Close()

	fmt.Println("  -> GET /submissions 成功")
}
