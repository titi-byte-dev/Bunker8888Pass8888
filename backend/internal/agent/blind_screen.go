package agent

import (
	"regexp"
	"strings"
)

// Padrões de campos que a triagem às cegas deve ocultar (RGPD / não discriminação).
// Didático: género e etnia não influenciam competências — removemos antes do agente ver.
var blindFieldPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(g[eé]nero|gender|sexo|sex)\s*[:=]\s*.+`),
	regexp.MustCompile(`(?i)(etnia|ethnicity|ra[cç]a|race|cor\s+da\s+pele)\s*[:=]\s*.+`),
	regexp.MustCompile(`(?i)(estado\s+civil|marital\s+status)\s*[:=]\s*.+`),
	regexp.MustCompile(`(?i)(idade|age|data\s+de\s+nascimento|date\s+of\s+birth|nascimento)\s*[:=]\s*.+`),
}

const blindRedacted = "[oculto — triagem às cegas]"

// IsRecruitmentEmail detecta e-mails de candidatura (AGENT-007).
func IsRecruitmentEmail(subject, body string) bool {
	text := strings.ToLower(subject + " " + body)
	keywords := []string{
		"candidatura", "candidato", "candidata", "curriculum", "currículo", "curriculo",
		" cv ", "cv:", "recrutamento", "recruitment", "application", "candidacy", "candidature",
	}
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return strings.Contains(strings.ToLower(subject), "cv")
}

// BlindScreenCV remove ou mascara campos sensíveis para não discriminação.
func BlindScreenCV(text string) string {
	lines := strings.Split(text, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := line
		for _, re := range blindFieldPatterns {
			trimmed = re.ReplaceAllString(trimmed, blindRedacted)
		}
		out = append(out, trimmed)
	}
	return strings.Join(out, "\n")
}
