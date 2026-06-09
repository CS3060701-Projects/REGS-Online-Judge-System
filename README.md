# REGS - Remote Evaluation and Grading System

**REGS** 是一個線上程式評測系統（Online Judge）後端，專為跨平台使用、自動化處理與環境隔離而設計。本系統結合 **Docker 容器技術** 與 **CMake/Ninja 編譯工具鏈**，並以 JWT 進行權限控管。

---

## 核心功能

- **自動化編譯管線**：支援 `.zip` 格式上傳，自動執行 CMake 配置與 Ninja 編譯
- **環境隔離與安全**：編譯容器使用 `regs-build-net`，執行容器使用 `--network none` 完全斷網
- **非同步評測隊列**：接收提交後立即回傳 `operatorId`，背景透過 Job Queue 與 Semaphore 控制併發
- **多維度狀態判定**：AC、WA、CE、RE、SE、TLE
- **分層權限管理 (RBAC)**：Admin、User、Guest 三種角色，ECDSA JWT 認證
- **分階段日誌查詢**：`configure.log`、`compile.log`、`output.log`
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

### 啟動步驟

1. **啟動資料庫**

```bash
docker compose up -d
```

> PostgreSQL 預設映射至 `localhost:5433`。

2. **建立評測用 Docker 映像**

```bash
docker build -t regs-judger -f dockerfile .
```

3. **建立管理員帳號**

```bat
scripts\register_admin.bat
```

4. **啟動後端伺服器**

```bash
go run ./cmd/server
```

或使用 `Server.bat` / `start_all.bat`。

5. **訪問服務**

- API：`http://localhost:8081`
- Swagger：`http://localhost:8081/swagger/index.html`

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
| GET | /api/submissions/{operatorId}/logs/{type} | User |
| POST | /api/submissions/{operatorId}/rerun | User |
| GET | /api/stats/problems/{problem_id} | Guest |
| GET | /api/stats/users/{user_id} | Guest |

> 需要授權的請求請帶入 `Authorization: Bearer <token>`。

---

## 專案結構

```
regs-backend/
├── cmd/                # 應用程式進入點
│   ├── server/         # Web 伺服器主程式
│   └── seed/           # 建立管理員的獨立工具
├── docs/               # Swagger、ERD、操作說明
├── internal/           # 內部套件
│   ├── api/            # handlers、middleware
│   ├── database/       # 資料庫連線
│   ├── judge/          # 評測沙盒邏輯
│   └── models/         # GORM 資料模型
├── pkg/                # 公共套件（jwt、utils）
├── scripts/            # 輔助腳本
├── storage/            # 提交與工作目錄（執行期生成）
├── testdata/           # 題目測資
├── docker-compose.yml
├── dockerfile
└── go.mod
```

## 文件

- `docs/swagger.yaml`：OpenAPI 3.0 API 文件
- `docs/erd.md`：資料庫 ERD
- `docs/OPERATION.md`：專案操作說明

## 輔助腳本

- `scripts/register_admin.bat`：建立管理員帳號
- `scripts/reset_database.bat`：**(危險操作)** 重設資料庫

---

## 授權

本專案為 NTUST CS3060701 課程期末專案，請遵守學術誠信規範。
