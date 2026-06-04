# Diagramas — Visão de Sistema

Diagramas em [Mermaid](https://mermaid.js.org/) (renderizam no GitHub e no Cursor).

## Topologia geral (BYOD → VPN → VPS)

```mermaid
flowchart LR
    subgraph Dev["Dispositivo do funcionário (BYOD)"]
        APP["App AegisPass<br/>(Svelte + Capacitor)"]
    end
    subgraph Net["Perímetro"]
        FW["Firewall<br/>(só UDP WireGuard aberta)"]
    end
    subgraph VPS["VPS / Data Center privado"]
        API["Go API (Guardião)"]
        PG["PostgreSQL (RLS)"]
        MAIL["Relay SMTP / Aliases"]
        AG["Orquestrador de Agentes"]
    end
    Ext["Google Workspace / Open Banking / Breach Data API"]

    APP -->|túnel cifrado| FW --> API
    API --> PG
    API --> MAIL
    API --> AG
    API -->|proxy cifrado| Ext
```

## Camadas de defesa (Zero-Trust)

```mermaid
flowchart TB
    R1["1. Túnel WireGuard (rede)"]
    R2["2. Auth Hash + sessão (identidade)"]
    R3["3. Turno + Geofencing (contexto)"]
    R4["4. RLS multi-tenant (dados)"]
    R5["5. Zero-Knowledge (cifragem client-side)"]
    R1 --> R2 --> R3 --> R4 --> R5
```

## Fluxo de evento no ERP (Fase 3)

```mermaid
sequenceDiagram
    participant CRM as Agente CRM
    participant BUS as Event Bus
    participant FIN as Agente Financeiro
    participant OPS as Agente Operações
    participant HR as Agente RH
    CRM->>BUS: deal_closed (€5000)
    BUS->>FIN: gera fatura + link seguro
    FIN->>BUS: payment_detected (Open Banking)
    BUS->>OPS: aloca serviço / atualiza inventário
    BUS->>HR: calcula comissão (ficha cifrada)
```

> Mais diagramas específicos vivem junto do tema relevante (ex: o fluxo
> Zero-Knowledge está em [`../zero-knowledge.md`](../zero-knowledge.md)).
