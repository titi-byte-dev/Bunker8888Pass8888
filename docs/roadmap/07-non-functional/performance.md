# Requisitos Não-Funcionais — Performance

A escolha Go + Svelte existe precisamente para correr muito numa VPS modesta.

## Metas (SLOs)

> 💡 **Conceito — SLO (Service Level Objective):** meta mensurável de qualidade
> de serviço (ex: "99% dos logins em < 300ms"). Guia decisões e alarmes.

| Métrica | Meta |
|---|---|
| Sincronização de alteração (WebSocket) | < 1s |
| Revogação de acesso refletida no cliente | < 1s |
| Latência de login (p95) | < 300ms (excl. Argon2id) |
| Derivação Argon2id (cliente) | 0.5–1s (deliberadamente lento) |
| Tenants por VPS modesta | centenas (multi-tenancy) |
| Arranque da API (cold start) | < 2s |

> ⚠️ **Nota:** o Argon2id ser "lento" é uma *feature* de segurança, não um bug.
> O alvo é ~0.5–1s no dispositivo do utilizador — suficientemente lento para
> travar força bruta, suficientemente rápido para não irritar.

## Porque Go + Svelte ajudam

- **Goroutines:** concorrência barata → muitas ligações WebSocket simultâneas.
- **Binários estáticos:** arranque rápido, baixo footprint de memória.
- **Svelte:** compila para JS sem runtime pesado (vs React/Angular) → cripto
  client-side corre à velocidade máxima.

```go
// Exemplo: servir muitos clientes em simultâneo. Cada ligação corre na sua
// goroutine; o `select` espera por vários canais ao mesmo tempo sem bloquear.
func (h *Hub) servirCliente(c *Cliente) {
    for {
        select {
        case msg := <-c.envio: // mensagem para enviar a este cliente
            c.conn.WriteJSON(msg)
        case <-c.fechar: // sinal de desligar (ex: revogação)
            return
        }
    }
}
```

## Estratégias

- [ ] Connection pooling para o PostgreSQL.
- [ ] Índices adequados (incluindo blind indexing onde aplicável).
- [ ] Caching em memória de dados não sensíveis (com TTL).
- [ ] Testes de carga antes de cada release maior.

## Observabilidade

- [ ] Métricas (latência, throughput, erros) — Prometheus.
- [ ] Tracing distribuído — OpenTelemetry.
- [ ] Alarmes ligados aos SLOs acima.
