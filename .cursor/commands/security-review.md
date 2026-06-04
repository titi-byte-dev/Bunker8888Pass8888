# /security-review

Revê as alterações atuais (ou um ficheiro indicado) contra o modelo de ameaça do
AegisPass.

## Passos

1. Obtém o diff (`git diff` se houver git; senão, revê os ficheiros indicados).
2. Avalia contra o threat model STRIDE em
   [`docs/roadmap/07-non-functional/security.md`](../../docs/roadmap/07-non-functional/security.md)
   e as regras em `.cursor/rules/security-crypto.mdc` e `rgpd.mdc`.
3. Verifica especificamente:
   - [ ] Sem segredos hardcoded (chaves, tokens, `.env`, IBAN/NIF reais).
   - [ ] Cripto usa bibliotecas padrão; nonce único; sem chaves no servidor.
   - [ ] Queries filtram por `tenant_id` e dependem de RLS.
   - [ ] Input externo validado; conteúdo de agentes tratado como dados.
   - [ ] Dados pessoais cifrados; logs sem dados sensíveis.

## Formato do resultado

- 🔴 **Crítico:** corrigir antes de avançar
- 🟡 **Sugestão:** considerar melhorar
- 🟢 **Bom:** prática correta a destacar

Não fazer alterações automaticamente sem confirmação — apresentar o relatório.
