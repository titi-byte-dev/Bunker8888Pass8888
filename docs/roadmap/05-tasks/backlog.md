# Backlog — Tasks Granulares

Cada task tem um **ID rastreável** (`PREFIXO-NNN`), uma estimativa relativa
(*story points* / tamanho) e o seu estado. Liga ao epic e à fase correspondentes.

> 💡 **Conceito — Story points:** estimativa de *esforço relativo* (não horas).
> Usamos a escala simples: **S** (pequeno), **M** (médio), **L** (grande),
> **XL** (muito grande, candidato a ser dividido).

> Legenda de estado: ⚪ todo · 🟡 em curso · 🟢 done · 🔵 bloqueado

## Como adicionar uma task

1. Escolhe o prefixo do epic (`VAULT`, `HR`, `FIN`, `SHARE`, `DW`, `MAIL`, `AGENT`, `INFRA`, `UI`).
2. Usa o próximo número livre.
3. Formato: `| ID | Descrição | Fase | Tamanho | Estado | Depende de |`

---

## INFRA — Infraestrutura & Setup

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| INFRA-001 | Provisionar VPS (Debian/Ubuntu) + hardening SSH (só chaves) | 1 | M | 🟡 | — |
| INFRA-002 | Configurar WireGuard (servidor) + firewall (só UDP aberta) | 1 | M | 🟡 | INFRA-001 |
| INFRA-003 | Docker + Docker Compose; esqueleto de serviços | 1 | M | 🟢 | INFRA-001 |
| INFRA-004 | PostgreSQL em container + backups cifrados | 1 | M | 🟢 | INFRA-003 |
| INFRA-005 | Pipeline CI (lint, test, build Go + Svelte) | 1 | M | 🟢 | — |
| INFRA-006 | Scaffold monorepo (`/backend` Go, `/frontend` Svelte, `/cli`) | 1 | S | 🟢 | — |
| CI-001 | Fix svelte-check (reconcile.test mocks) + actions @v5/Node 24 | 1 | S | 🟢 | INFRA-005 |

## VAULT — The Vault

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| VAULT-001 | Implementar Argon2id (derivação Master Key + Auth Hash) | 1 | M | 🟢 | INFRA-006 |
| VAULT-002 | Cifragem AES-GCM-256 (cliente) + test vectors | 1 | M | 🟢 | VAULT-001 |
| VAULT-003 | Modelo de dados de itens do cofre (blobs cifrados) | 1 | S | 🟢 | INFRA-004 |
| VAULT-004 | API Go: login por Auth Hash + sessões | 1 | M | 🟢 | VAULT-001 |
| VAULT-005 | CRUD de itens (logins, notas, cartões) | 1 | M | 🟢 | VAULT-003 |
| VAULT-006 | Sincronização em tempo real (WebSockets) | 1 | M | 🟢 | VAULT-005 |
| VAULT-007 | Gerador de palavras-passe | 1 | S | 🟢 | — |
| VAULT-008 | Score de higiene (fraca/reutilizada) calculado no cliente | 1 | M | 🟢 | VAULT-005 |
| VAULT-009 | 2FA / TOTP (RFC 6238) integrado | 1 | M | 🟢 | VAULT-005 |
| VAULT-010 | Acesso por turnos (validação NTP + expurgo de chave) | 1 | L | 🟢 | VAULT-004 |
| VAULT-011 | Geofencing (IP/GPS) na validação de acesso | 1 | M | 🟢 | VAULT-010 |
| VAULT-012 | Remote wipe de emergência (push WebSocket) | 1 | L | 🟢 | VAULT-006 |
| VAULT-013 | Browser isolado / sandbox (injeção sem revelar password) | 1 | XL | 🟢 | VAULT-005 |
| VAULT-014 | Suporte a passkeys | 1 | M | 🟢 | VAULT-004 |
| VAULT-015 | Importação de palavras-passe | 1 | M | 🟢 | VAULT-005 |
| VAULT-019 | Playground dev UI (login + cofre + import CSV) | 1 | M | 🟢 | VAULT-005 |
| VAULT-016 | Acesso de Emergência (herdeiro digital, período de espera) | 1 | L | 🟢 | VAULT-004 |
| VAULT-017 | CLI em Go (injeção de segredos via mTLS) | 1 | L | 🟢 | VAULT-004 |
| VAULT-018 | Chave de recuperação (mitiga perda de Master Password) | 1 | M | 🟢 | VAULT-002 |

## UI — Design & Experiência de Produto

