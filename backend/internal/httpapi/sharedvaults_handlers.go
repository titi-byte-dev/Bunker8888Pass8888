package httpapi

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sharedvaults"
)

// --- Tipos de pedido/resposta (tudo opaco: blobs em base64) -----------------

type createVaultRequest struct {
	NameBlob        string `json:"name_blob"`         // nome cifrado com a chave do cofre
	Algorithm       string `json:"algorithm"`         // marcador de esquema
	WrappedVaultKey string `json:"wrapped_vault_key"` // chave do cofre cifrada p/ a minha PK
}

type addMemberRequest struct {
	UserID          string `json:"user_id"`
	Role            string `json:"role"`
	WrappedVaultKey string `json:"wrapped_vault_key"` // chave do cofre cifrada p/ a PK do convidado
}

type vaultItemRequest struct {
	Type string `json:"type"`
	Blob string `json:"blob"`
}

// mapSharedVaultError traduz erros de domínio em códigos HTTP.
func mapSharedVaultError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, sharedvaults.ErrNotFound):
		writeError(w, http.StatusNotFound, "cofre não encontrado")
	case errors.Is(err, sharedvaults.ErrForbidden):
		writeError(w, http.StatusForbidden, "sem permissão para esta operação")
	case errors.Is(err, sharedvaults.ErrInvalidRole):
		writeError(w, http.StatusBadRequest, "papel inválido (admin, member ou viewer)")
	case errors.Is(err, sharedvaults.ErrOwnerImmutable):
		writeError(w, http.StatusForbidden, "o dono do cofre não pode ser alterado por aqui")
	default:
		writeError(w, http.StatusInternalServerError, "falha na operação do cofre partilhado")
	}
	return true
}

func vaultJSON(vm *sharedvaults.VaultForMember) map[string]any {
	return map[string]any{
		"id":                vm.ID,
		"owner_id":          vm.OwnerID,
		"name_blob":         base64.StdEncoding.EncodeToString(vm.NameBlob),
		"algorithm":         vm.Algorithm,
		"role":              vm.Role,
		"wrapped_vault_key": base64.StdEncoding.EncodeToString(vm.WrappedVaultKey),
		"created_at":        vm.CreatedAt,
		"updated_at":        vm.UpdatedAt,
	}
}

// --- Handlers de cofres ------------------------------------------------------

func handleCreateSharedVault(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req createVaultRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		nameBlob, err := base64.StdEncoding.DecodeString(req.NameBlob)
		if err != nil || len(nameBlob) == 0 {
			writeError(w, http.StatusBadRequest, "name_blob base64 inválido")
			return
		}
		wrapped, err := base64.StdEncoding.DecodeString(req.WrappedVaultKey)
		if err != nil || len(wrapped) == 0 {
			writeError(w, http.StatusBadRequest, "wrapped_vault_key base64 inválido")
			return
		}
		v, err := repo.CreateVault(r.Context(), userID, nameBlob, strings.TrimSpace(req.Algorithm), wrapped)
		if mapSharedVaultError(w, err) {
			return
		}
		vm := &sharedvaults.VaultForMember{Vault: *v, Role: sharedvaults.RoleOwner, WrappedVaultKey: wrapped}
		writeJSON(w, http.StatusCreated, vaultJSON(vm))
	}
}

func handleListSharedVaults(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		vaults, err := repo.ListForUser(r.Context(), userID)
		if mapSharedVaultError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(vaults))
		for i := range vaults {
			out = append(out, vaultJSON(&vaults[i]))
		}
		writeJSON(w, http.StatusOK, map[string]any{"vaults": out})
	}
}

func handleGetSharedVault(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		vm, err := repo.GetForUser(r.Context(), userID, r.PathValue("id"))
		if mapSharedVaultError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, vaultJSON(vm))
	}
}

func handleDeleteSharedVault(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := repo.DeleteVault(r.Context(), r.PathValue("id"), userID); mapSharedVaultError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// --- Handlers de membros -----------------------------------------------------

func handleListSharedVaultMembers(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		members, err := repo.ListMembers(r.Context(), r.PathValue("id"), userID)
		if mapSharedVaultError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(members))
		for _, m := range members {
			out = append(out, map[string]any{
				"user_id":    m.UserID,
				"email":      m.Email,
				"role":       m.Role,
				"created_at": m.CreatedAt,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"members": out})
	}
}

func handleAddSharedVaultMember(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req addMemberRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		newUserID := strings.TrimSpace(req.UserID)
		if newUserID == "" {
			writeError(w, http.StatusBadRequest, "user_id em falta")
			return
		}
		wrapped, err := base64.StdEncoding.DecodeString(req.WrappedVaultKey)
		if err != nil || len(wrapped) == 0 {
			writeError(w, http.StatusBadRequest, "wrapped_vault_key base64 inválido")
			return
		}
		err = repo.AddMember(r.Context(), r.PathValue("id"), userID, newUserID, strings.TrimSpace(req.Role), wrapped)
		if mapSharedVaultError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleRemoveSharedVaultMember(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		target := r.PathValue("userId")
		if err := repo.RemoveMember(r.Context(), r.PathValue("id"), userID, target); mapSharedVaultError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// --- Handlers de itens -------------------------------------------------------

func handleListSharedVaultItems(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		items, err := repo.ListItems(r.Context(), r.PathValue("id"), userID)
		if mapSharedVaultError(w, err) {
			return
		}
		out := make([]map[string]any, 0, len(items))
		for _, it := range items {
			out = append(out, map[string]any{
				"id":         it.ID,
				"type":       it.ItemType,
				"blob":       base64.StdEncoding.EncodeToString(it.Blob),
				"created_by": it.CreatedBy,
				"created_at": it.CreatedAt,
				"updated_at": it.UpdatedAt,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": out})
	}
}

func handleCreateSharedVaultItem(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		var req vaultItemRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if strings.TrimSpace(req.Type) == "" {
			writeError(w, http.StatusBadRequest, "type em falta")
			return
		}
		blob, err := base64.StdEncoding.DecodeString(req.Blob)
		if err != nil || len(blob) == 0 {
			writeError(w, http.StatusBadRequest, "blob base64 inválido")
			return
		}
		it, err := repo.CreateItem(r.Context(), r.PathValue("id"), userID, strings.TrimSpace(req.Type), blob)
		if mapSharedVaultError(w, err) {
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"id":         it.ID,
			"type":       it.ItemType,
			"created_by": it.CreatedBy,
			"created_at": it.CreatedAt,
			"updated_at": it.UpdatedAt,
		})
	}
}

func handleDeleteSharedVaultItem(repo *sharedvaults.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, _ := r.Context().Value(userIDKey).(string)
		if err := repo.DeleteItem(r.Context(), r.PathValue("id"), userID, r.PathValue("itemId")); mapSharedVaultError(w, err) {
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
