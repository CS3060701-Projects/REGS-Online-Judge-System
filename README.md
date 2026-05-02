# Remote-Evaluation-and-Grading-System

REGS 是一套線上程式評測系統（Online Judge），包含：

- `backend`：Go + Gin + PostgreSQL + Docker sandbox 程式測試
- `frontend`：React + TypeScript + Vite 使用者介面

此專案目前為「前後端分離」架構：前端呼叫後端 API 進行認證、題目查詢、提交與評測結果查詢。

---

## 專案結構

### backend
- Go 後端服務
### frontend
- React 前端服務（含 Start-Frontend.bat）
### Start-All.bat
- 一鍵同時啟動 backend / frontend

---

## 功能概覽

### 前端

- 使用者註冊 / 登入 / 登出
- 題目列表與搜尋
- 題目詳情與統計資訊
- 題目範例測資顯示（最多兩筆 input/output）
- 提交 zip 程式碼
- 我的提交紀錄 / 單筆 operatorId 查詢
- 管理員題目管理（建立、更新、刪除、上傳測資）

### 後端

- JWT 驗證與 RBAC 權限控管（User / Admin）
- 非同步評測 Queue
- Docker 隔離編譯與執行
- 題目、提交、統計 API
- Swagger 文件

---

### 本機開發啟動

## 一鍵同時啟動（Windows）

在專案根目錄直接執行，會同時開啟兩個視窗分別啟動 backend / frontend。：

```bat
Start-All.bat
```

---

## 環境變數

前端可用 `frontend/.env` 設定 API 位置：

```env
VITE_API_BASE_URL=http://localhost:8081/api
```

---

## API 注意事項（目前前端有使用）

- `GET /api/problems/:id`：取得題目敘述
- `GET /api/problems/:id/examples`：取得範例測資（新增 API，不影響既有 API）
- `GET /api/stats/problems/:problem_id`：題目統計
- `POST /api/users/login`：登入
- `POST /api/submissions`：提交

---

## 相關文件

- 後端詳細說明：`backend/README.md`
- 前端詳細說明：`frontend/README.md`

