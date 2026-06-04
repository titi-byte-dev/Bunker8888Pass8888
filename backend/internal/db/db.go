// Package db trata da ligação ao PostgreSQL e das migrações de esquema.
package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect abre um pool de ligações ao PostgreSQL e confirma que responde.
//
// Didático: um *pgxpool.Pool gere VÁRIAS ligações reutilizáveis (em vez de
// abrir/fechar uma por pedido). É seguro para uso concorrente por goroutines.
func Connect(ctx context.Context, url string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("criar pool: %w", err)
	}
	// Ping verifica que a BD está mesmo acessível (não só que a string é válida).
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping à BD: %w", err)
	}
	return pool, nil
}
