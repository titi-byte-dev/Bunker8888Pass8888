---
name: scaffold-go-service
description: Cria a estrutura base de um serviço/módulo Go no backend do AegisPass, seguindo a arquitetura por camadas (handler → service → store) e as regras de segurança. Usar ao iniciar um novo módulo backend.
disable-model-invocation: true
---

# Scaffold de Serviço Go — AegisPass

Cria um módulo backend coeso, com separação de responsabilidades e segurança
desde o início.

## Estrutura por módulo

```
backend/internal/<modulo>/
├── handler.go     # camada HTTP/WS: parse, validação, resposta
├── service.go     # lógica de negócio (sem detalhes de transporte nem SQL)
├── store.go       # acesso a dados (PostgreSQL), sempre com tenant_id + RLS
├── types.go       # structs do domínio
└── service_test.go / store_test.go  # testes (ver skill table-driven-test)
```

## Regras ao gerar

1. **Camadas:** o `handler` não fala com a BD diretamente; chama o `service`,
   que chama o `store`.
2. **context.Context** como primeiro parâmetro em funções de I/O.
3. **tenant_id:** toda a função do `store` recebe/usa o tenant e confia na RLS.
4. **Erros** embrulhados com `%w`; nada de `panic` em fluxo normal.
5. **Segredos/cripto:** seguir `.cursor/rules/security-crypto.mdc`.
6. **Sem dados sensíveis em logs.**

## Esqueleto do store (exemplo)

```go
package <modulo>

import "context"

type Store struct{ db DB }

func (s *Store) PorID(ctx context.Context, tenantID, id string) (Item, error) {
    // SET LOCAL app.tenant_id é definido no início da transação (middleware),
    // por isso a RLS já restringe as linhas a este tenant.
    // ...
}
```

## Passos

1. Confirmar o nome do módulo e a que epic pertence.
2. Criar os ficheiros acima com stubs comentados (estilo didático).
3. Criar os testes correspondentes.
4. Indicar o(s) ID(s) de task do backlog que este módulo cobre.
