# Arquitectura — Pipeline de e-mail (MAIL)

Visão técnica do percurso SMTP → inbox → relay → CRM, com vários tipos de
diagrama para aprendizagem.

## Componentes

| Componente | Papel |
|---|---|
| **Mailpit** (dev) / **Postfix** (prod) | Recebe SMTP na porta 25/1025 |
| **Webhook** | Notifica o backend (`POST /api/mail/webhook/mailpit`) |
| **IngestService** | Resolve alias, grava inbox, dispara relay |
| **RelayService** | SMTP de saída para `destination` do alias |
| **InboxRepo** | `mail_inbox_messages` por `owner_id` |
| **Prospection** (AGENT-003) | inbox → rascunhos → CRM cifrado |

## Pipeline completo (flowchart)

```mermaid
flowchart LR
    subgraph Entrada
        SMTP[SMTP recebido]
        WH[Webhook]
    end
    subgraph Backend
        ING[Ingest]
        RL[Relay MAIL-004]
        INB[(inbox BD)]
        LIM[Rate limit MAIL-005]
    end
    subgraph Saida
        DST[E-mail real destination]
        CRM[Agente prospection]
    end

    SMTP --> WH --> LIM --> ING
    ING --> INB
    ING --> RL --> DST
    INB --> CRM
```

## Relay vs compose (MAIL-004)

```mermaid
flowchart TB
    subgraph Inbound["Entrada (automático)"]
        I1[E-mail para alias] --> I2[Ingest + inbox]
        I2 --> I3[Relay → destination]
    end
    subgraph Outbound["Saída (utilizador)"]
        O1[POST /api/mail/compose] --> O2[Verifica dono do alias]
        O2 --> O3[SMTP From: alias@aegis.email]
        O3 --> O4[Destinatário externo]
    end
```

## Anti open-relay (MAIL-005)

```mermaid
flowchart TD
    R[Pedido SMTP] --> A{Destino é alias @aegis.email registado?}
    A -->|não| X[Rejeitar — não somos open relay]
    A -->|sim| B{Rate limit OK?}
    B -->|não| Y[429 Too Many Requests]
    B -->|sim| C[Processar ingest/relay]
```

> ⚠️ **Segurança:** só aceitamos entrega para aliases que existem na BD do
> tenant. Relay de saída exige sessão autenticada + propriedade do alias.

## Ligação ao epic

Ver [`../03-epics/epic-aliases-email.md`](../03-epics/epic-aliases-email.md) e
journeys em [`../04-user-journeys/`](../04-user-journeys/).
