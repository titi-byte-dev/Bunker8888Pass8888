package httpapi

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/auth"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/geofence"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/shifts"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sentinel"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
)

// accessPolicyDeps agrupa validações de contexto de acesso (turno + geofence).
type accessPolicyDeps struct {
	Shifts    *shifts.Repo
	Geofence  *geofence.Repo
}

type geofencePolicyRequest struct {
	Enabled      bool     `json:"enabled"`
	AllowedCIDRs []string `json:"allowed_cidrs"`
	GPSEnabled   bool     `json:"gps_enabled"`
	GPSLat       *float64 `json:"gps_lat"`
	GPSLon       *float64 `json:"gps_lon"`
	GPSRadiusM   float64  `json:"gps_radius_m"`
}

func handleGetAccessGeofence(geoRepo *geofence.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		p, err := geoRepo.Get(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		geo, geoErr := geofence.ParseClientGeo(r.Header.Get("X-Geo-Latitude"), r.Header.Get("X-Geo-Longitude"))
		if geoErr != nil {
			writeError(w, http.StatusBadRequest, "coordenadas inválidas")
			return
		}
		allowed, err := geofence.IsAllowed(p, clientIP(r), geo)
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
			"within_fence":  allowed,
			"client_ip":     clientIP(r),
		})
	}
}

func handleAdminSetAccessGeofence(userRepo *users.Repo, geoRepo *geofence.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetID := r.PathValue("id")
		if targetID == "" {
			writeError(w, http.StatusBadRequest, "user id em falta")
			return
		}
		if _, err := userRepo.ByID(r.Context(), targetID); err != nil {
			if err == users.ErrNotFound {
				writeError(w, http.StatusNotFound, "utilizador não encontrado")
				return
			}
			writeError(w, http.StatusInternalServerError, "erro interno")
			return
		}
		var req geofencePolicyRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		p := geofence.Policy{
			UserID:       targetID,
			Enabled:      req.Enabled,
			AllowedCIDRs: req.AllowedCIDRs,
			GPSEnabled:   req.GPSEnabled,
			GPSLat:       req.GPSLat,
			GPSLon:       req.GPSLon,
			GPSRadiusM:   req.GPSRadiusM,
		}
		if err := geoRepo.Upsert(r.Context(), p); err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao gravar geofence")
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleLoginWithAccessPolicy(svc *auth.Service, ap accessPolicyDeps, sent *sentinel.Service, userRepo *users.Repo) http.HandlerFunc {
	if ap.Shifts == nil && ap.Geofence == nil && sent == nil {
		return handleLogin(svc)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		authHash, err := decodeAuthHash(req.AuthHash)
		if err != nil {
			writeError(w, http.StatusBadRequest, "base64 inválido")
			return
		}

		userID, err := svc.ValidateCredentials(r.Context(), req.Email, authHash)
		if errors.Is(err, auth.ErrInvalidCredentials) {
			if sent != nil && userRepo != nil {
				lc := loginContextFromRequest(r, "", req.Email)
				if u, e := userRepo.ByEmail(r.Context(), req.Email); e == nil {
					lc.UserID = u.ID
					lc.Email = u.Email
				}
				_ = sent.RecordFailure(r.Context(), lc)
			}
			writeError(w, http.StatusUnauthorized, "credenciais inválidas")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha no login")
			return
		}

		if err := writeAccessPolicyError(w, assertUserAccessPolicy(r.Context(), r, userID, ap)); err != nil {
			return
		}

		finishLoginWithSentinel(w, r, userID, req.Email, svc, sent)
	}
}

func requireAuthWithAccessPolicy(svc *auth.Service, ap accessPolicyDeps, next http.HandlerFunc) http.Handler {
	if ap.Shifts == nil && ap.Geofence == nil {
		return requireAuth(svc, next)
	}
	return requireAuth(svc, func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := writeAccessPolicyError(w, assertUserAccessPolicy(r.Context(), r, userID, ap)); err != nil {
			return
		}
		next.ServeHTTP(w, r)
	})
}

func assertUserAccessPolicy(ctx context.Context, r *http.Request, userID string, ap accessPolicyDeps) error {
	if ap.Shifts != nil {
		if err := assertUserWithinShift(ctx, ap.Shifts, userID); err != nil {
			return err
		}
	}
	if ap.Geofence != nil {
		if err := assertUserWithinGeofence(ctx, r, ap.Geofence, userID); err != nil {
			return err
		}
	}
	return nil
}

func assertUserWithinGeofence(ctx context.Context, r *http.Request, geoRepo *geofence.Repo, userID string) error {
	p, err := geoRepo.Get(ctx, userID)
	if err != nil {
		return err
	}
	geo, err := geofence.ParseClientGeo(r.Header.Get("X-Geo-Latitude"), r.Header.Get("X-Geo-Longitude"))
	if err != nil {
		return err
	}
	return geofence.AssertAllowed(p, clientIP(r), geo)
}

func writeAccessPolicyError(w http.ResponseWriter, err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, shifts.ErrOutsideShift):
		writeError(w, http.StatusForbidden, "fora do horário de turno")
	case errors.Is(err, geofence.ErrOutsideFence):
		writeError(w, http.StatusForbidden, "fora da zona geográfica permitida")
	case errors.Is(err, geofence.ErrGPSRequired):
		writeError(w, http.StatusForbidden, "localização GPS necessária")
	default:
		writeError(w, http.StatusInternalServerError, "acesso negado")
	}
	return err
}

// clientIP extrai o IP do cliente (suporta X-Forwarded-For atrás de proxy).
//
// ⚠️ Segurança: em produção o proxy deve sobrescrever X-Forwarded-For; confiar
// apenas em proxies controlados (INFRA-001).
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return strings.TrimSpace(r.RemoteAddr)
	}
	return host
}
