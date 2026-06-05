package httpapi

import (
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/geofence"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/security"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/shifts"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
)

// requireAdminKey valida X-Admin-Key antes de executar handlers administrativos.
//
// ⚠️ Segurança: a chave admin é um segredo de operador — nunca vai para o
// frontend em localStorage permanente; a UI guarda-a só em sessionStorage.
func requireAdminKey(adminKey string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if adminKey == "" {
			writeError(w, http.StatusServiceUnavailable, "admin desactivado")
			return
		}
		if r.Header.Get("X-Admin-Key") != adminKey {
			writeError(w, http.StatusForbidden, "chave admin inválida")
			return
		}
		next(w, r)
	}
}

func handleAdminListUsers(userRepo *users.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := userRepo.ListSummaries(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao listar utilizadores")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"users": list})
	}
}

func handleAdminGetUser(userRepo *users.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			writeError(w, http.StatusBadRequest, "user id em falta")
			return
		}
		u, err := userRepo.ByID(r.Context(), id)
		if err != nil {
			if err == users.ErrNotFound {
				writeError(w, http.StatusNotFound, "utilizador não encontrado")
				return
			}
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"id":    u.ID,
			"email": u.Email,
		})
	}
}

func handleAdminGetAccessShift(shiftsRepo *shifts.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetID := r.PathValue("id")
		if targetID == "" {
			writeError(w, http.StatusBadRequest, "user id em falta")
			return
		}
		p, err := shiftsRepo.Get(r.Context(), targetID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"enabled":                p.Enabled,
			"timezone":               p.Timezone,
			"schedule":               p.Schedule,
			"max_clock_skew_seconds": p.MaxClockSkewSecs,
		})
	}
}

func handleAdminGetAccessGeofence(geoRepo *geofence.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetID := r.PathValue("id")
		if targetID == "" {
			writeError(w, http.StatusBadRequest, "user id em falta")
			return
		}
		p, err := geoRepo.Get(r.Context(), targetID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"enabled":       p.Enabled,
			"allowed_cidrs": p.AllowedCIDRs,
			"gps_enabled":   p.GPSEnabled,
			"gps_lat":       p.GPSLat,
			"gps_lon":       p.GPSLon,
			"gps_radius_m":  p.GPSRadiusM,
		})
	}
}

func handleAdminListWipeAudit(deps Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if deps.Pool == nil {
			writeError(w, http.StatusServiceUnavailable, "auditoria indisponível")
			return
		}
		events, err := security.ListWipeEvents(r.Context(), deps.Pool, 50)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao listar auditoria")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"events": events})
	}
}

func registerAdminRoutes(mux *http.ServeMux, deps Deps) {
	if deps.AdminKey == "" {
		return
	}
	guard := func(h http.HandlerFunc) http.HandlerFunc {
		return requireAdminKey(deps.AdminKey, h)
	}

	if deps.Users != nil {
		mux.HandleFunc("GET /api/admin/users", guard(handleAdminListUsers(deps.Users)))
		mux.HandleFunc("GET /api/admin/users/{id}", guard(handleAdminGetUser(deps.Users)))
	}
	if deps.Shifts != nil {
		mux.HandleFunc("GET /api/admin/users/{id}/access-shift", guard(handleAdminGetAccessShift(deps.Shifts)))
	}
	if deps.Geofence != nil {
		mux.HandleFunc("GET /api/admin/users/{id}/access-geofence", guard(handleAdminGetAccessGeofence(deps.Geofence)))
	}
	if deps.Shifts != nil && deps.Users != nil {
		mux.HandleFunc(
			"PUT /api/admin/users/{id}/access-shift",
			guard(handleAdminSetAccessShift(deps.Users, deps.Shifts)),
		)
	}
	if deps.Geofence != nil && deps.Users != nil {
		mux.HandleFunc(
			"PUT /api/admin/users/{id}/access-geofence",
			guard(handleAdminSetAccessGeofence(deps.Users, deps.Geofence)),
		)
	}
	if deps.Wipe != nil && deps.Users != nil {
		mux.HandleFunc(
			"POST /api/admin/users/{id}/remote-wipe",
			guard(handleAdminRemoteWipe(deps.Users, deps.Wipe)),
		)
	}
	if deps.Pool != nil {
		mux.HandleFunc("GET /api/admin/audit/wipe-events", guard(handleAdminListWipeAudit(deps)))
	}
}
