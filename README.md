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

O plano completo do projeto vive em [`docs/roadmap/`](docs/roadmap/README.md):

- **Visão & glossário** — [`00-overview`](docs/roadmap/00-overview.md)
- **Arquitetura** — [`01-architecture`](docs/roadmap/01-architecture/)
- **Fases** — [`02-phases`](docs/roadmap/02-phases/)
- **Épicos** — [`03-epics`](docs/roadmap/03-epics/)
- **User journeys** — [`04-user-journeys`](docs/roadmap/04-user-journeys/)
- **Backlog** — [`05-tasks/backlog.md`](docs/roadmap/05-tasks/backlog.md)
- **Testes** — [`06-testing`](docs/roadmap/06-testing/)
- **Requisitos não-funcionais** — [`07-non-functional`](docs/roadmap/07-non-functional/)
- **AI tooling** — [`08-ai-tooling`](docs/roadmap/08-ai-tooling/README.md)

Para assistentes de IA e novos contribuidores, ler primeiro [`AGENTS.md`](AGENTS.md).

## Estado

🟡 **Planeamento.** O código ainda não foi iniciado; o foco atual é a Fundação
(Fase 1) — ver [`02-phases/phase-1-foundation.md`](docs/roadmap/02-phases/phase-1-foundation.md).

## Licença

A definir.
