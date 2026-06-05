---
title: Panorama competitivo
slug: competitive-landscape
category: competitive
order: 1
audience: [developer, admin]
layer: [product]
level: 2
in_app: true
summary: Concorrentes, abordagens tecnológicas e ideias validadas para o roadmap.
related: [vault, hr-rgpd]
---

:::summary
O AegisPass cruza **password manager**, **IAM BYOD**, **RH** e (futuro) **ERP/CRM**.
Esta análise alimenta decisões — cada ideia aceite deve virar task no backlog.
:::

:::level{level=2 title="Categorias e players"}
| Categoria | Exemplos | O que fazem bem |
|---|---|---|
| Password manager | 1Password, Bitwarden, Proton Pass | UX de cofre, TOTP, partilha familiar/team |
| IAM / SSO | Okta, Microsoft Entra | Políticas enterprise, SCIM |
| HR BYOD | Rippling, Deel | Onboarding 1-clique, payroll |
| ZTNA / VPN | Tailscale, Cloudflare Zero Trust | Rede zero-trust (nós: WireGuard) |
:::

:::level{level=2 title="Matriz de diferenciação (AegisPass)"}
| Capacidade | AegisPass | Típico password manager |
|---|---|---|
| Zero-Knowledge E2EE | ✅ core | ✅ |
| BYOD + wipe selectivo | ✅ | ❌ / MDM pesado |
| Turnos + geofence | ✅ | ❌ |
| RH + RGPD campo-a-campo | ✅ | ❌ |
| Browser sandbox inject | ✅ (roadmap) | Extensão expõe DOM |
| ERP/CRM + agentes IA | 🔜 fase 2–3 | ❌ |
:::

:::level{level=3 title="Pipeline: observação → backlog"}
1. **Observar** feature ou stack num concorrente
2. **Validar** com utilizador piloto ou spike técnico
3. **Decidir** incluir / adiar / rejeitar (registar em `validated-ideas`)
4. **Criar task** com ID (`VAULT-*`, `HR-*`, …) em `backlog.md`

:::concept{id="validated-example" title="Exemplo validado" level=3}
**k-anonymity para breach check** (estilo HIBP) — adoptado em `DW-001` após confirmar
que não envia passwords em claro ao serviço externo.
:::

:::level{level=3 title="Ideias em observação (não validadas)"}
| Ideia | Fonte inspiração | Estado |
|---|---|---|
| Google proxy ZK (Drive/Sheets) | Interno + custo storage | ⚪ GOOGLE-* |
| Cartões virtuais efémeros | Brex, Stripe Issuing | ⚪ FIN-004 |
| Command palette universal | Linear, Raycast | 🟢 UI-006 |
| SCIM provisioning | Okta | ⚪ futuro |
:::
