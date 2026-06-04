package emergency

import (
	"testing"
	"time"
)

func TestPromoteIfExpired(t *testing.T) {
	unlock := time.Date(2026, 6, 10, 12, 0, 0, 0, time.UTC)

	if got := PromoteIfExpired(unlock.Add(-time.Hour), StatusWaiting, unlock); got != StatusWaiting {
		t.Fatalf("antes do prazo: got %s", got)
	}
	if got := PromoteIfExpired(unlock, StatusWaiting, unlock); got != StatusReady {
		t.Fatalf("no prazo: got %s", got)
	}
	if got := PromoteIfExpired(unlock.Add(time.Hour), StatusWaiting, unlock); got != StatusReady {
		t.Fatalf("após prazo: got %s", got)
	}
	if got := PromoteIfExpired(unlock.Add(time.Hour), StatusRejected, unlock); got != StatusRejected {
		t.Fatalf("rejected inalterado: got %s", got)
	}
}
