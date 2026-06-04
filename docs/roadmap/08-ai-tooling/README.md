# 08 — AI Tooling (configuração de assistentes de IA)

Esta secção descreve **como configuramos os assistentes de IA** para desenvolver
o produto. É deliberadamente **tool-agnostic** (independente do fornecedor):
hoje usamos o Cursor, mas os mesmos conceitos aplicam-se a qualquer agente — e,
no futuro, ao **nosso próprio agente** (alinhado com a visão de ERP/CRM com
agentes em [`../03-epics/epic-agents.md`](../03-epics/epic-agents.md)).

## Filosofia

Um bom assistente precisa de 5 tipos de "configuração". Estes conceitos existem
em quase todas as plataformas de agentes com nomes diferentes:

| Conceito (genérico) | O que é | Cursor (hoje) | Forma genérica |
|---|---|---|---|
| **Memória / Contexto** | Orientação sempre presente | `AGENTS.md` | ficheiro de contexto do projeto |
| **Regras / Convenções** | Padrões por tipo de ficheiro | `.cursor/rules/*.mdc` | regras / settings |
| **Comandos** | Ações rápidas repetíveis | `.cursor/commands/*.md` | comandos (slash) |
| **Hooks / Gates** | Automação em eventos | `.cursor/hooks.json` | hooks de eventos |
| **Skills** | Workflows reutilizáveis | `.cursor/skills/*` | skills |
| **Integrações (MCP)** | Acesso a ferramentas externas | `.cursor/mcp.json` | conectores MCP |

> 💡 **Princípio:** preferir configuração **versionada e legível** (texto/markdown)
> a conhecimento "na cabeça" de quem programa. Assim qualquer humano **ou** IA
> que entre no projeto fica produtivo e seguro de imediato.

## O que existe hoje neste repositório

### Memória
- [`/AGENTS.md`](../../../AGENTS.md) — orientação + regras de ouro.

### Regras (`.cursor/rules/`)
- `didactic-style.mdc` — comentários que ensinam (sempre).
- `security-crypto.mdc` — segurança e cripto (sempre).
- `go-conventions.mdc` — Go (`**/*.go`).
- `svelte-conventions.mdc` — Svelte/TS (`**/*.svelte`, `**/*.ts`).
- `rgpd.mdc` — RGPD (ficheiros de RH/compliance).

### Comandos (`.cursor/commands/`)
- `/new-task`, `/new-epic`, `/security-review`, `/crypto-check`.

### Hooks (`.cursor/hooks.json`)
- Gate de escrita: bloqueia segredos e protege `_private/`.
- Gate de shell: pede confirmação em comandos destrutivos/perigosos.
- *(prompt hooks — multiplataforma, sem dependências de binários)*

### Skills (`.cursor/skills/`)
- `table-driven-test`, `scaffold-go-service`, `add-vault-item-type`.

## Princípios de design destes mecanismos

1. **Seguro por omissão:** hooks e regras assumem a escolha mais segura.
2. **Mínimo e legível:** cada ficheiro é curto e focado num só assunto.
3. **Portável:** evitar dependências de SO/binários (daí os *prompt hooks*).
4. **Ensina enquanto faz:** alinhado com o estilo didático do projeto.

## Próximos passos

Ver a lista priorizada de melhorias em
[`future-enhancements.md`](future-enhancements.md).
