@echo off
title REGS - Create Admin User

echo =========================================
echo       REGS System - Create Admin User
echo =========================================
echo.
echo This script will check the database and create the first admin user if needed.
echo Please make sure the PostgreSQL container in Docker is running.
echo.

@echo off
cd /d "%~dp0.."

go run cmd/seed/main.go

pause

echo.
echo Operation finished.
pause