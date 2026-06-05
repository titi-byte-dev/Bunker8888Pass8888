# Arquitectura — Orquestrador multi-agente (AGENT-005)

Coordena workers que reagem ao Event Bus sem acoplamento directo entre módulos.

> 💡 **Conceito — Orquestrador:** regista agentes (`Worker`), subscreve tipos de
> evento e delega `Handle`. Não substitui o Guardião — só encadeia reacções seguras.

## Workers v1

| Worker | Escuta | Publica |
|---|---|---|
| `prospection` | `mail.inbox.received` | `orchestrator.action.suggested` (`run_prospection`) |

> ⚠️ **Human-in-the-loop:** prospeção **não** corre automaticamente — o cliente
> precisa da Master Key. O worker só *sugere* a acção no feed.

## Fluxo mail → sugestão (sequence)

```mermaid
sequenceDiagram
    participant M as Mail ingest
    participant Bus as Event Bus
    participant Orq as Orchestrator
    participant PW as Worker prospection
    participant U as Utilizador /crm

    M->>Bus: mail.inbox.received
    Bus->>PW: dispatch
    PW->>Bus: orchestrator.action.suggested
    U->>U: vê sugestão no feed
    U->>U: Correr prospeção (manual)
```

## Registo de workers (flowchart)

```mermaid
flowchart TD
    START[Arranque servidor] --> RUN[agentBus.Run]
    RUN --> ORC[orchestrator.Start]
    ORC --> LOOP[Para cada Worker]
    LOOP --> SUB[bus.Subscribe handles]
    SUB --> READY[Pronto para eventos]
```

## Estados de uma sugestão (stateDiagram)

```mermaid
stateDiagram-v2
    [*] --> Sugerida: orchestrator.action.suggested
    Sugerida --> Ignorada: utilizador não actua
    Sugerida --> Executada: run_prospection manual
    Executada --> [*]
    Ignorada --> [*]
```

## API

```
GET /api/agent/orchestrator/status  → lista workers registados
GET /api/agent/events               → feed inclui sugestões
```

Evolução: mais workers (FIN, RH), filas persistentes, NATS (AGENT-004 roadmap).
