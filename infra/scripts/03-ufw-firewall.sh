#!/usr/bin/env bash
# INFRA-002 — UFW: deny por omissão; abre WireGuard (+ SMTP opcional)
set -euo pipefail

WG_PORT="${WG_PORT:-51820}"
ALLOW_SMTP=false
SSH_VIA_WG_ONLY=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --wg-port) WG_PORT="$2"; shift 2 ;;
    --allow-smtp) ALLOW_SMTP=true; shift ;;
    --ssh-via-wg-only) SSH_VIA_WG_ONLY=true; shift ;;
    *) echo "Opção desconhecida: $1" >&2; exit 1 ;;
  esac
done

export DEBIAN_FRONTEND=noninteractive
apt-get install -y -qq ufw

ufw --force reset
ufw default deny incoming
ufw default allow outgoing

# WireGuard (única porta UDP obrigatória na Internet pública)
ufw allow "${WG_PORT}/udp" comment 'WireGuard'

if [[ "${ALLOW_SMTP}" == "true" ]]; then
  ufw allow 25/tcp comment 'SMTP inbound'
  ufw allow 587/tcp comment 'SMTP submission'
fi

if [[ "${SSH_VIA_WG_ONLY}" == "true" ]]; then
  # SSH só na interface wg0 — requer túnel activo para administrar
  ufw allow in on wg0 to any port 22 proto tcp comment 'SSH via VPN'
else
  ufw allow 22/tcp comment 'SSH bootstrap'
fi

ufw --force enable
ufw status verbose

echo "OK: firewall activo. Com --ssh-via-wg-only, confirma SSH via 10.8.0.1 antes de fechar sessão pública."
