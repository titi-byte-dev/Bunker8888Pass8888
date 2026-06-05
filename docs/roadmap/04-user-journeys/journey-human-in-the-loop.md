# Journey: Aprovar sugestão do orquestrador (AGENT-009)

> **Ator:** Utilizador comercial · **Epics:** `AGENT`, `CRM`

O feed em `/crm` mostra sugestões com botões **Aprovar** / **Rejeitar**. Só após
aprovação é que a prospeção corre no cliente (Master Key necessária).

## Pré-condições

- Orquestrador activo (AGENT-005).
- E-mail ingerido ou simulado → sugestão `run_prospection` no feed.

## Fluxo principal

```mermaid
sequenceDiagram
    participant U as Utilizador
    participant CRM as /crm
    participant API as API
    participant Bus as Event Bus

    CRM->>API: GET events (sugestão pending)
    U->>CRM: Aprovar
    CRM->>API: POST .../approve
    API->>Bus: orchestrator.action.approved
    CRM->>API: POST prospection/run
    Note over CRM: rascunhos + import ZK
```

## Passo-a-passo

1. Simula ou recebe e-mail no alias.
2. Abre `/crm` — vês **Sugestão: correr prospeção** com botões.
3. Clica **Aprovar** — decisão fica auditada.
4. A prospeção corre automaticamente (precisas do cofre desbloqueado).
5. Importa rascunhos para o funil.

Alternativa: **Rejeitar** — sugestão fica marcada como rejeitada, sem prospeção.

## Pós-condições

- Eventos `approved` / `rejected` em `agent_events`.
- Nenhuma tool ZK executada sem clique humano.
