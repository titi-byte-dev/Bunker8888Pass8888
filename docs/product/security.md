---
title: Segurança e confiança
slug: security
category: product
order: 3
audience: [user]
layer: [product]
feature: vault
level: 1
in_app: true
summary: Higiene, dispositivos, turnos, geofencing, Sentinel e emergência.
related: [vault, glossary, journey-remote-wipe]
---

:::summary
A área **Segurança** reúne tudo o que te ajuda a perceber e reforçar a postura da
conta — sem alarmismo, com dados acionáveis.
:::

:::concept{id="hygiene-score" title="Score de higiene" level=1}
Pontuação calculada **no teu browser** a partir dos logins decifrados: passwords
fracas, reutilizadas ou com aparência em fugas públicas (via k-anonymity). O servidor
não precisa de ver as passwords para te dar recomendações.
:::

:::concept{id="sentinel" title="Sentinel Mode" level=2}
Deteta logins **geograficamente impossíveis** (ex.: Lisboa e Tokyo em minutos). Pode
exigir *step-up* (passkey ou confirmação) antes de permitir acesso sensível.
:::

:::level{level=1 title="Páginas principais"}
| Rota | Função |
|---|---|
| `/security/hygiene` | Score, passwords fracas, acções |
| `/security/devices` | Sessões, passkeys, revogação |
| `/security/sentinel` | Alertas de viagem impossível |
| `/security/emergency` | Acesso de emergência (herdeiro digital) |
| `/work/shifts` | Turnos e geofence |
:::

:::level{level=2 title="Aprofundar: turnos e geofencing"}
Fora do **horário de turno** a Master Key é expurgada da memória — mesmo com sessão
HTTP válida, o cofre fica inacessível até novo unlock dentro das regras.

O **geofencing** valida IP/GPS nos pedidos; combina com WireGuard para Zero-Trust
por camadas.
:::

:::level{level=3 title="Técnico: remote wipe"}
O admin pode acionar **remote wipe** via WebSocket: invalida sessão, revoga chaves
de desencriptação locais e apaga caches corporativos **sem tocar em dados pessoais**
do dispositivo BYOD.

Ver percurso completo: **Percursos → Remote Wipe**.
:::
