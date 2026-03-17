param(
    [int]$AppPort = 0
)

$rootDir = Split-Path $PSScriptRoot -Parent
$envFile = Join-Path $rootDir ".env"
$envExample = Join-Path $rootDir ".env.example"

if (-not (Test-Path $envFile)) {
    Copy-Item $envExample $envFile
    Write-Host "Created .env from .env.example"
}

if (Test-Path $envFile) {
    Get-Content $envFile | Where-Object { $_ -match '^[^#]+=.+' } | ForEach-Object {
        $parts = $_ -split '=', 2
        [System.Environment]::SetEnvironmentVariable($parts[0].Trim(), $parts[1].Trim())
    }
}

$appKey = $env:APP_KEY
if (-not $appKey -or $appKey.Length -ne 32) {
    $letters = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
    $keyBytes = New-Object byte[] 32
    [System.Security.Cryptography.RandomNumberGenerator]::Create().GetBytes($keyBytes)
    $newKey = -join ($keyBytes | ForEach-Object { $letters[$_ % $letters.Length] })
    $envContent = Get-Content $envFile -Raw
    $envContent = $envContent -replace '(?m)^APP_KEY=.*', "APP_KEY=$newKey"
    Set-Content -Path $envFile -Value $envContent -NoNewline
    [System.Environment]::SetEnvironmentVariable("APP_KEY", $newKey)
    Write-Host "Generated APP_KEY"
} else {
    Write-Host "APP_KEY already set — skipping"
}

if ($AppPort -eq 0) {
    $AppPort = if ($env:APP_PORT) { [int]$env:APP_PORT } else { 8080 }
}

$NgrokUrl        = $env:NGROK_URL
$NgrokAuthtoken  = $env:NGROK_AUTHTOKEN
$AppmaxAppIdUUID = $env:APPMAX_APP_ID_UUID
$BaseUrl    = "http://localhost:$AppPort"
$HealthUrl  = "$BaseUrl/health"
$ComposeFiles = @("-f", "docker-compose.yml")

function Wait-ForHealth {
    for ($i = 1; $i -le 30; $i++) {
        try {
            $r = Invoke-WebRequest -Uri $HealthUrl -UseBasicParsing -TimeoutSec 3 -ErrorAction Stop
            if ($r.StatusCode -eq 200) {
                return $true
            }
        } catch {}
        Start-Sleep -Seconds 3
    }
    return $false
}

function Test-Endpoints {
    if (-not $NgrokAuthtoken) {
        Write-Host "  [FAIL] NGROK_AUTHTOKEN is empty."
        Write-Host "  Create a ngrok account at https://dashboard.ngrok.com/signup and set NGROK_AUTHTOKEN in .env"
        return $false
    }

    if (-not $AppmaxAppIdUUID) {
        Write-Host "  [FAIL] APPMAX_APP_ID_UUID is empty."
        Write-Host "  Set APPMAX_APP_ID_UUID in .env with your app's UUID from the Appmax AppStore."
        return $false
    }

    try {
        $r = Invoke-WebRequest -Uri "$BaseUrl/install/start" -UseBasicParsing -TimeoutSec 5 -ErrorAction Stop -MaximumRedirection 0
    } catch {
        $code = $_.Exception.Response.StatusCode.Value__
        if (-not $code -or $code -ge 500) {
            Write-Host "  [FAIL] GET /install/start (HTTP $code)"
            return $false
        }
    }

    $NgrokHealthBase = ""
    if ($NgrokUrl) {
        $NgrokHealthBase = if ($NgrokUrl -match '^https?://') { $NgrokUrl } else { "https://$NgrokUrl" }
        for ($i = 1; $i -le 20; $i++) {
            try {
                Invoke-WebRequest -Uri "$NgrokHealthBase/health" -UseBasicParsing -TimeoutSec 3 -ErrorAction Stop | Out-Null
                break
            } catch { Start-Sleep -Seconds 2 }
        }
    }

    $activeUrl = ""
    for ($i = 1; $i -le 40; $i++) {
        try {
            $json = & docker compose @ComposeFiles exec -T ngrok sh -lc "wget -qO- http://127.0.0.1:4040/api/tunnels 2>/dev/null" 2>$null
            if ($json -match '"public_url":"([^"]+)"') {
                $activeUrl = $Matches[1]
                break
            }
        } catch {}
        Start-Sleep -Seconds 2
    }

    if (-not $activeUrl) {
        Write-Host "  [FAIL] ngrok active tunnel URL not found in ngrok API"
        return $false
    }

    $frontendChecks = @(
        @{ Label = "Frontend URL"; Url = "$activeUrl/" },
        @{ Label = "Health URL"; Url = "$activeUrl/health" },
        @{ Label = "Callback URL"; Url = "$activeUrl/integrations/appmax/callback/install" }
    )
    foreach ($check in $frontendChecks) {
        try {
            $r = Invoke-WebRequest -Uri $check.Url -UseBasicParsing -TimeoutSec 6 -ErrorAction Stop -Headers @{ Accept = "text/html" }
            $code = [int]$r.StatusCode
            $ct = if ($r.Headers["Content-Type"]) { "$($r.Headers["Content-Type"])" } else { "" }
            if ($code -ne 200 -or $ct -notmatch "text/html") {
                Write-Host "  [FAIL] $($check.Label) did not return frontend HTML: $($check.Url) (HTTP $code, Content-Type: $ct)"
                return $false
            }
        } catch {
            $code = $_.Exception.Response.StatusCode.Value__
            Write-Host "  [FAIL] $($check.Label) did not return frontend HTML: $($check.Url) (HTTP $code)"
            return $false
        }
    }

    Write-Host "  Frontend URL: $activeUrl/"
    Write-Host "  Health URL: $activeUrl/health"
    Write-Host "  Callback URL: $activeUrl/integrations/appmax/callback/install"

    return $true
}

Push-Location $rootDir

Write-Host "==> Removing old containers and volumes..."
& docker compose @ComposeFiles down -v --remove-orphans
if ($LASTEXITCODE -ne 0) { Pop-Location; exit $LASTEXITCODE }

Write-Host "==> Starting stack [air hot reload]..."
& docker compose @ComposeFiles up -d --build
if ($LASTEXITCODE -ne 0) { Pop-Location; exit $LASTEXITCODE }

Write-Host "==> Running migrations..."
& docker compose @ComposeFiles exec app ./tmp/server artisan migrate
if ($LASTEXITCODE -ne 0) { Pop-Location; exit $LASTEXITCODE }

if (-not (Wait-ForHealth)) {
    Write-Host "Healthcheck failed after 30 attempts."
    Pop-Location
    exit 1
}

Write-Host "==> Validating endpoints..."
if (-not (Test-Endpoints)) {
    Pop-Location
    exit 1
}

Write-Host "Stack is ready."
Pop-Location
