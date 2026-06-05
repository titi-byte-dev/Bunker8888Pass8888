package agent

import (
	"encoding/json"
	"fmt"
)

// DecodeInput faz unmarshal para T e devolve erro de validação uniforme.
// Didático: cada tool define um struct Go — validação explícita sem inventar JSON Schema à mão.
func DecodeInput[T any](raw json.RawMessage, dest *T) error {
	if len(raw) == 0 {
		return fmt.Errorf("%w: body vazio", ErrInvalidToolInput)
	}
	if err := json.Unmarshal(raw, dest); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidToolInput, err)
	}
	return nil
}
