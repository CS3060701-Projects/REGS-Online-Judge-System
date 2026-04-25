# **REGS 後端評測系統 (REGS Backend)**

## **專案概述**

本專案為 CS3060701 期末專案的後端評測系統。系統採用 Go 語言開發，結合 Gin 框架與 PostgreSQL 資料庫，並利用 Docker 容器技術建立安全的隔離沙盒。系統提供非同步、自動化且安全的程式碼編譯與評測環境，支援 C++ 專案的 CMake 自動化建置與標準測資比對。

## **目前開發進度**

目前已完成專案規格書中，自動化編譯的 30 分標準。

### **1\. 核心基礎建設**

* **資料庫與資料模型**：完成 User, Problem, Submission 資料表的建立與關聯 (GORM)。

* **非同步任務隊列**：實作基於 Goroutine 與 Channel 的 Worker Pool 架構。接收提交請求後立即回傳 operatorId，並將任務放入佇列，嚴格控制最高併發數量，確保伺服器穩定性。

### **2\. 自動化編譯**

* **檔案處理**：支援接收 ZIP 壓縮檔，自動建立獨立工作空間 (storage/workspaces/{operatorId}) 並進行解壓縮。  
* **分段編譯機制**：  
  * **Configure 階段**：使用 Ninja 產生器 (cmake \-G Ninja \-B build)，若失敗則判定為 SE (Setup Error)。  
  * **Build 階段**：執行編譯 (cmake \--build build \--verbose)，若失敗則判定為 CE (Compile Error)。  
* **實體日誌儲存**：各階段的標準輸出與錯誤輸出已實體化，分別儲存為 configure.log, compile.log, output.log，供後續 API 查詢。

### **3\. 安全隔離與判題邏輯**

* **沙盒環境**：使用指定的 Docker 映像檔 (yhlib/cs3060701) 執行所有外部程式碼。  
* **資源限制與網路隔離**：執行評測時套用 \--network none 中斷網路，並限制 CPU (--cpus 1.0)、記憶體 (--memory 256m) 與處理程序數量 (--pids-limit 50)。  
* **狀態判定**：包含 AC (Accepted), WA (Wrong Answer), CE (Compile Error), SE (Setup Error), RE (Runtime Error)，以及TLE (Time Limit Exceeded) 超時控制。

## ---

**系統需求**

在開始測試前，請確保開發環境已安裝以下組件：

* Go (1.20 或以上版本)  
* Docker Desktop 或 Docker Engine  
* PostgreSQL (設定為 Port 5433\)  
* Postman (或其他 API 測試工具)

## ---

**環境配置與啟動**

### **1\. 準備測試資料 (Testcases)**

在專案根目錄下建立測資資料夾，系統目前預設測試題目為 p1001：

Plaintext

test\_data/  
└── p1001/  
    ├── 1.in  
    └── 1.out

### **3\. 資料庫配置**

host=localhost
user=regs_user
password=regs_passwod
dbname=regs_db
port=5433

### **4\. 啟動伺服器**

1.開啟 Docker Desktop 應用程式
2.在regs-backend底下執行 docker-compose up -d
3.執行server.bat啟動伺服器：

伺服器啟動時，終端機將顯示資料庫連線成功、資料表遷移完成，以及 Worker Pool 啟動的日誌。

## ---

**測試指南**

### **1\. 準備提交檔案 (Submission)**

建立一個供測試的 C++ 專案，必須包含 CMakeLists.txt 與原始碼檔案。

**CMakeLists.txt 範例**：

cmake\_minimum\_required(VERSION 3.10)  
project(Submission)  
add\_executable(main main.cpp)

將該目錄下的所有檔案（不包含外層資料夾本身）壓縮為 submission.zip。

### **2\. 使用 Postman 進行提交**

* **Method**: POST  
* **URL**: http://localhost:8081/api/submissions  
* **Body Type**: multipart/form-data  
* **欄位範例**:  
  * problem\_id (Text): p1001  
  * user\_id (Text): 0
  * file (File): 選擇您剛剛建立的 submission.zip

### **3\. 驗證結果**

1. **API 回應**：點擊 Send 後，應立即收到 HTTP 200 回應，包含生成的 operatorId。  

2. **伺服器日誌**：觀察 Go 伺服器的終端機輸出，將看到背景 Worker 接手任務，並依序印出 Configure、Build 與各測資的比對結果，最終輸出狀態 (例如 AC 或 WA)。  

3. **日誌檔案檢查**：前往 storage/workspaces/\<operatorId\>/ 目錄，檢查是否成功生成 configure.log, compile.log, 以及 output.log。 (若已實作自動清理機制，可暫時將清理邏輯註解以檢視檔案)。