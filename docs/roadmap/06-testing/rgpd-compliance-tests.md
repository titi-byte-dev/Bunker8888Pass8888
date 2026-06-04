# Testes — Conformidade RGPD

Casos de teste que provam, de forma automatizada, que o produto cumpre as
exigências de proteção de dados. Servem também de evidência para auditoria.

## Mapa: artigo RGPD → teste

| Artigo | Exigência | Teste |
|---|---|---|
| Art. 5 (minimização) | Recolher só o necessário | Verificar que o servidor não persiste dados em claro |
| Art. 17 (esquecimento) | Apagar a pedido | `erasure` torna dados irrecuperáveis |
| Art. 25 (privacy by design) | Privacidade por desenho | Master Key nunca chega ao servidor |
| Art. 30 (registo de tratamento) | Logs de acesso | Logs imutáveis íntegros |
| Art. 32 (segurança) | Cifragem adequada | Vetores cripto + cifragem em repouso |

## Casos de teste

### CT-RGPD-01 — Servidor nunca vê dados em claro

```gherkin
Dado que um funcionário guarda o salário "2500 €"
Quando o dado é enviado para a API
Então a base de dados contém apenas um blob cifrado
E nenhum log do servidor contém "2500"
```

### CT-RGPD-02 — Direito ao esquecimento (crypto-shredding)

```gherkin
Dado um funcionário com ficha cifrada
Quando o admin executa o erasure (Art. 17)
Então a chave do funcionário é destruída
E os dados pessoais tornam-se indecifráveis
E é gerado um certificado de eliminação verificável
E a cadeia de logs imutáveis permanece íntegra
```

### CT-RGPD-03 — Integridade dos logs de auditoria

```gherkin
Dado um conjunto de entradas de log encadeadas por hash
Quando se altera uma entrada antiga
Então a verificação da cadeia deteta a adulteração
```

### CT-RGPD-04 — Triagem às cegas (não discriminação)

```gherkin
Dado um currículo com dados de género/etnia
Quando é apresentado ao agente de recrutamento
Então esses campos estão ocultos/mascarados
```

> 💡 **Conceito — Gherkin (Given/When/Then):** forma legível de escrever testes
> de comportamento. "Dado" (contexto), "Quando" (ação), "Então" (resultado
> esperado). Ferramentas como Cucumber/godog executam isto diretamente.

## Evidência para auditoria

- [ ] Relatório de conformidade gerado automaticamente (ver `HR-008`).
- [ ] Resultados dos CT-RGPD-* exportáveis como prova.
- [ ] Registo de DPIA (avaliação de impacto) mantido à parte.

> ⚠️ **Nota legal:** estes testes ajudam a *demonstrar* conformidade técnica,
> mas não substituem aconselhamento jurídico nem um DPO. A conformidade RGPD é
> também organizacional, não só técnica.
