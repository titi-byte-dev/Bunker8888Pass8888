package burnnotes

import (
	"context"
	"time"
)

// StartReaper corre um ciclo em segundo plano que remove notas expiradas a cada
// `interval`, para que notas por ler nao fiquem em RAM mais do que o necessario
// (mesmo que ninguem as volte a tocar). Termina quando o context é cancelado.
// Pensado para ser chamado uma vez, com `go` no arranque.
func (s *Store) StartReaper(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.Reap(time.Now())
		}
	}
}
