@echo off
title REGS - Reset Database

echo =========================================
echo      REGS System - Reset Database
echo =========================================
echo.
echo WARNING: This will completely and IRREVERSIBLY delete all data
echo from the PostgreSQL database, including all users, problems,
echo and submissions.
echo.

set /p "confirm=Are you sure you want to continue? (y/n): "

if /i not "%confirm%"=="y" (
    echo Operation cancelled.
    pause
    exit /b
)

echo.
echo Stopping database container and removing data volume...
cd /d "%~dp0..\"
docker-compose down -v

echo.
echo Database has been successfully reset.
echo You can now restart it using 'docker-compose up -d'.
echo.
pause