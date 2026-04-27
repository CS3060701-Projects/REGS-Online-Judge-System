# REGS - A Modern Online Judge System

**REGS** 是一個現代化的線上評測系統（Online Judge），專為跨平台使用、自動化處理與環境隔離而設計。本系統結合了 **Docker 容器技術** 與 **CMake/Ninja 編譯工具鏈**，確保使用者提交的程式碼在完全隔離的沙盒環境中運行，並提供即時的評測結果與多階段日誌查詢。

---

## 核心功能

*   **自動化編譯管線**：支援 `.zip` 格式上傳，自動執行 CMake 配置與 Ninja 編譯。
*   **環境隔離與安全**：所有編譯與執行階段均在 Docker 容器內完成，並使用 `--network none` 實施完全斷網隔離。
*   **非同步評測隊列**：系統接收提交後立即回傳 `operatorId`，評測邏輯在背景透過 Job Queue 執行，確保 API 高響應性。
*   **多維度狀態判定**： **AC** (Accepted), **WA** (Wrong Answer), **CE** (Compile Error), **RE** (Runtime Error), **SE** (Setup Error) 與 **TLE** (Time Limit Exceeded)。
*   **分層權限管理 (RBAC)**：定義了 `Admin`, `User` 與 `Guest` 三種權限角色，並以 ECDSA 簽署的 JWT 進行身份驗證。
*   **完整的 API 文件**：透過 Swagger UI 提供互動式的 API 文件。

---

## 技術使用

*   **後端**: Go, Gin
*   **資料庫**: PostgreSQL, GORM
*   **容器化**: Docker
*   **編譯工具**: CMake, Ninja
*   **API 文件**: Swagger

---

## 快速開始

### 1. 環境準備

*   安裝 Go (建議版本 1.22 或以上)。
*   安裝 Docker。

### 2. 啟動步驟

1.  **啟動資料庫**
    在專案根目錄執行以下指令，背景啟動 PostgreSQL 服務。
    ```bash
    docker-compose up -d
    ```
2.  **啟動後端伺服器**
    執行根目錄的 `Server.bat` 來編譯並啟動後端服務。

3.  **訪問服務**
    *   **API 服務**: `http://localhost:8081`
    *   **API 文件**: `http://localhost:8081/swagger/index.html`

---

## API 文件與測試

### 查看 API 文件

在伺服器啟動後，您可以透過瀏覽器訪問以下網址來查看完整的互動式 API 文件：

*   [http://localhost:8081/swagger/index.html](http://localhost:8081/swagger/index.html)

### 測試 API

1.  在 API 文件頁面中，點開您想測試的任何一個 API 端點。
2.  點擊右上角的 **"Try it out"** 按鈕。
3.  填寫必要的參數（例如，請求內文、路徑參數）。
4.  對於需要授權的 API，請先透過 `/api/users/login` 取得 `token`，然後點擊頁面右上角的 **"Authorize"** 按鈕，在彈出的視窗中輸入 `Bearer <您的token>`。
5.  點擊 **"Execute"** 按鈕即可發送請求並查看回應。

---

## 輔助腳本

專案的script目錄中包含了方便開發的工具：

*   `create_admin.bat`: 建立管理員帳號。
*   `reset_database.bat`: **(危險操作)** 徹底清空資料庫，用於開發和測試。執行前會要求確認。

---

## 專案結構

```
regs-backend/
├── cmd/                # 應用程式進入點
│   ├── server/         # Web 伺服器主程式
│   └── seed/           # 建立管理員的獨立工具
├── docs/               # Swagger API 文件
├── internal/           # 內部套件 (不應被外部引用)
│   ├── api/            # API 處理邏輯 (handlers, middleware)
│   ├── database/       # 資料庫連線與遷移
│   ├── judge/          # 核心評測沙盒邏輯
│   └── models/         # GORM 資料模型
├── pkg/                # 可被外部引用的公共套件
│   ├── jwt/            # JWT 產生與驗證
│   └── utils/          # 通用工具函式 (如解壓縮)
├── storage/            # 執行期間生成的檔案 (submissions, workspaces)
├── testdata/           # 題目測資
├── Server.bat          # 啟動伺服器腳本
├── docker-compose.yml  # Docker 服務編排
└── go.mod              # Go 模組定義
```

---

## 授權

本專案為 NTUST CS3060701 課程期末專案，相關內容遵循學術誠信規範。