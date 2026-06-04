package vault

import (
	"errors"
	"testing"
)

func TestValidateType(t *testing.T) {
	for _, tc := range []struct {
		tipo string
		ok   bool
	}{
		{TypeLogin, true},
		{TypeNote, true},
		{TypeCard, true},
		{"", false},
		{"ssh", false},
	} {
		err := ValidateType(tc.tipo)
		if tc.ok && err != nil {
			t.Fatalf("%q devia passar: %v", tc.tipo, err)
		}
		if !tc.ok && !errors.Is(err, ErrInvalidType) {
			t.Fatalf("%q devia falhar com ErrInvalidType", tc.tipo)
		}
	}
}
