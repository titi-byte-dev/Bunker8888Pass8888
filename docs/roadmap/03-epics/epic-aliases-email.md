# Epic: Aliases & E-mail "Ghost" — `MAIL`

> **Fase:** 1 (aliases básicos) → 2 (relay + domínio próprio) · **Prioridade:** 🟡 Média

## Objetivo

Proteger a identidade dos funcionários com aliases de e-mail ilimitados e um
relay que mascara remetentes, com suporte a domínio próprio.

## Valor de negócio

Anonimato em pesquisas/compras e proteção da caixa de entrada real contra spam e
fugas. Funcionalidade de "esconder e-mail" totalmente self-hosted.

## Funcionalidades

- **Aliases Hide-my-Email ilimitados** (10 no plano base, ilimitados acima)
- **Domínio personalizado** para aliases (`projeto-x@suaempresa.com`)
- **Caixas de correio adicionais** para aliases
- **Iniciar e-mail a partir do alias** (compor na interface Svelte como remetente)
- Reencaminhamento para o e-mail real; o destinatário nunca vê o endereço original
- Assente num servidor de e-mail open-source self-hosted (SMTP/IMAP)

## Critérios de aceitação

- [ ] Alias gerado encaminha corretamente para o e-mail real.
- [ ] Resposta a partir do alias mantém o endereço real oculto.
- [ ] Domínio próprio configurável (SPF/DKIM/DMARC corretos).
- [ ] Desativar um alias bloqueia imediatamente o reencaminhamento.

## Conceitos didáticos

> 💡 **SPF / DKIM / DMARC:** três mecanismos que provam que um e-mail é
> legítimo. **SPF** diz que servidores podem enviar pelo teu domínio; **DKIM**
> assina o e-mail criptograficamente; **DMARC** define o que fazer se falharem.
> Sem eles, os e-mails de aliases caem em spam.

> ⚠️ **Segurança:** um relay mal configurado pode tornar-se um *open relay*
> (usado por spammers). É obrigatório autenticar todos os envios e limitar taxas.

## Dependências

- Infra de e-mail open-source (SMTP) na VPS.
- Necessário para o agente de prospeção (CRM) na Fase 2.

## Tasks

Ver prefixo `MAIL-*` em [`../05-tasks/backlog.md`](../05-tasks/backlog.md).
