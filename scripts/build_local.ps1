$ErrorActionPreference = 'Stop'

$repo = Split-Path -Parent $MyInvocation.MyCommand.Path | Split-Path -Parent
Set-Location $repo

Write-Host '==> Build frontend'
Push-Location frontend
try {
  if (Test-Path package-lock.json) {
    npm ci
  } else {
    npm install
  }
} catch {
  Write-Host 'npm install failed. If this is EPERM/locked file, close any running Vite/dev server or antivirus scan, then retry.' -ForegroundColor Yellow
  if (Test-Path node_modules) {
    try {
      Remove-Item -Recurse -Force node_modules
    } catch {
      Write-Host 'Failed to remove node_modules (maybe locked). Please close processes and retry.' -ForegroundColor Yellow
      throw
    }
  }
  npm install
}

try {
  npm run build
} catch {
  Write-Host 'Build failed. If you see "run-p not recognized", ensure npm install completed successfully.' -ForegroundColor Yellow
  throw
}
Pop-Location

Write-Host '==> Build Go binary'
$dist = Join-Path $repo 'dist'
New-Item -ItemType Directory -Force -Path $dist | Out-Null
$exe = Join-Path $dist 'wrzapi.exe'

go build -trimpath -ldflags "-s -w" -o $exe ./cmd/server

Write-Host "Done: $exe"
