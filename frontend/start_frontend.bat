@echo off
setlocal

set "FRONTEND_DIR=%~dp0"

if not exist "%FRONTEND_DIR%" (
  echo [Error] frontend directory not found: "%FRONTEND_DIR%"
  pause
  exit /b 1
)

echo Starting frontend in a new window...
start "REGS Frontend" cmd /k "pushd ""%FRONTEND_DIR%"" && if not exist node_modules (npm install) && npm run dev"

echo Done.
echo Frontend URL: http://localhost:5173
echo You can close this window.

endlocal
