#!/usr/bin/env bash
# INFRA-001 — Hardening SSH em Debian/Ubuntu
# Executar como root na VPS nova. Requer DEPLOY_USER e DEPLOY_PUBKEY no ambiente.
set -euo pipefail

DEPLOY_USER="${DEPLOY_USER:-titi}"
DEPLOY_PUBKEY="${DEPLOY_PUBKEY:-}"

if [[ -z "${DEPLOY_PUBKEY}" ]]; then
  echo "ERRO: define DEPLOY_PUBKEY com o conteúdo da chave pública (.pub)" >&2
  exit 1
fi

export DEBIAN_FRONTEND=noninteractive
apt-get update -qq
apt-get upgrade -y -qq
apt-get install -y -qq fail2ban unattended-upgrades sudo

# Utilizador de deploy (não-root) com sudo sem password — só para automação controlada.
if ! id "${DEPLOY_USER}" &>/dev/null; then
  useradd -m -s /bin/bash "${DEPLOY_USER}"
fi
usermod -aG sudo "${DEPLOY_USER}"
echo "${DEPLOY_USER} ALL=(ALL) NOPASSWD:ALL" > "/etc/sudoers.d/99-${DEPLOY_USER}"
chmod 440 "/etc/sudoers.d/99-${DEPLOY_USER}"

install -d -m 700 -o "${DEPLOY_USER}" -g "${DEPLOY_USER}" "/home/${DEPLOY_USER}/.ssh"
echo "${DEPLOY_PUBKEY}" > "/home/${DEPLOY_USER}/.ssh/authorized_keys"
chmod 600 "/home/${DEPLOY_USER}/.ssh/authorized_keys"
chown -R "${DEPLOY_USER}:${DEPLOY_USER}" "/home/${DEPLOY_USER}/.ssh"

# sshd: só chaves, sem root por password
install -d /etc/ssh/sshd_config.d
cat > /etc/ssh/sshd_config.d/99-aegispass.conf <<'EOF'
# AegisPass INFRA-001 — gerado por 01-ssh-hardening.sh
PasswordAuthentication no
KbdInteractiveAuthentication no
ChallengeResponseAuthentication no
PermitRootLogin prohibit-password
PubkeyAuthentication yes
MaxAuthTries 3
X11Forwarding no
AllowTcpForwarding yes
EOF

systemctl enable fail2ban unattended-upgrades
systemctl restart ssh
systemctl restart fail2ban

echo "OK: utilizador ${DEPLOY_USER} criado. Testa login noutra sessão antes de sair."
