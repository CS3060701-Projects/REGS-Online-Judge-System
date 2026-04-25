# **CS3060701 期末專案 \- REGS (Online Judge System)**

## **專案概述**

**REGS** 是一個現代化的線上評測系統（Online Judge），專為跨平台使用、自動化處理與環境隔離而設計。本系統結合了 **Docker 容器技術** 與 **CMake/Ninja 編譯工具鏈**，確保使用者提交的程式碼在完全隔離的沙盒環境中運行，並提供即時的評測結果與多階段日誌查詢。

## ---

**核心功能**

* **自動化編譯管線**：支援 .zip 格式上傳，自動執行 CMake 配置與 Ninja 編譯。  
* **環境隔離與安全**：所有編譯與執行階段均在 Docker 容器內完成，並使用 \--network none 實施完全斷網隔離。  
* **非同步評測隊列**：系統接收提交後立即回傳 operatorId，評測邏輯在背景透過 Job Queue 執行。  
* **五種狀態判定**：準確回傳 **AC** (Accepted)、**WA** (Wrong Answer)、**CE** (Compile Error)、**RE** (Runtime Error)、**SE** (Setup Error) 與 **TLE** (Time Limit Exceeded)。  
* **分層權限管理 (RBAC)**：定義了 Admin、User 與 Guest 三種權限角色，並以 JWT 進行身份驗證。

## ---

**環境需求**

* **Go**: 1.2x 或以上版本  
* **Docker**: 需具備執行容器之權限，並預先下載映像檔：docker pull yhlib/cs3060701  
* **PostgreSQL**: 作為主要的資料庫系統  
* **CMake & Ninja**: 容器內部預裝的編譯工具鏈

## ---

**快速開始**

### **1\. 資料庫配置**

請在 PostgreSQL 中建立名為 regs 的資料庫，並確保你的連接字串（DSN）正確配置於 internal/database/database.go。

### **2\. 啟動伺服器**

1.啟動 Docker Desktop 應用程式。
2.於 regs-backend 啟動資料庫容器
```Bash
docker-compose up -d
```
3.啟動 Server.bat

### **3\. 目錄結構準備**

系統啟動後會自動在根目錄建立以下資料夾，用於存放評測數據：

* storage/submissions/：儲存使用者上傳的原始原始碼 zip 檔。  
* storage/workspaces/：評測時的臨時工作空間與日誌檔。  
* storage/testdata/：存放題目測資，支援 Admin 打包下載。

## ---

**API 權限表摘要**

| 功能 | 方法 | URL | 權限要求 |
| :---- | :---- | :---- | :---- |
| 註冊/登入 | POST | /api/users/register, /api/users/login | Guest |
| 題目清單 | GET | /api/problems | Guest |
| 提交程式碼 | POST | /api/submissions | User |
| 下載原始碼 | GET | /api/submissions/{id}/source | User (本人/Admin) |
| 建立題目 | PUT | /api/problems | Admin |
| 下載測資 | GET | /api/problems/{id}/testcases | Admin |

## ---

**評測流程說明**

1. **任務受理**：伺服器接收壓縮檔後，分配 operatorId 並回傳。  
2. **預處理**：檢查專案根目錄是否含有 CMakeLists.txt。  
3. **配置與編譯**：在斷網容器中執行 cmake \-G Ninja 與 cmake \--build。若失敗則分別記錄為 SE 或 CE。  
4. **執行與比對**：逐一執行測資點，監控 Exit Code（非 0 則為 RE）與資源消耗（TLE/MLE），並進行逐行結果比對。  
5. **日誌產出**：所有階段的輸出將分別存為 configure.log、compile.log 與 output.log 供查詢。

## ---

**授權說明**

本專案為 NTUST CS3060701 課程期末專案，相關內容遵循學術誠信規範。