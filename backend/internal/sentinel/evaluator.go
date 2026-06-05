// Package sentinel implementa o Sentinel Mode (DW-004): deteção de login
// geograficamente impossível e step-up com passkey.
package sentinel

import (
	"fmt"
	"math"
	"time"
)

// MaxTravelSpeedKmh é a velocidade máxima plausível (voo comercial + margem).
const MaxTravelSpeedKmh = 900.0

// MinTravelDistanceM evita falsos positivos quando o utilizador está parado
// mas o GPS oscila ligeiramente.
const MinTravelDistanceM = 50_000

// Point representa uma localização observada num login.
type Point struct {
	Lat float64
	Lon float64
	At  time.Time
}

// Assessment resume se o login actual parece impossível face ao anterior.
type Assessment struct {
	Suspicious bool
	Reason     string
	Detail     string
}

// AssessTravel compara dois logins com GPS válido.
//
// Didático: dividimos distância (Haversine) pelo tempo decorrido. Se a
// velocidade implícita exceder ~900 km/h, é fisicamente implausível —
// sinal clássico de sessão roubada ou credencial partilhada.
func AssessTravel(previous, current Point) Assessment {
	if previous.At.IsZero() || current.At.IsZero() {
		return Assessment{}
	}
	elapsed := current.At.Sub(previous.At)
	if elapsed <= 0 {
		return Assessment{}
	}

	distM := haversineMeters(previous.Lat, previous.Lon, current.Lat, current.Lon)
	if distM < MinTravelDistanceM {
		return Assessment{}
	}

	hours := elapsed.Hours()
	if hours <= 0 {
		return Assessment{Suspicious: true, Reason: "impossible_travel", Detail: "Intervalo de tempo inválido"}
	}

	speedKmh := (distM / 1000) / hours
	if speedKmh > MaxTravelSpeedKmh {
		return Assessment{
			Suspicious: true,
			Reason:     "impossible_travel",
			Detail: fmt.Sprintf(
				"~%.0f km em %.0f min (%.0f km/h)",
				distM/1000,
				elapsed.Minutes(),
				speedKmh,
			),
		}
	}
	return Assessment{}
}

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
