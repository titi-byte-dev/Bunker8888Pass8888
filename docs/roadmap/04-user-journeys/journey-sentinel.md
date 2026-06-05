# Journey: Sentinel deteta login impossível

> **Ator:** Funcionário (login) · Admin (revisão) · **Epics:** `DW`, `VAULT`

Quando o padrão de acesso é geograficamente ou temporalmente incoerente, o
Sentinel Mode bloqueia o login até **step-up** adicional — sem desligar o
Zero-Knowledge.

## Pré-condições

- Sentinel activo no tenant (DW-004).
- Último login registado com IP/GPS e timestamp fiável (NTP).

## Fluxo principal

```mermaid
sequenceDiagram
    participant U as Funcionário
    participant App as App (Svelte)
    participant G as Go API
    participant Adm as Admin

    U->>App: login (Auth Hash) de localização nova
    App->>G: POST /api/auth/login + contexto geo
    G->>G: avalia viagem impossível / horário anómalo
    alt padrão suspeito
        G-->>App: 403 sentinel_step_up + challenge_id
        App-->>U: banner: confirma identidade (passkey/GPS)
        U->>App: completa step-up (WebAuthn ou geo)
        App->>G: POST step-up/finish
        G-->>App: token de sessão
    else padrão normal
        G-->>App: token de sessão
    end
    App-->>U: sessão activa; cofre exige unlock separado
    G->>G: regista evento em /security/sentinel
    Adm->>App: revê alertas no painel
```

## Passo-a-passo

1. O funcionário tenta login a partir de um IP/GPS inconsistente com o último
   acesso (ex.: Lisboa há 10 min, Tokyo agora).
2. A API recusa o token completo e devolve `sentinel_step_up` com um
   `challenge_id`.
3. A app mostra um banner em `/auth/login` pedindo **prova adicional** —
   passkey registada ou confirmação de geolocalização.
4. Após step-up bem-sucedido, emite-se o token HTTP; o cofre continua a exigir
   Master Password (ZK intacto).
5. O admin vê o evento em `/security/sentinel` com motivo legível.

## Conceito didático

> 💡 **Step-up authentication:** a sessão "normal" não basta para contextos de
> risco. Pedimos um segundo factor **só quando o risco sobe** — melhor UX que
> 2FA obrigatório em cada clique.

> ⚠️ **Segurança:** o Sentinel não substitui turnos nem geofencing — acrescenta
> uma camada quando o padrão de login parece **impossível** mesmo dentro das regras
> base.

## Fluxos alternativos

- **Step-up falha:** login negado; tentativa fica no log; admin pode forçar
  revogação de sessões.
- **Dispositivo sem passkey:** fallback para confirmação GPS ou contacto admin.

## Pós-condições

- Sessão só emitida após step-up em caso de alerta.
- Evento auditável sem expor Master Key nem passwords.
