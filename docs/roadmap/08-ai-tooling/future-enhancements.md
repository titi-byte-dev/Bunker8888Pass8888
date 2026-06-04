# AI Tooling — Listagem de Melhorias Futuras

Inventário do que **poderíamos ter mais** para tornar o desenvolvimento assistido
por IA melhor. Inspirado em padrões conhecidos de agentic tooling, reimplementado
de raiz e tool-agnostic. Serve de backlog desta área.

> Legenda: 🟢 alto valor / baixo esforço · 🟡 médio · 🔵 estratégico (mais tarde)

## A. Regras adicionais

| ID | Melhoria | Valor |
|---|---|---|
| AIT-001 | Regra de convenções de **commits** (Conventional Commits) | 🟢 |
| AIT-002 | Regra de **SQL/migrations** (sempre reversíveis, RLS em tabelas novas) | 🟢 |
| AIT-003 | Regra de **API design** (versionamento, erros, paginação) | 🟡 |
| AIT-004 | Regra de **testes** (cobertura mínima por área) | 🟡 |
| AIT-005 | Regra de **acessibilidade** (WCAG) no frontend | 🔵 |

## B. Comandos adicionais

| ID | Melhoria | Valor |
|---|---|---|
| AIT-010 | `/commit` — gera mensagem no formato Conventional Commits | 🟢 |
| AIT-011 | `/review` — code review completo segundo as regras do projeto | 🟢 |
| AIT-012 | `/rgpd-check` — valida tratamento de dados pessoais num diff | 🟡 |
| AIT-013 | `/new-journey` — cria um user journey com diagrama Mermaid | 🟡 |
| AIT-014 | `/threat-model` — gera/atualiza STRIDE para uma feature | 🔵 |

## C. Hooks adicionais

| ID | Melhoria | Valor | Nota |
|---|---|---|---|
| AIT-020 | **Formatação automática** (`gofmt`/`prettier`) em `afterFileEdit` | 🟢 | requer toolchain instalado |
| AIT-021 | **Secret scanner por script** (ex: gitleaks) em vez de prompt | 🟡 | mais determinístico; requer binário |
| AIT-022 | `beforeSubmitPrompt` — avisar se o prompt do utilizador contém segredos | 🟡 | |
| AIT-023 | `afterShellExecution` — auditar output de comandos sensíveis | 🔵 | |
| AIT-024 | `subagentStop` — encadear loops de revisão automática | 🔵 | |

## D. Skills adicionais

| ID | Melhoria | Valor |
|---|---|---|
| AIT-030 | `scaffold-svelte-component` (componente + store + teste) | 🟢 |
| AIT-031 | `write-migration` (migration PostgreSQL reversível + RLS) | 🟢 |
| AIT-032 | `add-agent-tool` (nova Tool para um agente, com schema + validação) | 🟡 |
| AIT-033 | `e2e-flow` (esqueleto de teste Playwright para um journey) | 🟡 |
| AIT-034 | `incident-runbook` (gera runbook de resposta a incidente) | 🔵 |

## E. Integrações (MCP)

| ID | Melhoria | Valor |
|---|---|---|
| AIT-040 | MCP de **PostgreSQL** (explorar schema/queries em dev) | 🟡 |
| AIT-041 | MCP de **gestão de issues** (sincronizar backlog ↔ tracker) | 🟡 |
| AIT-042 | MCP de **observabilidade** (consultar métricas/logs em debug) | 🔵 |

## F. O "nosso próprio agente" (visão de longo prazo)

> 💡 Ligação ao produto: o AegisPass terá agentes (Fase 3). A experiência de
> configurar agentes de *desenvolvimento* (esta secção) informa o desenho dos
> agentes de *produto*.

| ID | Melhoria | Valor |
|---|---|---|
| AIT-050 | Especificação tool-agnostic dos nossos agentes (formato próprio) | 🔵 |
| AIT-051 | Camada de **permissões/Guardião** partilhada dev↔produto | 🔵 |
| AIT-052 | Catálogo de **Tools** reutilizáveis (schema único) | 🔵 |
| AIT-053 | Telemetria/auditoria de ações de agentes (logs imutáveis) | 🔵 |

## Como priorizar

1. Fazer primeiro os 🟢 (alto valor, baixo esforço): AIT-001, AIT-002, AIT-010,
   AIT-011, AIT-030, AIT-031.
2. Adicionar hooks determinísticos (AIT-020/021) **quando** o toolchain existir.
3. Os 🔵 estratégicos acompanham a Fase 3 (agentes de produto).
