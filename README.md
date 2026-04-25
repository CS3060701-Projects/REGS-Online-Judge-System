# **REGS 後端評測系統**

## **專案概述**

本專案為 CS3060701 期末專案之後端系統。基於 Go (Gin) 與 PostgreSQL 開發，透過 Docker 實作隔離沙盒，支援 C++ 專案之 CMake 自動化建置與非同步測資比對。

## **核心功能**

* **非同步任務隊列**：基於 Worker Pool 實作，控制最高併發數以維持系統穩定。  
* **分段自動化編譯**：  
  * Configure 階段：cmake \-G Ninja \-B build (失敗判定 SE)。  
  * Build 階段：cmake \--build build \--verbose (失敗判定 CE)。  
* **沙盒隔離 (Docker)**：採用 yhlib/cs3060701 映像檔，套用 \--network none 斷網，並限制 CPU (1.0)、Memory (256m) 及 PIDs (50)。  
* **判題與日誌**：支援 AC, WA, CE, SE, RE, TLE 狀態判定，並實體化生成 configure.log, compile.log, output.log 供後續查詢。

## **環境需求**

* Go 1.20+  
* Docker Desktop 或 Docker Engine  
* PostgreSQL (Port 5433\)

## ---

**環境配置與啟動**

### **1\. 準備測試資料**

於專案根目錄配置測資結構 (以 p1001 為例)：

Plaintext

test\_data/  
└── p1001/  
    ├── 1.in  
    └── 1.out

### **2\. 資料庫配置**

確認 PostgreSQL 環境設定符合以下連線參數：

Plaintext

host=localhost  
user=regs\_user  
password=regs\_password  
dbname=regs\_db  
port=5433

### **3\. 系統啟動順序**

依序執行以下步驟啟動系統：

1. 啟動 Docker Desktop 應用程式。  
2. 於專案根目錄啟動資料庫容器：  
   Bash  
   docker-compose up \-d

3. 啟動後端伺服器：  
   Bash  
   server.bat

*(伺服器啟動後，終端機將顯示資料庫連線、資料表遷移及 Worker Pool 啟動日誌。)*

## ---

**測試指南**

### **1\. 準備提交檔案**

建立 C++ 專案，確保根目錄包含以下 CMakeLists.txt。將原始碼與 CMake 設定檔一起壓縮為 submission.zip (請勿包含外層資料夾)。

CMake

cmake\_minimum\_required(VERSION 3.10)  
project(Submission)  
add\_executable(main main.cpp)

### **2\. 發送 API 請求**

使用 Postman 或其他工具發送提交請求：

* **URL**: http://localhost:8081/api/submissions  
* **Method**: POST  
* **Body**: multipart/form-data  
  * problem\_id (Text): p1001  
  * user\_id (Text): 0  
  * file (File): 上傳 submission.zip

### **3\. 驗證結果**

1. **API 回應**：確認收到 HTTP 200 回應，並取得系統配發的 operatorId。  
2. **終端機日誌**：檢查伺服器後台，確認 Worker 依序執行 Configure、Build 與測資比對，並印出最終判定狀態 (AC / WA 等)。  
3. **實體日誌檢查**：前往 storage/workspaces/\<operatorId\>/，確認是否成功寫入 .log 檔案。