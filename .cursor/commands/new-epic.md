# /new-epic

Cria um novo epic em `docs/roadmap/03-epics/` seguindo o template padrão.

## Passos

1. Pergunta (se não foi dado): nome do epic, prefixo de ID e fase.
2. Cria `docs/roadmap/03-epics/epic-<nome>.md` com as secções:
   - Título + Fase + Prioridade
   - **Objetivo**
   - **Valor de negócio**
   - **Funcionalidades**
   - **Critérios de aceitação** (checklist)
   - **Conceitos didáticos** (`> 💡` e `> ⚠️` conforme a regra de estilo)
   - **Dependências**
   - **Tasks** (aponta para o prefixo no backlog)
3. Adiciona o epic à tabela do README do roadmap e cria a secção do prefixo no
   `05-tasks/backlog.md` se ainda não existir.

## Referência

Usa [`epic-vault.md`](../../docs/roadmap/03-epics/epic-vault.md) como modelo de
estrutura e tom.
