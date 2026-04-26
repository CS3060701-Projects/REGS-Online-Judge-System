# REGS Frontend

這是對應 `backend` API 的前端系統（React + TypeScript + Vite）。

## 功能

- 使用者註冊/登入/登出
- 題目列表與題目詳情
- 題目統計資訊顯示
- 上傳 `.zip` 提交程式
- 查看個人提交紀錄與單筆 `operatorId` 查詢
- 管理員題目管理（建立/更新/刪除題目、上傳測資）

## 啟動

1. 先確保後端在 `http://localhost:8081` 運行。
2. 複製環境變數檔：

```bash
cp .env.example .env
```

3. 安裝與啟動：

```bash
npm install
npm run dev
```

若你的系統出現 `npm` 指令不存在，請先修正 Node.js/npm 的 PATH，或重新安裝 Node.js LTS。

## Windows 快速啟動

在 `frontend` 目錄內執行：

```bat
Start-Frontend.bat
```
