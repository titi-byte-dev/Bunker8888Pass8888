package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/invoicing"
)

// Faturacao (FIN-006). O conteudo do documento (cliente, linhas, valores) chega
// ja cifrado num blob base64; o servidor so atribui o numero legal e gere o
// estado. Nunca decifra.

type issueInvoiceRequest struct {
	DocType      string `json:"doc_type"`
	SourceLeadID string `json:"source_lead_id"`
	Blob         string `json:"blob"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func mapInvoicingError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, invoicing.ErrNotFound):
		writeError(w, http.StatusNotFound, "documento não encontrado")
	case errors.Is(err, invoicing.ErrInvalidType):
		writeError(w, http.StatusBadRequest, "tipo de documento inválido")
	case errors.Is(err, invoicing.ErrInvalidStatus):
		writeError(w, http.StatusBadRequest, "estado inválido")
	default:
		writeError(w, http.StatusInternalServerError, "falha na operação de faturação")
	}
	return true
}

func invoiceJSON(d *invoicing.Document) map[string]any {
	return map[string]any{
		"id":             d.ID,
		"doc_type":       d.DocType,
		"year":           d.Year,
		"seq":            d.Seq,
		"number":         d.Number,
		"status":         d.Status,
		"source_lead_id": d.SourceLeadID,
		"blob":           base64.StdEncoding.EncodeToString(d.Blob),
		"created_at":     d.CreatedAt,
		"updated_at":     d.UpdatedAt,
	}
}

func handleListInvoices(repo *invoicing.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		docs, err := repo.List(r.Context(), userID)
		if mapInvoicingError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(docs))
		for i := range docs {
			out = append(out, invoiceJSON(&docs[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"invoices": out})
	}
}

func handleIssueInvoice(repo *invoicing.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req issueInvoiceRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		blob, err := base64.StdEncoding.DecodeString(req.Blob)
		if err != nil || len(blob) == 0 {
			writeError(w, http.StatusBadRequest, "blob base64 inválido")
			return
		}
		d, err := repo.Issue(r.Context(), userID, req.DocType, req.SourceLeadID, blob)
		if mapInvoicingError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, invoiceJSON(d))
	}
}

func handleUpdateInvoiceStatus(repo *invoicing.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req updateStatusRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		d, err := repo.UpdateStatus(r.Context(), userID, r.PathValue("id"), req.Status)
		if mapInvoicingError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, invoiceJSON(d))
	}
}
