# Epic: Agentes de IA (ERP/CRM) — `AGENT`

> **Fase:** 2 (1º agente) → 3 (orquestrador completo) · **Prioridade:** 🟣 Estratégica

## Objetivo

Construir o "staff de IA" que opera o ERP/CRM autonomamente, sob a supervisão do
Guardião e com *human-in-the-loop* nas decisões críticas.

## Valor de negócio

Quebra a rigidez dos ERPs tradicionais e monolíticos. O utilizador supervisiona
um feed de ações em vez de preencher dezenas de formulários.

## Arquitetura

Ver [`../01-architecture/agents-architecture.md`](../01-architecture/agents-architecture.md)
(EDA, Guardião, e a nota de propriedade intelectual).

## Agentes

| Agente | Módulo | Função resumida |
|---|---|---|
| Prospeção & Enriquecimento | CRM | Lê e-mails/leads, preenche fichas |
| Negociação (Sales Copilot) | CRM | Sugere margens, redige propostas |
| Reconciliação Bancária | Financeiro | Concilia pagamentos (Open Banking) |
| Otimização Fiscal | Financeiro | Categoriza despesas para dedução |
| Recrutamento | RH | Triagem às cegas + agenda entrevistas |
| Onboarding | RH | Cria credenciais, e-mail, contrato |
| Compras & Inventário | Operações | Prevê stock, gera ordens de compra |

## Critérios de aceitação

- [ ] Cada agente declara as suas *tools* com schema e permissões explícitas.
- [ ] O Guardião só decifra os campos pedidos (menor privilégio) — auditável.
- [ ] Ações financeiras/legais exigem aprovação humana (1 clique).
- [ ] Eventos publicados/consumidos são rastreáveis (log imutável).
- [ ] LLM corre isolado; nenhum dado sensível em claro chega ao modelo sem o
      Guardião.

## Conceitos didáticos

> 💡 **Function Calling / Tools:** dá-se ao LLM uma lista de "ferramentas"
> (funções com schema). Em vez de o LLM "adivinhar", ele pede para chamar uma
> tool com argumentos; o nosso código Go executa-a com segurança e devolve o
> resultado. O LLM nunca toca diretamente na base de dados.

```go
// Padrão de "Tool" reimplementado de raiz (inspirado em padrões conhecidos de
// agentic tooling, SEM copiar código de terceiros). Valida o input antes de executar.
type Tool interface {
    Name() string
    // Validate corre ANTES de Execute: falha cedo se o input for inválido,
    // evitando efeitos colaterais com dados malformados.
    Validate(input json.RawMessage) error
    Execute(ctx context.Context, input json.RawMessage) (any, error)
}
```

> ⚠️ **Segurança — prompt injection:** um e-mail malicioso pode tentar "dar
> ordens" ao agente que o lê. Mitigação: tratar conteúdo externo como **dados,
> nunca instruções**; o agente só pode agir através de tools validadas e com
> aprovação humana em ações sensíveis.

## Dependências

- Fases 1 e 2 (vault, multi-tenancy, proxy Google, aliases).
- Orquestrador/Event Bus em Go.

## Tasks

Ver prefixo `AGENT-*` em [`../05-tasks/backlog.md`](../05-tasks/backlog.md).
