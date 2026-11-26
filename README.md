# OIPivot 行业开发 SDK

`github.com/zhixingdx/public` 是一个 Go 模块。公开协议与 Go 绑定从 `public/api/...` 引入，运行时库从 `public/pkg/...` 引入，Ent 注解和驱动从 `public/ent/...` 引入。`SDK_VERSION` 记录版本和来源提交。

业务项目直接锁定一个公共仓版本。首个正式标签发布前，可将下列命令的版本换为公共仓预览提交。

```sh
go get github.com/zhixingdx/public@v0.1.0
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

在项目 go.mod 中加入上述模块依赖，并锁定 Wire 工具（`go get -tool github.com/google/wire/cmd/wire@v0.7.0`）。在 `backend/edu/course/` 维护该服务的 `tables/` 和 `config.yaml`。配置中的 `serviceType: edu`、`serviceName: course` 决定服务身份，`package`、`protoGoModule` 和 `protoImportPath` 决定业务 Go 包与 PB 的导入路径。生成器使用项目自己的 PB 包，公共 Proto 来自 go.mod 选中的 SDK。首次执行通过 Go 模块代理获取锁定的 Buf 和插件，已有相同版本的原生 Buf 可直接复用。

```sh
go mod tidy
cd backend/edu/course
yimoka gen --backend-only -c config.yaml
go mod tidy
go tool wire ./cmd/app
go build ./cmd/app
```

[教育课程服务配置、业务表与 Docker 部署](https://github.com/zhixingdx/demo)
