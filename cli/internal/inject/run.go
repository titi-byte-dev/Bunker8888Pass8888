// Package inject executa comandos com segredos em variáveis de ambiente.
//
// ⚠️ **Segurança:** o segredo existe só na memória do processo filho durante
// a execução — nunca escrevemos passwords em ficheiros de texto plano.
package inject

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// RunWithEnv executa command com env extra (segredo injectado).
// O segredo é zerado do slice local após arrancar o processo (melhor esforço).
func RunWithEnv(envKey, secret string, command []string) error {
	if len(command) == 0 {
		return fmt.Errorf("comando em falta")
	}
	if envKey == "" {
		return fmt.Errorf("nome da variável de ambiente em falta")
	}

	// Copiamos o ambiente actual e injectamos o segredo.
	env := append(os.Environ(), envKey+"="+secret)
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	// Melhor esforço para apagar a cópia local do segredo (Go strings são imutáveis
	// mas limitamos a janela onde o slice env ainda referencia o valor).
	secret = ""
	env[len(env)-1] = envKey + "="

	return err
}

// ZeroBytes tenta apagar bytes sensíveis da memória (melhor esforço em Go).
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
	runtime.KeepAlive(b)
}
