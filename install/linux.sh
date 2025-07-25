#!/usr/bin/env bash
set -e

# 固定下载地址（仅 Linux x86-64）
BASE_URL="https://gitea.licnoc.top/adm/gitsod/releases/latest/download"
BIN_NAME="gitsod-linux-amd64"
DOWNLOAD_URL="${BASE_URL}/${BIN_NAME}"

# 安装目录
INSTALL_DIR="/usr/local/bin"
TARGET="${INSTALL_DIR}/gitsod"

# 自动 sudo
[[ ! -w "$INSTALL_DIR" ]] && SUDO="sudo" || SUDO=""

echo "正在下载 $DOWNLOAD_URL ..."
$SUDO curl -L -o "$TARGET" "$DOWNLOAD_URL"
$SUDO chmod +x "$TARGET"

echo
echo "===== 安装完成 ====="
echo "已安装到 $TARGET"
echo "现在可直接执行：gitsod --help"