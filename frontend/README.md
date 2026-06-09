# REGS Frontend

前端應用程式使用 React + TypeScript + Vite，對應 `backend` 提供的 REST API，並支援使用者註冊、登入、題目瀏覽與提交紀錄查詢。

## 功能

- 使用者註冊、登入、登出
- 題目列表瀏覽與搜尋
- 題目詳情與統計資訊
- `.zip` 程式碼上傳提交
- 查詢個人提交紀錄與指定 `operatorId`
- 管理員題目管理：建立、更新、刪除題目、上傳測資

## 環境需求

- Node.js
- npm
- 後端服務應運行於 `http://localhost:8081`

## 啟動步驟

### 1. 複製環境變數檔

在 `frontend` 目錄中執行：

```bat
copy .env.example .env
```

若您使用 Bash，也可以執行：

```bash
cp .env.example .env
```

### 2. 安裝相依套件

```bash
npm install
```

### 3. 啟動開發伺服器

```bash
npm run dev
```

開啟後，前端預設網址為：

`http://localhost:5173`

## 設定檔

`frontend/.env.example` 內容：

```env
VITE_API_BASE_URL=http://localhost:8081/api
```

如果後端 API 位置改變，請更新 `.env` 中的 `VITE_API_BASE_URL`。

## Windows 快速啟動

在 `frontend` 目錄內執行：

```bat
start_frontend.bat
```

此腳本會在新的命令提示字元視窗中啟動 Vite 開發伺服器。

## 注意事項

- 前端與後端開發伺服器需要同時啟動。
- 當後端部署在不同位置時，請調整 `VITE_API_BASE_URL`。
- 如果 `npm` 指令不存在，請安裝 Node.js 並確認 PATH 設定。

