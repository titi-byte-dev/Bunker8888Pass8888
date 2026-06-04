package geofence

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repo acede à tabela access_geofence.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositório de geofencing.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// Get devolve a política. Sem linha → enabled=false (sem restrição).
func (r *Repo) Get(ctx context.Context, userID string) (Policy, error) {
	var p Policy
	var raw []byte
	var lat, lon *float64
	err := r.pool.QueryRow(ctx, `
		SELECT user_id::text, enabled, allowed_cidrs, gps_enabled, gps_lat, gps_lon, gps_radius_m
		FROM access_geofence WHERE user_id = $1`, userID,
	).Scan(&p.UserID, &p.Enabled, &raw, &p.GPSEnabled, &lat, &lon, &p.GPSRadiusM)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Policy{UserID: userID, Enabled: false, GPSRadiusM: 500}, nil
		}
		return Policy{}, err
	}
	p.AllowedCIDRs, err = ParseCIDRsJSON(raw)
	if err != nil {
		return Policy{}, err
	}
	p.GPSLat = lat
	p.GPSLon = lon
	if p.GPSRadiusM <= 0 {
		p.GPSRadiusM = 500
	}
	return p, nil
}

// Upsert grava ou actualiza a política de geofencing.
func (r *Repo) Upsert(ctx context.Context, p Policy) error {
	raw, err := MarshalCIDRsJSON(p.AllowedCIDRs)
	if err != nil {
		return err
	}
	radius := p.GPSRadiusM
	if radius <= 0 {
		radius = 500
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO access_geofence (user_id, enabled, allowed_cidrs, gps_enabled, gps_lat, gps_lon, gps_radius_m)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id) DO UPDATE SET
			enabled = EXCLUDED.enabled,
			allowed_cidrs = EXCLUDED.allowed_cidrs,
			gps_enabled = EXCLUDED.gps_enabled,
			gps_lat = EXCLUDED.gps_lat,
			gps_lon = EXCLUDED.gps_lon,
			gps_radius_m = EXCLUDED.gps_radius_m,
			updated_at = now()`,
		p.UserID, p.Enabled, raw, p.GPSEnabled, p.GPSLat, p.GPSLon, radius,
	)
	return err
}
