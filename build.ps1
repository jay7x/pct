#!/usr/bin/env pwsh

[CmdletBinding()]
param (
  [Parameter()]
  [ValidateSet('build', 'quick', 'package')]
  [string]
  $Target = 'build'
)
$arch = go env GOHOSTARCH
$platform = go env GOHOSTOS
$binPath = Join-Path $PSScriptRoot "dist" "pct_${platform}_${arch}"

$amd64 = go env GOAMD64
if ($amd64) {
	$binPath = "${binPath}_${amd64}"
}

switch ($Target) {
  'build' {
    # Set goreleaser to build for current platform only
    goreleaser build --snapshot --clean --single-target
    git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/baker-round (Join-Path $binPath "templates")
		Get-ChildItem -Path (Join-Path $binPath "templates")
  }
  'quick' {
    If ($Env:OS -match '^Windows') {
      go build -o "$binPath/pct.exe"
    } else {
      go build -o "$binPath/pct"
    }
  }
  'package' {
    git clone -b main --depth 1 --single-branch https://github.com/puppetlabs/baker-round "templates"
    goreleaser --skip-publish --snapshot --clean
  }
}
