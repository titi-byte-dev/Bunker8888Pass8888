// Package sharedvaults implementa os Cofres Partilhados (SHARE-002): coleções de
// itens cifrados sob uma chave de cofre própria, partilhada com vários membros.
//
// ⚠️ Modelo Zero-Knowledge: o servidor trata todos os blobs como bytes opacos.
//   - name_blob / item.blob  → cifrados com a chave do cofre (AES-GCM no cliente);
//   - wrapped_vault_key       → a chave do cofre cifrada para a chave pública de
//     cada membro (RSA-OAEP). O servidor nunca decifra nada disto.
//
// O servidor só faz cumprir o controlo de acesso: quem é membro, com que papel,
// e quem pode gerir membros ou escrever itens.
package sharedvaults

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Papéis possíveis num cofre partilhado, do mais para o menos privilegiado.
const (
	RoleOwner  = "owner"  // controlo total: gere membros, apaga o cofre
	RoleAdmin  = "admin"  // gere membros (exceto owner), lê/escreve itens
	RoleMember = "member" // lê/escreve itens
	RoleViewer = "viewer" // só lê itens
)

// Erros de domínio (mapeados para códigos HTTP nos handlers).
var (
	// ErrNotFound: cofre inexistente OU o utilizador não é membro (mesma
	// resposta, para não revelar a existência de cofres alheios).
	ErrNotFound = errors.New("sharedvaults: cofre não encontrado")
	// ErrForbidden: é membro, mas o papel não permite a operação.
	ErrForbidden = errors.New("sharedvaults: sem permissão")
	// ErrInvalidRole: papel desconhecido ou não atribuível (ex.: 'owner').
	ErrInvalidRole = errors.New("sharedvaults: papel inválido")
	// ErrOwnerImmutable: não se remove nem rebaixa o dono do cofre.
	ErrOwnerImmutable = errors.New("sharedvaults: o dono não pode ser removido")
)

