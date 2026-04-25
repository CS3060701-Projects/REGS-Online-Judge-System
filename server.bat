@echo off
chcp 65001 >nul
title REGS 後端伺服器

echo =========================================
echo         REGS 評測系統 - 伺服器啟動
echo =========================================
echo.

cd /d "%~dp0"

echo [1/2] 正在編譯最新程式碼...
go build -o server.exe ./cmd/server

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

server.exe

echo.
pause