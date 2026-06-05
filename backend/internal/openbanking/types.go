package openbanking

import "time"

// ConnectionStatus reflecte o ciclo de consentimento PSD2/Open Banking.
type ConnectionStatus string

const (
	StatusPending   ConnectionStatus = "pending"
	StatusConnected ConnectionStatus = "connected"
	StatusExpired   ConnectionStatus = "expired"
)

// Connection metadados persistidos — sem tokens bancários em claro.
type Connection struct {
	ID               string
	OwnerID          string
	Provider         string
	Status           ConnectionStatus
	ConsentExpiresAt *time.Time
	LastSyncAt       *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// Transaction representa um movimento devolvido pelo provider (mock em dev).
// Em produção o cliente cifra estes campos antes de guardar localmente.
type Transaction struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	BookedAt    time.Time `json:"booked_at"`
	Description string    `json:"description"`
	MerchantRef string    `json:"merchant_ref,omitempty"`
}
