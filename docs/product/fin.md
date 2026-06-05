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
summary: Monitorização de subscrições, custos mensais e alertas de licenças sem uso.
related: [vault, journey-saas-costs]
---

:::summary
A secção **Finanças** (`/fin`) ajuda a perceber quanto a empresa gasta em SaaS e
quais licenças estão pagas mas sem uso — cruzando metadados com itens do cofre.
:::

:::concept{id="saas-subscription" title="Subscrição SaaS" level=1}
Registo de um serviço pago: nome, custo, moeda, ciclo de facturação e data de
último uso. Pode ligar-se a um login do cofre pelo ID do item (sem enviar a password).
:::

:::concept{id="unused-license-alert" title="Alerta de licença sem uso" level=2}
Regra calculada **no cliente**: subscrição activa com `last_used_at` antigo ou
em falta. Sugere rever acesso ou desactivar a licença no offboarding.
:::

:::level{level=1 title="Onde encontrar na app"}
- `/fin` — painel de subscrições, totais e alertas
- Associar cada SaaS a um item do cofre para rastreio de uso
:::

:::level{level=2 title="Aprofundar: offboarding"}
Quando um funcionário sai, desactivar subscrições ligadas ao mesmo tempo que
revogar cofres evita pagar licenças órfãs. Ver percurso **SaaS e licenças**.
:::

:::level{level=3 title="Técnico"}
Frontend: `frontend/src/lib/fin/`, rota `(app)/fin`. Backend: API de subscrições
com filtro `tenant_id`. Epic: `docs/roadmap/03-epics/epic-fintech.md`.
Open Banking e cartões efémeros são Fase 3 (`FIN-003+`).
:::
