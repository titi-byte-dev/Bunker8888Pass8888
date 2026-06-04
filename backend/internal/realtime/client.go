package realtime

import (
	"sync"
	"time"
)

const (
	SendBuf       = 16
	WriteWait     = 10 * time.Second
	PongWait      = 60 * time.Second
	PingPeriod    = (PongWait * 9) / 10
	maxMessageLen = 512
)

// Client representa uma ligação WebSocket de um dispositivo.
type Client struct {
	send chan []byte
	once sync.Once
}

// NewClient cria um cliente com fila de envio.
func NewClient() *Client {
	return &Client{send: make(chan []byte, SendBuf)}
}

// Send devolve o canal onde o hub escreve mensagens JSON.
func (c *Client) Send() <-chan []byte {
	return c.send
}

// Close fecha o canal de envio (idempotente).
func (c *Client) Close() {
	c.once.Do(func() { close(c.send) })
}

// trySend tenta enfileirar sem bloquear; descarta se a fila estiver cheia.
func (c *Client) trySend(msg []byte) {
	select {
	case c.send <- msg:
	default:
		// Cliente lento — ignoramos este evento (defesa contra backpressure).
	}
}
