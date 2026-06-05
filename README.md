# AegisPass

> Identity, HR & Zero-Trust Vault — uma plataforma de gestão de identidade,
> acesso efémero, controlo financeiro de SaaS e Recursos Humanos para empresas
> em modelo **BYOD** (*Bring Your Own Device*), com cifragem **Zero-Knowledge** e
> conformidade **RGPD por desenho**.

## Visão

Muitas empresas não conseguem dar um computador a cada colaborador, mas precisam
de garantir que os dados corporativos ficam isolados e seguros nos dispositivos
pessoais. O AegisPass é a ponte segura que separa a vida profissional da pessoal
no mesmo aparelho — e que, a longo prazo, evolui para um **ERP + CRM com agentes
de IA**.

## Pilares

- 🔒 **Zero-Knowledge / E2EE** — o servidor nunca vê dados em claro nem a chave.
- 🧱 **Zero-Trust & multi-tenant** — isolamento por empresa, acesso validado por
  rede, identidade, turno e geofencing.
- 📜 **RGPD por desenho** — cifragem campo-a-campo, direito ao esquecimento e
  logs imutáveis.
- ⚡ **Performance** — backend concorrente em Go, frontend leve em Svelte.

## Stack

| Camada | Tecnologia |
|---|---|
| Backend | Go (Golang) |
| Frontend | Svelte + TypeScript |
| Base de dados | PostgreSQL (Row-Level Security) |
| Rede | WireGuard |
| Infraestrutura | Docker / Docker Compose |

100% open-source, sem licenças proprietárias.

## Documentação

### Na app (utilizador e programador)

Em **Definições → Documentação** (`/settings/docs`): conteúdo didático por níveis
(Essencial / Intermédio / Técnico), conceitos em dropdowns, fonte única em:

- [`docs/product/`](docs/product/) — como usar funcionalidades
- [`docs/concepts/`](docs/concepts/) — glossário embutido
- [`docs/developer/`](docs/developer/) — front, back, API
- [`docs/competitive/`](docs/competitive/) — panorama e ideias para o roadmap
- [`docs/roadmap/04-user-journeys/`](docs/roadmap/04-user-journeys/) — percursos

Regenerar após editar Markdown: `npm run docs:build` (raiz) ou `make docs-build`.

### Roadmap interno

Plano completo em [`docs/roadmap/`](docs/roadmap/README.md). Para IA e contribuidores:
[`AGENTS.md`](AGENTS.md).

## Estado

🟡 **Fase 1 em curso** — cofre Zero-Knowledge, shell UI, RH, partilha e
documentação in-app. Ver [`02-phases/phase-1-foundation.md`](docs/roadmap/02-phases/phase-1-foundation.md).

## Licença

A definir.
