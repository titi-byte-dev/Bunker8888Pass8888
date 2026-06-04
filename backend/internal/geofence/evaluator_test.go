package geofence

import (
	"testing"
)

func TestIPInAllowedCIDRs(t *testing.T) {
	ok, err := ipInAllowedCIDRs("10.0.0.5", []string{"10.0.0.0/24", "192.168.1.0/24"})
	if err != nil || !ok {
		t.Fatalf("esperado dentro da rede: ok=%v err=%v", ok, err)
	}
	ok, err = ipInAllowedCIDRs("8.8.8.8", []string{"10.0.0.0/24"})
	if err != nil || ok {
		t.Fatalf("esperado fora: ok=%v err=%v", ok, err)
	}
}

func TestHaversineMeters(t *testing.T) {
	// Lisboa ~ Porto ~ 275 km (ordem de grandeza)
	d := haversineMeters(38.7223, -9.1393, 41.1579, -8.6291)
	if d < 250_000 || d > 320_000 {
		t.Fatalf("distância inesperada: %f m", d)
	}
}

func TestIsAllowed_Disabled(t *testing.T) {
	ok, err := IsAllowed(Policy{Enabled: false}, "1.2.3.4", ClientGeo{})
	if err != nil || !ok {
		t.Fatal("disabled devia permitir")
	}
}

func TestIsAllowed_IPOnly(t *testing.T) {
	p := Policy{Enabled: true, AllowedCIDRs: []string{"127.0.0.1/32"}}
	ok, err := IsAllowed(p, "127.0.0.1", ClientGeo{})
	if err != nil || !ok {
		t.Fatalf("localhost devia passar: ok=%v err=%v", ok, err)
	}
	ok, err = IsAllowed(p, "10.0.0.1", ClientGeo{})
	if err != nil || ok {
		t.Fatalf("IP externo devia falhar: ok=%v err=%v", ok, err)
	}
}

func TestIsAllowed_GPSOnly(t *testing.T) {
	lat, lon := 38.7223, -9.1393
	p := Policy{
		Enabled:    true,
		GPSEnabled: true,
		GPSLat:     &lat,
		GPSLon:     &lon,
		GPSRadiusM: 500,
	}
	ok, err := IsAllowed(p, "8.8.8.8", ClientGeo{Lat: lat, Lon: lon, Ok: true})
	if err != nil || !ok {
		t.Fatalf("dentro do raio: ok=%v err=%v", ok, err)
	}
	ok, err = IsAllowed(p, "8.8.8.8", ClientGeo{Lat: 41.0, Lon: -8.0, Ok: true})
	if err != nil || ok {
		t.Fatalf("longe devia falhar: ok=%v err=%v", ok, err)
	}
}

func TestAssertAllowed_GPSRequired(t *testing.T) {
	lat, lon := 38.0, -9.0
	p := Policy{Enabled: true, GPSEnabled: true, GPSLat: &lat, GPSLon: &lon, GPSRadiusM: 100}
	if err := AssertAllowed(p, "", ClientGeo{}); err != ErrGPSRequired {
		t.Fatalf("esperado ErrGPSRequired, got %v", err)
	}
}

func TestParseClientGeo(t *testing.T) {
	g, err := ParseClientGeo("38.7223", "-9.1393")
	if err != nil || !g.Ok {
		t.Fatalf("parse falhou: %v", err)
	}
	if _, err := ParseClientGeo("bad", "-9"); err == nil {
		t.Fatal("devia rejeitar lat inválida")
	}
}
