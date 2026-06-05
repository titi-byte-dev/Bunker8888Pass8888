package httpapi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/clidevices"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/vault"
)

type cliDeps struct {
	ca      *clidevices.CA
	devices *clidevices.Repo
	certTTL time.Duration
}

type registerDeviceRequest struct {
	Name   string `json:"name"`
	CSRPEM string `json:"csr_pem"`
}

// handleRegisterCLIDevice regista um dispositivo CLI (requer sessão Bearer).
// O cliente envia um CSR; o servidor assina e devolve o certificado + CA.
func handleRegisterCLIDevice(d cliDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req registerDeviceRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.Name == "" || req.CSRPEM == "" {
			writeError(w, http.StatusBadRequest, "name e csr_pem são obrigatórios")
			return
		}
		if d.ca == nil {
			writeError(w, http.StatusServiceUnavailable, "mTLS não configurado no servidor")
			return
		}

		certPEM, cert, err := d.ca.SignCSR([]byte(req.CSRPEM), d.certTTL)
		if errors.Is(err, clidevices.ErrInvalidCSR) || errors.Is(err, clidevices.ErrUnsupportedKey) {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao assinar certificado")
			return
		}

		fp := clidevices.Fingerprint(cert)
		if err := d.devices.Register(r.Context(), userID, req.Name, fp); err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao registar dispositivo")
			return
		}

		caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: d.ca.CertDER()})
		writeJSON(w, http.StatusCreated, map[string]string{
			"cert_pem": string(certPEM),
			"ca_pem":   string(caPEM),
		})
	}
}

// handleCLIListDevices lista dispositivos CLI do utilizador (sessão Bearer).
func handleCLIListDevices(devices *clidevices.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		list, err := devices.ListByUser(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao listar dispositivos")
			return
		}
		out := make([]map[string]string, 0, len(list))
		for _, d := range list {
			out = append(out, map[string]string{
				"id":         d.ID,
				"name":       d.Name,
				"created_at": d.CreatedAt,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"devices": out})
	}
}

// handleCLIListVaultItems lista metadados do cofre via mTLS (sem Bearer).
func handleCLIListVaultItems(repo *vault.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		items, err := repo.ListByUser(r.Context(), userID, r.URL.Query().Get("type"))
		if mapVaultError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(items))
		for _, it := range items {
			out = append(out, itemJSON(&it, false))
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": out})
	}
}

// handleCLIGetVaultItem devolve um item com blob cifrado via mTLS.
func handleCLIGetVaultItem(repo *vault.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		item, err := repo.GetByID(r.Context(), userID, r.PathValue("id"))
		if mapVaultError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, itemJSON(item, true))
	}
}

// requireMTLSDevice valida o certificado cliente TLS e resolve o user_id.
func requireMTLSDevice(devices *clidevices.Repo, next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil || len(r.TLS.PeerCertificates) == 0 {
			writeError(w, http.StatusUnauthorized, "certificado cliente em falta")
			return
		}
		cert := r.TLS.PeerCertificates[0]
		fp := clidevices.Fingerprint(cert)
		userID, err := devices.LookupActiveByFingerprint(r.Context(), fp)
		if errors.Is(err, clidevices.ErrNotFound) {
			writeError(w, http.StatusUnauthorized, "dispositivo não registado ou revogado")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha na autenticação mTLS")
			return
		}
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// NewMTLSRouter devolve rotas exclusivas da CLI com autenticação por certificado.
func NewMTLSRouter(deps Deps) http.Handler {
	mux := http.NewServeMux()
	if deps.Devices == nil || deps.Vault == nil {
		return mux
	}
	mux.Handle("GET /api/cli/vault/items", requireMTLSDevice(deps.Devices, handleCLIListVaultItems(deps.Vault)))
	mux.Handle("GET /api/cli/vault/items/{id}", requireMTLSDevice(deps.Devices, handleCLIGetVaultItem(deps.Vault)))
	return mux
}

// MTLSTLSConfig constrói a configuração TLS do servidor mTLS.
func MTLSTLSConfig(serverCert, serverKey, caCert []byte) (*tls.Config, error) {
	cert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("CA PEM inválido")
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    pool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// registerCLIRoutes adiciona rotas de gestão de dispositivos ao router HTTP normal.
func registerCLIRoutes(mux *http.ServeMux, deps Deps) {
	if deps.Auth == nil || deps.Devices == nil {
		return
	}
	cd := cliDeps{ca: deps.CLIca, devices: deps.Devices, certTTL: deps.CLICertTTL}
	if deps.CLICertTTL == 0 {
		cd.certTTL = 365 * 24 * time.Hour
	}
	mux.Handle("POST /api/cli/devices", requireAuth(deps.Auth, handleRegisterCLIDevice(cd)))
	mux.Handle("GET /api/cli/devices", requireAuth(deps.Auth, handleCLIListDevices(deps.Devices)))
	mux.Handle("DELETE /api/cli/devices/{id}", requireAuth(deps.Auth, handleRevokeCLIDevice(deps.Devices)))
}
