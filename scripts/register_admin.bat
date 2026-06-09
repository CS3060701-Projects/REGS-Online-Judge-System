@echo off
chcp 65001 >nul
title REGS - Create Admin User

echo =========================================
echo       REGS System - Create Admin User
echo =========================================
echo.

call "%~dp0ensure_docker.bat"
if %errorlevel% neq 0 (
    pause
    exit /b 1
)

cd /d "%~dp0.."

go run ./cmd/seed
if %errorlevel% neq 0 (
    echo.
    echo [錯誤] 建立管理員失敗。
    pause
    exit /b %errorlevel%
)

echo.
echo Operation finished.
pause
