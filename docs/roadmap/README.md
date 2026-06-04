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
| 04 | [`04-user-journeys/`](04-user-journeys/) | Fluxos passo-a-passo do utilizador |
| 05 | [`05-tasks/`](05-tasks/) | Backlog granular com IDs rastreáveis |
| 06 | [`06-testing/`](06-testing/) | Estratégia de testes, segurança e RGPD |
| 07 | [`07-non-functional/`](07-non-functional/) | Segurança, performance, conformidade |
| 08 | [`08-ai-tooling/`](08-ai-tooling/README.md) | Configuração de IA + melhorias futuras |

## Convenções

- **Idioma:** Português (Portugal) no conteúdo; termos técnicos e IDs em inglês.
- **IDs:** cada epic tem um prefixo (`VAULT`, `HR`, `FIN`, `SHARE`, `DW`, `MAIL`,
  `AGENT`) e as tasks numeram-se a partir daí (`VAULT-001`). Ver [`05-tasks/backlog.md`](05-tasks/backlog.md).
- **Estilo didático:** o código e docs incluem comentários que ensinam a
  linguagem (ver regra em `.cursor/rules/didactic-style.mdc`).
- **Notas pessoais:** em [`_private/`](_private/) (fora do git por defeito).

## Estado

| Fase | Estado |
|---|---|
| Fase 1 — Fundação & Identidade | 🟡 Planeamento |
| Fase 2 — Pipeline de Vendas (CRM) | ⚪ Por iniciar |
| Fase 3 — Operações (ERP + Agentes) | ⚪ Por iniciar |

> Legenda: ⚪ por iniciar · 🟡 em planeamento/curso · 🟢 concluído
