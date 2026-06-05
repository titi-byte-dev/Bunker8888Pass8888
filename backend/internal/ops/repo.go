// Package ops gere inventário operacional (AGENT-008).
//
// Didático: nomes e SKUs são metadados de stock — não substituem cifragem ZK
// para dados pessoais. O orquestrador sugere ordens de compra quando o stock
// desce ao nível de reordenação.
package ops

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("ops: artigo não encontrado")

// Item espelha ops_inventory_items.
type Item struct {
	ID           string
	OwnerID      string
	Name         string
	SKU          string
	Quantity     int
	ReorderLevel int
	Unit         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Repo acede ao inventário por owner_id.
type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// IsLowStock indica se o stock atingiu o limiar de reordenação.
func IsLowStock(qty, reorder int) bool {
	return qty <= reorder
}

type CreateInput struct {
	Name         string
	SKU          string
	Quantity     int
	ReorderLevel int
	Unit         string
}

func (r *Repo) Create(ctx context.Context, ownerID string, in CreateInput) (*Item, error) {
	if in.Name == "" {
		return nil, errors.New("ops: nome obrigatório")
	}
	if in.Quantity < 0 || in.ReorderLevel < 0 {
		return nil, errors.New("ops: quantidades inválidas")
	}
	unit := in.Unit
	if unit == "" {
		unit = "un"
	}
	it := &Item{OwnerID: ownerID, Name: in.Name, SKU: in.SKU, Quantity: in.Quantity, ReorderLevel: in.ReorderLevel, Unit: unit}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO ops_inventory_items (owner_id, name, sku, quantity, reorder_level, unit)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id::text, created_at, updated_at`,
		ownerID, in.Name, in.SKU, in.Quantity, in.ReorderLevel, unit,
	).Scan(&it.ID, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return it, nil
}

func (r *Repo) List(ctx context.Context, ownerID string) ([]Item, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, owner_id::text, name, sku, quantity, reorder_level, unit, created_at, updated_at
		FROM ops_inventory_items WHERE owner_id = $1 ORDER BY name ASC`, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(&it.ID, &it.OwnerID, &it.Name, &it.SKU, &it.Quantity, &it.ReorderLevel, &it.Unit, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

func (r *Repo) Get(ctx context.Context, ownerID, id string) (*Item, error) {
	it := &Item{ID: id, OwnerID: ownerID}
	err := r.pool.QueryRow(ctx, `
		SELECT name, sku, quantity, reorder_level, unit, created_at, updated_at
		FROM ops_inventory_items WHERE id = $1 AND owner_id = $2`,
		id, ownerID,
	).Scan(&it.Name, &it.SKU, &it.Quantity, &it.ReorderLevel, &it.Unit, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return it, nil
}

type UpdateInput struct {
	Name         string
	SKU          string
	Quantity     int
	ReorderLevel int
	Unit         string
}

func (r *Repo) Update(ctx context.Context, ownerID, id string, in UpdateInput) (*Item, error) {
	if in.Name == "" {
		return nil, errors.New("ops: nome obrigatório")
	}
	if in.Quantity < 0 || in.ReorderLevel < 0 {
		return nil, errors.New("ops: quantidades inválidas")
	}
	unit := in.Unit
	if unit == "" {
		unit = "un"
	}
	it := &Item{ID: id, OwnerID: ownerID, Name: in.Name, SKU: in.SKU, Quantity: in.Quantity, ReorderLevel: in.ReorderLevel, Unit: unit}
	err := r.pool.QueryRow(ctx, `
		UPDATE ops_inventory_items
		SET name = $1, sku = $2, quantity = $3, reorder_level = $4, unit = $5, updated_at = now()
		WHERE id = $6 AND owner_id = $7
		RETURNING created_at, updated_at`,
		in.Name, in.SKU, in.Quantity, in.ReorderLevel, unit, id, ownerID,
	).Scan(&it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return it, nil
}

// Adjust altera quantity por delta (pode ser negativo).
func (r *Repo) Adjust(ctx context.Context, ownerID, id string, delta int) (*Item, error) {
	it := &Item{ID: id, OwnerID: ownerID}
	err := r.pool.QueryRow(ctx, `
		UPDATE ops_inventory_items
		SET quantity = GREATEST(0, quantity + $1), updated_at = now()
		WHERE id = $2 AND owner_id = $3
		RETURNING name, sku, quantity, reorder_level, unit, created_at, updated_at`,
		delta, id, ownerID,
	).Scan(&it.Name, &it.SKU, &it.Quantity, &it.ReorderLevel, &it.Unit, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return it, nil
}

func (r *Repo) Delete(ctx context.Context, ownerID, id string) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM ops_inventory_items WHERE id = $1 AND owner_id = $2`, id, ownerID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
