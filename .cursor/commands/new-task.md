# /new-task

Adiciona uma nova task ao backlog do roadmap, com o ID e formato corretos.

## Passos

1. Lê [`docs/roadmap/05-tasks/backlog.md`](../../docs/roadmap/05-tasks/backlog.md).
2. Identifica o epic/prefixo correto (`VAULT`, `HR`, `FIN`, `SHARE`, `DW`,
   `MAIL`, `AGENT`, `GOOGLE`, `INFRA`).
3. Atribui o **próximo número livre** desse prefixo.
4. Acrescenta uma linha à tabela respetiva no formato:
   `| ID | Descrição | Fase | Tamanho (S/M/L/XL) | ⚪ | Depende de |`
5. Se a descrição da task implicar dependências, preenche a coluna "Depende de".
6. Confirma com um resumo do ID criado.

## Notas

- Não renumerar tasks existentes.
- Manter a descrição curta e acionável (uma linha).
- Se o utilizador não indicar fase/tamanho, infere a partir do epic e propõe.
