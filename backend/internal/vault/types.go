// Tipos de item suportados pelo cofre (VAULT-005).
package vault

import (
	"errors"
	"fmt"
)

const (
	TypeLogin = "login"
	TypeNote  = "note"
	TypeCard  = "card"
)

var ErrInvalidType = errors.New("vault: tipo de item inválido")

var ValidTypes = map[string]struct{}{
	TypeLogin: {},
	TypeNote:  {},
	TypeCard:  {},
}

func ValidateType(itemType string) error {
	if _, ok := ValidTypes[itemType]; !ok {
		return fmt.Errorf("%w: %q (permitidos: login, note, card)", ErrInvalidType, itemType)
	}
	return nil
}
