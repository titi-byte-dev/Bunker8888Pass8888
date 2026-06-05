# FIN-003 — Open Banking real (mTLS)

> **Objetivo:** activar provider PSD2 real em vez do mock de dev, com certificados
> mTLS e variáveis `AEGIS_OB_*` já suportadas pelo backend.

> 💡 **Conceito:** *mTLS* (mutual TLS) exige que **cliente e servidor** apresentem
> certificados — o banco/TPP confirma que és um participante registado.

## Pré-requisitos

- INFRA-002 (WireGuard) — chamadas ao TPP só a partir da VPS ou túnel
- Conta junto de um *TPP* (Third Party Provider) PSD2 registado
- Certificados eIDAS QWAC / QSeal (conforme jurisdição)

## Variáveis (já no código)

Definidas em `backend/internal/openbanking/mtls_config.go`:

| Variável | Conteúdo |
|---|---|
| `AEGIS_OB_BASE_URL` | URL base da API PSD2 do TPP |
| `AEGIS_OB_CLIENT_CERT` | Certificado cliente PEM |
| `AEGIS_OB_CLIENT_KEY` | Chave privada PEM |
| `AEGIS_OB_CA_CERT` | CA do TPP (opcional) |

Se `AEGIS_OB_BASE_URL` estiver vazio, o servidor usa **só** `mock` provider (dev).

## 1. Instalar certificados na VPS

```bash
sudo mkdir -p /run/secrets/ob
sudo install -m 600 ob-client.pem /run/secrets/ob/client.pem
sudo install -m 600 ob-client.key /run/secrets/ob/client.key
sudo install -m 644 ob-ca.pem /run/secrets/ob/ca.pem
```

No `.env`:

```
AEGIS_OB_BASE_URL=https://api.tpp-registado.eu/psd2/v1
AEGIS_OB_CLIENT_CERT=/run/secrets/ob/client.pem
AEGIS_OB_CLIENT_KEY=/run/secrets/ob/client.key
AEGIS_OB_CA_CERT=/run/secrets/ob/ca.pem
```

## 2. Fluxo na UI

1. `/fin/banking` → **Ligar banco** (OAuth redirect do TPP, se aplicável)
2. **Sync** → `POST /api/fin/banking/sync` — débitos entram cifrados
3. Agente financeiro reconcilia com subscrições SaaS

Journey: [`journey-finance-agent-reconcile.md`](../04-user-journeys/journey-finance-agent-reconcile.md)

## 3. Verificação

```bash
# Com túnel WireGuard activo
curl -s http://10.8.0.1:8080/api/fin/banking/status -H "Cookie: ..."
# provider: "mtls" (não "mock")
```

Testes de integração com sandbox do TPP **antes** de ligar contas reais.

> ⚠️ **Segurança:** chaves mTLS são credenciais de sistema — rotação periódica,
> nunca em logs, permissões 600, backup offline cifrado.
