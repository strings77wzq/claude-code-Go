# install.ps1 — One-command installer for go-code on Windows
# Usage: irm https://raw.githubusercontent.com/strings77wzq/claude-code-Go/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$Repo = "strings77wzq/claude-code-Go"
$BinaryName = "go-code.exe"

Write-Host "========================================" -ForegroundColor Green
Write-Host "  go-code Installer (Windows)" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""

# Detect architecture
$Arch = $env:PROCESSOR_ARCHITECTURE
if ($Arch -eq "AMD64" -or $Arch -eq "x86_64") {
    $Arch = "amd64"
} else {
    Write-Host "Unsupported architecture: $Arch" -ForegroundColor Red
    exit 1
}

$Binary = "go-code-windows-$Arch.exe"
$DownloadUrl = "https://github.com/$Repo/releases/latest/download/$Binary"

Write-Host "Detected: windows/$Arch" -ForegroundColor Yellow
Write-Host "Downloading: $Binary" -ForegroundColor Yellow

# Download
$TempDir = [System.IO.Path]::GetTempPath()
$DestPath = Join-Path $TempDir $Binary
Invoke-WebRequest -Uri $DownloadUrl -OutFile $DestPath -UseBasicParsing

# Install to user's bin directory
$UserBin = Join-Path $env:USERPROFILE ".local\bin"
if (-not (Test-Path $UserBin)) {
    New-Item -ItemType Directory -Path $UserBin -Force | Out-Null
}

$InstallPath = Join-Path $UserBin $BinaryName
Copy-Item $DestPath $InstallPath -Force
Remove-Item $DestPath -Force

# Add to PATH if not already there
$CurrentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
if ($CurrentPath -notlike "*$UserBin*") {
    [Environment]::SetEnvironmentVariable("PATH", "$CurrentPath;$UserBin", "User")
    $env:PATH = "$env:PATH;$UserBin"
}

Write-Host ""
Write-Host "Binary installed to: $InstallPath" -ForegroundColor Green
Write-Host ""

# Run setup wizard
Write-Host "Starting setup wizard..." -ForegroundColor Yellow
Write-Host ""
try {
    & $InstallPath --setup
    if ($LASTEXITCODE -ne 0) {
        throw "Setup exited with code $LASTEXITCODE"
    }
} catch {
    Write-Host ""
    Write-Host "Setup wizard failed: $_" -ForegroundColor Yellow
    Write-Host "You can configure manually:"
    Write-Host "  mkdir `$env:USERPROFILE\.go-code"
    Write-Host "  echo '{`"apiKey`": `"sk-ant-...`"}' > `$env:USERPROFILE\.go-code\settings.json"
    Write-Host ""
    Write-Host "Or run setup again:"
    Write-Host "  go-code --setup"
}
