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
	case errors.Is(err, hr.ErrAlreadyShredded):
		writeError(w, http.StatusConflict, "campo já eliminado")
	case errors.Is(err, hr.ErrTooLargeContract):
		writeError(w, http.StatusRequestEntityTooLarge, "contrato demasiado grande")
	case errors.Is(err, hr.ErrNoSigningIdentity):
		writeError(w, http.StatusNotFound, "sem identidade de assinatura")
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

// --- Crypto-shredding (HR-003) + certificados de eliminação (HR-004) ---------

func certificateJSON(c *hr.Certificate) map[string]any {
	return map[string]any{
		"id":           c.ID,
		"owner_id":     c.OwnerID,
		"record_id":    c.RecordID,
		"field_name":   c.FieldName,
		"value_digest": c.ValueDigest,
		"shredded_at":  c.ShreddedAt,
		"cert_hash":    c.CertHash,
		"issued_at":    c.IssuedAt,
	}
}

// handleShredEmployeeField destrói a chave de um campo (crypto-shred) e devolve
// o certificado de eliminação emitido.
func handleShredEmployeeField(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		cert, err := repo.ShredField(r.Context(), userID, r.PathValue("id"), r.PathValue("field"))
		if mapHRError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, certificateJSON(cert))
	}
}

// handleShredEmployeeRecord elimina todos os campos ainda com chave de uma ficha,
// devolvendo a lista de certificados emitidos.
func handleShredEmployeeRecord(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		certs, err := repo.ShredRecord(r.Context(), userID, r.PathValue("id"))
		if mapHRError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(certs))
		for i := range certs {
			out = append(out, certificateJSON(&certs[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"certificates": out})
	}
}

// handleListErasureCertificates lista os certificados de eliminação do utilizador.
func handleListErasureCertificates(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		certs, err := repo.ListCertificates(r.Context(), userID)
		if mapHRError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(certs))
		for i := range certs {
			out = append(out, certificateJSON(&certs[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"certificates": out})
	}
}

// --- Logs imutáveis (HR-002) + relatório de conformidade (HR-008) ------------

func auditEntryJSON(e *hr.AuditEntry) map[string]any {
	return map[string]any{
		"owner_id":    e.OwnerID,
		"seq":         e.Seq,
		"action":      e.Action,
		"detail":      e.Detail,
		"occurred_at": e.OccurredAt,
		"prev_hash":   e.PrevHash,
		"entry_hash":  e.EntryHash,
	}
}

// handleListAuditLog devolve a cadeia de auditoria do utilizador (verificável
// no cliente recalculando os hashes encadeados).
func handleListAuditLog(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		entries, err := repo.ListAudit(r.Context(), userID)
		if mapHRError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(entries))
		for i := range entries {
			out = append(out, auditEntryJSON(&entries[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"entries": out})
	}
}

// --- Contratos (HR-005) + assinatura digital (HR-006) ------------------------

type signingIdentityRequest struct {
	PublicKey         string `json:"public_key"`          // SPKI base64
	WrappedPrivateKey string `json:"wrapped_private_key"` // PKCS8 cifrado base64
}

type contractRequest struct {
	MetaBlob   string `json:"meta_blob"`
	DataBlob   string `json:"data_blob"`
	WrappedKey string `json:"wrapped_key"`
}

type signContractRequest struct {
	ContentDigest string `json:"content_digest"`
	Signature     string `json:"signature"` // ECDSA base64
}

func contractMetaJSON(c *hr.Contract) map[string]any {
	m := map[string]any{
		"id":          c.ID,
		"record_id":   c.RecordID,
		"meta_blob":   base64.StdEncoding.EncodeToString(c.MetaBlob),
		"wrapped_key": base64.StdEncoding.EncodeToString(c.WrappedKey),
		"byte_size":   c.ByteSize,
		"signed":      len(c.Signature) > 0,
		"created_by":  c.CreatedBy,
		"created_at":  c.CreatedAt,
	}
	if len(c.Signature) > 0 {
		m["content_digest"] = c.ContentDigest
		m["signature"] = base64.StdEncoding.EncodeToString(c.Signature)
		m["signed_by"] = c.SignedBy
		if c.SignedAt != nil {
			m["signed_at"] = *c.SignedAt
		}
	}
	return m
}

