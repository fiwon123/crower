@echo off
REM Get the folder of this batch file
SET "SCRIPT_DIR=%~dp0"
REM Remove trailing backslash
SET "SCRIPT_DIR=%SCRIPT_DIR:~0,-1%"

REM Check if already in PATH
echo %PATH% | findstr /I /C:"%SCRIPT_DIR%" >nul
IF %ERRORLEVEL% EQU 0 (
    echo Folder %SCRIPT_DIR% already in PATH
) ELSE (
    REM Add to user PATH permanently
    setx PATH "%PATH%;%SCRIPT_DIR%"
    echo Folder %SCRIPT_DIR% added to PATH
    echo Close and reopen cmd to see changes
)

pause
