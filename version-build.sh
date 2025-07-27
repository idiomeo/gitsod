#!/usr/bin/env bash
set -euo pipefail

# 1. 询问版本（浮点数）
read -rp "请输入新的软件版本号（例如 1.2）: " NEW_VER

# 2. 替换 main.go 中整行 currentVer = <浮点数>
#    匹配浮点字面量：整数或小数点
sed -i.bak '/^[[:space:]]*currentVer[[:space:]]*float64[[:space:]]*=[[:space:]]*[0-9]\+\(\.[0-9]\+\)\?/c\	currentVer float64 = '"$NEW_VER" main.go
echo "已更新 main.go → currentVer = $NEW_VER"

# 3. 更新 ./install/VERSION（纯文本）
printf "%s" "$NEW_VER" > install/VERSION
echo "已更新 install/VERSION → $NEW_VER"

# 4. 更新 wix/gitsod.wxs 中 <Package … Version="…"
sed -i.bak -E 's/(<Package[^>]*Version=")[^"]*"/\1'"$NEW_VER"'"/' wix/gitsod.wxs
echo "已更新 wix/gitsod.wxs → Version=\"$NEW_VER\""

echo "全部版本号已同步完成！"