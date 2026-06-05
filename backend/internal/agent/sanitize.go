package agent

import (
	"regexp"
	"strings"
)

// Padrões comuns de prompt injection em e-mails / páginas externas.
// Didático: não tentamos "entender" o ataque — removemos frases que pedem ao
// modelo para ignorar regras ou revelar segredos.
var injectionPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)ignore\s+(all\s+)?(previous|prior)\s+instructions`),
	regexp.MustCompile(`(?i)disregard\s+(all\s+)?(previous|prior)\s+instructions`),
	regexp.MustCompile(`(?i)you\s+are\s+now\s+`),
	regexp.MustCompile(`(?i)system\s*:\s*`),
	regexp.MustCompile(`(?i)reveal\s+(the\s+)?(master\s+)?(key|password|secret)`),
	regexp.MustCompile(`(?i)jailbreak`),
}

// SanitizeExternalContent remove linhas suspeitas de conteúdo lido por agentes.
// AGENT-010: dados externos nunca são instruções — só contexto delimitado.
func SanitizeExternalContent(body string) string {
	lines := strings.Split(body, "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" {
			out = append(out, line)
			continue
		}
		skip := false
		for _, re := range injectionPatterns {
			if re.MatchString(trim) {
				skip = true
				break
			}
		}
		if !skip {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

// WrapExternalData envolve texto externo para o LLM tratar como dados, não ordens.
func WrapExternalData(source, body string) string {
	clean := SanitizeExternalContent(body)
	if len(clean) > 8000 {
		clean = clean[:8000] + "\n…[truncado]"
	}
	return "<external_data source=\"" + escapeXMLAttr(source) + "\">\n" + clean + "\n</external_data>"
}

func escapeXMLAttr(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	return s
}

// RejectIfLooksLikeInstruction falha cedo se o input parece uma ordem ao agente.
func RejectIfLooksLikeInstruction(text string) error {
	lower := strings.ToLower(strings.TrimSpace(text))
	if strings.HasPrefix(lower, "system:") || strings.HasPrefix(lower, "assistant:") {
		return ErrExternalAsCommand
	}
	for _, re := range injectionPatterns {
		if re.MatchString(text) {
			return ErrExternalAsCommand
		}
	}
	return nil
}
