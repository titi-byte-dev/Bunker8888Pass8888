# Roadmap AegisPass

> Identity, HR & Zero-Trust Vault — plataforma de gestão de identidade, acesso
> efémero, controlo financeiro de SaaS e RH para empresas em modelo **BYOD**
> (*Bring Your Own Device*), evoluindo para um **ERP + CRM com agentes de IA**.

Este `roadmap/` é o plano de trabalho navegável do projeto: fases, épicos,
*user journeys*, tasks e testes. A orientação rápida para humanos e IA está em
[`/AGENTS.md`](../../AGENTS.md).

## Como ler

A numeração indica a ordem de leitura sugerida (do "porquê" para o "como").

| # | Pasta / Ficheiro | O que contém |
|---|---|---|
| 00 | [`00-overview.md`](00-overview.md) | Visão, problema, proposta de valor, glossário |
| 01 | [`01-architecture/`](01-architecture/) | Stack, Zero-Knowledge, multi-tenancy, agentes, proxy Google |
| 02 | [`02-phases/`](02-phases/) | As 3 fases de desenvolvimento e marcos |
| 03 | [`03-epics/`](03-epics/) | Épicos por módulo (Vault, RH, FinTech, etc.) |
| 04 | [`04-user-journeys/`](04-user-journeys/README.md) | Fluxos passo-a-passo do utilizador (índice) |
| 05 | [`05-tasks/`](05-tasks/) | Backlog granular com IDs rastreáveis |
| 06 | [`06-testing/`](06-testing/) | Estratégia de testes, segurança e RGPD |
| 07 | [`07-non-functional/`](07-non-functional/) | Segurança, performance, conformidade |
| 08 | [`08-ai-tooling/`](08-ai-tooling/README.md) | Configuração de IA + melhorias futuras |
| 09 | [`09-design/`](09-design/README.md) | Visão UX/UI, design system, referências |
| 10 | [`10-production/`](10-production/README.md) | VPS, WireGuard, deploy, integrações reais |
| — | [`../../docs/`](../..) | **Documentação in-app** (product, concepts, developer, competitive) — ver `DOC-*` no backlog |

## Convenções

- **Idioma:** Português (Portugal) no conteúdo; termos técnicos e IDs em inglês.
- **IDs:** cada epic tem um prefixo (`VAULT`, `HR`, `FIN`, `SHARE`, `DW`, `MAIL`,
  `AGENT`, `UI`) e as tasks numeram-se a partir daí (`VAULT-001`). Ver [`05-tasks/backlog.md`](05-tasks/backlog.md).
- **Estilo didático:** o código e docs incluem comentários que ensinam a
  linguagem (ver regra em `.cursor/rules/didactic-style.mdc`).
- **Notas pessoais:** em [`_private/`](_private/) (fora do git por defeito).

## Estado (Jun 2026)

| Fase | Código | DoD dev local | Produção real |
|---|---|---|---|
| Fase 1 — Fundação | 🟢 maior parte no backlog | 🟡 parcial | INFRA VPS pendente |
| Fase 2 — CRM + 1º agente | 🟢 | 🟢 [`development-without-vps`](08-dev-environment/development-without-vps.md) | Google/Mail reais |
| Fase 3 — ERP + orquestrador | 🟢 | 🟢 idem | Open Banking TPP |

Matriz detalhada: [`10-production/implementation-status.md`](10-production/implementation-status.md)

> Legenda: ⚪ por iniciar · 🟡 em curso / stub · 🟢 concluído (no âmbito indicado)
