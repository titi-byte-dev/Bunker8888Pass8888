---
title: Finanças SaaS
slug: fin
category: product
order: 5
audience: [user, admin]
layer: [frontend, backend]
feature: fin
level: 1
in_app: true
summary: Monitorização de subscrições, custos, alertas, fiscal, faturação e Open Banking.
related: [vault, journey-saas-costs, journey-fiscal-categorization, production-path]
---

:::summary
O módulo **Finanças** cobre custos SaaS, classificação fiscal, faturação ERP,
comissões e reconciliação bancária — com dados sensíveis sempre cifrados no cliente.
:::

:::concept{id="saas-subscription" title="Subscrição SaaS" level=1}
Registo de um serviço pago: nome, custo, moeda, ciclo, categoria, código fiscal
opcional e data de último uso. Pode ligar-se a um login do cofre pelo ID do item.
:::

:::concept{id="fiscal-code" title="Código fiscal (FIN-005)" level=2}
Classificação IRC (ex.: dedutível 100%, investimento) guardada **dentro do blob
cifrado** da subscrição. Totais dedutíveis calculam-se só no browser.
:::

:::level{level=1 title="Rotas na app"}

```mermaid
sequenceDiagram
participant User as Utilizador
participant Fin as /fin
participant Vault as Cofre ZK
User->>Fin: abre hub Finanças
Fin->>Vault: desbloqueia Master Key
Vault-->>Fin: blobs cifrados
```

| Rota | Função |
|---|---|
| `/fin` | Subscrições, alertas licenças sem uso, feed agente |
| `/fin/fiscal` | Resumo dedutível mensal + classificação |
| `/fin/banking` | Open Banking (mock em dev) + reconciliação |
| `/fin/invoices` | Pro-forma, faturas, recibos (FIN-006) |
| `/fin/commissions` | Comissões sobre faturas pagas (FIN-007) |
:::

:::level{level=2 title="Alertas de licença sem uso"}

```mermaid
flowchart LR
  Sub[Subscrição activa] --> Check{last_used antigo?}
  Check -->|sim| Alert[Alerta FIN-001]
  Check -->|não| OK[Sem alerta]
  Alert --> Vault[Liga login cofre]
```

Regra no cliente: subscrição activa com `last_used_at` antigo ou em falta.
Cruza com logins do cofre para sugerir corte de custos.
:::

:::level{level=2 title="Agente financeiro"}

```mermaid
sequenceDiagram
participant Fin as Agente FIN
participant Bus as Event bus
participant Orq as Orquestrador
participant User as Aprovador
Fin->>Bus: fin.subscription.stale
Bus->>Orq: worker finance
Orq->>User: human-in-the-loop
User-->>Orq: aprova / rejeita
```

Reporta licenças stale e movimentos bancários ao orquestrador. Acções sensíveis
passam por **Human-in-the-loop** no feed de eventos.
:::

:::level{level=3 title="Técnico"}
Frontend: `frontend/src/lib/fin/`. Backend: blobs em `saas_subscriptions`, APIs
`/api/fin/*`, Open Banking em `internal/openbanking/`. Epic:
`docs/roadmap/03-epics/epic-fintech.md`.
:::
