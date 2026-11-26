# OIPivot 行业开发 SDK

`SDK_VERSION` 依次记录 SDK 版本和平台来源提交。根目录的 `pkg/` 提供运行时库，`ent/` 提供建模注解和驱动，`api/` 是独立的 Go 与 Proto 契约模块。公共 tag `vX.Y.Z` 和 `api/vX.Y.Z` 指向同一提交。Go 源与 Proto 沿用平台的类型和导入名。

在行业项目的 go.mod 中锁定两个模块的实际下载位置。以下示例使用 v0.1.0，升级时一起修改版本。

```go
require (
    github.com/zhixingdx/oipivot/backend v0.1.0
    github.com/zhixingdx/oipivot-api v0.1.0
)
replace github.com/zhixingdx/oipivot/backend => github.com/zhixingdx/public v0.1.0
replace github.com/zhixingdx/oipivot-api => github.com/zhixingdx/public/api v0.1.0
```

## 安装 Yimoka

在本仓对应版本的检出目录运行安装脚本。脚本自动识别系统和 CPU，从同版本 Release 下载压缩的可执行文件，校验 SHA256 和 CLI 自报版本后安装。支持 macOS、Linux、Windows 的 amd64 和 arm64。已安装程序可直接复制到离线开发机。

```sh
# macOS / Linux
sh install.sh v0.1.0
export PATH="$HOME/.local/bin:$PATH"
yimoka version
```

```powershell
# Windows PowerShell
.\install.ps1 -Version v0.1.0
$env:Path += ";$env:LOCALAPPDATA\Yimoka\bin"
yimoka version
```

设置 `YIMOKA_BIN_DIR` 可自选安装目录。本地演练时把压缩包及 SHA256SUMS 放到 `dist/`，设置 `YIMOKA_ASSET_DIR` 为其绝对路径后运行同一安装命令。

## 开发业务服务

在项目 go.mod 中加入上述两个依赖和替换，并锁定 Wire 工具（`go get -tool github.com/google/wire/cmd/wire@v0.7.0`）。在项目根目录执行 `yimoka init shared/<name>`，维护新服务的 `tables/` 后运行以下命令。生成器使用项目自己的 PB 包，公共 Proto 来自 go.mod 选中的 SDK。首次执行通过 Go 模块代理获取锁定的 Buf 和插件，已有相同版本的原生 Buf 可直接复用。

```sh
go mod tidy
cd backend/shared/<name>
yimoka gen --backend-only -c config.yaml
go mod tidy
go tool wire ./cmd/app
go build ./cmd/app
```

[完整的维修工单示例与 Docker 部署](https://github.com/zhixingdx/demo)