> Visão completa: [`../09-design/product-ui-vision.md`](../09-design/product-ui-vision.md)

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| UI-001 | Design tokens (CSS vars, light/dark, tipografia, spacing) | 1 | M | 🟢 | VAULT-019 |
| UI-002 | App shell SvelteKit (routing, layouts, sidebar/tab bar) | 1 | L | 🟢 | UI-001 |
| UI-003 | Fluxos Auth (login, register, unlock, recovery, passkey) | 1 | L | 🟢 | UI-002, VAULT-014, VAULT-018 |
| UI-004 | Vault list + detail + CRUD (substituir playground) | 1 | L | 🟢 | UI-003, VAULT-005 |
| UI-005 | Motion system (GSAP presets, reduced motion) | 1 | M | 🟢 | UI-001 |
| UI-006 | Command palette (⌘K) | 1 | M | 🟢 | UI-002 |
| UI-007 | Painel higiene + preview Dark Web (score cliente) | 2 | M | 🟢 | UI-004, VAULT-008, DW-001 |
| UI-008 | Admin: utilizadores, políticas, audit, remote wipe | 1 | L | 🟢 | UI-002, VAULT-010, VAULT-012 |
| UI-009 | Shell mobile Capacitor + biometria | 1 | L | 🟡 | UI-004 |
| UI-010 | Catálogo componentes (Storybook ou página dev) | 1 | M | 🟢 | UI-001 |
| UI-011 | `ROUTE_TREE` + sidebar hierárquica + breadcrumbs | 1 | L | 🟢 | UI-002 |
| UI-012 | `lib/ui/` — PageShell, Panel, Button, Field, Eyebrow, EmptyState, HubLinks, StatusBanner, Breadcrumbs | 1 | M | 🟢 | UI-001, UI-010 |
| UI-013 | Paletas múltiplas (Aegis, Aurora, Midnight, Paper) + picker `/settings`; `data-palette` ⊥ `data-theme`; resolução por tenant (white-label); contraste AA testado | 1 | M | 🟢 | UI-001 |
| UI-014 | Hubs `/fin`, `/team` via `HubLinks` derivados da `ROUTE_TREE`; custos SaaS movidos p/ `/fin/costs` | 1 | S | 🟢 | UI-011 |
| UI-015 | DataTable, MetricCard, ListRow | 2 | M | 🟢 | UI-012 |
| UI-016 | Migração páginas para `PageShell` (`/fin/*` ✅, `/crm` ✅, hubs `/work` `/security` ✅; pendentes RH, mail, vault…) | 2 | L | 🟡 | UI-011, UI-012, UI-015, UI-017 |
| UI-017 | Toast, Skeleton, ConfirmDialog | 1 | M | 🟢 | UI-012 |

> Plano detalhado: [`../09-design/design-system.md`](../09-design/design-system.md)

## HR — Recursos Humanos & RGPD

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| HR-001 | Ficha de empregado com cifragem campo-a-campo | 1 | L | 🟢 | VAULT-002 |
| HR-002 | Logs imutáveis com hashing encadeado | 1 | M | 🟢 | INFRA-004 |
| HR-003 | Crypto-shredding + erasure RGPD (Art. 17) | 1 | L | 🟢 | HR-001 |
| HR-004 | Certificado criptográfico de eliminação | 1 | M | 🟢 | HR-003 |
| HR-005 | Gestão de contratos (object storage, chave por ficheiro) | 1 | L | 🟢 | VAULT-002 |
| HR-006 | Assinatura digital de contratos | 1 | L | 🟢 | HR-005 |
| HR-007 | Dashboard de onboarding (1 clique) | 1 | M | 🟢 | VAULT-005, MAIL-001 |
| HR-008 | Relatório de conformidade RGPD (PDF) | 2 | M | 🟢 | HR-002 |

## MAIL — Aliases & E-mail

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| MAIL-001 | Geração de aliases + reencaminhamento | 1 | M | 🟢 | INFRA-003 |
| MAIL-002 | Servidor de e-mail open-source (SMTP/IMAP) na VPS | 1 | L | 🟡 | INFRA-003 |
| MAIL-003 | SPF/DKIM/DMARC + domínio personalizado | 2 | M | 🟡 | MAIL-002 |
| MAIL-004 | Compor/iniciar e-mail a partir do alias (relay) | 2 | M | 🟢 | MAIL-002 |
| MAIL-005 | Rate limiting / anti open-relay | 2 | M | 🟢 | MAIL-002 |

## SHARE — Partilha

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| SHARE-001 | Chaves assimétricas por utilizador | 2 | M | 🟢 | VAULT-002 |
| SHARE-002 | Cofres partilhados (Shared Vaults) + permissões | 2 | L | 🟢 | SHARE-001 |
| SHARE-003 | Secret links efémeros (servidos via RAM) | 2 | L | 🟢 | SHARE-001 |
| SHARE-004 | Anexos cifrados por ficheiro | 2 | M | 🟢 | VAULT-002 |
| SHARE-005 | Notas temporárias auto-destrutivas | 2 | M | 🟢 | SHARE-003 |

## DW — Dark Web & Auditoria

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| DW-001 | Verificação de fugas via k-anonymity (breach data API) | 2 | M | 🟢 | VAULT-005 |
| DW-002 | Forçar alteração de password em exposição | 2 | S | 🟢 | DW-001 |
| DW-003 | Painel de saúde de segurança (score) | 2 | M | 🟢 | VAULT-008 |
| DW-004 | Sentinel Mode (deteção de login impossível) | 2 | L | 🟢 | VAULT-004 |

