---
title: Google Workspace (proxy)
slug: google-workspace
category: product
order: 6
audience: [user, admin]
layer: [frontend, backend]
feature: google_workspace
level: 1
in_app: true
summary: Layer de controlo sobre Drive, Sheets e Gmail — cifragem ZK e mascaramento sem expor PII à Google.
related: [journey-google-dev-stub, vault]
---

:::summary
O AegisPass permite continuar a usar **Google Workspace** com Drive/Sheets/Gmail,
mas a Google deixa de ver dados sensíveis em claro — ficheiros cifrados e células
mascaradas com tokens.
:::

:::concept{id="google-proxy" title="Proxy de ofuscação" level=1}
A app intercepta dados **antes** de irem para a Google. A Google fica como
armazenamento/transporte; a inteligência e a cifragem são do AegisPass.
:::

:::level{level=1 title="Onde na app"}
| Rota | Modo | Task |
|---|---|---|
| `/work/google` | Estado do provider (mock ou Service Account) | GOOGLE-001 |
| `/work/google-dev` | Simulação local sem OAuth | DoD Fase 2 dev |
:::

:::level{level=2 title="Drive & Docs (GOOGLE-002)"}
Ficheiros são cifrados no browser com a Master Key antes do upload. Na Drive da
empresa só existem blobs opacos. Abrir no AegisPass decifra em memória.
:::

:::level{level=2 title="Sheets (GOOGLE-003)"}
NIF, IBAN e campos sensíveis são substituídos por tokens (`TOKEN_NIF_*`). O mapa
token→valor fica cifrado (em dev: cliente; em prod: PostgreSQL ZK).
:::

:::level{level=3 title="Produção"}
Configura Service Account e domain-wide delegation — ver
`docs/roadmap/10-production/google-001-oauth.md`. Variáveis: `AEGIS_GOOGLE_SA_JSON`,
`AEGIS_GOOGLE_DELEGATED_USER`, `AEGIS_GOOGLE_ENABLED=true`.
:::
