# Epic: Dark Web & Auditoria — `DW`

> **Fase:** 2 · **Prioridade:** 🟡 Média

## Objetivo

Vigilância proativa: detetar credenciais expostas em fugas de dados e medir a
"saúde de segurança" da empresa sem o servidor ver as passwords reais.

## Valor de negócio

Transforma segurança reativa em proativa. Dá ao RH/gestão um painel de risco —
argumento forte de venda e retenção.

## Funcionalidades

- **Monitorização da Dark Web:** verifica periodicamente bases de dados de fugas
  (via APIs de *breach data*) e, se um e-mail da empresa aparecer, força a
  alteração da password
- **Audit de higiene de passwords:** painel com % de passwords fracas,
  reutilizadas ou expostas — score calculado **no cliente**
- **Proteção avançada (Sentinel Mode):** deteta padrões de login impossíveis
  (horário/geografia incoerente) e exige prova adicional

## Critérios de aceitação

- [ ] Verificação de fugas usa modelo *k-anonymity* (não envia a password/hash
      completa para o serviço externo).
- [ ] Score de higiene calculado no cliente; servidor só recebe o número.
- [ ] Sentinel sinaliza login geograficamente impossível e exige re-verificação.

## Conceitos didáticos

> 💡 **k-anonymity:** para verificar se uma password vazou **sem a revelar**,
> envia-se só os **primeiros 5 caracteres** do hash SHA-1. O serviço de breach
> data devolve todos os hashes com esse prefixo; a comparação final faz-se
> **localmente**. O servidor externo nunca sabe qual password verificámos.

```go
// Hash SHA-1 da password, em hexadecimal e MAIÚSCULAS (formato comum das APIs
// de breach data baseadas em k-anonymity).
sum := sha1.Sum([]byte(password))
hash := strings.ToUpper(hex.EncodeToString(sum[:]))
prefix, suffix := hash[:5], hash[5:] // só o prefixo sai do dispositivo
// GET {breach_api}/range/{prefix}
// → comparar `suffix` na lista devolvida, localmente.
```

> ⚠️ **Segurança:** nunca enviar a password (nem o hash completo) para serviços
> externos. O k-anonymity preserva a privacidade mesmo numa funcionalidade que,
> à primeira vista, exigiria partilhar a credencial.

## Dependências

- Vault (`VAULT`) para inventário de credenciais.

## Tasks

Ver prefixo `DW-*` em [`../05-tasks/backlog.md`](../05-tasks/backlog.md).
