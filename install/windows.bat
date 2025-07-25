# 临时允许脚本执行
Set-ExecutionPolicy Bypass -Scope Process -Force

# 下载并运行安装脚本
$tempFile = "$env:TEMP\install_gitsod.ps1"
Invoke-WebRequest "https://gitee.com/idiomeo/gitsod/raw/master/install/windows.bat" -OutFile $tempFile

# 修改后的安装脚本 (修复PATH设置问题)
@"
@echo off
setlocal

set BASE_URL=https://gitee.com/idiomeo/gitsod/raw/master/install/bin/
set BIN_NAME=gitsod.exe
set DOWNLOAD_URL=%BASE_URL%/%BIN_NAME%

set INSTALL_DIR=%USERPROFILE%\bin
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"
set TARGET=%INSTALL_DIR%\gitsod.exe

echo 正在下载 %DOWNLOAD_URL% ...
powershell -Command "Invoke-WebRequest -Uri '%DOWNLOAD_URL%' -OutFile '%TARGET%' -UseBasicParsing"

if not exist "%TARGET%" (
    echo 下载失败，请检查网络或下载地址。
    pause
    exit /b
)

:: 仅添加新路径（避免重复）
powershell -Command "$newPath = '%INSTALL_DIR%'; $currentPath = [Environment]::GetEnvironmentVariable('Path', 'User'); if(-not ($currentPath -split ';' -contains $newPath)) { [Environment]::SetEnvironmentVariable('Path', "$currentPath;$newPath", 'User') }"

echo.
echo ===== 安装成功 =====
echo 文件位置: %TARGET%
echo 请关闭所有终端窗口后重新打开
echo 测试命令: gitsod --version
pause
"@ | Set-Content $tempFile -Encoding ASCII

# 执行安装
Start-Process cmd.exe -ArgumentList "/c $tempFile" -Wait

# 刷新当前会话的PATH
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "User") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "Machine")

# 验证安装
if (Get-Command gitsod -ErrorAction SilentlyContinue) {
    Write-Host "安装成功! 版本信息:" -ForegroundColor Green
    gitsod --version
} else {
    Write-Host "安装失败，请手动检查 $env:USERPROFILE\bin" -ForegroundColor Red
}