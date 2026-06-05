package ops

import "testing"

func TestIsLowStock(t *testing.T) {
	tests := []struct {
		qty, reorder int
		want         bool
	}{
		{0, 5, true},
		{5, 5, true},
		{6, 5, false},
		{10, 0, false},
	}
	for _, tc := range tests {
		if got := IsLowStock(tc.qty, tc.reorder); got != tc.want {
			t.Fatalf("IsLowStock(%d,%d)=%v want %v", tc.qty, tc.reorder, got, tc.want)
		}
	}
}
