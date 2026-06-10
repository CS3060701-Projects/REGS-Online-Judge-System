# REGS - Remote Evaluation and Grading System

**REGS** 是一個線上程式評測系統（Online Judge）後端，專為跨平台使用、自動化處理與環境隔離而設計。本系統結合 **Docker 容器技術** 與 **CMake/Ninja 編譯工具鏈**，並以 JWT 進行權限控管。

---

## 核心功能

- **自動化編譯管線**：支援 `.zip` 格式上傳，自動執行 CMake 配置與 Ninja 編譯
- **環境隔離與安全**：編譯容器使用 `regs-build-net`，執行容器使用 `--network none` 完全斷網
- **非同步評測隊列**：接收提交後立即回傳 `operatorId`，背景透過 Job Queue 與 Semaphore 控制併發
- **多維度狀態判定**：AC、WA、CE、RE、SE、TLE
- **分層權限管理 (RBAC)**：Admin、User、Guest 三種角色，ECDSA JWT 認證
- **Swagger API 文件**：互動式 API 測試介面

## 技術使用

- 後端：Go + Gin
- 資料庫：PostgreSQL + GORM
- 編譯工具：CMake + Ninja（Docker 映像 `regs-judger`）
- 認證方式：JWT (ECDSA P-256)
- 容器化：Docker

---

## 快速開始

### 環境需求

- Go（建議 1.22+）
- Docker

### 啟動與管理 (跨平台支援)

本專案提供跨平台統一腳本，無論是 **Windows** 或 **Ubuntu** 皆可直接執行，不需設定任何執行權限。

- **啟動伺服器** (自動開啟 DB 並編譯)
  ```bash
  go run cmd/task/main.go server
  ```
- **建立管理員帳號**
  ```bash
  go run cmd/task/main.go seed-admin
  ```
- **重置資料庫** (刪除所有資料)
  ```bash
  go run cmd/task/main.go reset-db
  ```

---

### 手動啟動步驟 (進階)

若不想使用自動腳本，手動執行步驟如下：

1. **啟動資料庫**

```bash
docker compose up -d
```

> PostgreSQL 預設映射至 `localhost:5433`。

2. **建立評測用 Docker 映像**

```bash
docker build -t regs-judger -f dockerfile .
```

3. 確認 `private.pem` / `public.pem` 存在（或使用 openssl 重新生成）
```bash
openssl genpkey -algorithm EC -pkeyopt ec_paramgen_curve:P-256 -out private.pem
openssl pkey -in private.pem -pubout -out public.pem
```



## 提交評測流程

1. 使用者以 JWT 登入後，POST `/api/submissions` 上傳 `.zip` 與 `problem_id`
2. 伺服器解壓縮、檢查 `CMakeLists.txt`、替換官方 `entrypoint.cpp`
3. 立即回傳 `operatorId`，背景 Worker 透過 Semaphore 控制併發
4. 以 GET `/api/submissions/{operatorId}` 輪詢狀態

## 題目管理（Admin）

1. PUT `/api/problems` 建立或更新題目（含一次性上傳測資 `.zip`，並自動讀取 `settings.yaml` 覆蓋資源限制）
2. GET `/api/problems/{id}/testcases` 下載測資
3. DELETE `/api/problems/{id}` 刪除題目

題目測資預設放在 `testdata/{problem_id}/`，需包含根目錄 `CMakeLists.txt` 與 `solution/`、`spec/` 等結構。

5. **API文件**

- API：`http://localhost:8081`
- OpenAPI 3.0：`http://localhost:8081/openapi.yaml`
- Redoc：`http://localhost:8081/docs`
- Swagger UI：`http://localhost:8081/swagger/index.html`

---

## 核心 API 概覽

| 方法 | URL | 權限 |
|------|-----|------|
| POST | /api/users/register | Guest |
| POST | /api/users/login | Guest |
| POST | /api/users/logout | User |
| GET | /api/users/me | User |
| GET | /api/users/{user_id}/submissions | Guest |
| GET | /api/problems | Guest |
| GET | /api/problems/{id} | Guest |
| PUT | /api/problems | Admin |
| DELETE | /api/problems/{id} | Admin |
| GET | /api/problems/{id}/testcases | Admin |
| POST | /api/submissions | User |
| GET | /api/submissions | User |
| GET | /api/submissions/{operatorId} | User |
| GET | /api/submissions/{operatorId}/source | User |
| GET | /api/stats/problems/{problem_id} | Guest |
| GET | /api/stats/users/{user_id} | Guest |

> 需要授權的請求請帶入 `Authorization: Bearer <token>`。

---

## 專案結構

```
regs-backend/
├── cmd/                # 應用程式進入點
│   ├── server/         # Web 伺服器主程式
│   ├── seed/           # 建立管理員的獨立工具
│   └── task/           # 跨平台統一工作腳本
├── docs/               # Swagger、ERD、操作說明
├── internal/           # 內部套件
│   ├── api/            # handlers、middleware
│   ├── database/       # 資料庫連線
│   ├── judge/          # 評測沙盒邏輯
│   └── models/         # GORM 資料模型
├── pkg/                # 公共套件（jwt、utils）
├── storage/            # 提交與工作目錄（執行期生成）
├── testdata/           # 題目測資
├── docker-compose.yml
├── dockerfile
└── go.mod
```

## 文件

- `docs/openapi.yaml`：OpenAPI 3.0 API 文件
- `docs/swagger.yaml`：Swagger 2.0（相容用）
- `ERD實體關聯圖.md`：資料庫 ERD 與 RBAC 說明
- `docs/OPERATION.md`：專案操作說明

## 跨平台任務腳本

統一透過 `go run cmd/task/main.go [command]` 來執行：
- `server`：檢查 Docker 後啟動 Web 伺服器
- `seed-admin`：建立預設管理員帳號
- `reset-db`：清空所有資料庫數據

---
