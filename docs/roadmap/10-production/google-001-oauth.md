# GOOGLE-001 — OAuth2 + Service Account (Workspace)

> **Objetivo:** substituir o stub `/work/google-dev` por integração real com Google
> Workspace (Drive, Sheets, Gmail relay — GOOGLE-002/003/004).

> 💡 **Conceito:** *Service Account* é uma identidade de máquina (não humana) que
> acede a APIs Google com delegação de domínio (*domain-wide delegation*).

## Pré-requisitos

- Google Workspace (não conta Gmail pessoal)
- Acesso de super-admin ao Admin Console
- INFRA-006 já concluído (monorepo com módulos Google no frontend)

## 1. Criar projecto GCP

1. [Google Cloud Console](https://console.cloud.google.com/) → novo projecto `aegispass-prod`
2. APIs & Services → Enable:
   - Google Drive API
   - Google Sheets API
   - Gmail API (para GOOGLE-004)

## 2. Service Account

1. IAM → Service Accounts → Create
2. Nome: `aegispass-workspace`
3. Criar chave JSON → guardar em cofre ( **nunca** no git )
4. Na VPS: `/run/secrets/google-sa.json` (perm 600, root)

## 3. Domain-wide delegation

1. Admin Console → Security → API controls → Domain-wide delegation
2. Client ID da service account → scopes:

```
https://www.googleapis.com/auth/drive
https://www.googleapis.com/auth/spreadsheets
https://www.googleapis.com/auth/gmail.send
```

> ⚠️ **Segurança:** princípio do menor privilégio — remove scopes que não fores
> usar na primeira entrega.

## 4. Variáveis de ambiente (produção)

Adiciona ao `.env` (ver `.env.production.example`):

```
AEGIS_GOOGLE_SA_JSON=/run/secrets/google-sa.json
AEGIS_GOOGLE_DELEGATED_USER=admin@suaempresa.com
AEGIS_GOOGLE_ENABLED=true
```

## 5. Substituir stub no frontend

| Dev | Produção |
|---|---|
| `/work/google-dev` | `/work/google` (a implementar com GOOGLE-002/003) |
| `driveDevStore.ts` (localStorage) | Drive API + blobs ZK do cliente |
| `masking.ts` local | Sheets API + tokens de masking |

A cifragem ZK **mantém-se no cliente** — a service account só vê blobs opacos
e metadados não sensíveis, igual ao modelo Vault.

## 6. Verificação

1. Lista ficheiros Drive com utilizador delegado
2. Upload de blob cifrado (test vector VAULT-002)
3. Abrir Sheet com masking activo — células sensíveis substituídas por tokens

Journey de dev (referência): [`journey-google-dev-stub.md`](../04-user-journeys/journey-google-dev-stub.md)
