param([Parameter(Mandatory=$true)][ValidatePattern('^v[01]\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$')][string]$Version)
$ErrorActionPreference = 'Stop'
$arch = switch ([Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString()) {
    'X64'   { 'amd64' }
    'Arm64' { 'arm64' }
    default { throw '不支持当前 CPU 架构' }
}
$name = "yimoka-windows-$arch.gz"
$dest = if ($env:YIMOKA_BIN_DIR) { $env:YIMOKA_BIN_DIR } else { Join-Path $env:LOCALAPPDATA 'Yimoka/bin' }
$work = Join-Path $dest ([guid]::NewGuid().ToString())
New-Item -ItemType Directory -Path $work -Force | Out-Null
try {
    foreach ($file in @($name, 'SHA256SUMS')) {
        $path = Join-Path $work $file
        if ($env:YIMOKA_ASSET_DIR) { Copy-Item (Join-Path $env:YIMOKA_ASSET_DIR $file) $path }
        else { Invoke-WebRequest -UseBasicParsing "https://github.com/zhixingdx/public/releases/download/$Version/$file" -OutFile $path }
    }
    $matchesForFile = @(Get-Content (Join-Path $work 'SHA256SUMS') | Where-Object { ($_ -split '\s+')[1] -eq "./$name" })
    $hash = (Get-FileHash (Join-Path $work $name) -Algorithm SHA256).Hash
    if ($matchesForFile.Count -ne 1 -or ($matchesForFile[0] -split '\s+')[0] -ne $hash) { throw 'SHA256 校验失败' }
    $exe = Join-Path $work 'yimoka.exe'
    $inputFile = [IO.File]::OpenRead((Join-Path $work $name))
    try {
        $gzip = [IO.Compression.GZipStream]::new($inputFile, [IO.Compression.CompressionMode]::Decompress)
        $outputFile = [IO.File]::Create($exe)
        try { $gzip.CopyTo($outputFile) } finally { $outputFile.Dispose(); $gzip.Dispose() }
    } finally { $inputFile.Dispose() }
    $info = & $exe version
    if ($LASTEXITCODE -ne 0 -or !$info.StartsWith("yimoka version $Version (")) { throw 'CLI 版本不匹配' }
    Move-Item $exe (Join-Path $dest 'yimoka.exe') -Force
    Write-Output "$info — $dest"
} finally { Remove-Item -Recurse -Force $work }
