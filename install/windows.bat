@echo off
setlocal

:: 下载地址
set BASE_URL=https://gitee.com/idiomeo/gitsod/raw/master/install/bin/
set BIN_NAME=gitsod.exe
set DOWNLOAD_URL=%BASE_URL%/%BIN_NAME%

:: 安装目录
set INSTALL_DIR=%USERPROFILE%\bin
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
set TARGET=%INSTALL_DIR%\gitsod.exe

:: 下载并覆盖
echo 正在下载 %DOWNLOAD_URL% ...
powershell -Command "& {Invoke-WebRequest -Uri '%DOWNLOAD_URL%' -OutFile '%TARGET%'}"

if not exist "%TARGET%" (
    echo 下载失败，请检查网络或下载地址。
    pause
    exit /b
)

:: 添加到用户 PATH
setx PATH "%PATH%;%INSTALL_DIR%" >nul 2>&1

echo.
echo ===== 安装完成 =====
echo 已安装到 %TARGET%
echo 关闭并重新打开终端后，即可使用 gitsod 命令。
pause