@echo off
:: 設定終端機為 UTF-8 編碼，避免中文字變成亂碼
chcp 65001 >nul
title REGS 後端伺服器

echo =========================================
echo         REGS 評測系統 - 伺服器啟動
echo =========================================
echo.

:: 切換到這個 bat 檔所在的絕對目錄，確保路徑正確
cd /d "%~dp0"

echo [1/2] 正在編譯最新程式碼...
go build -o server.exe ./cmd/server

:: 檢查編譯是否成功 (Error Level 不為 0 代表有錯)
if %errorlevel% neq 0 (
    echo.
    echo [錯誤] 編譯失敗！請檢查上方的錯誤訊息。
    pause
    exit /b %errorlevel%
)

echo [2/2] 編譯成功！準備啟動伺服器...
echo -----------------------------------------
echo 提示：若要關閉伺服器，請直接關閉此視窗，或按 Ctrl+C
echo -----------------------------------------
echo.

:: 執行編譯好的伺服器
server.exe

:: 如果伺服器意外關閉，暫停視窗讓你能夠看到錯誤訊息
echo.
pause