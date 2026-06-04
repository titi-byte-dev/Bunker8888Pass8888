# Testes — Segurança

Testes específicos para garantir que o modelo de ameaça se mantém. Complementa a
[`strategy.md`](strategy.md).

## 1. Vetores de teste criptográficos

> 💡 **Conceito — Test vectors:** valores de entrada/saída conhecidos e
> publicados (ex: NIST) para um algoritmo. Se a nossa implementação produz o
> mesmo output para o mesmo input, está conforme a especificação.

- [ ] Argon2id validado contra vetores oficiais (RFC 9106).
- [ ] AES-GCM validado contra vetores do NIST.
- [ ] TOTP validado contra vetores do RFC 6238.

## 2. Testes de isolamento multi-tenant

- [ ] Tentar ler dados de outro `tenant_id` → **deve falhar** (RLS).
- [ ] Esquecer o `WHERE tenant_id` numa query → RLS continua a proteger.
- [ ] Token de sessão do tenant A não dá acesso a recursos do tenant B.

> ⚠️ **Segurança:** este é um teste *negativo* — confirmamos que algo **NÃO**
> acontece. Testes negativos são tão importantes como os positivos em segurança.

## 3. Fuzzing

> 💡 **Conceito — Fuzzing:** alimentar a função com inputs aleatórios/malformados
> em massa para encontrar crashes ou comportamentos inesperados. Go tem fuzzing
> nativo desde a 1.18.

```go
// FuzzDecifrar atira bytes aleatórios ao decifrador. NUNCA deve causar panic;
// no máximo deve devolver um erro controlado. `f.Fuzz` corre milhares de vezes.
func FuzzDecifrar(f *testing.F) {
    key := make([]byte, 32)
    f.Add([]byte("input inicial")) // "semente" para o fuzzer
    f.Fuzz(func(t *testing.T, data []byte) {
        // Ignoramos o resultado; só queremos garantir que não há panic
        // e que input inválido devolve erro em vez de comportamento estranho.
        _, _ = Decifrar(key, data)
    })
}
```

- [ ] Fuzz dos parsers (decifragem, import de passwords, masking regex).
- [ ] Fuzz das entradas da API.

## 4. Testes de abuso / antifuncionais

| Cenário | Resultado esperado |
|---|---|
| Login fora do turno | Negado + sem chave em memória |
| Acesso sem VPN | API inalcançável |
| Reutilização de nonce (injetado) | Rejeitado / impossível por design |
| Secret link após expirar | 404 / sem dados em disco |
| Prompt injection num e-mail lido por agente | Tratado como dados, não ordens |
| Open relay (envio não autenticado) | Bloqueado |

## 5. Auditoria externa (antes de produção)

- [ ] Pen-test por terceiro independente.
- [ ] Revisão de cripto por especialista.
- [ ] Análise de dependências (SCA) e SBOM.

> ⚠️ **Regra de ouro:** **nunca** implementar criptografia "à mão". Usar sempre
> as bibliotecas padrão auditadas (`crypto/*` de Go, WebCrypto no browser).
