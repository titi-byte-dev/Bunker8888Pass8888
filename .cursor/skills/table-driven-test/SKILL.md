---
name: table-driven-test
description: Escreve testes table-driven idiomáticos em Go para o AegisPass, com sub-testes nomeados e casos de borda. Usar ao criar ou pedir testes de funções Go.
disable-model-invocation: true
---

# Table-Driven Test (Go)

Padrão idiomático de testes em Go: uma tabela de casos + iteração com `t.Run`.

## Template

```go
func TestNomeDaFuncao(t *testing.T) {
    casos := []struct {
        nome    string
        entrada TipoEntrada
        querErr bool
        espera  TipoSaida
    }{
        {nome: "caso normal", entrada: ..., espera: ...},
        {nome: "caso de borda (vazio)", entrada: ..., espera: ...},
        {nome: "input inválido", entrada: ..., querErr: true},
    }
    for _, c := range casos {
        // t.Run cria um sub-teste nomeado → output claro sobre QUAL caso falhou.
        t.Run(c.nome, func(t *testing.T) {
            got, err := NomeDaFuncao(c.entrada)
            if (err != nil) != c.querErr {
                t.Fatalf("erro = %v, querErr = %v", err, c.querErr)
            }
            if c.querErr {
                return // esperávamos erro; não comparar saída
            }
            // bytes.Equal para []byte; reflect.DeepEqual para structs/slices.
            if got != c.espera {
                t.Errorf("esperado %v, obtido %v", c.espera, got)
            }
        })
    }
}
```

## Regras

- Incluir sempre pelo menos: caso normal, caso de borda, e caso de erro.
- Para cripto, incluir um caso com **test vectors** conhecidos.
- Comparar `[]byte` com `bytes.Equal` (o operador `==` não funciona em slices).
- Nomes de casos descritivos (aparecem no output do teste).
