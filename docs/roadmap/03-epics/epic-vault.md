# Epic: The Vault — `VAULT`

> **Fase:** 1 · **Prioridade:** 🔴 Crítica (fundação)

## Objetivo

Cofre Zero-Knowledge para guardar e sincronizar credenciais, notas e cartões,
com acesso efémero (turnos), 2FA, e capacidade de remote wipe.

## Valor de negócio

É a fundação da confiança do produto. Sem um cofre seguro e usável, nada do
resto tem credibilidade. Resolve a dor "funcionário usa PC pessoal".

## Funcionalidades

- Inícios de sessão, notas seguras e cartões de crédito **ilimitados**
- Dispositivos ilimitados + apps (browser, telemóvel, computador)
- Gerador de palavras-passe
- Alertas para passwords fracas e reutilizadas (score calculado no cliente)
- Chaves de acesso (passkeys) suportadas
- Importação fácil de palavras-passe
- **Acesso por turnos:** chave expurgada da memória fora do horário
- **Browser isolado (sandbox):** injeta credenciais sem revelar a password,
  bloqueia copy-paste
- **2FA / TOTP integrado** (RFC 6238)
- **Remote wipe de emergência**
- **Acesso de Emergência** (herdeiro digital, com período de espera)
- CLI em Go (injeta segredos em scripts sem os escrever em texto plano)

## Critérios de aceitação

- [ ] A Master Key nunca é transmitida nem persistida no servidor.
- [ ] Itens cifrados com AES-GCM-256; cada cifragem usa nonce único.
- [ ] Sincronização entre dispositivos em < 1s após alteração (WebSockets).
- [ ] Revogação de acesso reflete-se no dispositivo em < 1s.
- [ ] Fora do turno, a chave local é apagada (verificável).
- [ ] No sandbox browser, a password real nunca é exibida nem copiável.
- [ ] Score de higiene calculado no cliente; servidor só recebe o número.

## Conceitos didáticos

> 💡 **TOTP:** gera um código de 6 dígitos a partir de um segredo partilhado + a
> hora atual (em janelas de 30s). Como ambos os lados conhecem o segredo e a
> hora, geram o mesmo código sem comunicar.

> ⚠️ **Segurança — expurgar chave da memória:** em JS/TS não há "apagar memória"
> garantido (o *garbage collector* decide), mas minimizamos a janela: nunca
> guardar a chave em `localStorage`; mantê-la só em memória volátil e
> descartar a referência ao bloquear/fim de turno.

## Dependências

- Núcleo Zero-Knowledge ([`zero-knowledge.md`](../01-architecture/zero-knowledge.md))
- Multi-tenancy ([`multi-tenancy.md`](../01-architecture/multi-tenancy.md))

## Tasks

Ver prefixo `VAULT-*` em [`../05-tasks/backlog.md`](../05-tasks/backlog.md).
