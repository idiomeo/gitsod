#!/usr/bin/env bash
set -e

# ========= 可修改的下载根地址 =========
BASE_URL="https://gitea.licnoc.top/adm/gitsod/releases/latest/download"

# ========= 判断系统与架构 =========
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64)  ARCH=amd64 ;;
    aarch64|arm64) ARCH=arm64 ;;
    *) echo "暂不支持架构 $ARCH"; exit 1 ;;
esac

BIN_NAME="gitsod-${OS}-${ARCH}"
DOWNLOAD_URL="${BASE_URL}/${BIN_NAME}"

# ========= 安装目录 =========
INSTALL_DIR="/usr/local/bin"
TARGET="${INSTALL_DIR}/gitsod"

# ========= 需要 sudo 时提示 =========
[[ ! -w "$INSTALL_DIR" ]] && SUDO="sudo" || SUDO=""

echo "正在下载 $DOWNLOAD_URL ..."
$SUDO curl -L -o "$TARGET" "$DOWNLOAD_URL"
$SUDO chmod +x "$TARGET"

echo
echo "===== 安装完成 ====="
echo "已安装到 $TARGET"
echo "现在可直接执行：gitsod --help"