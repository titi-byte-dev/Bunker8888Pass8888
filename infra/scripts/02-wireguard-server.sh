#!/usr/bin/env bash
# INFRA-002 — Instala WireGuard e gera wg0.conf + cliente titi
set -euo pipefail

WG_PORT="${WG_PORT:-51820}"
WG_NET="${WG_NET:-10.8.0.0/24}"
WG_SERVER_IP="${WG_SERVER_IP:-10.8.0.1/24}"
CLIENT_NAME="${CLIENT_NAME:-titi}"
CLIENT_IP="${CLIENT_IP:-10.8.0.2/32}"

export DEBIAN_FRONTEND=noninteractive
apt-get update -qq
apt-get install -y -qq wireguard qrencode

install -d -m 700 /etc/wireguard
install -d -m 700 /etc/wireguard/clients

umask 077
wg genkey | tee /etc/wireguard/server.key | wg pubkey > /etc/wireguard/server.pub
wg genkey | tee "/etc/wireguard/clients/${CLIENT_NAME}.key" | wg pubkey > "/etc/wireguard/clients/${CLIENT_NAME}.pub"

SERVER_PRIV=$(cat /etc/wireguard/server.key)
SERVER_PUB=$(cat /etc/wireguard/server.pub)
CLIENT_PRIV=$(cat "/etc/wireguard/clients/${CLIENT_NAME}.key")
CLIENT_PUB=$(cat "/etc/wireguard/clients/${CLIENT_NAME}.pub")

# Interface pública (primeira não-loopback)
PUB_IF=$(ip -o -4 route show to default | awk '{print $5}' | head -1)
PUB_IP=$(curl -4 -s ifconfig.me || hostname -I | awk '{print $1}')

cat > /etc/wireguard/wg0.conf <<EOF
# AegisPass INFRA-002
[Interface]
Address = ${WG_SERVER_IP}
ListenPort = ${WG_PORT}
PrivateKey = ${SERVER_PRIV}
PostUp = iptables -A FORWARD -i wg0 -j ACCEPT; iptables -t nat -A POSTROUTING -o ${PUB_IF} -j MASQUERADE
PostDown = iptables -D FORWARD -i wg0 -j ACCEPT; iptables -t nat -D POSTROUTING -o ${PUB_IF} -j MASQUERADE

[Peer]
# ${CLIENT_NAME}
PublicKey = ${CLIENT_PUB}
AllowedIPs = ${CLIENT_IP}
EOF

cat > "/etc/wireguard/clients/${CLIENT_NAME}.conf" <<EOF
[Interface]
PrivateKey = ${CLIENT_PRIV}
Address = ${CLIENT_IP//\/32/"/24"}
DNS = 1.1.1.1

[Peer]
PublicKey = ${SERVER_PUB}
Endpoint = ${PUB_IP}:${WG_PORT}
AllowedIPs = ${WG_NET}
PersistentKeepalive = 25
EOF

chmod 600 /etc/wireguard/wg0.conf "/etc/wireguard/clients/${CLIENT_NAME}.conf"

# IP forwarding para NAT
sysctl -w net.ipv4.ip_forward=1
grep -q 'net.ipv4.ip_forward=1' /etc/sysctl.d/99-aegispass.conf 2>/dev/null || \
  echo 'net.ipv4.ip_forward=1' >> /etc/sysctl.d/99-aegispass.conf

systemctl enable wg-quick@wg0
systemctl start wg-quick@wg0

echo "OK: cliente em /etc/wireguard/clients/${CLIENT_NAME}.conf"
echo "Importa esse ficheiro no WireGuard do portátil."
