package shifts

import (
	"testing"
	"time"
)

func mustLoc(t *testing.T, name string) *time.Location {
	t.Helper()
	loc, err := time.LoadLocation(name)
	if err != nil {
		t.Fatal(err)
	}
	return loc
}

func TestIsWithinShift_DisabledAlwaysAllows(t *testing.T) {
	p := Policy{Enabled: false}
	ok, err := IsWithinShift(time.Now(), p)
	if err != nil || !ok {
		t.Fatalf("enabled=false devia permitir: ok=%v err=%v", ok, err)
	}
}

func TestIsWithinShift_WeekdayWindow(t *testing.T) {
	loc := mustLoc(t, "Europe/Lisbon")
	// Quarta-feira 2026-06-03 10:30 em Lisboa
	now := time.Date(2026, 6, 3, 10, 30, 0, 0, loc)

	p := Policy{
		Enabled:  true,
		Timezone: "Europe/Lisbon",
		Schedule: WeeklySchedule{
			"wed": {{Start: "09:00", End: "17:00"}},
		},
	}
	ok, err := IsWithinShift(now, p)
	if err != nil || !ok {
		t.Fatalf("dentro do turno: ok=%v err=%v", ok, err)
	}

	now = time.Date(2026, 6, 3, 18, 0, 0, 0, loc)
	ok, err = IsWithinShift(now, p)
	if err != nil || ok {
		t.Fatalf("fora do turno: ok=%v err=%v", ok, err)
	}
}

func TestIsWithinShift_OvernightWindow(t *testing.T) {
	loc := mustLoc(t, "UTC")
	now := time.Date(2026, 6, 3, 23, 30, 0, 0, loc)
	p := Policy{
		Enabled:  true,
		Timezone: "UTC",
		Schedule: WeeklySchedule{
			"wed": {{Start: "22:00", End: "06:00"}},
		},
	}
	ok, err := IsWithinShift(now, p)
	if err != nil || !ok {
		t.Fatalf("overnight dentro: ok=%v err=%v", ok, err)
	}
}

func TestAssertWithinShift_ReturnsError(t *testing.T) {
	loc := mustLoc(t, "UTC")
	now := time.Date(2026, 6, 3, 8, 0, 0, 0, loc)
	p := Policy{
		Enabled:  true,
		Timezone: "UTC",
		Schedule: WeeklySchedule{
			"wed": {{Start: "09:00", End: "17:00"}},
		},
	}
	if err := AssertWithinShift(now, p); err != ErrOutsideShift {
		t.Fatalf("esperado ErrOutsideShift, got %v", err)
	}
}
