# ============================================================================
# backup-postgres.ps1 — backup cifrado do PostgreSQL (INFRA-004)
# ----------------------------------------------------------------------------
# Uso (PowerShell, na raiz do repo):
#   $env:AEGIS_BACKUP_KEY = "..."
#   .\scripts\backup-postgres.ps1
# ============================================================================

$ErrorActionPreference = "Stop"

$root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $root

if (-not $env:AEGIS_BACKUP_KEY) {
    Write-Error "AEGIS_BACKUP_KEY nao definida. Gera com: make backup-gen-key"
}

$user = if ($env:POSTGRES_USER) { $env:POSTGRES_USER } else { "aegis" }
$db   = if ($env:POSTGRES_DB)   { $env:POSTGRES_DB }   else { "aegis" }

New-Item -ItemType Directory -Force -Path "backups" | Out-Null
$stamp = Get-Date -Format "yyyyMMdd-HHmmss"
$out   = (Resolve-Path "backups").Path + "\aegis-$stamp.sql.enc"

Write-Host "A fazer pg_dump de $db..."
Push-Location backend
try {
    docker compose exec -T db pg_dump -U $user $db | go run ./cmd/backup encrypt -o $out
    if ($LASTEXITCODE -ne 0) { throw "backup falhou" }
} finally {
    Pop-Location
}

Write-Host "Backup cifrado: $out"
