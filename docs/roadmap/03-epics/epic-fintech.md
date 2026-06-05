# Epic: Controlo Financeiro & FinTech — `FIN`

> **Fase:** 2 (monitorização) → 3 (transacional) · **Prioridade:** 🟠 Média-Alta

## Objetivo

Dar à empresa controlo sobre custos de SaaS e capacidade de emitir cartões
virtuais efémeros, integrando com Open Banking.

## Valor de negócio

"As empresas adoram controlo de custos." Cruzar passwords de SaaS com gastos e
permitir revogar acessos/cartões instantaneamente é valor financeiro direto.

## Funcionalidades

- **Monitorização de custos SaaS:** cruza subscrições ativas com passwords
  guardadas; alerta licenças pagas sem uso
- **Cartões virtuais efémeros:** gerar cartão com plafond para uma compra
  específica; autodestrói-se após uso
- **Revogação instantânea:** se o funcionário sai, acesso e cartões são cortados
- **Categorização fiscal** (base para o agente de otimização fiscal)

## Critérios de aceitação

- [x] Painel mostra gasto por licença e sinaliza licenças sem uso (`FIN-001/002`, `/fin`).
- [x] Categorização fiscal com códigos IRC no blob ZK (`FIN-005`, `/fin/fiscal`).
- [ ] Cartão efémero respeita limite e expira (uso único ou prazo) — `FIN-004`.
- [x] Open Banking scaffold com mock em dev; [ ] mTLS TPP em produção (`FIN-003`).

## Conceitos didáticos

> 💡 **Open Banking:** regulação que obriga os bancos a expor APIs seguras para
> terceiros autorizados acederem a dados/iniciarem pagamentos **com consentimento
> do cliente**. Permite-nos conciliar pagamentos sem "raspar" o site do banco.

> 💡 **mTLS (TLS mútuo):** no TLS normal só o servidor prova a sua identidade
> (cadeado do browser). No mTLS, **ambos** os lados apresentam certificado — o
> cliente também prova quem é. Usado entre a CLI/serviços e a API.

> ⚠️ **Segurança:** dados de pagamento e tokens bancários são segredos de alto
> valor → cifrados em repouso e nunca expostos a agentes de IA sem o Guardião.

## Dependências

- Fase 1 (vault) para guardar credenciais SaaS.
- Layer de proxy e agentes (para categorização automática) na Fase 3.

## Tasks

Ver prefixo `FIN-*` em [`../05-tasks/backlog.md`](../05-tasks/backlog.md).
