#!/bin/sh
# 从 GitHub Release 或本地同版制品安装到用户目录。
set -eu
version=${1:?用法: sh install.sh vX.Y.Z}
printf '%s\n' "$version" | grep -Eq '^v[01]\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$' || exit 2
case "$(uname -s)" in Darwin) os=darwin ;; Linux) os=linux ;; *) echo 'Windows 请使用 install.ps1' >&2; exit 2 ;; esac
case "$(uname -m)" in arm64|aarch64) arch=arm64 ;; x86_64|amd64) arch=amd64 ;; *) echo '不支持当前 CPU 架构' >&2; exit 2 ;; esac
asset=yimoka-$os-$arch.gz
dest=${YIMOKA_BIN_DIR:-"$HOME/.local/bin"}
mkdir -p "$dest"
work=$(mktemp -d "$dest/.yimoka.XXXXXXXX")
trap 'rm -rf "$work"' EXIT
for file in "$asset" SHA256SUMS; do
  if [ -n "${YIMOKA_ASSET_DIR:-}" ]; then
    cp "$YIMOKA_ASSET_DIR/$file" "$work/$file"
  else
    curl -fsSL --retry 3 "https://github.com/zhixingdx/public/releases/download/$version/$file" -o "$work/$file"
  fi
done
(cd "$work"
  awk -v asset="./$asset" '$2 == asset {print; n++} END {if (n != 1) exit 1}' SHA256SUMS > check
  if command -v sha256sum >/dev/null; then sha256sum -c check; else shasum -a 256 -c check; fi
  gzip -dc "$asset" > yimoka
  chmod 755 yimoka
  case "$(./yimoka version)" in "yimoka version $version ("*) ;; *) echo 'CLI 版本不匹配' >&2; exit 1 ;; esac
)
mv -f "$work/yimoka" "$dest/yimoka"
"$dest/yimoka" version
echo "安装位置 $dest/yimoka"
