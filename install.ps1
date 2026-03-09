# SnapGit installer for Windows
# Usage: irm https://raw.githubusercontent.com/ovair/snapgit/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

$repo = "ovair/snapgit"
$installDir = "$env:LOCALAPPDATA\SnapGit\bin"

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { "arm64" } else { "amd64" }
} else {
    Write-Error "SnapGit requires a 64-bit operating system."
    return
}

# Get latest release tag
$release = Invoke-RestMethod -Uri "https://api.github.com/repos/$repo/releases/latest"
$version = $release.tag_name
Write-Host "Installing SnapGit $version ..." -ForegroundColor Cyan

# Download
$asset = "sg_$($version.TrimStart('v'))_windows_$arch.zip"
$url = "https://github.com/$repo/releases/download/$version/$asset"
$tmp = Join-Path $env:TEMP $asset

Write-Host "Downloading $asset ..."
Invoke-WebRequest -Uri $url -OutFile $tmp -UseBasicParsing

# Extract
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
}
Expand-Archive -Path $tmp -DestinationPath $installDir -Force
Remove-Item $tmp

# Add to PATH if not already there
$userPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($userPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$userPath;$installDir", "User")
    Write-Host "Added $installDir to your PATH." -ForegroundColor Yellow
    Write-Host "Restart your terminal for PATH changes to take effect." -ForegroundColor Yellow
}

Write-Host "SnapGit $version installed successfully!" -ForegroundColor Green
Write-Host "Run 'sg help' to get started."
