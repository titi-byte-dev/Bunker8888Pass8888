# Journey: Enviar segredo com Secret Link efémero

> **Ator:** Utilizador (remetente) · **Epics:** `SHARE`

Partilhar uma password ou nota sensível **uma vez**, sem deixar rasto permanente
no servidor — o segredo vive em RAM até expirar ou ser lido.

## Gatilhos

- Onboarding rápido (credencial de projeto para freelancer).
- Partilha pontual de IBAN, API key ou código de um só uso.
- Alternativa segura ao "cola isto no Slack".

## Fluxo principal

```mermaid
sequenceDiagram
    participant R as Remetente
    participant App as App (Svelte)
    participant G as Go API
    participant D as Destinatário

    R->>App: cria secret link (TTL ou 1 clique)
    App->>App: gera chave efémera + cifra payload (AES-GCM)
    App->>G: envia blob cifrado + metadados (sem plaintext)
    G->>G: guarda em RAM / TTL agressivo
    G-->>App: URL do link
    App-->>R: copia link para enviar (canal à parte)
    D->>G: abre link (GET único)
    G->>G: serve payload da RAM + invalida
    G-->>D: blob cifrado (uma vez)
    D->>App: decifra no browser (chave na URL fragment)
    App-->>D: mostra segredo; link fica morto
```

## Passo-a-passo

1. Em `/team/links`, o remetente define o conteúdo, TTL (minutos) ou política de
   **um clique**.
2. O cliente cifra o segredo; a chave de desencriptação pode ir no *fragment*
   da URL (`#...`) — nunca chega ao servidor nos logs de acesso.
3. O servidor mantém o blob só em **memória** com expiração; após leitura ou
   timeout, o link responde 404.
4. O destinatário abre o link, vê o segredo no browser e o link deixa de funcionar.

## Conceito didático

> 💡 **Fragment URL:** tudo depois de `#` na URL **não é enviado ao servidor**
> no pedido HTTP. É um truque clássico para passar chaves de desencriptação só
> ao destinatário, sem o servidor as ver.

> ⚠️ **Segurança — RAM only:** links efémeros não devem ser persistidos em disco
> nem em logs. TTL curto + invalidação após primeiro GET reduzem a janela de
> exposição se a RAM for despejada tarde.

## Fluxos alternativos

- **Link expirado:** o destinatário vê mensagem genérica; o remetente gera novo link.
- **Tentativa de reabrir:** segundo GET falha — política de clique único.

## Pós-condições

- Segredo não recuperável no servidor após expiração/leitura.
- Remetente pode auditar *que* criou um link, sem o conteúdo em claro.
