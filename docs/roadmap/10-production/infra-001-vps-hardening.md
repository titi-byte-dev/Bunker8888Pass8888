# INFRA-001 — VPS + hardening SSH

> **Objetivo:** servidor Debian 12 ou Ubuntu 24.04 LTS acessível **apenas** por chave
> SSH, com utilizador não-root para deploy.

> 💡 **Conceito:** *hardening* reduz a superfície de ataque — desactiva serviços e
> métodos de login que um atacante poderia abusar (password SSH, root directo).

## 1. Escolher e provisionar VPS

| Critério | Recomendação |
|---|---|
| SO | Debian 12 ou Ubuntu 24.04 LTS |
| RAM | ≥ 4 GB (PostgreSQL + Postfix + backend) |
| Disco | ≥ 40 GB SSD |
| Região | UE (RGPD) — ex. Frankfurt, Paris |

Regista o IP público como `AEGIS_VPS_IP` (só para bootstrap; depois o acesso
administrativo passa pelo WireGuard — ver INFRA-002).

## 2. Chave SSH local (no teu portátil)

```bash
# Ed25519 — curva moderna, chave curta e segura
ssh-keygen -t ed25519 -C "titi@aegispass" -f ~/.ssh/aegispass_deploy
```

Copia a chave pública para a VPS (primeiro login, ainda com password do painel):

```bash
ssh-copy-id -i ~/.ssh/aegispass_deploy.pub root@AEGIS_VPS_IP
```

## 3. Executar script de hardening

Na VPS, como `root`:

```bash
git clone https://github.com/titi-byte-dev/Bunker8888Pass8888.git /opt/aegispass
cd /opt/aegispass/infra/scripts
chmod +x 01-ssh-hardening.sh
DEPLOY_USER=titi DEPLOY_PUBKEY="$(cat ~/.ssh/aegispass_deploy.pub)" ./01-ssh-hardening.sh
```

O script:

1. Actualiza pacotes (`apt upgrade`)
2. Cria utilizador `titi` com `sudo` (sem password no sudoers — só chave)
3. Instala `fail2ban`, `unattended-upgrades`
4. Configura `sshd`: `PasswordAuthentication no`, `PermitRootLogin prohibit-password`
5. Activa só autenticação por chave pública

> ⚠️ **Segurança:** abre **outra** sessão SSH e confirma login com `titi` **antes**
> de fechar a sessão root. Se falhar, ainda podes reverter via consola do provider.

## 4. Verificação

```bash
ssh -i ~/.ssh/aegispass_deploy titi@AEGIS_VPS_IP
sudo whoami   # deve responder "root"
grep PasswordAuthentication /etc/ssh/sshd_config.d/99-aegispass.conf
# PasswordAuthentication no
```

## 5. Próximo passo

[`infra-002-wireguard-firewall.md`](infra-002-wireguard-firewall.md) — depois disso o SSH
público pode ficar restrito ao interface WireGuard.
