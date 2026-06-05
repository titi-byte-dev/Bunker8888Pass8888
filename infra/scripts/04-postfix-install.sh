#!/usr/bin/env bash
# MAIL-002 — Instala Postfix + pipe aegis-ingest na VPS
# Executar como root (ou sudo). Lê segredos de /etc/default/aegis-mail ou ambiente.
#
# Uso:
#   AEGIS_MAIL_DOMAIN=mail.seudominio.com \
#   AEGIS_MAIL_WEBHOOK_SECRET="$(grep AEGIS_MAIL_WEBHOOK_SECRET /opt/aegispass/.env | cut -d= -f2)" \
#   ./04-postfix-install.sh
#
# Pré-requisito: backend Docker a escutar em 127.0.0.1:8080 (05-deploy-stack.sh)
set -euo pipefail

AEGIS_MAIL_DOMAIN="${AEGIS_MAIL_DOMAIN:-aegis.email}"
AEGIS_INGEST_URL="${AEGIS_INGEST_URL:-http://127.0.0.1:8080/api/mail/ingest}"
AEGIS_MAIL_WEBHOOK_SECRET="${AEGIS_MAIL_WEBHOOK_SECRET:-}"
REPO_DIR="${REPO_DIR:-/opt/aegispass}"

if [[ -z "${AEGIS_MAIL_WEBHOOK_SECRET}" ]]; then
  if [[ -f "${REPO_DIR}/.env" ]]; then
    # shellcheck disable=SC1091
    AEGIS_MAIL_WEBHOOK_SECRET="$(grep -E '^AEGIS_MAIL_WEBHOOK_SECRET=' "${REPO_DIR}/.env" | cut -d= -f2- | tr -d '\r' || true)"
  fi
fi

if [[ -z "${AEGIS_MAIL_WEBHOOK_SECRET}" ]]; then
  echo "ERRO: define AEGIS_MAIL_WEBHOOK_SECRET (igual ao .env do backend)" >&2
  exit 1
fi

export DEBIAN_FRONTEND=noninteractive
apt-get update -qq
apt-get install -y -qq postfix mailutils curl python3

# Variáveis do pipe — lidas pelo aegis-ingest.sh via ambiente do master process
cat > /etc/default/aegis-mail <<EOF
# Gerado por 04-postfix-install.sh — NÃO commitar
AEGIS_INGEST_URL=${AEGIS_INGEST_URL}
AEGIS_MAIL_WEBHOOK_SECRET=${AEGIS_MAIL_WEBHOOK_SECRET}
EOF
chmod 600 /etc/default/aegis-mail

install -m 755 "${REPO_DIR}/infra/postfix/aegis-ingest.sh" /usr/local/bin/aegis-ingest.sh

# Wrapper garante env mesmo que Postfix não propague /etc/default/aegis-mail
cat > /usr/local/bin/aegis-ingest-wrapper.sh <<'WRAP'
#!/usr/bin/env bash
set -euo pipefail
# shellcheck disable=SC1091
[[ -f /etc/default/aegis-mail ]] && source /etc/default/aegis-mail
exec /usr/local/bin/aegis-ingest.sh "$@"
WRAP
chmod 755 /usr/local/bin/aegis-ingest-wrapper.sh

# Transport map: todo o tráfego @domínio → pipe
cat > /etc/postfix/aegis-ingest.cf <<EOF
# Domínio de aliases AegisPass (MAIL-002)
${AEGIS_MAIL_DOMAIN}    aegis-ingest:
EOF
postmap /etc/postfix/aegis-ingest.cf

# Fragmento Postfix — include no fim de main.cf
cat > /etc/postfix/aegispass.cf <<EOF
# AegisPass MAIL-002 — gerado por 04-postfix-install.sh
relay_domains = ${AEGIS_MAIL_DOMAIN}
transport_maps = hash:/etc/postfix/aegis-ingest.cf
# ⚠️ Segurança: só aceita destinos autorizados (relay_domains + mydestination)
smtpd_recipient_restrictions = permit_mynetworks, reject_unauth_destination
EOF

if ! grep -q 'aegispass.cf' /etc/postfix/main.cf 2>/dev/null; then
  echo "" >> /etc/postfix/main.cf
  echo "# AegisPass" >> /etc/postfix/main.cf
  echo "include /etc/postfix/aegispass.cf" >> /etc/postfix/main.cf
fi

# Transport pipe no master.cf (idempotente)
if ! grep -q '^aegis-ingest' /etc/postfix/master.cf; then
  cat >> /etc/postfix/master.cf <<'EOF'

# AegisPass — entrega via HTTP ingest (MAIL-002)
aegis-ingest unix  -       n       n       -       -       pipe
  flags=R user=nobody argv=/usr/local/bin/aegis-ingest-wrapper.sh ${recipient}
EOF
fi

postfix check
systemctl enable postfix
systemctl reload postfix

echo "OK: Postfix configurado para ${AEGIS_MAIL_DOMAIN}"
echo "     Teste: swaks --to ALIAS@${AEGIS_MAIL_DOMAIN} --server 127.0.0.1"
echo "     Próximo: DNS MX/SPF/DKIM — ver mail-002-003-production.md"
