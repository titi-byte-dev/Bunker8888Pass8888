# Journey — Categorização fiscal (FIN-005)

> **Ator:** Gestor financeiro · **Pré-requisito:** subscrições SaaS em `/fin`

## Fluxo

```mermaid
flowchart LR
    SUB[Subscrições ZK] --> DEC[Decifrar no cliente]
    DEC --> SUG[Sugestão heurística / agente]
    SUG --> HITL[Revisão humana]
    HITL --> BLOB[Re-cifrar com fiscalCode]
    BLOB --> SUM[Resumo dedutível /fin/fiscal]
```

## Passos

1. Em `/fin`, cada subscrição pode ter **código fiscal** no blob cifrado.
2. Em `/fin/fiscal`, vê totais mensais e parcela **dedutível estimada**.
3. **Sugerir automaticamente** aplica regras (SaaS → `dedutivel_100`, hardware → `investimento`).
4. Ajusta manualmente no dropdown — a classificação volta a cifrar-se.

> ⚠️ **Segurança:** totais e códigos fiscais **nunca** passam pelo servidor em claro.
> O agente futuro (Guardião) só verá campos estritamente necessários.

> 💡 **Conceito:** *dedução fiscal* reduz o lucro tributável — nem todas as despesas
> SaaS são 100% dedutíveis; hardware pode ser amortizado como investimento.

## Verificação

- [ ] Subscrição «Figma» classificada como `dedutivel_100`
- [ ] Resumo mostra dedutível = 100% do custo mensal dessa linha
- [ ] Servidor continua a ver só blob opaco (inspecionar rede)
