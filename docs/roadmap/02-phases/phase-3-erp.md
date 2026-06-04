# Fase 3 — Fecho do Círculo Operacional (O ERP)

> **Objetivo:** transformar o AegisPass no "sistema operativo" da empresa —
> faturação, contabilidade, banca e um **orquestrador multi-agente** autónomo.

## Âmbito (in scope)

- Módulo de Faturação e Contabilidade
- Integração bancária (Open Banking) + cartões virtuais efémeros — ver [`epic-fintech`](../03-epics/epic-fintech.md)
- **Orquestrador de Agentes (multi-agente)** em Go — ver [`epic-agents`](../03-epics/epic-agents.md)
- "Modo Autónomo": empresa define regras e os agentes operam rotinas diárias
- Agentes especializados: CRM, Financeiro, RH, Operações/Logística

## Modelo de orquestração

Ver [`agents-architecture`](../01-architecture/agents-architecture.md) para o
detalhe técnico (EDA, Guardião, segurança com LLMs).

```mermaid
flowchart TB
    REGRAS["Regras de negócio definidas pela empresa"]
    REGRAS --> ORQ["Orquestrador (Go)"]
    ORQ --> A1["Agente CRM"]
    ORQ --> A2["Agente Financeiro"]
    ORQ --> A3["Agente RH"]
    ORQ --> A4["Agente Operações"]
    A1 & A2 & A3 & A4 --> FEED["Dashboard Svelte<br/>(centro de comando, WebSockets)"]
```

## Definition of Done (Fase 3)

- [ ] Faturação completa (pro-forma → fatura → recibo) com numeração legal.
- [ ] Reconciliação bancária automática via Open Banking.
- [ ] Orquestrador executa o fluxo "deal_closed → fatura → pagamento → comissão"
      ponta-a-ponta com aprovação humana nos pontos críticos.
- [ ] Guardião comprovadamente nunca expõe a Master Key aos agentes (auditoria).
- [ ] Relatório de conformidade RGPD gerado automaticamente (PDF).

## Riscos principais

| Risco | Mitigação |
|---|---|
| Agente de IA toma decisão errada/cara | "Human-in-the-loop" obrigatório em ações financeiras; limites por política |
| LLM expõe dados sensíveis | Guardião com menor privilégio; modelos locais isolados |
| Complexidade de ERP descontrola scope | Entregar módulo a módulo; cada um com valor isolado |

## Visão de longo prazo (pós-Fase 3)

- Mais agentes verticais (jurídico, marketing, suporte).
- Marketplace de skills/conectores.
- Selo "AegisPass Compliant" + seguro cibernético integrado.