// Vault espelha uma linha de "shared_vaults".
type Vault struct {
	ID        string
	OwnerID   string
	NameBlob  []byte
	Algorithm string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// VaultForMember é um cofre na perspetiva de um membro: inclui o seu papel e a
// chave do cofre cifrada para a sua chave pública.
type VaultForMember struct {
	Vault
	Role            string
	WrappedVaultKey []byte
}

// Membership espelha uma linha de "shared_vault_members" (com o email do membro).
type Membership struct {
	VaultID         string
	UserID          string
	Email           string
	Role            string
	WrappedVaultKey []byte
	CreatedAt       time.Time
}

// Item espelha uma linha de "shared_vault_items".
type Item struct {
	ID        string
	VaultID   string
	ItemType  string
	Blob      []byte
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Repo acede às tabelas shared_vault*.
type Repo struct {
	pool *pgxpool.Pool
}

// NewRepo cria o repositório.
func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// isAssignableRole valida papéis que um gestor pode atribuir (não 'owner').
func isAssignableRole(role string) bool {
	switch role {
	case RoleAdmin, RoleMember, RoleViewer:
		return true
	default:
		return false
	}
}

// canManageMembers indica se o papel pode convidar/remover membros.
func canManageMembers(role string) bool {
	return role == RoleOwner || role == RoleAdmin
}

// canWriteItems indica se o papel pode criar/apagar itens (tudo menos viewer).
func canWriteItems(role string) bool {
	return role == RoleOwner || role == RoleAdmin || role == RoleMember
}

// roleOf devolve o papel do utilizador no cofre, ou ErrNotFound se não for membro.
func (r *Repo) roleOf(ctx context.Context, vaultID, userID string) (string, error) {
	var role string
	err := r.pool.QueryRow(ctx,
		`SELECT role FROM shared_vault_members WHERE vault_id = $1 AND user_id = $2`,
		vaultID, userID,
	).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	return role, nil
}

// CreateVault cria um cofre e inscreve o criador como owner, numa só transação.
// ownerWrappedKey é a chave do cofre cifrada para a chave pública do próprio.
func (r *Repo) CreateVault(ctx context.Context, ownerID string, nameBlob []byte, algorithm string, ownerWrappedKey []byte) (*Vault, error) {
	if len(nameBlob) == 0 || len(ownerWrappedKey) == 0 {
		return nil, errors.New("sharedvaults: nome ou chave do cofre em falta")
	}
	if algorithm == "" {
		algorithm = "AES-GCM-256+RSA-OAEP-3072"
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	v := &Vault{OwnerID: ownerID, NameBlob: nameBlob, Algorithm: algorithm}
	if err := tx.QueryRow(ctx, `
		INSERT INTO shared_vaults (owner_id, name_blob, algorithm)
		VALUES ($1, $2, $3)
		RETURNING id::text, created_at, updated_at`,
		ownerID, nameBlob, algorithm,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt); err != nil {
		return nil, err
	}

	if _, err := tx.Exec(ctx, `
		INSERT INTO shared_vault_members (vault_id, user_id, role, wrapped_vault_key)
		VALUES ($1, $2, $3, $4)`,
		v.ID, ownerID, RoleOwner, ownerWrappedKey,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return v, nil
}

// ListForUser devolve os cofres de que o utilizador é membro (com papel e chave).
func (r *Repo) ListForUser(ctx context.Context, userID string) ([]VaultForMember, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT v.id::text, v.owner_id::text, v.name_blob, v.algorithm,
		       v.created_at, v.updated_at, m.role, m.wrapped_vault_key
		FROM shared_vault_members m
		JOIN shared_vaults v ON v.id = m.vault_id
		WHERE m.user_id = $1
		ORDER BY v.created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []VaultForMember
	for rows.Next() {
		var vm VaultForMember
		if err := rows.Scan(
			&vm.ID, &vm.OwnerID, &vm.NameBlob, &vm.Algorithm,
			&vm.CreatedAt, &vm.UpdatedAt, &vm.Role, &vm.WrappedVaultKey,
		); err != nil {
			return nil, err
		}
		out = append(out, vm)
	}
	return out, rows.Err()
}

// GetForUser devolve um cofre se o utilizador for membro, senão ErrNotFound.
func (r *Repo) GetForUser(ctx context.Context, userID, vaultID string) (*VaultForMember, error) {
	vm := &VaultForMember{}
	err := r.pool.QueryRow(ctx, `
		SELECT v.id::text, v.owner_id::text, v.name_blob, v.algorithm,
		       v.created_at, v.updated_at, m.role, m.wrapped_vault_key
		FROM shared_vault_members m
		JOIN shared_vaults v ON v.id = m.vault_id
		WHERE m.user_id = $1 AND m.vault_id = $2`, userID, vaultID,
	).Scan(
		&vm.ID, &vm.OwnerID, &vm.NameBlob, &vm.Algorithm,
		&vm.CreatedAt, &vm.UpdatedAt, &vm.Role, &vm.WrappedVaultKey,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return vm, nil
}

// ListMembers devolve os membros do cofre. Requer que o requerente seja membro.
func (r *Repo) ListMembers(ctx context.Context, vaultID, userID string) ([]Membership, error) {
	if _, err := r.roleOf(ctx, vaultID, userID); err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `
		SELECT m.vault_id::text, m.user_id::text, u.email, m.role, m.created_at
		FROM shared_vault_members m
		JOIN users u ON u.id = m.user_id
		WHERE m.vault_id = $1
		ORDER BY (m.role = 'owner') DESC, m.created_at ASC`, vaultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Membership
	for rows.Next() {
		var m Membership
		if err := rows.Scan(&m.VaultID, &m.UserID, &m.Email, &m.Role, &m.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// AddMember convida (ou actualiza) um membro. O requerente tem de poder gerir
// membros; o papel atribuído não pode ser 'owner'. wrappedKey é a chave do cofre
// cifrada para a chave pública do novo membro.
func (r *Repo) AddMember(ctx context.Context, vaultID, actorID, newUserID, role string, wrappedKey []byte) error {
	if !isAssignableRole(role) {
		return ErrInvalidRole
	}
	if len(wrappedKey) == 0 {
		return errors.New("sharedvaults: chave do cofre cifrada em falta")
	}
	actorRole, err := r.roleOf(ctx, vaultID, actorID)
	if err != nil {
		return err
	}
	if !canManageMembers(actorRole) {
		return ErrForbidden
	}
	// Nunca alterar o dono por esta via (mesmo que o alvo já seja owner).
	if existing, err := r.roleOf(ctx, vaultID, newUserID); err == nil && existing == RoleOwner {
		return ErrOwnerImmutable
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO shared_vault_members (vault_id, user_id, role, wrapped_vault_key)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (vault_id, user_id) DO UPDATE SET
			role              = EXCLUDED.role,
			wrapped_vault_key = EXCLUDED.wrapped_vault_key`,
		vaultID, newUserID, role, wrappedKey,
	)
	return err
}

// RemoveMember retira um membro (revogação imediata). O requerente tem de poder
// gerir membros; o dono nunca é removido por esta via.
func (r *Repo) RemoveMember(ctx context.Context, vaultID, actorID, targetID string) error {
	actorRole, err := r.roleOf(ctx, vaultID, actorID)
	if err != nil {
		return err
	}
	if !canManageMembers(actorRole) {
		return ErrForbidden
	}
	targetRole, err := r.roleOf(ctx, vaultID, targetID)
	if err != nil {
		return err
	}
	if targetRole == RoleOwner {
		return ErrOwnerImmutable
	}
	_, err = r.pool.Exec(ctx,
		`DELETE FROM shared_vault_members WHERE vault_id = $1 AND user_id = $2`,
		vaultID, targetID,
	)
	return err
}

// DeleteVault apaga o cofre (e em cascata membros e itens). Só o dono o pode fazer.
func (r *Repo) DeleteVault(ctx context.Context, vaultID, actorID string) error {
	role, err := r.roleOf(ctx, vaultID, actorID)
	if err != nil {
		return err
	}
	if role != RoleOwner {
		return ErrForbidden
	}
	_, err = r.pool.Exec(ctx, `DELETE FROM shared_vaults WHERE id = $1`, vaultID)
	return err
}

// CreateItem adiciona um item cifrado ao cofre. Exige papel com escrita.
func (r *Repo) CreateItem(ctx context.Context, vaultID, userID, itemType string, blob []byte) (*Item, error) {
	if itemType == "" || len(blob) == 0 {
		return nil, errors.New("sharedvaults: tipo ou blob em falta")
	}
	role, err := r.roleOf(ctx, vaultID, userID)
	if err != nil {
		return nil, err
	}
	if !canWriteItems(role) {
		return nil, ErrForbidden
	}
	it := &Item{VaultID: vaultID, ItemType: itemType, Blob: blob, CreatedBy: userID}
	err = r.pool.QueryRow(ctx, `
		INSERT INTO shared_vault_items (vault_id, item_type, blob, created_by)
		VALUES ($1, $2, $3, $4)
		RETURNING id::text, created_at, updated_at`,
		vaultID, itemType, blob, userID,
	).Scan(&it.ID, &it.CreatedAt, &it.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return it, nil
}

// ListItems devolve os itens do cofre. Qualquer membro (incl. viewer) pode ler.
func (r *Repo) ListItems(ctx context.Context, vaultID, userID string) ([]Item, error) {
	if _, err := r.roleOf(ctx, vaultID, userID); err != nil {
		return nil, err
	}
	rows, err := r.pool.Query(ctx, `
		SELECT id::text, vault_id::text, item_type, blob,
		       COALESCE(created_by::text, ''), created_at, updated_at
		FROM shared_vault_items
		WHERE vault_id = $1
		ORDER BY created_at DESC`, vaultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Item
	for rows.Next() {
		var it Item
		if err := rows.Scan(
			&it.ID, &it.VaultID, &it.ItemType, &it.Blob,
			&it.CreatedBy, &it.CreatedAt, &it.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, it)
	}
	return out, rows.Err()
}

// DeleteItem remove um item. Exige papel com escrita (tudo menos viewer).
func (r *Repo) DeleteItem(ctx context.Context, vaultID, userID, itemID string) error {
	role, err := r.roleOf(ctx, vaultID, userID)
	if err != nil {
		return err
	}
	if !canWriteItems(role) {
		return ErrForbidden
	}
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM shared_vault_items WHERE id = $1 AND vault_id = $2`, itemID, vaultID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
