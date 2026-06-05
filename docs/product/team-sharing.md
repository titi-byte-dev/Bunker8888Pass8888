---
title: Equipa e partilha
slug: team-sharing
category: product
order: 4
audience: [user]
layer: [frontend, backend]
feature: share
level: 1
in_app: true
summary: Cofres partilhados, secret links e notas temporárias.
related: [vault, journey-admin-onboarding]
---

:::summary
A secção **Equipa** permite partilhar segredos com colegas sem expor passwords em
chat ou e-mail — sempre com cifragem e permissões explícitas.
:::

:::concept{id="shared-vault" title="Shared Vault" level=1}
Cofre com **vários membros** e papéis (ler, editar, administrar). Cada utilizador tem
chaves assimétricas; o conteúdo re-cifrado por membro mantém Zero-Knowledge face ao
servidor.
:::

:::concept{id="secret-link" title="Secret link efémero" level=2}
Link de uso único ou temporário que serve o segredo **da RAM** — não persiste em claro
em disco no servidor. Ideal para credenciais de projeto ou onboarding rápido.
:::

:::level{level=1 title="Onde encontrar na app"}
- `/team/vaults` — cofres partilhados
- `/team/links` — secret links
- `/team/notes` — notas auto-destrutivas
:::

:::level{level=2 title="Aprofundar: notas temporárias"}
Notas com **TTL** (*time-to-live*): após expirar, o conteúdo deixa de estar disponível.
Útil para partilhar IBAN ou códigos de um só uso sem rasto permanente.
:::

:::level{level=3 title="Técnico"}
Backend: `backend/internal/sharing/`. Frontend: `frontend/src/routes/(app)/team/`.
Epic: `docs/roadmap/03-epics/epic-sharing.md`.
:::
