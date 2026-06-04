// Package emergency implementa acesso de emergência com herdeiro digital (VAULT-016).
//
// Didático: o servidor gere o workflow (pedido, espera, aprovar/rejeitar) mas
// nunca vê a Master Key — só armazena um blob cifrado pelo titular no cliente.
package emergency

import (
	"errors"
	"time"
)

var (
	ErrNotConfigured   = errors.New("emergency: herdeiro não configurado")
	ErrNotFound        = errors.New("emergency: pedido não encontrado")
	ErrForbidden       = errors.New("emergency: operação não permitida")
	ErrActiveRequest   = errors.New("emergency: já existe pedido activo")
	ErrNotReady        = errors.New("emergency: acesso ainda não disponível")
	ErrSelfHeir        = errors.New("emergency: não podes ser o teu próprio herdeiro")
	ErrHeirMismatch    = errors.New("emergency: email não corresponde ao herdeiro configurado")
	ErrNoEncryptedBlob = errors.New("emergency: blob de emergência não configurado")
)

// Status do pedido de acesso de emergência.
type Status string

const (
	StatusWaiting  Status = "waiting"
	StatusRejected Status = "rejected"
	StatusReady    Status = "ready"
	StatusConsumed Status = "consumed"
)

// Config guarda o herdeiro designado pelo titular.
type Config struct {
	HeirEmail string    `json:"heir_email"`
	WaitDays  int       `json:"wait_days"`
	HasBlob   bool      `json:"has_blob"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Request representa um pedido de acesso iniciado pelo herdeiro.
type Request struct {
	ID          string     `json:"id"`
	OwnerUserID string     `json:"owner_user_id"`
	HeirUserID  string     `json:"heir_user_id"`
	HeirEmail   string     `json:"heir_email"`
	Status      Status     `json:"status"`
	RequestedAt time.Time  `json:"requested_at"`
	UnlocksAt   time.Time  `json:"unlocks_at"`
	RejectedAt  *time.Time `json:"rejected_at,omitempty"`
	ConsumedAt  *time.Time `json:"consumed_at,omitempty"`
}

// PromoteIfExpired actualiza waiting→ready quando o período de espera terminou.
func PromoteIfExpired(now time.Time, status Status, unlocksAt time.Time) Status {
	if status == StatusWaiting && !now.Before(unlocksAt) {
		return StatusReady
	}
	return status
}
