# Requisitos Não-Funcionais — Conformidade

Para além da segurança técnica, o AegisPass tem de cumprir obrigações legais e
organizacionais — é parte da proposta de valor.

## RGPD (foco principal)

| Princípio | Como o produto cumpre |
|---|---|
| Licitude e transparência | Consentimento explícito (OAuth2, termos) |
| Minimização de dados | Guardar só o necessário; cifrar o resto |
| Privacidade por desenho (Art. 25) | Zero-Knowledge desde o início |
| Direito ao esquecimento (Art. 17) | Crypto-shredding + certificado |
| Segurança (Art. 32) | Cifragem, RLS, logs imutáveis |
| Registo de tratamento (Art. 30) | Logs de auditoria imutáveis |

Ver testes em [`../06-testing/rgpd-compliance-tests.md`](../06-testing/rgpd-compliance-tests.md).

> 💡 **Conceito — DPO e DPIA:** o **DPO** (Encarregado de Proteção de Dados) é a
> pessoa responsável pela conformidade; a **DPIA** (Avaliação de Impacto sobre a
> Proteção de Dados) é a análise de risco obrigatória para tratamentos de alto
> risco. São requisitos *organizacionais* que acompanham os técnicos.

## Soberania de dados

- [ ] VPS sob jurisdição clara (idealmente UE para clientes UE).
- [ ] Documentar onde os dados (cifrados) residem fisicamente.
- [ ] Subprocessadores (Google, banco) listados e com base legal.

## Licenciamento e propriedade intelectual

> ⚠️ **Importante:** o produto baseia-se em stack **open-source** sem licenças
> proprietárias. Não incorporamos código proprietário de terceiros — ver
> [`../01-architecture/agents-architecture.md`](../01-architecture/agents-architecture.md).

- [ ] Manter inventário de licenças das dependências (compatibilidade).
- [ ] Verificar licenças de bibliotecas cripto e de e-mail.
- [ ] Garantir que nenhum código de terceiros sem licença entra no produto.

## Acessibilidade e qualidade (futuro)

- [ ] Frontend conforme WCAG 2.1 AA (objetivo de médio prazo).
- [ ] Internacionalização (i18n) — PT/EN inicialmente.

## Checklist de "pronto para vender a empresas"

- [ ] Relatório de conformidade RGPD automático (`HR-008`).
- [ ] Auditoria de segurança externa concluída.
- [ ] Política de retenção e backup documentada.
- [ ] Termos de serviço + acordo de processamento de dados (DPA).
