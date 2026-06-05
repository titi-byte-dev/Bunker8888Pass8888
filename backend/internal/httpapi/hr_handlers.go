package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/hr"
)

// Fichas de Empregado com cifragem campo-a-campo (HR-001). Tudo opaco: o cliente
// cifra cada campo (value_blob) com uma chave própria e envia também essa chave
// cifrada com a Master Key (wrapped_key), ambas em base64. O servidor só guarda
// e faz cumprir a posse da ficha — nunca decifra.

type putFieldRequest struct {
	ValueBlob  string `json:"value_blob"`  // AES-GCM(chave_campo, valor)
	WrappedKey string `json:"wrapped_key"` // AES-GCM(master_key, chave_campo)
}

// mapHRError traduz erros de domínio em códigos HTTP.
func mapHRError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, hr.ErrNotFound):
		writeError(w, http.StatusNotFound, "ficha não encontrada")
	case errors.Is(err, hr.ErrInvalidField):
		writeError(w, http.StatusBadRequest, "campo inválido (valor ou chave em falta)")
	default:
		writeError(w, http.StatusInternalServerError, "falha na operação da ficha")
	}
	return true
}

func recordJSON(rec *hr.Record) map[string]any {
	return map[string]any{
		"id":         rec.ID,
		"owner_id":   rec.OwnerID,
		"created_at": rec.CreatedAt,
		"updated_at": rec.UpdatedAt,
	}
}

func fieldJSON(f *hr.Field) map[string]any {
	m := map[string]any{
		"id":         f.ID,
		"field_name": f.FieldName,
		"value_blob": base64.StdEncoding.EncodeToString(f.ValueBlob),
		"shredded":   f.WrappedKey == nil,
		"created_at": f.CreatedAt,
		"updated_at": f.UpdatedAt,
	}
	// wrapped_key só vai quando existe (um campo crypto-shredded não tem chave).
	if f.WrappedKey != nil {
		m["wrapped_key"] = base64.StdEncoding.EncodeToString(f.WrappedKey)
	}
	if f.ShreddedAt != nil {
		m["shredded_at"] = *f.ShreddedAt
	}
	return m
}

func handleCreateEmployeeRecord(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		rec, err := repo.CreateRecord(r.Context(), userID)
		if mapHRError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, recordJSON(rec))
	}
}

func handleListEmployeeRecords(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		recs, err := repo.ListRecords(r.Context(), userID)
		if mapHRError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(recs))
		for i := range recs {
			out = append(out, recordJSON(&recs[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"records": out})
	}
}

func handleGetEmployeeRecord(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		rec, fields, err := repo.GetRecord(r.Context(), userID, r.PathValue("id"))
		if mapHRError(w, err) {
			return
		}
		fs := make([]map[string]any, 0, len(fields))
		for i := range fields {
			fs = append(fs, fieldJSON(&fields[i]))
		}
		out := recordJSON(rec)
		out["fields"] = fs
		writeJSON(w, http.StatusOK, out)
	}
}

func handleDeleteEmployeeRecord(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := repo.DeleteRecord(r.Context(), userID, r.PathValue("id")); mapHRError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handlePutEmployeeField(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req putFieldRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		value, err := base64.StdEncoding.DecodeString(req.ValueBlob)
		if err != nil || len(value) == 0 {
			writeError(w, http.StatusBadRequest, "value_blob base64 inválido")
			return
		}
		wrapped, err := base64.StdEncoding.DecodeString(req.WrappedKey)
		if err != nil || len(wrapped) == 0 {
			writeError(w, http.StatusBadRequest, "wrapped_key base64 inválido")
			return
		}
		f, err := repo.PutField(r.Context(), userID, r.PathValue("id"), r.PathValue("field"), value, wrapped)
		if mapHRError(w, err) {
			return
		}
		f.FieldName = r.PathValue("field")
		f.ValueBlob = value
		f.WrappedKey = wrapped
		writeJSON(w, http.StatusOK, fieldJSON(f))
	}
}

func handleDeleteEmployeeField(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := repo.DeleteField(r.Context(), userID, r.PathValue("id"), r.PathValue("field")); mapHRError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
