@echo off
setlocal

set "ROOT_DIR=%~dp0"

echo Starting REGS backend server...

start "REGS Backend" cmd /k "pushd ""%ROOT_DIR%"" && go run ./cmd/server"

echo Done.
echo You can close this window.

endlocal
