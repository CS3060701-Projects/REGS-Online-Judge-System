# REGS 專案操作說明

## 環境準備

1. 安裝 Go、Docker Desktop
2. 在專案根目錄啟動 PostgreSQL：`docker compose up -d`
3. 建立評測映像：`docker build -t regs-judger -f dockerfile .`
4. 確認 `private.pem` / `public.pem` 存在（或使用 openssl 重新生成）
5. 建立管理員：`scripts\register_admin.bat`
6. 啟動伺服器：`go run ./cmd/server`

## 評測容器說明

- **編譯容器**：掛載學生上傳目錄與題目目錄，使用 `regs-build-net` 網路
- **執行容器**：使用 `--network none` 完全斷網，防止網路攻擊
- 映像需包含：CMake、Ninja、build-essential

## 提交評測流程

1. 使用者以 JWT 登入後，POST `/api/submissions` 上傳 `.zip` 與 `problem_id`
2. 伺服器解壓縮、檢查 `CMakeLists.txt`、替換官方 `entrypoint.cpp`
3. 立即回傳 `operatorId`，背景 Worker 透過 Semaphore 控制併發
4. 以 GET `/api/submissions/{operatorId}` 輪詢狀態
5. 以 GET `/api/submissions/{operatorId}/logs/{type}` 查詢日誌（`configure` / `compile` / `output`）
6. 可 POST `/api/submissions/{operatorId}/rerun` 重新執行評測

## 題目管理（Admin）

1. PUT `/api/problems` 建立或更新題目
2. POST `/api/problems/{id}/testdata` 上傳測資 `.zip`
3. GET `/api/problems/{id}/testcases` 下載測資
4. DELETE `/api/problems/{id}` 刪除題目

題目測資預設放在 `testdata/{problem_id}/`，需包含根目錄 `CMakeLists.txt` 與 `solution/`、`spec/` 等結構。

## 測試資料

`backend/testdata/` 內含 `113final001` ~ `113final006` 範例題目，可作為本地評測測試。

## 常見問題

- **Configure 失敗 (SE)**：檢查學生上傳的 CMake 專案結構，或查看 `configure.log`
- **編譯失敗 (CE)**：查看 `compile.log` 中的語法錯誤
- **Docker 映像不存在**：執行 `docker build -t regs-judger -f dockerfile .`
- **資料庫連線失敗**：確認 `docker compose up -d` 且埠號為 `5433`
