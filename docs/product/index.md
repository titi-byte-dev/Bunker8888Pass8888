---
title: Começar no AegisPass
slug: start
category: product
order: 1
audience: [user]
layer: [product]
level: 1
in_app: true
summary: Primeiros passos — sessão, desbloqueio do cofre e navegação da app.
related: [glossary, vault, security]
---

:::summary
O AegisPass separa **autenticação no servidor** (quem és) do **desbloqueio local**
(o cofre). Esta distinção é o primeiro hábito a internalizar.
:::

:::concept{id="login-vs-unlock" title="Login ≠ Desbloquear cofre" level=1}
**Login** (ou passkey) autentica a tua sessão HTTP — o servidor sabe que és tu.
**Unlock** deriva a Master Key localmente para decifrar itens. Podes estar autenticado
e o cofre continuar bloqueado até introduzires a Master Password.
:::

:::level{level=1 title="Fluxo recomendado (primeira vez)"}
1. **Registo** — defines Master Password; o cliente envia só o Auth Hash.
2. **Login** — sessão activa; opcionalmente registas uma **passkey** em Definições.
3. **Unlock** — desbloqueias o cofre com a Master Password (ou recovery key).
4. **Cofre** — consultas, crias e editas logins, notas e cartões.
5. **Segurança** — revê higiene, dispositivos e políticas do teu turno.
:::

:::level{level=2 title="Navegação da app"}
| Secção | Para quê |
|---|---|
| **Cofre** | Credenciais, notas, cartões, importação |
| **Segurança** | Higiene, breaches, dispositivos, Sentinel |
| **Trabalho** | Turnos, sandbox browser, CLI |
| **Equipa** | Cofres partilhados, secret links |
| **RH** | Fichas, contratos, onboarding, RGPD |
| **Definições** | Conta, tema, passkeys, **esta documentação** |
:::

:::level{level=3 title="Para programadores: arranque local"}
```bash
docker compose up --build   # API + PostgreSQL
cd frontend && npm run dev  # SvelteKit em :5173
```

A API expõe `GET /healthz`. O proxy Vite encaminha `/api/*` para `localhost:8080`.
Documentação técnica: secção **Programador** neste mesmo painel.
:::
