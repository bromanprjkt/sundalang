$ErrorActionPreference = "Stop"

function Get-TermWidth {
    try { return $Host.UI.RawUI.WindowSize.Width } catch { return 80 }
}

function Write-Center {
    param([string]$Text, [string]$Color = "White")
    $width = Get-TermWidth
    $padding = [Math]::Max(0, [Math]::Floor(($width - $Text.Length) / 2))
    Write-Host (" " * $padding + $Text) -ForegroundColor $Color
}

function Write-Box {
    param([string[]]$Lines, [string]$BorderColor = "Cyan", [string]$TextColor = "White")
    $width = Get-TermWidth
    $maxLen = ($Lines | Measure-Object -Maximum -Property Length).Maximum
    $boxWidth = [Math]::Min($maxLen + 6, $width - 4)
    $innerWidth = $boxWidth - 4
    $padding = [Math]::Max(0, [Math]::Floor(($width - $boxWidth) / 2))
    $pad = " " * $padding
    
    Write-Host ($pad + [char]0x256D + ([string][char]0x2500 * ($boxWidth - 2)) + [char]0x256E) -ForegroundColor $BorderColor
    foreach ($line in $Lines) {
        $textPad = [Math]::Max(0, [Math]::Floor(($innerWidth - $line.Length) / 2))
        $leftPad = " " * $textPad
        $rightPad = " " * ($innerWidth - $textPad - $line.Length)
        Write-Host ($pad + [char]0x2502 + " " + $leftPad + $line + $rightPad + " " + [char]0x2502) -ForegroundColor $BorderColor -NoNewline
        Write-Host ""
    }
    Write-Host ($pad + [char]0x2570 + ([string][char]0x2500 * ($boxWidth - 2)) + [char]0x256F) -ForegroundColor $BorderColor
}

function Write-Status {
    param([string]$Status, [string]$Text, [string]$Color)
    $width = Get-TermWidth
    $msg = "[$Status] $Text"
    $padding = [Math]::Max(0, [Math]::Floor(($width - $msg.Length) / 2))
    Write-Host (" " * $padding + $msg) -ForegroundColor $Color
}

function Write-Divider {
    $width = Get-TermWidth
    $divLen = [Math]::Min(50, $width - 10)
    $padding = [Math]::Max(0, [Math]::Floor(($width - $divLen) / 2))
    Write-Host (" " * $padding + ([string][char]0x2500 * $divLen)) -ForegroundColor DarkGray
}

[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Clear-Host
Write-Host ""
Write-Box @(
    "",
    "   SUNDALANG   ",
    "",
    "Bahasa Pemrograman Sunda Pandeglang",
    "Installer for Windows",
    ""
) -BorderColor Cyan -TextColor White
Write-Host ""

$InstallDir = Join-Path $env:USERPROFILE ".sundalang\bin"
$BinaryPath = Join-Path $InstallDir "sundalang.exe"

if (-not (Test-Path $InstallDir)) {
    Write-Status "INFO" "Nyieun folder instalasi..." Yellow
    New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    Write-Status "OK" "Folder instalasi dibuat" Green
}
Write-Host ""

Write-Divider
Write-Host ""
Write-Status "INFO" "Nyari versi terbaru..." Yellow
Write-Host ""

try {
    $ReleaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/bromanprjkt/sundalang/releases/latest"
    $Version = $ReleaseInfo.tag_name
    
    $Asset = $ReleaseInfo.assets | Where-Object { $_.name -eq "sundalang.exe" }
    
    if (-not $Asset) {
        Write-Host ""
        Write-Status "ERROR" "Gagal manggihan binary Windows" Red
        exit 1
    }
    
    $DownloadUrl = $Asset.browser_download_url
    Write-Status "OK" "Kapanggih: SundaLang $Version" Green
    
} catch {
    Write-Host ""
    Write-Status "ERROR" "Gagal nyambung ka GitHub" Red
    Write-Center "Pastikeun koneksi internet aktif" Red
    exit 1
}

Write-Host ""
Write-Divider
Write-Host ""
Write-Status "INFO" "Ngeundeur SundaLang $Version..." Yellow
Write-Host ""

try {
    $TempFile = Join-Path $env:TEMP "sundalang.exe"
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempFile
    
    Move-Item -Path $TempFile -Destination $BinaryPath -Force
    Write-Status "OK" "Hasil ngundeur" Green
    
} catch {
    Write-Host ""
    Write-Status "ERROR" "Gagal ngundeur binary" Red
    exit 1
}

Write-Host ""
Write-Divider
Write-Host ""
Write-Status "INFO" "Ngatur PATH..." Yellow
$UserPath = [Environment]::GetEnvironmentVariable("Path", "User")

if ($UserPath -notlike "*$InstallDir*") {
    $NewPath = "$UserPath;$InstallDir"
    [Environment]::SetEnvironmentVariable("Path", $NewPath, "User")
    Write-Status "OK" "PATH geus diupdate" Green
    Write-Host ""
    Write-Center "PENTING: Tutup tur buka deui terminal" Yellow
    Write-Center "pikeun nerapkeun parobahan PATH" Yellow
} else {
    Write-Status "OK" "PATH geus aya SundaLang" Green
}

Write-Host ""
Write-Host ""
Write-Box @(
    "",
    "INSTALASI SUKSES!",
    ""
) -BorderColor Green -TextColor Green
Write-Host ""

Write-Center "Kumaha carana make:" Cyan
Write-Host ""
Write-Center "1. Tutup terminal ieu tur buka deui" White
Write-Center "2. Jalankeun: sundalang --version" White
Write-Center "3. Jalankeun file .sl: sundalang namafile.sl" White
Write-Host ""
Write-Divider
Write-Host ""
Write-Center "Lokasi: $BinaryPath" DarkGray
Write-Host ""
Write-Center "Uninstall:" DarkGray
Write-Center "irm https://raw.githubusercontent.com/bromanprjkt/sundalang/main/uninstall.ps1 | iex" DarkGray
Write-Host ""
