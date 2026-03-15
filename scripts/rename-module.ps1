param(
    [string]$NewPath = ""
)

$rootDir = Split-Path $PSScriptRoot -Parent
Push-Location $rootDir

$oldPath = (Select-String -Path go.mod -Pattern '^module\s+(\S+)').Matches[0].Groups[1].Value

if (-not $NewPath) {
    $envFile = Join-Path $rootDir ".env"
    if (Test-Path $envFile) {
        $line = Get-Content $envFile | Where-Object { $_ -match '^MODULE_PATH=' } | Select-Object -First 1
        if ($line) { $NewPath = ($line -split '=', 2)[1].Trim() }
    }
}

if (-not $NewPath) {
    Write-Host "Usage: .\rename-module.ps1 -NewPath <new-module-path>"
    Write-Host "       or set MODULE_PATH in .env"
    Pop-Location; exit 1
}

if ($oldPath -eq $NewPath) {
    Write-Host "Module path already is '$NewPath' — nothing to do."
    Pop-Location; exit 0
}

Write-Host "Renaming module:"
Write-Host "  old: $oldPath"
Write-Host "  new: $NewPath"

go mod edit -module $NewPath

Get-ChildItem -Recurse -Filter "*.go" |
    Where-Object { $_.FullName -notmatch '\.git|\.gocache|\.gomodcache|vendor' } |
    ForEach-Object {
        $content = Get-Content $_.FullName -Raw
        $updated = $content -replace [regex]::Escape("`"$oldPath"), "`"$NewPath"
        if ($content -ne $updated) {
            Set-Content -Path $_.FullName -Value $updated -NoNewline
        }
    }

Write-Host "Done. Run 'go mod tidy' to verify."
Pop-Location
