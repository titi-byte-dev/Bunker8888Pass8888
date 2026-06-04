# Requisitos Não-Funcionais — Segurança

A segurança é o requisito não-funcional **número um** do AegisPass.

## Modelo de ameaça (resumo STRIDE)

> 💡 **Conceito — STRIDE:** checklist de categorias de ameaça: **S**poofing,
> **T**ampering, **R**epudiation, **I**nformation disclosure, **D**enial of
> service, **E**levation of privilege.

| Ameaça | Exemplo | Mitigação |
|---|---|---|
| Spoofing | Fingir-se de outro utilizador | Auth Hash, mTLS, 2FA |
| Tampering | Alterar dados/logs | AES-GCM (autenticado), hash encadeado |
| Repudiation | Negar uma ação | Logs imutáveis assinados |
| Info disclosure | Fuga de dados | Zero-Knowledge, RLS, cifragem em repouso |
| Denial of service | Sobrecarga | Rate limiting, firewall, só VPN exposta |
| Elevation of privilege | Ganhar acesso indevido | Menor privilégio, Guardião, RLS |

## Princípios

1. **Zero-Trust:** negar por omissão; validar todos os pedidos (rede, identidade,
   contexto, dados).
2. **Defesa em profundidade:** 5 barreiras independentes (ver
   [`../04-user-journeys/journey-employee-byod.md`](../04-user-journeys/journey-employee-byod.md)).
3. **Menor privilégio:** cada componente/agente acede só ao mínimo necessário.
4. **Zero-Knowledge:** o servidor nunca tem dados em claro nem a Master Key.
5. **Não inventar cripto:** usar bibliotecas padrão auditadas.

## Requisitos de superfície de ataque

- [ ] Só a porta UDP do WireGuard exposta à internet.
- [ ] SSH só por chave (password login desativado).
- [ ] Segredos fora do código (variáveis de ambiente / cofre), nunca no git.
- [ ] Dependências verificadas (SCA) + SBOM mantido.
- [ ] Headers de segurança no frontend (CSP, HSTS, etc.).

> ⚠️ **Segurança:** o `.gitignore` já bloqueia `.env`, `*.pem`, `*.key` e
> `secrets/`. Mesmo assim, validar com um *secret scanner* no CI.

## Resposta a incidentes

- [ ] Procedimento de revogação rápida (chaves, sessões, cartões).
- [ ] Capacidade de remote wipe em massa.
- [ ] Logs imutáveis permitem investigação forense fiável.
