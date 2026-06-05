package sentinel

import (
	"testing"
	"time"
)

func TestAssessTravel_normal(t *testing.T) {
	prev := Point{Lat: 38.7223, Lon: -9.1393, At: time.Now().Add(-24 * time.Hour)}
	curr := Point{Lat: 41.1579, Lon: -8.6291, At: time.Now()} // Porto — ~300 km
	a := AssessTravel(prev, curr)
	if a.Suspicious {
		t.Fatalf("viagem plausível marcada como suspeita: %s", a.Detail)
	}
}

func TestAssessTravel_impossible(t *testing.T) {
	prev := Point{Lat: 38.7223, Lon: -9.1393, At: time.Now().Add(-30 * time.Minute)} // Lisboa
	curr := Point{Lat: 35.6762, Lon: 139.6503, At: time.Now()}                      // Tóquio
	a := AssessTravel(prev, curr)
	if !a.Suspicious || a.Reason != "impossible_travel" {
		t.Fatalf("esperava viagem impossível, got %+v", a)
	}
}

func TestAssessTravel_shortDistanceIgnored(t *testing.T) {
	prev := Point{Lat: 38.72, Lon: -9.14, At: time.Now().Add(-5 * time.Minute)}
	curr := Point{Lat: 38.73, Lon: -9.13, At: time.Now()}
	a := AssessTravel(prev, curr)
	if a.Suspicious {
		t.Fatal("distância curta não devia disparar sentinel")
	}
}
