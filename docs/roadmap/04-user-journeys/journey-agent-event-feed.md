# Journey: Feed de acções dos agentes (AGENT-004)

> **Ator:** Utilizador / supervisor · **Epics:** `AGENT`, `MAIL`, `CRM`

Supervisiona o que os agentes e o pipeline de mail fazem — sem formulários
monolíticos, através de um **feed de eventos** auditável.

## Pré-condições

- Sessão activa.
- Event Bus ligado (arranque com base de dados).

## Fluxo principal

```mermaid
sequenceDiagram
    participant M as Mailpit
    participant API as API Go
    participant Bus as Event Bus
    participant CRM as /crm
    participant U as Utilizador

    M->>API: webhook e-mail
    API->>Bus: mail.inbox.received
    U->>CRM: abre feed de eventos
    CRM->>API: GET /api/agent/events
    API-->>CRM: lista cronológica
    U->>CRM: Correr prospeção
    API->>Bus: crm.prospection.run
```

## Mapa mental do feed (flowchart)

```mermaid
flowchart TD
    E[Evento no feed] --> T{Tipo?}
    T -->|mail.inbox.received| A[Novo e-mail — corre prospeção?]
    T -->|crm.prospection.run| B[Rascunhos gerados — importar?]
    T -->|agent.tool.executed| C[Tool usada — ver auditoria]
    A --> CRM[/crm prospeção/]
    B --> CRM
    C --> AUD[/api/agent/audit/]
```

## Passo-a-passo

1. Recebe e-mail no alias → evento `mail.inbox.received` no feed.
2. Abre `/crm` — painel «Actividade dos agentes» mostra eventos recentes.
3. Clica **Correr prospeção** → novo evento `crm.prospection.run`.
4. Importa leads — funil actualiza (cifragem local).
5. Opcional: `GET /api/agent/audit` para detalhe de tools.

## Pós-condições

- Histórico em `agent_events` (append-only).
- Base para orquestrador AGENT-005 (subscritores automáticos).
