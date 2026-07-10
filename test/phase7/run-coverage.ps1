param(
    [string]$OutputPath = ""
)

$ErrorActionPreference = "Stop"

if ([string]::IsNullOrWhiteSpace($OutputPath)) {
    $OutputPath = Join-Path $PSScriptRoot "coverage.out"
}

$env:GOCACHE = Join-Path $PSScriptRoot ".gocache"

go test ./test/phase7 "-coverpkg=akadia/internal/payment/handler" "-coverprofile=$OutputPath"
go tool cover -func $OutputPath
