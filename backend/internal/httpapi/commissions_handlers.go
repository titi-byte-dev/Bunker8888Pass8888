package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/commissions"
)

// Comissoes (FIN-007). O conteudo (beneficiario, percentagem, valor) chega ja
// cifrado num blob base64; o servidor so guarda a ligacao a fatura e o estado
// de liquidacao. Nunca decifra.

type createCommissionRequest struct {
	InvoiceID string `json:"invoice_id"`
	Blob      string `json:"blob"`
}

func mapCommissionsError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, commissions.ErrNotFound):
		writeError(w, http.StatusNotFound, "comissão não encontrada")
	case errors.Is(err, commissions.ErrInvalidStatus):
		writeError(w, http.StatusBadRequest, "estado inválido")
	default:
		writeError(w, http.StatusInternalServerError, "falha na operação de comissões")
	}
	return true
}

func commissionJSON(c *commissions.Commission) map[string]any {
	return map[string]any{
		"id":         c.ID,
		"invoice_id": c.InvoiceID,
		"status":     c.Status,
		"blob":       base64.StdEncoding.EncodeToString(c.Blob),
		"created_at": c.CreatedAt,
		"updated_at": c.UpdatedAt,
	}
}

func handleListCommissions(repo *commissions.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		items, err := repo.List(r.Context(), userID)
		if mapCommissionsError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(items))
		for i := range items {
			out = append(out, commissionJSON(&items[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"commissions": out})
	}
}

func handleCreateCommission(repo *commissions.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req createCommissionRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		blob, err := base64.StdEncoding.DecodeString(req.Blob)
		if err != nil || len(blob) == 0 {
			writeError(w, http.StatusBadRequest, "blob base64 inválido")
			return
		}
		c, err := repo.Create(r.Context(), userID, req.InvoiceID, blob)
		if mapCommissionsError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, commissionJSON(c))
	}
}

func handleUpdateCommissionStatus(repo *commissions.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req updateStatusRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		c, err := repo.UpdateStatus(r.Context(), userID, r.PathValue("id"), req.Status)
		if mapCommissionsError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, commissionJSON(c))
	}
}
