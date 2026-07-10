param(
    [string]$OutputPath = ""
)

$ErrorActionPreference = "Stop"

$repoRoot = Resolve-Path (Join-Path $PSScriptRoot "..\..")
Set-Location $repoRoot

if ([string]::IsNullOrWhiteSpace($OutputPath)) {
    $OutputPath = Join-Path $PSScriptRoot "coverage.out"
}

$cachePath = Join-Path $PSScriptRoot ".gocache"
if (-not (Test-Path $cachePath)) {
    New-Item -ItemType Directory -Path $cachePath | Out-Null
}
$env:GOCACHE = $cachePath

$coverPkg = "akadia/internal/payment/service,akadia/internal/shared"

& go test ./test/phase1 "-coverpkg=$coverPkg" "-coverprofile=$OutputPath"
& go tool cover -func $OutputPath
