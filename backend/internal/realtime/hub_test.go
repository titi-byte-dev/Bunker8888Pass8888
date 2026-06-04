package realtime

import (
	"testing"
	"time"
)

func TestHub_NotifyDelivered(t *testing.T) {
	h := NewHub()
	c := NewClient()
	defer c.Close()
	h.Register("user-1", c)

	h.Notify("user-1", Event{
		Type:     EventCreated,
		ItemID:   "abc",
		ItemType: "login",
	})

	select {
	case msg := <-c.Send():
		if string(msg) == "" {
			t.Fatal("mensagem vazia")
		}
	case <-time.After(time.Second):
		t.Fatal("evento não entregue")
	}

	h.Unregister("user-1", c)
	if h.ClientCount("user-1") != 0 {
		t.Fatal("cliente devia estar desregistado")
	}
}

func TestHub_NotifyWipeDelivered(t *testing.T) {
	h := NewHub()
	c := NewClient()
	defer c.Close()
	h.Register("user-1", c)

	h.NotifyWipe("user-1", "offboarding")

	select {
	case msg := <-c.Send():
		if string(msg) == "" {
			t.Fatal("mensagem wipe vazia")
		}
	case <-time.After(time.Second):
		t.Fatal("wipe não entregue")
	}
}

func TestHub_NotifyOtherUserIgnored(t *testing.T) {
	h := NewHub()
	c := NewClient()
	defer c.Close()
	h.Register("user-a", c)

	h.Notify("user-b", Event{Type: EventDeleted, ItemID: "x"})

	select {
	case <-c.Send():
		t.Fatal("user-a não devia receber evento de user-b")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestClient_trySend_DropWhenFull(t *testing.T) {
	c := NewClient()
	defer c.Close()
	// Enchemos a fila (SendBuf=16).
	for i := 0; i < SendBuf; i++ {
		c.trySend([]byte("x"))
	}
	// A 17ª não deve bloquear nem panic.
	c.trySend([]byte("y"))
}
