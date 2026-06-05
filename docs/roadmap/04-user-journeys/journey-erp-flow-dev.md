# Journey: Fluxo ERP em desenvolvimento (DoD Fase 3)

> **Ator:** Gestor comercial/financeiro · **Epics:** CRM, FIN, HR, AGENT

Ciclo completo com orquestrador e human-in-the-loop — **sem** banco nem VPS real.

## Pré-condições

- Master Key desbloqueada.
- Backend com orquestrador activo (`docker compose`).

## Fluxo principal

```mermaid
sequenceDiagram
    participant U as Utilizador
    participant CRM as /crm
    participant INV as /fin/invoices
    participant COM as /fin/commissions
    participant HR as /hr/compliance
    participant Orq as Orquestrador

    U->>CRM: lead → Ganho
    CRM->>Orq: crm.deal_closed
    Orq->>CRM: issue_proforma suggested
    U->>CRM: Aprovar
    CRM->>INV: emite PF (ZK)

    U->>INV: Converter em fatura + Marcar pago
    INV->>Orq: fin.invoice.paid
    Orq->>INV: calculate_commission suggested
    U->>INV: Aprovar
    INV->>COM: cria comissão

    COM->>Orq: hr.compliance.requested
    Orq->>COM: generate_rgpd_report suggested
    U->>COM: Aprovar
    U->>HR: PDF conformidade
```

## Passo-a-passo

1. **CRM** — cria lead; move para **Ganho**; aprova **emitir pro-forma**.
2. **Faturas** — converte PF→FT; preenche linhas se necessário; **Marcar pago**.
3. Aprova **calcular comissão** (10% vendedor por omissão em dev).
4. **Comissões** — confirma registo; aprova **relatório RGPD**.
5. **Conformidade** — imprime PDF (`Ctrl+P` → Guardar como PDF).

## Recibo (RC)

Em `/fin/invoices`, em fatura **paga**, clica **Emitir recibo**.

## Reconciliação (opcional)

`/fin/banking` → consentimento mock → sync → reconcilia com subscrições SaaS.

## Pós-condições

- Eventos auditados em `agent_events`.
- Nenhuma acção financeira automática sem aprovação humana.
