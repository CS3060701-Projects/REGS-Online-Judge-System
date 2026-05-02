@echo off
setlocal

set "ROOT_DIR=%~dp0"
set "BACKEND_DIR=%ROOT_DIR%backend"
set "FRONTEND_DIR=%ROOT_DIR%frontend"

if not exist "%BACKEND_DIR%" (
  echo [Error] backend directory not found: "%BACKEND_DIR%"
  pause
  exit /b 1
)

if not exist "%FRONTEND_DIR%" (
  echo [Error] frontend directory not found: "%FRONTEND_DIR%"
  pause
  exit /b 1
)

echo Starting backend and frontend in separate windows...

REM Start backend server
start "REGS Backend" cmd /k "pushd ""%BACKEND_DIR%"" && go run ./cmd/server"

REM Start frontend dev server
start "REGS Frontend" cmd /k "pushd ""%FRONTEND_DIR%"" && npm run dev"

echo Done. You should get:
echo - Backend:  http://localhost:8081
echo - Frontend: http://localhost:5173
echo.
echo You can close this window.

endlocal
