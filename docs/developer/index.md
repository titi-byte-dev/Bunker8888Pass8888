---
title: Guia do programador
slug: developer-index
category: developer
order: 1
audience: [developer]
layer: [fullstack]
level: 2
in_app: true
summary: Estrutura do monorepo, convenções e onde começar a contribuir.
related: [developer-crypto, developer-api, glossary]
---

:::summary
Monorepo com **backend Go**, **frontend SvelteKit**, **CLI** e **docs** como fonte única.
Lê `AGENTS.md` e o roadmap antes de decisões de arquitectura.
:::

:::level{level=2 title="Estrutura do repositório"}
```
backend/     API Go, migrações SQL, pkg/crypto
frontend/    SvelteKit, WebCrypto, vault UI
cli/         Injeção mTLS de segredos
docs/        Documentação SSOT (esta app lê daqui)
scripts/     build-docs.mjs, backups
```
:::

:::level{level=2 title="Convenções obrigatórias"}
- **Zero-Knowledge:** nunca enviar Master Key ao servidor
- **Multi-tenant:** filtrar `tenant_id`; confiar na RLS como rede de segurança
- **Cripto:** só `crypto/*` (Go) e WebCrypto (browser) — nunca inventar algoritmos
- **Didático:** comentários que explicam o *porquê* (ver `.cursor/rules/`)
:::

:::level{level=3 title="Comandos úteis"}
| Comando | Efeito |
|---|---|
| `make backend-test` | Testes Go incl. cripto |
| `make frontend-test` | Vitest no frontend |
| `make test` | Ambos |
| `node scripts/build-docs.mjs` | Regenera docs na app |
| `docker compose up` | Ambiente local completo |
:::

:::level{level=3 title="Backlog e épicos"}
Tasks com IDs (`VAULT-001`, `UI-004`) em `docs/roadmap/05-tasks/backlog.md`.
Épicos em `docs/roadmap/03-epics/`.
:::
