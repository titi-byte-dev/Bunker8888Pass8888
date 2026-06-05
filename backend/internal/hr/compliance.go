package hr

import (
	"context"
	"time"
)

// Relatorio de conformidade RGPD (HR-008).
//
// Reune metadados NAO-secretos sobre as fichas e a auditoria do utilizador: o
// servidor nunca decifra nada, so conta e atesta a integridade da cadeia de
// logs (HR-002). O cliente formata isto num relatorio imprimivel/PDF.

// ComplianceReport e a fotografia de conformidade num dado instante.
type ComplianceReport struct {
	GeneratedAt        time.Time
	RecordCount        int64
	ActiveFieldCount   int64 // campos ainda cifrados (com chave)
	ShreddedFieldCount int64 // campos crypto-shredded (sem chave)
	CertificateCount   int64 // certificados de eliminacao emitidos
	AuditEntryCount    int64 // entradas na cadeia de auditoria
	AuditChainValid    bool  // cadeia intacta (recalculada server-side)
	AuditBrokenSeq     int64 // primeira entrada partida (0 se intacta)
}

// ComplianceReportFor reune o relatorio para o utilizador.
func (r *Repo) ComplianceReportFor(ctx context.Context, ownerID string) (*ComplianceReport, error) {
	rep := &ComplianceReport{GeneratedAt: time.Now().UTC()}

	if err := r.pool.QueryRow(ctx,
		`SELECT count(*) FROM employee_records WHERE owner_id = $1`, ownerID,
	).Scan(&rep.RecordCount); err != nil {
		return nil, err
	}

	// Campos activos vs eliminados, restritos as fichas do dono.
	if err := r.pool.QueryRow(ctx, `
		SELECT
			count(*) FILTER (WHERE f.wrapped_key IS NOT NULL),
			count(*) FILTER (WHERE f.wrapped_key IS NULL)
		FROM employee_fields f
		JOIN employee_records e ON e.id = f.record_id
		WHERE e.owner_id = $1`, ownerID,
	).Scan(&rep.ActiveFieldCount, &rep.ShreddedFieldCount); err != nil {
		return nil, err
	}

	if err := r.pool.QueryRow(ctx,
		`SELECT count(*) FROM erasure_certificates WHERE owner_id = $1`, ownerID,
	).Scan(&rep.CertificateCount); err != nil {
		return nil, err
	}

	if err := r.pool.QueryRow(ctx,
		`SELECT count(*) FROM audit_log WHERE owner_id = $1`, ownerID,
	).Scan(&rep.AuditEntryCount); err != nil {
		return nil, err
	}

	valid, brokenSeq, err := r.VerifyAudit(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	rep.AuditChainValid = valid
	rep.AuditBrokenSeq = brokenSeq

	return rep, nil
}
