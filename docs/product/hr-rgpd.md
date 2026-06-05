---
title: RH e conformidade RGPD
slug: hr-rgpd
category: product
order: 5
audience: [user, admin]
layer: [frontend, backend]
feature: hr
level: 1
in_app: true
summary: Fichas cifradas, contratos, onboarding e direito ao esquecimento.
related: [glossary, journey-rgpd-erasure, journey-admin-onboarding]
---

:::summary
O módulo **RH** trata dados pessoais de colaboradores com **cifragem campo-a-campo** e
auditoria imutável — conformidade técnica alinhada com o RGPD por desenho.
:::

:::concept{id="field-encryption" title="Cifragem campo-a-campo" level=1}
Campos sensíveis (NIF, IBAN, salário, saúde) cifram-se no cliente **antes** de
chegarem à API. A base de dados guarda ciphertext; só quem tem a chave adequada
decifra no browser autorizado.
:::

:::concept{id="crypto-shredding" title="Crypto-shredding (Art. 17)" level=2}
Para o **direito ao esquecimento**, destruímos a chave do titular. Os blobs tornam-se
irrecuperáveis sem reescrever backups históricos — gera-se um **certificado
criptográfico** de eliminação para auditoria.
:::

:::level{level=1 title="Funcionalidades na app"}
- `/hr` — fichas de empregado
- `/hr/onboarding` — wizard de integração
- `/hr/compliance` — relatório RGPD (PDF)
:::

:::level{level=2 title="Aprofundar: logs imutáveis"}
Acessos a dados sensíveis registam-se com **hashing encadeado** — alterar uma linha
quebra a cadeia e é detectável. Os logs nunca contêm o dado em si, só metadados.
:::

:::level{level=3 title="Técnico"}
`backend/internal/hr/`, migrações `employees`, `audit_logs`. Testes RGPD em
`docs/roadmap/06-testing/rgpd-compliance-tests.md`.
:::
