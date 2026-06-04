// Package geofence valida acesso por IP (CIDR) e GPS (VAULT-011).
//
// Didático: geofencing restringe ONDE o pedido pode originar-se. IP identifica
// a rede (útil com VPN corporativa); GPS identifica o dispositivo (BYOD no
// escritório). Quando ambos estão configurados, exigimos os dois (defesa em
// profundidade).
package geofence

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/netip"
	"strconv"
	"strings"
)

// ErrOutsideFence indica que IP ou GPS estão fora da zona permitida.
var ErrOutsideFence = errors.New("geofence: fora da zona permitida")

// ErrGPSRequired indica que a política exige GPS mas o cliente não enviou coords.
var ErrGPSRequired = errors.New("geofence: coordenadas GPS em falta")

// Policy agrupa regras de geofencing de um utilizador.
type Policy struct {
	UserID       string
	Enabled      bool
	AllowedCIDRs []string
	GPSEnabled   bool
	GPSLat       *float64
	GPSLon       *float64
	GPSRadiusM   float64
}

// ClientGeo são coordenadas enviadas pelo cliente (headers X-Geo-*).
type ClientGeo struct {
	Lat float64
	Lon float64
	Ok  bool
}

// IsAllowed devolve true se o pedido cumpre a política activa.
func IsAllowed(p Policy, clientIP string, geo ClientGeo) (bool, error) {
	if !p.Enabled {
		return true, nil
	}

	ipRequired := len(p.AllowedCIDRs) > 0
	gpsRequired := p.GPSEnabled && p.GPSLat != nil && p.GPSLon != nil

	if !ipRequired && !gpsRequired {
		// Política activada mas sem regras concretas — não bloqueamos (dev/admin).
		return true, nil
	}

	if ipRequired {
		ok, err := ipInAllowedCIDRs(clientIP, p.AllowedCIDRs)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}

	if gpsRequired {
		if !geo.Ok {
			return false, nil
		}
		dist := haversineMeters(geo.Lat, geo.Lon, *p.GPSLat, *p.GPSLon)
		if dist > p.GPSRadiusM {
			return false, nil
		}
	}

	return true, nil
}

// AssertAllowed devolve ErrOutsideFence ou ErrGPSRequired quando aplicável.
func AssertAllowed(p Policy, clientIP string, geo ClientGeo) error {
	if !p.Enabled {
		return nil
	}

	gpsRequired := p.GPSEnabled && p.GPSLat != nil && p.GPSLon != nil
	if gpsRequired && !geo.Ok {
		return ErrGPSRequired
	}

	ok, err := IsAllowed(p, clientIP, geo)
	if err != nil {
		return err
	}
	if !ok {
		return ErrOutsideFence
	}
	return nil
}

func ipInAllowedCIDRs(ipStr string, cidrs []string) (bool, error) {
	ipStr = strings.TrimSpace(ipStr)
	if ipStr == "" {
		return false, nil
	}
	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return false, fmt.Errorf("IP inválido %q: %w", ipStr, err)
	}
	for _, cidr := range cidrs {
		prefix, err := netip.ParsePrefix(strings.TrimSpace(cidr))
		if err != nil {
			return false, fmt.Errorf("CIDR inválido %q: %w", cidr, err)
		}
		if prefix.Contains(addr) {
			return true, nil
		}
	}
	return false, nil
}

// haversineMeters calcula distância entre dois pontos WGS84 em metros.
func haversineMeters(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusM = 6_371_000
	rad := math.Pi / 180
	dLat := (lat2 - lat1) * rad
	dLon := (lon2 - lon1) * rad
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusM * c
}

// ParseCIDRsJSON desserializa allowed_cidrs da BD.
func ParseCIDRsJSON(raw []byte) ([]string, error) {
	if len(raw) == 0 {
		return nil, nil
	}
	var cidrs []string
	if err := json.Unmarshal(raw, &cidrs); err != nil {
		return nil, err
	}
	return cidrs, nil
}

// MarshalCIDRsJSON serializa CIDRs para JSONB.
func MarshalCIDRsJSON(cidrs []string) ([]byte, error) {
	if cidrs == nil {
		cidrs = []string{}
	}
	return json.Marshal(cidrs)
}

// ParseClientGeo lê headers X-Geo-Latitude / X-Geo-Longitude.
func ParseClientGeo(latHdr, lonHdr string) (ClientGeo, error) {
	latHdr = strings.TrimSpace(latHdr)
	lonHdr = strings.TrimSpace(lonHdr)
	if latHdr == "" && lonHdr == "" {
		return ClientGeo{}, nil
	}
	if latHdr == "" || lonHdr == "" {
		return ClientGeo{}, fmt.Errorf("latitude e longitude devem vir em par")
	}
	lat, err1 := strconv.ParseFloat(latHdr, 64)
	lon, err2 := strconv.ParseFloat(lonHdr, 64)
	if err1 != nil || err2 != nil {
		return ClientGeo{}, fmt.Errorf("coordenadas inválidas")
	}
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return ClientGeo{}, fmt.Errorf("coordenadas fora de intervalo")
	}
	return ClientGeo{Lat: lat, Lon: lon, Ok: true}, nil
}
