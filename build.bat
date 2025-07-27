@echo off
setlocal enabledelayedexpansion

echo [1/5] Building Windows executable...
set CGO_ENABLED=0
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o gitsod.exe .

echo [2/5] Building MSI installer...
pushd wix
wix build gitsod.wxs -o gitsod.msi
popd

echo [3/5] Copying MSI to install/bin...
if not exist "install\bin" mkdir install\bin
copy /Y wix\gitsod.msi install\bin\gitsod.msi

echo [4/5] Building Linux executable...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o gitsod .

echo [5/5] Moving Linux binary to install/bin...
move /Y gitsod install\bin\gitsod

echo.
echo === Build Complete ===
echo Windows EXE : gitsod.exe
echo MSI         : install\bin\gitsod.msi
echo Linux BIN   : install\bin\gitsod
pause