# Journey: Registar e usar passkey

> **Ator:** Utilizador · **Epics:** `VAULT`, `UI`

A passkey autentica a **sessão HTTP** com WebAuthn (biometria ou PIN do
dispositivo) — mas a **Master Key** continua a vir só da Master Password
(Zero-Knowledge).

## Pré-condições

- Browser ou SO com suporte WebAuthn (Secure Enclave / TPM).
- Utilizador com conta criada e pelo menos um login por password bem-sucedido.

## Fluxo principal

```mermaid
sequenceDiagram
    participant U as Utilizador
    participant App as App (Svelte)
    participant SE as Secure Enclave
    participant G as Go API

    Note over U,G: Registo (Definições)
    U->>App: "Registar passkey" + nome
    App->>G: GET opções de criação WebAuthn
    G-->>App: challenge + user id
    App->>SE: navigator.credentials.create()
    SE-->>App: par de chaves (privada fica no dispositivo)
    App->>G: POST credencial pública
    G-->>App: passkey registada

    Note over U,G: Login
    U->>App: login com passkey
    App->>G: GET opções de autenticação
    G-->>App: challenge
    App->>SE: navigator.credentials.get()
    SE-->>App: assinatura do challenge
    App->>G: POST verificação
    G-->>App: token de sessão
    App-->>U: redirecciona; cofre pede unlock (Master Password)
```

## Passo-a-passo

1. Em **Definições**, o utilizador regista uma passkey com nome descritivo
   (ex.: "MacBook Touch ID").
2. O servidor envia um *challenge*; o dispositivo cria chave privada no Secure
   Enclave e devolve só a parte pública.
3. No login, o utilizador escolhe passkey; o dispositivo assina o challenge sem
   revelar a chave privada.
4. Com token de sessão, a app pede **unlock** do cofre — Argon2id com Master
   Password, como sempre.

## Conceito didático

> 💡 **WebAuthn / passkey:** autenticação por **par de chaves assimétricas** ligado
> ao dispositivo. Resiste a phishing melhor que passwords — o browser só assina
> pedidos do domínio correcto.

> ⚠️ **Segurança — duas camadas:** passkey = "és tu para o servidor". Master
> Password = "consegues decifrar o cofre". Comprometer uma não compromete a outra
> automaticamente.

## Fluxos alternativos

- **Sem WebAuthn:** login clássico por Auth Hash continua disponível.
- **Sentinel step-up:** a mesma passkey pode servir de segundo factor no login
  suspeito (ver journey Sentinel).
- **Revogação:** em Definições ou `/security/devices`, remove passkeys perdidas.

## Pós-condições

- Sessão HTTP activa sem enviar Master Password ao servidor.
- Master Key só em memória após unlock explícito.
