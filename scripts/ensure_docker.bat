@echo off
cd /d "%~dp0.."

echo [INFO] Ensuring PostgreSQL container is running...
docker compose up -d
if errorlevel 1 (
    echo [ERROR] Failed to start Docker. Is Docker Desktop running?
    exit /b 1
)

echo [INFO] Waiting for database to be ready...
ping 127.0.0.1 -n 4 >nul
exit /b 0
