# Arquitectura — Event Bus de agentes (AGENT-004)

> 💡 **Conceito — Event Bus:** canal onde componentes **publicam** factos
> (`mail.inbox.received`) e outros **subscrevem** reacções, sem acoplamento
> directo. Em Go começamos com `chan Event`; NATS para cluster multi-nó.

## Tipos de evento (v1)

| Tipo | Origem | Payload (metadados) |
|---|---|---|
| `mail.inbox.received` | Webhook Mailpit | `inbox_id`, `alias`, `relayed` |
| `crm.prospection.run` | POST prospection | `draft_count` |
| `agent.tool.executed` | POST tool/run | `tool`, `agent_id`, `success` |

> ⚠️ **Segurança:** payloads sem PII de leads — blobs CRM nunca no bus.

## Fluxo publish/subscribe (sequence)

```mermaid
sequenceDiagram
    participant P as Publicador
    participant Bus as Event Bus
    participant DB as agent_events
    participant H1 as Subscritor A
    participant H2 as Subscritor B

    P->>Bus: Publish(type, user_id, payload)
    Bus->>DB: Append (auditoria)
    Bus->>H1: dispatch async
    Bus->>H2: dispatch async
    H1-->>Bus: ok / log erro
```

## Pipeline mail → CRM (flowchart)

```mermaid
flowchart LR
    WH[Webhook mail] --> P1[Publish mail.inbox.received]
    P1 --> FE[Feed /api/agent/events]
    P1 --> FUT[Subscritores futuros AGENT-005]
    U[Utilizador] --> PR[POST prospection/run]
    PR --> P2[Publish crm.prospection.run]
    P2 --> FE
```

## Ciclo de vida de um evento (stateDiagram)

```mermaid
stateDiagram-v2
    [*] --> Publicado: Publish()
    Publicado --> Persistido: PGStore.Append
    Publicado --> NaFila: chan buffer
    NaFila --> Entregue: dispatch handlers
    Entregue --> [*]
    NaFila --> Descartado: buffer cheio
    Descartado --> [*]
```

## Evolução (AGENT-005)

```mermaid
flowchart TB
    subgraph Hoje["AGENT-004"]
        CH[channels Go]
        PG[(PostgreSQL)]
    end
    subgraph Futuro["Escala"]
        NATS[NATS JetStream]
        ORQ[Orquestrador]
    end
    CH --> PG
    CH -.-> NATS
    NATS --> ORQ
```

Ver também [`agents-architecture.md`](agents-architecture.md) e
[`../04-user-journeys/journey-agent-event-feed.md`](../04-user-journeys/journey-agent-event-feed.md).