func handlePutSigningIdentity(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req signingIdentityRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		pub, err1 := base64.StdEncoding.DecodeString(req.PublicKey)
		priv, err2 := base64.StdEncoding.DecodeString(req.WrappedPrivateKey)
		if err1 != nil || err2 != nil || len(pub) == 0 || len(priv) == 0 {
			writeError(w, http.StatusBadRequest, "chaves base64 inválidas")
			return
		}
		if err := repo.PutSigningIdentity(r.Context(), userID, pub, priv); mapHRError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func handleGetSigningIdentity(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		si, err := repo.GetSigningIdentity(r.Context(), userID)
		if mapHRError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"public_key":          base64.StdEncoding.EncodeToString(si.PublicKey),
			"wrapped_private_key": base64.StdEncoding.EncodeToString(si.WrappedPrivateKey),
			"created_at":          si.CreatedAt,
		})
	}
}

// handleGetSignerPublicKey devolve só a chave pública de um signatário (verificar).
func handleGetSignerPublicKey(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pk, err := repo.GetSignerPublicKey(r.Context(), r.PathValue("userId"))
		if mapHRError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"public_key": base64.StdEncoding.EncodeToString(pk),
		})
	}
}

func handleListContracts(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		cs, err := repo.ListContracts(r.Context(), userID, r.PathValue("id"))
		if mapHRError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(cs))
		for i := range cs {
			out = append(out, contractMetaJSON(&cs[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"contracts": out})
	}
}

func handleAddContract(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req contractRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		meta, e1 := base64.StdEncoding.DecodeString(req.MetaBlob)
		data, e2 := base64.StdEncoding.DecodeString(req.DataBlob)
		wrapped, e3 := base64.StdEncoding.DecodeString(req.WrappedKey)
		if e1 != nil || e2 != nil || e3 != nil || len(meta) == 0 || len(data) == 0 || len(wrapped) == 0 {
			writeError(w, http.StatusBadRequest, "blobs base64 inválidos")
			return
		}
		c, err := repo.AddContract(r.Context(), userID, r.PathValue("id"), meta, data, wrapped)
		if mapHRError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, contractMetaJSON(c))
	}
}

func handleGetContract(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		c, err := repo.GetContract(r.Context(), userID, r.PathValue("id"), r.PathValue("cid"))
		if mapHRError(w, err) {
			return
		}
		out := contractMetaJSON(c)
		out["data_blob"] = base64.StdEncoding.EncodeToString(c.DataBlob)
		writeJSON(w, http.StatusOK, out)
	}
}

func handleSignContract(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req signContractRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		sig, err := base64.StdEncoding.DecodeString(req.Signature)
		if err != nil || len(sig) == 0 || req.ContentDigest == "" {
			writeError(w, http.StatusBadRequest, "assinatura inválida")
			return
		}
		c, err := repo.SignContract(r.Context(), userID, r.PathValue("id"), r.PathValue("cid"), req.ContentDigest, sig)
		if mapHRError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, contractMetaJSON(c))
	}
}

func handleDeleteContract(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := repo.DeleteContract(r.Context(), userID, r.PathValue("id"), r.PathValue("cid")); mapHRError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// handleComplianceReport devolve o relatório de conformidade RGPD (metadados +
// estado da cadeia de auditoria). O cliente formata-o para PDF.
func handleComplianceReport(repo *hr.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		rep, err := repo.ComplianceReportFor(r.Context(), userID)
		if mapHRError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"generated_at":         rep.GeneratedAt,
			"record_count":         rep.RecordCount,
			"active_field_count":   rep.ActiveFieldCount,
			"shredded_field_count": rep.ShreddedFieldCount,
			"certificate_count":    rep.CertificateCount,
			"audit_entry_count":    rep.AuditEntryCount,
			"audit_chain_valid":    rep.AuditChainValid,
			"audit_broken_seq":     rep.AuditBrokenSeq,
		})
	}
}