## FIN — FinTech

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| FIN-001 | Monitorização de custos SaaS (cruza com vault) | 2 | M | 🟢 | VAULT-005 |
| FIN-002 | Alertas de licenças sem uso | 2 | S | 🟢 | FIN-001 |
| FIN-003 | Integração Open Banking (mTLS) | 3 | XL | 🟡 | INFRA-002 |
| FIN-004 | Cartões virtuais efémeros | 3 | XL | ⚪ | FIN-003 |
| FIN-005 | Categorização fiscal automática | 3 | M | 🟢 | FIN-001, AGENT-002 |
| FIN-006 | Faturação (pro-forma/fatura/recibo + numeração legal) | 3 | M | 🟢 | CRM-001 |
| FIN-007 | Comissões de vendas sobre faturas pagas | 3 | M | 🟢 | FIN-006 |

## CRM — Pipeline de vendas

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| CRM-001 | Modelo leads cifrados + API CRUD | 2 | M | 🟢 | VAULT-005 |
| CRM-002 | UI funil básico (`/crm`) | 2 | M | 🟢 | CRM-001, UI-002 |
| CRM-003 | Fluxo `deal_closed` → pro-forma a partir de lead ganho | 3 | S | 🟢 | CRM-001, FIN-006 |
| CRM-004 | Conversão pro-forma → fatura definitiva | 3 | S | 🟢 | FIN-006, CRM-003 |

## AGENT — Agentes de IA

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| AGENT-001 | Interface `Tool` + validação (reimplementada de raiz) | 2 | M | 🟢 | INFRA-006 |
| AGENT-002 | Guardião (decifra só o necessário, menor privilégio) | 2 | L | 🟢 | VAULT-002 |
| AGENT-003 | Primeiro agente: prospeção (lê e-mails → leads) | 2 | L | 🟢 | AGENT-001, MAIL-001 |
| AGENT-004 | Event Bus (channels Go / NATS) | 3 | L | 🟢 | INFRA-003 |
| AGENT-005 | Orquestrador multi-agente | 3 | XL | 🟢 | AGENT-004 |
| AGENT-006 | Agente Financeiro (reconciliação) | 3 | L | 🟡 | AGENT-005, FIN-003 |
| AGENT-007 | Agente RH (recrutamento às cegas + onboarding) | 3 | L | 🟢 | AGENT-005, HR-007 |
| AGENT-008 | Agente Operações (compras/inventário) | 3 | L | 🟢 | AGENT-005 |
| AGENT-009 | Human-in-the-loop p/ ações sensíveis | 3 | M | 🟢 | AGENT-005 |
| AGENT-010 | Mitigação de prompt injection (dados ≠ instruções) | 2 | M | 🟢 | AGENT-001 |

## DOC — Documentação in-app

> SSOT em `docs/`; a app lê JSON gerado por `scripts/build-docs.mjs`.

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| DOC-001 | Convenção Markdown + script `build-docs.mjs` + manifest | 1 | M | 🟢 | — |
| DOC-002 | UI `/settings/docs` com níveis e dropdowns didáticos | 1 | L | 🟢 | DOC-001, UI-002 |
| DOC-003 | Conteúdo product + concepts (glossary, vault, security…) | 1 | M | 🟢 | DOC-001 |
| DOC-004 | Conteúdo developer (crypto, API) + competitive landscape | 1 | M | 🟢 | DOC-001 |
| DOC-005 | Journeys do roadmap expostos na app (04-user-journeys) | 1 | S | 🟢 | DOC-002 |
| DOC-006 | Command palette + link em Definições | 1 | S | 🟢 | DOC-002 |
| DOC-007 | Journeys em falta (SHARE, DW, emergency…) | 2 | M | 🟢 | DOC-005 |
| DOC-008 | Renderer Mermaid nos percursos | 2 | S | 🟢 | DOC-005 |
| DOC-009 | FlowPlayer animado (passo-a-passo + reproduzir) | 2 | M | 🟢 | DOC-008 |
| DOC-010 | Pesquisa full-text + glossário inline nas páginas | 2 | M | 🟢 | DOC-002 |
| DOC-011 | POC Svelte Flow (`@xyflow/svelte`) com highlight por passo | 2 | M | 🟢 | DOC-009 |
| DOC-012 | Converter sequenceDiagram → graph no build-docs | 2 | S | 🟢 | DOC-011 |
| DOC-013 | Links contextuais «Como funciona?» nas rotas da app | 2 | S | 🟢 | DOC-002 |
| DOC-014 | Docs produção + matriz dev/prod + índice journeys | 2 | S | 🟢 | DOC-001 |

## GOOGLE — Layer de Proxy (transversal)

| ID | Descrição | Fase | Tamanho | Estado | Depende de |
|---|---|---|---|---|---|
| GOOGLE-001 | OAuth2 + Service Account (Workspace) | 2 | L | 🟡 | INFRA-006 |
| GOOGLE-002 | Cifragem ZK de ficheiros Drive/Docs | 2 | L | 🟡 | VAULT-002, GOOGLE-001 |
| GOOGLE-003 | Data masking dinâmico em Sheets (regex + tokens) | 2 | L | 🟡 | GOOGLE-001 |
| GOOGLE-004 | Relay Gmail com aliases + PGP | 2 | M | ⚪ | MAIL-004, GOOGLE-001 |
