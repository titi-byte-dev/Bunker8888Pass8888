# Smoke test MAIL-002 - simula Postfix pipe -> POST /api/mail/ingest (dev local)
# Uso: .\scripts\smoke-mail-ingest.ps1 -AliasAddress "abc123@aegis.email"
# Nota: confirma que nada mais ocupa a porta 8080 (docker compose backend).
param(
    [string]$BaseUrl = "http://127.0.0.1:8080",
    [string]$Secret = $env:AEGIS_MAIL_WEBHOOK_SECRET,
    [Parameter(Mandatory = $true)]
    [string]$AliasAddress,
    [string]$From = "prospect@example.com",
    [string]$Subject = "Smoke test AegisPass ingest"
)

if (-not $Secret) {
    $envFile = Join-Path (Join-Path $PSScriptRoot "..") ".env"
    if (Test-Path $envFile) {
        $line = Get-Content $envFile | Where-Object { $_ -match '^AEGIS_MAIL_WEBHOOK_SECRET=' } | Select-Object -First 1
        if ($line) { $Secret = ($line -split '=', 2)[1].Trim() }
    }
}
if (-not $Secret) { $Secret = "dev-mail-webhook-secret" }

$payload = @{
    message_id = [guid]::NewGuid().ToString("N")
    from       = $From
    to         = @($AliasAddress)
    subject    = $Subject
    body       = "Mensagem de teste gerada por smoke-mail-ingest.ps1"
} | ConvertTo-Json -Compress

$url = "$BaseUrl/api/mail/ingest?secret=$Secret"
Write-Host "POST $url (alias: $AliasAddress)"

try {
    $resp = Invoke-WebRequest -Uri $url -Method POST -ContentType "application/json" -Body $payload -UseBasicParsing
    Write-Host "HTTP $($resp.StatusCode): $($resp.Content)"
    if ($resp.StatusCode -eq 201) {
        Write-Host "OK: email ingerido - verifica /mail inbox e evento mail.inbox.received"
    }
    exit 0
}
catch {
    $status = $null
    $errBody = $_.Exception.Message
    if ($_.Exception.Response) {
        $status = [int]$_.Exception.Response.StatusCode
        $stream = $_.Exception.Response.GetResponseStream()
        if ($stream) {
            $reader = New-Object System.IO.StreamReader($stream)
            $errBody = $reader.ReadToEnd()
            $reader.Close()
        }
    }
    Write-Host "HTTP ${status}: $errBody"
    if ($status -eq 200 -and $errBody -match "ignored") {
        Write-Host "AVISO: alias desconhecido - cria o alias em /mail antes do teste"
        exit 0
    }
    exit 1
}
