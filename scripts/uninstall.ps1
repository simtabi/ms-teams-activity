# Remove mta on Windows: stops + removes the logon task / service and deletes the
# binary (delegates to `mta self uninstall`).
#
#   irm https://raw.githubusercontent.com/simtabi/ms-teams-activity/main/scripts/uninstall.ps1 | iex
#   ./scripts/uninstall.ps1 -Purge    # also delete config and runtime data

param([switch]$Purge)

$ErrorActionPreference = "Stop"

$cmd = Get-Command mta -ErrorAction SilentlyContinue
$mta = if ($cmd) { $cmd.Source } else { $null }
if (-not $mta) {
    $cand = Join-Path "$env:LOCALAPPDATA\Programs\mta" "mta.exe"
    if (Test-Path $cand) { $mta = $cand }
}
if (-not $mta) { throw "mta not found on PATH or in the usual install dir." }

Write-Host "Removing mta (service + binary)..."
if ($Purge) { & $mta self uninstall --yes --purge }
else { & $mta self uninstall --yes }
Write-Host "Done."
