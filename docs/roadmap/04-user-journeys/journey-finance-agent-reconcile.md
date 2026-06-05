# Journey: Agente financeiro reconcilia pagamentos (FIN-003 / AGENT-006)

> **Ator:** Responsável financeiro · **Epics:** `AGENT`, `FIN`

Liga uma conta bancária (mock em dev), sincroniza movimentos e cruza débitos com
subscrições SaaS. O orquestrador sugere reconciliação — sem alterar dados automaticamente.

> 💡 **Conceito:** Open Banking — regulação PSD2 que permite acesso a movimentos
> com consentimento do titular, via APIs mTLS entre TPP e banco.

## Pré-condições

- Orquestrador activo com worker `finance` (AGENT-005).
- Subscrições em `/fin` com custos mensais definidos.
- Master Key desbloqueada.

## Fluxo principal

```mermaid
sequenceDiagram
    participant U as Utilizador
    participant OB as /fin/banking
    participant API as API
    participant Orq as Orquestrador

    U->>OB: Simular consentimento PSD2
    OB->>API: POST /api/fin/banking/connect
    U->>OB: Sincronizar movimentos
    OB->>API: POST /api/fin/banking/sync
    OB->>OB: reconcileTransactions() no cliente
    OB->>API: POST /api/agent/finance/report-sync
    Note over API: só contagens (matched/unmatched)
    API->>Orq: fin.transactions.synced
    Orq->>API: orchestrator.action.suggested
    Note over OB: «Sugestão: reconciliar pagamentos»
    U->>OB: Aprovar (AGENT-009)
    OB->>OB: revê movimentos sem correspondência
```

## Passo-a-passo

1. Em `/fin`, abre **Open Banking**.
2. Clica **Simular consentimento PSD2** — estado passa a `connected`.
3. **Sincronizar movimentos** — recebe transacções mock (Netflix, Spotify, AWS…).
4. O cliente associa débitos a subscrições pelo valor mensal.
5. No feed, aparece **Sugestão: reconciliar pagamentos**.
6. **Aprova** — revê os movimentos «sem correspondência» manualmente.

## Pós-condições

- Metadados de ligação em `bank_connections` (sem tokens em claro).
- Evento `fin.transactions.synced` auditado; movimentos não persistidos no servidor.

> ⚠️ **Segurança:** provider real e certificados mTLS (`AEGIS_OB_*`) ficam para
> quando INFRA-002 (WireGuard) e credenciais de TPP estiverem disponíveis.
