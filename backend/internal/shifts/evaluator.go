// Package shifts implementa políticas de acesso por turnos (VAULT-010).
//
// Didático: o servidor usa o relógio UTC (sincronizado via NTP no SO) e converte
// para o fuso do utilizador antes de comparar com a janela do turno. O cliente
// também valida o desvio do relógio local para mitigar manipulação.
package shifts

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrOutsideShift indica que o instante actual está fora das janelas permitidas.
var ErrOutsideShift = errors.New("shifts: fora do horário de turno")

// Window é um intervalo dentro de um dia (formato 24h "HH:MM").
type Window struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// WeeklySchedule mapeia dias ISO (mon..sun) a listas de janelas.
type WeeklySchedule map[string][]Window

// Policy agrupa a configuração de turno de um utilizador.
type Policy struct {
	UserID            string
	Timezone          string
	Schedule          WeeklySchedule
	Enabled           bool
	MaxClockSkewSecs  int
}

// dayKey devolve a chave do dia (mon..sun) para um instante no fuso indicado.
func dayKey(t time.Time, loc *time.Location) string {
	weekday := t.In(loc).Weekday()
	switch weekday {
	case time.Monday:
		return "mon"
	case time.Tuesday:
		return "tue"
	case time.Wednesday:
		return "wed"
	case time.Thursday:
		return "thu"
	case time.Friday:
		return "fri"
	case time.Saturday:
		return "sat"
	case time.Sunday:
		return "sun"
	default:
		return "mon"
	}
}

// parseHM converte "HH:MM" em minutos desde meia-noite.
func parseHM(s string) (int, error) {
	s = strings.TrimSpace(s)
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("formato inválido %q", s)
	}
	var h, m int
	if _, err := fmt.Sscanf(parts[0], "%d", &h); err != nil || h < 0 || h > 23 {
		return 0, fmt.Errorf("hora inválida em %q", s)
	}
	if _, err := fmt.Sscanf(parts[1], "%d", &m); err != nil || m < 0 || m > 59 {
		return 0, fmt.Errorf("minuto inválido em %q", s)
	}
	return h*60 + m, nil
}

// IsWithinShift devolve true se `now` cai dentro de alguma janela do turno.
// Se enabled=false, devolve sempre true.
func IsWithinShift(now time.Time, p Policy) (bool, error) {
	if !p.Enabled {
		return true, nil
	}
	loc, err := time.LoadLocation(p.Timezone)
	if err != nil {
		return false, fmt.Errorf("timezone inválido %q: %w", p.Timezone, err)
	}
	local := now.In(loc)
	key := dayKey(local, loc)
	windows := p.Schedule[key]
	if len(windows) == 0 {
		return false, nil
	}

	minutes := local.Hour()*60 + local.Minute()
	for _, w := range windows {
		start, err1 := parseHM(w.Start)
		end, err2 := parseHM(w.End)
		if err1 != nil || err2 != nil {
			return false, fmt.Errorf("janela inválida: %v %v", err1, err2)
		}
		// Janela normal (ex: 09:00–17:00) ou overnight (ex: 22:00–06:00).
		if start <= end {
			if minutes >= start && minutes < end {
				return true, nil
			}
		} else if minutes >= start || minutes < end {
			return true, nil
		}
	}
	return false, nil
}

// AssertWithinShift devolve ErrOutsideShift se o acesso não for permitido.
func AssertWithinShift(now time.Time, p Policy) error {
	ok, err := IsWithinShift(now, p)
	if err != nil {
		return err
	}
	if !ok {
		return ErrOutsideShift
	}
	return nil
}

// ParseSchedule desserializa JSONB da BD.
func ParseSchedule(raw []byte) (WeeklySchedule, error) {
	if len(raw) == 0 || string(raw) == "{}" {
		return WeeklySchedule{}, nil
	}
	var s WeeklySchedule
	if err := json.Unmarshal(raw, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// MarshalSchedule serializa o schedule para JSONB.
func MarshalSchedule(s WeeklySchedule) ([]byte, error) {
	if s == nil {
		s = WeeklySchedule{}
	}
	return json.Marshal(s)
}
