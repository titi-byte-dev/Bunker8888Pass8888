package agent

import (
	"strings"
	"testing"
)

func TestBlindScreenCV(t *testing.T) {
	in := "Nome: Joana\nGénero: Feminino\nEtnia: X\nSkills: Go, SQL"
	out := BlindScreenCV(in)
	if strings.Contains(out, "Feminino") {
		t.Fatalf("género não foi ocultado: %q", out)
	}
	if !strings.Contains(out, "Joana") || !strings.Contains(out, "Go") {
		t.Fatalf("dados profissionais devem permanecer: %q", out)
	}
}

func TestIsRecruitmentEmail(t *testing.T) {
	if !IsRecruitmentEmail("Candidatura Engenheira", "Segue CV em anexo") {
		t.Fatal("devia detectar candidatura")
	}
	if IsRecruitmentEmail("Reunião amanhã", "Confirmar hora") {
		t.Fatal("não devia ser recrutamento")
	}
}
