# AGENTS.md — Orientação para assistentes de IA

> Este ficheiro é lido por assistentes de IA (Cursor e outros) no início de cada
> sessão. Dá o contexto mínimo para trabalhar bem no projeto. Mantém-no curto.

## O que é este projeto

**AegisPass** — plataforma de gestão de identidade, acesso efémero, controlo
financeiro de SaaS e RH para empresas em modelo **BYOD**, com cifragem
**Zero-Knowledge** e conformidade **RGPD por desenho**. Roadmap de longo prazo:
evoluir para um **ERP + CRM com agentes de IA**.

A documentação completa (visão, arquitetura, fases, épicos, journeys, tasks,
testes) está em [`docs/roadmap/`](docs/roadmap/README.md). **Lê-a antes de
decisões de arquitetura.**

## Stack

- **Backend:** Go (Golang) · **Frontend:** Svelte + TypeScript
- **BD:** PostgreSQL (com Row-Level Security) · **Rede:** WireGuard
- **Infra:** Docker / Docker Compose · **Runtime cripto cliente:** WebCrypto

## Regras de ouro (não-negociáveis)

1. 🔒 **Zero-Knowledge:** o servidor **nunca** vê dados sensíveis em claro nem a
   Master Key. Cifrar/decifrar acontece no cliente.
2. 🚫 **Nunca inventar criptografia.** Usar só bibliotecas padrão auditadas
   (`crypto/*` em Go, WebCrypto no browser). Nonce único por cifragem.
3. 🤐 **Nunca commitar segredos** (`.env`, `*.pem`, `*.key`, IBANs/credenciais
   reais). O `.gitignore` já bloqueia, mas confirma sempre.
4. 🧱 **Multi-tenant:** isolar dados por `tenant_id` + RLS. Testar isolamento.
5. 📚 **Estilo didático:** comentários que ensinam a linguagem (ver
   `.cursor/rules/didactic-style.mdc`).
6. ⚖️ **Nunca copiar código proprietário de terceiros.** Inspirar em padrões
   conhecidos, mas reimplementar tudo de raiz.

## Idioma

Conteúdo e explicações em **Português (Portugal)**; termos técnicos, nomes de
código e IDs em **inglês** (`goroutine`, `nonce`, `VAULT-001`).

## Como trabalhar aqui

- Antes de uma feature, vê o epic e as tasks em
  [`docs/roadmap/03-epics/`](docs/roadmap/03-epics/) e
  [`docs/roadmap/05-tasks/backlog.md`](docs/roadmap/05-tasks/backlog.md).
- Código novo traz testes (ver [`docs/roadmap/06-testing/`](docs/roadmap/06-testing/)).
- Notas pessoais vão para `docs/roadmap/_private/` (fora do git).

## Dev local vs produção

| Ambiente | Guia |
|---|---|
| Sem VPS (stubs) | [`docs/roadmap/08-dev-environment/development-without-vps.md`](docs/roadmap/08-dev-environment/development-without-vps.md) |
| VPS + integrações reais | [`docs/roadmap/10-production/README.md`](docs/roadmap/10-production/README.md) |
| Matriz dev/prod por task | [`implementation-status.md`](docs/roadmap/10-production/implementation-status.md) |

Documentação in-app: `docs/product/` + `docs/developer/` → `npm run docs:build`.

## Configuração de IA neste repositório

Ver [`docs/roadmap/08-ai-tooling/`](docs/roadmap/08-ai-tooling/README.md) para a
filosofia (tool-agnostic) e a lista de melhorias futuras. Resumo do que existe:

| Mecanismo | Onde | Para quê |
|---|---|---|
| Rules | `.cursor/rules/*.mdc` | Convenções (Go, Svelte, cripto, RGPD, didático) |
| Commands | `.cursor/commands/*.md` | Ações rápidas (`/new-task`, `/security-review`, …) |
| Hooks | `.cursor/hooks.json` | Gates de segurança automáticos |
| Skills | `.cursor/skills/*` | Workflows reutilizáveis |
| MCP | (utilizador) | Integrações externas |
