#!/usr/bin/env pwsh
# Run: Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass; .\install.ps1

$ErrorActionPreference = 'Stop'

$Org = 'jay7x'
$Repo = 'pct'
$AppName = 'pct'

function Install-Pct {
    [CmdletBinding(SupportsShouldProcess)]
    param()

    # Suppress progress bar in Windows PowerShell (Desktop edition)
    if ($PSVersionTable.PSEdition -eq 'Desktop') { $ProgressPreference = 'SilentlyContinue' }

    # Architecture: windows_amd64 or windows_arm64
    $arch = if ($env:PROCESSOR_ARCHITECTURE -match 'arm64|ARM64') { 'arm64' } else { 'amd64' }
    $dest = Join-Path $env:USERPROFILE '.puppetlabs' 'pct'

    # Resolve latest version from checksums.txt
    $namePattern = "${AppName}_[0-9.]+_windows_${arch}.zip"
    $text = Invoke-RestMethod "https://github.com/${Org}/${Repo}/releases/latest/download/checksums.txt"
    $line = ($text -split "`n") -match $namePattern | Select-Object -First 1
    if (-not $line) { throw "Could not find archive matching ${namePattern}" }
    $checksum, $filename = $line -split '\s+', 2
    $version = ($filename -split '_')[1]
    $url = "https://github.com/${Org}/${Repo}/releases/download/v${version}/${filename}"
    $tmp = Join-Path ([System.IO.Path]::GetTempPath()) $filename

    if ($PSCmdlet.ShouldProcess($dest, "Download and extract ${AppName} v${version}")) {
        try {
            Write-Host "Downloading and extracting ${AppName} v${version} to ${dest}..."
            Invoke-WebRequest $url -OutFile $tmp -UseBasicParsing

            $actual = (Get-FileHash $tmp -Algorithm SHA256).Hash.ToLower()
            if ($actual -ne $checksum) { throw "Checksum mismatch for $filename" }

            $null = New-Item -ItemType Directory -Path $dest -Force
            Expand-Archive $tmp -DestinationPath $dest -Force
        }
        finally {
            if (Test-Path $tmp) { Remove-Item $tmp -Force }
        }

        Write-Host 'Remember to add the pct app to your path:'
        Write-Host "`$env:Path += `"`$env:PATH;${dest}`""
    }
}

# Auto-run when executed directly (skip when dot-sourced)
if ($MyInvocation.InvocationName -ne '.') {
    Install-Pct
}
