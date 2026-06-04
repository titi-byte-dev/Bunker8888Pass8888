package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/emergency"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
)

type emergencyConfigRequest struct {
	HeirEmail string `json:"heir_email"`
	WaitDays  int    `json:"wait_days"`
	Blob      string `json:"blob,omitempty"`
}

type emergencyRequestBody struct {
	OwnerEmail string `json:"owner_email"`
}

func handlePutEmergencyConfig(svc *emergency.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req emergencyConfigRequest
		if err := decodeJSON(r, &req); err != nil || req.HeirEmail == "" {
			writeError(w, http.StatusBadRequest, "heir_email em falta")
			return
		}
		waitDays := req.WaitDays
		if waitDays == 0 {
			waitDays = 7
		}
		var blob []byte
		if req.Blob != "" {
			var err error
			blob, err = base64.StdEncoding.DecodeString(req.Blob)
			if err != nil {
				writeError(w, http.StatusBadRequest, "blob base64 inválido")
				return
			}
		}
		if err := svc.SetConfig(r.Context(), userID, req.HeirEmail, waitDays, blob); err != nil {
			mapEmergencyError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleGetEmergencyConfig(svc *emergency.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		cfg, err := svc.GetConfig(r.Context(), userID)
		if errors.Is(err, emergency.ErrNotConfigured) {
			writeJSON(w, http.StatusOK, map[string]any{"configured": false})
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"configured": true,
			"heir_email": cfg.HeirEmail,
			"wait_days":  cfg.WaitDays,
			"has_blob":   cfg.HasBlob,
			"updated_at": cfg.UpdatedAt.UTC().Format(timeRFC3339),
		})
	}
}

func handleDeleteEmergencyConfig(svc *emergency.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := svc.DeleteConfig(r.Context(), userID); err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleListEmergencyRequests(svc *emergency.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		reqs, err := svc.ListRequestsForOwner(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"requests": serializeRequests(reqs)})
	}
}

func handleRejectEmergencyRequest(svc *emergency.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := r.PathValue("id")
		if err := svc.RejectRequest(r.Context(), userID, id); err != nil {
			mapEmergencyError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
	}
}

func handleApproveEmergencyRequest(svc *emergency.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		id := r.PathValue("id")
		if err := svc.ApproveEarly(r.Context(), userID, id); err != nil {
			mapEmergencyError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}

func handleCreateEmergencyRequest(svc *emergency.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		heirID, _ := r.Context().Value(userIDKey).(string)
		var req emergencyRequestBody
		if err := decodeJSON(r, &req); err != nil || req.OwnerEmail == "" {
			writeError(w, http.StatusBadRequest, "owner_email em falta")
			return
		}
		created, err := svc.CreateRequest(r.Context(), heirID, req.OwnerEmail)
		if err != nil {
			mapEmergencyError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"request": serializeRequest(*created)})
	}
}

func handleGetEmergencyRequestStatus(svc *emergency.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		heirID, _ := r.Context().Value(userIDKey).(string)
		ownerEmail := r.URL.Query().Get("owner_email")
		if ownerEmail == "" {
			writeError(w, http.StatusBadRequest, "owner_email em falta")
			return
		}
		req, err := svc.GetHeirRequest(r.Context(), heirID, ownerEmail)
		if errors.Is(err, emergency.ErrNotFound) {
			writeJSON(w, http.StatusOK, map[string]any{"active": false})
			return
		}
		if err != nil {
			mapEmergencyError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"active":  true,
			"request": serializeRequest(*req),
		})
	}
}

func handleFetchEmergencyAccess(svc *emergency.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		heirID, _ := r.Context().Value(userIDKey).(string)
		ownerEmail := r.URL.Query().Get("owner_email")
		if ownerEmail == "" {
			writeError(w, http.StatusBadRequest, "owner_email em falta")
			return
		}
		blob, err := svc.FetchAccessBlob(r.Context(), heirID, ownerEmail)
		if err != nil {
			mapEmergencyError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{
			"blob": base64.StdEncoding.EncodeToString(blob),
		})
	}
}

const timeRFC3339 = "2006-01-02T15:04:05Z07:00"

func serializeRequests(reqs []emergency.Request) []map[string]any {
	out := make([]map[string]any, len(reqs))
	for i, req := range reqs {
		out[i] = serializeRequest(req)
	}
	return out
}

func serializeRequest(req emergency.Request) map[string]any {
	m := map[string]any{
		"id":           req.ID,
		"heir_email":   req.HeirEmail,
		"status":       req.Status,
		"requested_at": req.RequestedAt.UTC().Format(timeRFC3339),
		"unlocks_at":   req.UnlocksAt.UTC().Format(timeRFC3339),
	}
	if req.RejectedAt != nil {
		m["rejected_at"] = req.RejectedAt.UTC().Format(timeRFC3339)
	}
	if req.ConsumedAt != nil {
		m["consumed_at"] = req.ConsumedAt.UTC().Format(timeRFC3339)
	}
	return m
}

func mapEmergencyError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, emergency.ErrNotConfigured),
		errors.Is(err, emergency.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, emergency.ErrForbidden),
		errors.Is(err, emergency.ErrHeirMismatch):
		writeError(w, http.StatusForbidden, err.Error())
	case errors.Is(err, emergency.ErrActiveRequest):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, emergency.ErrNotReady):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, emergency.ErrSelfHeir):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, users.ErrNotFound):
		writeError(w, http.StatusNotFound, "titular não encontrado")
	default:
		writeError(w, http.StatusInternalServerError, "erro interno")
	}
}
