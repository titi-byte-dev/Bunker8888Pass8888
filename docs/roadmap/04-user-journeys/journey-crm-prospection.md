# Journey: Prospeção automática — e-mail → lead cifrado (AGENT-003)

> **Ator:** Comercial / gestor de vendas · **Epics:** `AGENT`, `CRM`, `MAIL`

O agente de prospeção lê e-mails pendentes na inbox, gera rascunhos de lead e o
utilizador importa para o funil CRM com cifragem Zero-Knowledge no browser.

## Pré-condições

- Cofre **desbloqueado** (Master Key em memória).
- Mensagens na inbox (`/mail`) com estado **pendente**.
- Agente `prospection` autorizado pelo Guardião (AGENT-002).

## Fluxo principal (sequence)

```mermaid
sequenceDiagram
    participant U as Utilizador
    participant CRM as /crm
    participant API as Go API
    participant AG as Agente prospection
    participant DB as PostgreSQL
    participant MK as Master Key (cliente)

    U->>CRM: Correr prospeção
    CRM->>API: POST /api/agent/prospection/run
    API->>DB: inbox WHERE processed_at IS NULL
    loop Por mensagem
        API->>AG: draft_lead_from_email
        AG-->>API: rascunho JSON (sanitizado)
    end
    API-->>CRM: drafts[]

    U->>CRM: Importar rascunho
    CRM->>MK: encryptLead(payload)
    MK-->>CRM: blob AES-GCM
    CRM->>API: POST /api/crm/leads { blob }
    CRM->>API: POST /api/mail/inbox/{id}/processed
    API-->>CRM: lead criado
```

## Fluxograma do Guardião e tools

```mermaid
flowchart TD
    A[POST prospection/run] --> B{agent_id = prospection?}
    B -->|não| X[403 Forbidden]
    B -->|sim| C[list_mail_inbox + draft_lead]
    C --> D{Permissões OK?}
    D -->|mail:read_metadata + crm:write_lead_draft| E[Executar tools]
    D -->|não| X
    E --> F{Body com prompt injection?}
    F -->|sim| G[Ignorar mensagem]
    F -->|não| H[Rascunho no JSON de resposta]
    G --> E
```

## Diagrama de estados (lead)

```mermaid
stateDiagram-v2
    [*] --> Rascunho: agente gera draft
    Rascunho --> Cifrado: utilizador importa
    Cifrado --> Novo: stage = new no funil
    Novo --> Contactado: move estágio
    Contactado --> Qualificado
    Qualificado --> Proposta
    Proposta --> Ganho
    Proposta --> Perdido
```

## Passo-a-passo

1. Garante e-mails na inbox (SMTP para alias ou simulação em dev).
2. Abre `/crm` com cofre desbloqueado.
3. Clica **Correr prospeção** — o servidor devolve rascunhos (não grava PII em claro no CRM).
4. Revê cada rascunho (e-mail, assunto, notas sanitizadas).
5. **Importar para o funil** — o cliente cifra e grava o blob.
6. A mensagem de inbox fica **processada**; nova prospeção não a repete.

## Conceito didático

> 💡 **Function calling:** o agente não escreve na BD directamente — pede tools
> validadas (`list_mail_inbox`, `draft_lead_from_email`) com schemas JSON.

> ⚠️ **Segurança:** conteúdo de e-mail passa por `SanitizeExternalContent` e
> `RejectIfLooksLikeInstruction` (AGENT-010). Leads no CRM são sempre blobs opacos.

## Fluxos alternativos

- **Cofre bloqueado:** importação falha — só metadados de rascunho visíveis até unlock.
- **Sem e-mails pendentes:** prospeção devolve lista vazia.
- **Rate limit (MAIL-005):** ingest bloqueada antes de chegar à inbox.

## Pós-condições

- Lead no funil CRM (cifrado).
- Entrada em `GET /api/agent/audit` para `prospection_run`.
- Mensagem de inbox marcada processada.
