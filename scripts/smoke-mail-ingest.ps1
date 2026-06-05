# Smoke test MAIL-002 — simula Postfix pipe → POST /api/mail/ingest (dev local)
# Uso: .\scripts\smoke-mail-ingest.ps1 -AliasAddress "abc123@aegis.email"
param(
    [string]$BaseUrl = "http://127.0.0.1:8080",
    [string]$Secret = $env:AEGIS_MAIL_WEBHOOK_SECRET,
    [Parameter(Mandatory = $true)]
    [string]$AliasAddress,
    [string]$From = "prospect@example.com",
    [string]$Subject = "Smoke test AegisPass ingest"
)

if (-not $Secret) {
    $envFile = Join-Path $PSScriptRoot ".." ".env"
    if (Test-Path $envFile) {
        $line = Get-Content $envFile | Where-Object { $_ -match '^AEGIS_MAIL_WEBHOOK_SECRET=' } | Select-Object -First 1
        if ($line) { $Secret = ($line -split '=', 2)[1].Trim() }
    }
}
if (-not $Secret) { $Secret = "dev-mail-webhook-secret" }

$body = @{
    message_id = [guid]::NewGuid().ToString("N")
    from       = $From
    to         = @($AliasAddress)
    subject    = $Subject
    body       = "Mensagem de teste gerada por smoke-mail-ingest.ps1"
} | ConvertTo-Json -Compress

Write-Host "→ POST $BaseUrl/api/mail/ingest (alias: $AliasAddress)"

try {
    $resp = Invoke-WebRequest -Uri "$BaseUrl/api/mail/ingest?secret=$Secret" `
        -Method POST -ContentType "application/json" -Body $body -UseBasicParsing
    Write-Host "← $($resp.StatusCode) $($resp.Content)"
    if ($resp.StatusCode -eq 201) {
        Write-Host "OK: e-mail ingerido — verifica /mail inbox e evento mail.inbox.received"
    }
} catch {
    $status = $_.Exception.Response.StatusCode.value__
    $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
    $errBody = $reader.ReadToEnd()
    Write-Host "← $status $errBody"
    if ($status -eq 200 -and $errBody -match "ignored") {
        Write-Host "AVISO: alias desconhecido — cria o alias em /mail antes do teste"
    }
    exit 1
}
