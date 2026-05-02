@echo off
chcp 65001 >nul
title REGS 後端伺服器

echo =========================================
echo         REGS 評測系統 - 伺服器啟動
echo =========================================
echo.

:: 切換到腳本所在的資料夾路徑
cd /d "%~dp0"

:: -----------------------------------------
:: [步驟 1/3] 檢查 Docker 容器環境
:: -----------------------------------------
echo [1/3] 檢查 Docker 容器環境...

:: 檢查是否有任何狀態為 running 的服務
docker compose ps --status running --format "{{.Service}}" | findstr . >nul

if %errorlevel% neq 0 (
    echo [提示] 檢測到 Docker 容器尚未啟動，正在執行 docker compose up -d...
    docker compose up -d
    
    if %errorlevel% neq 0 (
        echo.
        echo 請在 Docker 啟動後重新執行腳本。
        pause
        exit /b %errorlevel%
    )
    pause
    exit /b
)

echo [提示] Docker 容器已在運行中。
echo.

:: -----------------------------------------
:: [步驟 2/3] 編譯最新程式碼
:: -----------------------------------------
echo [2/3] 正在編譯最新程式碼...
if exist server.exe del server.exe

go build -o server.exe ./cmd/server

if %errorlevel% neq 0 (
    echo.
    echo [錯誤] 編譯失敗！請檢查上方的錯誤訊息。
    pause
    exit /b %errorlevel%
)

:: -----------------------------------------
:: [步驟 3/3] 啟動伺服器
:: -----------------------------------------
echo [3/3] 編譯成功！準備啟動伺服器...
echo -----------------------------------------
echo 提示：若要關閉伺服器，請直接關閉此視窗，或按 Ctrl+C
echo -----------------------------------------
echo.

:: 執行伺服器
server.exe

echo.
pause