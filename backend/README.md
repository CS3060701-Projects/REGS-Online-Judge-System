# REGS Backend

`backend` 是 REGS 的後端服務，使用 Go + Gin 實作 REST API、JWT 認證、PostgreSQL 資料庫與 Docker 隔離評測流程。

## 主要功能

- 使用者註冊、登入、登出、JWT 驗證
- 題目管理與題目查詢
- 提交 `.zip` 程式碼並進行非同步評測
- 查詢提交結果、原始專案與分階段日誌
- 管理員題目與測資上傳、刪除
- Swagger API 文件

## 環境需求

- Go
- Docker
- `docker compose`
- `openssl` (若重新生成 JWT 金鑰時使用)

## 快速啟動

### 1. 啟動資料庫

在 `backend` 目錄中執行：

```bat
docker compose up -d
```

此命令會啟動 PostgreSQL 容器，並將內部 `5432` 映射到 `localhost:5433`。

### 2. 檢查或生成 JWT 金鑰

```bash
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out private.pem
openssl pkey -in private.pem -pubout -out public.pem
```

### 3. 建立管理員帳號

在 `backend` 目錄中執行：

```bat
scripts\register_admin.bat
```

### 4. 啟動後端伺服器

執行：

```bat
start_backend.bat
```

此腳本會：

- 檢查 Docker 容器是否已啟動
- 編譯 `cmd/server` 並產生 `server.exe`
- 啟動後端伺服器

若你希望直接執行，也可以：

```bat
go run ./cmd/server
```

### 5. 服務位置

- API: `http://localhost:8081`
- Swagger 文件: `http://localhost:8081/swagger/index.html`

## 重要 API 路由

### 公開路由

- `POST /api/users/register`
- `POST /api/users/login`
- `GET /api/problems`
- `GET /api/problems/:id`
- `GET /api/problems/:id/examples`
- `GET /api/users/:user_id/submissions`
- `GET /api/stats/problems/:problem_id`
- `GET /api/stats/users/:user_id`

### 需要驗證的路由

- `POST /api/users/logout`
- `GET /api/users/me`
- `POST /api/submissions`
- `GET /api/submissions`
- `GET /api/submissions/:operatorId`
- `GET /api/submissions/:operatorId/source`
- `GET /api/submissions/:operatorId/logs/:type`

### 管理員專用路由

- `PUT /api/problems`
- `GET /api/problems/:id/testcases`
- `POST /api/problems/:id/testdata`
- `DELETE /api/problems/:id`

> 後端採用 `Authorization: Bearer <token>`，並以 JWT 進行權限驗證。

## 輔助腳本

`backend/scripts` 目錄包含：

- `register_admin.bat`：建立管理員帳號
- `reset_database.bat`：**(危險操作)** 重設 PostgreSQL 資料庫


## 其他注意事項

- `backend/cmd/server/main.go` 會初始化 3 個 judge worker：`handlers.InitJudger(3)`。
- CORS 允許來源為 `http://localhost:5173`，對應前端預設位置。

## 授權

本專案為 NTUST CS3060701 課程期末專案，請遵守學術誠信規範。
