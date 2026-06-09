@echo off
chcp 65001 >nul
title REGS - Reset Database

cd /d "%~dp0"

echo =========================================
echo       REGS System - Reset Database
echo =========================================
echo.
echo 警告：此操作將刪除上一層目錄中的資料庫數據！
echo.

set /p "confirm=確定要繼續嗎？ (y/n): "
if /i not "%confirm%"=="y" exit /b

echo.
echo [1/3] 正在停止容器並清除資料庫數據
docker-compose -f "..\docker-compose.yml" down -v

echo.
echo [2/3] 強制清理可能殘留的數據 (若有錯誤請檢查 docker-compose.yml 中的 volume 名稱)...
docker volume rm REGS_pgdata

echo.
echo [3/3] 正在重新啟動服務...
docker-compose -f "..\docker-compose.yml" up -d

echo.
echo -----------------------------------------
echo 資料庫重置完成！
echo -----------------------------------------
pause