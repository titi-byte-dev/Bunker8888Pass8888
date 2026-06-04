# Epic: Recursos Humanos & RGPD — `HR`

> **Fase:** 1 · **Prioridade:** 🔴 Alta (diferenciador de mercado)

## Objetivo

Módulo de RH com cifragem campo-a-campo e conformidade RGPD por desenho, que
torna o AegisPass um sistema core (não só um utilitário de segurança).

## Valor de negócio

As PMEs têm pavor de multas RGPD. Vender "conformidade automática + cifragem de
dados de RH" abre portas em departamentos jurídicos e de RH.

## Funcionalidades

- **Ficha de empregado encriptada** (cartão de cidadão, IBAN, salário, saúde) —
  encriptação campo-a-campo
- **Direito ao Esquecimento (RGPD Art. 17):** apaga todos os registos do
  funcionário com um clique + **certificado criptográfico de eliminação**
- **Logs imutáveis (auditoria):** cada visualização de dado sensível registada
  com *hashing encadeado* (estilo blockchain)
- **Gestão de contratos + assinatura digital:** upload para object storage
  privado, cada ficheiro cifrado com chave única por funcionário
- **Triagem às cegas** (oculta género/etnia) — base para o agente de recrutamento
- **Onboarding/Offboarding:** cria/revoga credenciais, alias, cofre da equipa

## Critérios de aceitação

- [ ] Campos sensíveis chegam à BD já cifrados (servidor nunca vê em claro).
- [ ] Erasure remove dados e produz certificado verificável.
- [ ] Logs de acesso são *append-only* e detetam adulteração via hash encadeado.
- [ ] Cada contrato tem a sua própria chave de ficheiro.

## Conceitos didáticos

> 💡 **Hashing encadeado:** cada registo de log inclui o hash do registo
> anterior. Alterar um registo antigo muda o seu hash, o que "parte a corrente"
> de todos os seguintes — tornando a adulteração detetável.

```go
// Cada entrada do log "aponta" para a anterior pelo hash. Isto cria uma cadeia
// inviolável: mudar uma entrada antiga invalida TODAS as posteriores.
type LogEntry struct {
    Timestamp time.Time
    Actor     string
    Action    string
    PrevHash  []byte // hash da entrada anterior
}

func (e LogEntry) Hash() []byte {
    // sha256.Sum256 devolve um array de tamanho fixo [32]byte; convertemos para
    // slice com [:] para ser mais fácil de usar/armazenar.
    data := fmt.Sprintf("%v|%s|%s|%x", e.Timestamp, e.Actor, e.Action, e.PrevHash)
    sum := sha256.Sum256([]byte(data))
    return sum[:]
}
```

> ⚠️ **Segurança — Art. 17 vs. logs imutáveis:** há tensão entre "apagar tudo" e
> "logs inalteráveis". Resolve-se cifrando os dados pessoais e, no erasure,
> **destruindo a chave** (*crypto-shredding*): o log mantém-se íntegro mas o
> conteúdo pessoal torna-se irrecuperável.

## Dependências

- Núcleo Zero-Knowledge, object storage cifrado.

## Tasks

Ver prefixo `HR-*` em [`../05-tasks/backlog.md`](../05-tasks/backlog.md).
