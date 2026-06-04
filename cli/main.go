// Command aegis é a interface de linha de comandos (CLI) do AegisPass.
//
// Por agora é um esqueleto: imprime a versão e uma ajuda mínima. Será expandido
// na task VAULT-017 para injetar segredos em scripts via mTLS, sem nunca os
// escrever em ficheiros de texto plano.
package main

import (
	"flag"
	"fmt"
	"os"
)

// version é preenchida no build com -ldflags "-X main.version=...".
// Didático: assim o binário sabe a sua própria versão sem a ter "hardcoded".
var version = "dev"

func main() {
	// flag é o parser de argumentos da biblioteca padrão. Define a flag -version.
	showVersion := flag.Bool("version", false, "mostra a versão e sai")
	flag.Parse()

	if *showVersion {
		fmt.Printf("aegis %s\n", version)
		return
	}

	fmt.Println("AegisPass CLI (esqueleto). Usa --version para ver a versão.")
	fmt.Println("Funcionalidades a chegar na task VAULT-017.")
	os.Exit(0)
}
