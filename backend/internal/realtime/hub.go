// Package realtime gere notificações push via WebSocket (VAULT-006).
//
// Didático: WebSocket mantém um canal bidirecional aberto entre cliente e
// servidor. Ao contrário do HTTP (pedido→resposta), o servidor pode enviar
// mensagens a qualquer momento — ideal para sincronizar o cofre entre
// dispositivos sem polling.
//
// ⚠️ Zero-Knowledge: os eventos só transportam METADADOS (id, tipo, data).
// Nunca enviamos blobs cifrados nem conteúdo em claro pelo WebSocket; o
// cliente faz GET /api/vault/{id} se precisar do blob actualizado.
package realtime

import (
	"encoding/json"
	"sync"
	"time"
)

// Tipos de evento do cofre (contrato com o frontend).
const (
	EventCreated = "vault.item.created"
	EventUpdated = "vault.item.updated"
	EventDeleted = "vault.item.deleted"
	// EventRemoteWipe — ordem de apagar cache local + descartar Master Key (VAULT-012).
	EventRemoteWipe = "security.remote_wipe"
)

// Event é a mensagem JSON enviada aos clientes ligados.
type Event struct {
	Type      string    `json:"type"`
	ItemID    string    `json:"item_id"`
	ItemType  string    `json:"item_type,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// WipeEvent é enviado aos dispositivos ligados para apagar dados locais da app.
type WipeEvent struct {
	Type     string    `json:"type"`
	Reason   string    `json:"reason,omitempty"`
	IssuedAt time.Time `json:"issued_at"`
}

// NotifyWipe envia ordem de remote wipe a todos os dispositivos do utilizador.
func (h *Hub) NotifyWipe(userID, reason string) {
	ev := WipeEvent{
		Type:     EventRemoteWipe,
		Reason:   reason,
		IssuedAt: time.Now().UTC(),
	}
	payload, err := json.Marshal(ev)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients[userID] {
		c.trySend(payload)
	}
}

// Hub regista ligações WebSocket por userID e faz broadcast de eventos.
type Hub struct {
	mu      sync.RWMutex
	clients map[string]map[*Client]struct{}
}

// NewHub cria um hub vazio.
func NewHub() *Hub {
	return &Hub{clients: make(map[string]map[*Client]struct{})}
}

// Register associa um cliente a um utilizador.
func (h *Hub) Register(userID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[userID] == nil {
		h.clients[userID] = make(map[*Client]struct{})
	}
	h.clients[userID][c] = struct{}{}
}

// Unregister remove um cliente (chamado ao fechar a ligação).
func (h *Hub) Unregister(userID string, c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[userID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.clients, userID)
		}
	}
}

// Notify envia um evento a todos os dispositivos ligados do utilizador.
func (h *Hub) Notify(userID string, ev Event) {
	payload, err := json.Marshal(ev)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients[userID] {
		c.trySend(payload)
	}
}

// ClientCount devolve quantos clientes estão ligados (útil em testes).
func (h *Hub) ClientCount(userID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[userID])
}
