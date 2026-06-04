// Command backup cifra/decifra dumps PostgreSQL (INFRA-004).
//
// Uso típico (com docker compose a correr):
//
//	docker compose exec -T db pg_dump -U aegis aegis | go run ./cmd/backup encrypt -o backups/dump.sql.enc
//	go run ./cmd/backup decrypt -i backups/dump.sql.enc | docker compose exec -T db psql -U aegis aegis
//
// A chave vem de AEGIS_BACKUP_KEY (32 bytes em base64). Gera uma nova com:
//
//	go run ./cmd/backup gen-key
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/pkg/backup"
)

func main() {
	log.SetFlags(0)

	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "encrypt":
		runEncrypt(os.Args[2:])
	case "decrypt":
		runDecrypt(os.Args[2:])
	case "gen-key":
		runGenKey()
	default:
		usage()
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, `AegisPass — backup cifrado da base de dados

Subcomandos:
  encrypt [-o ficheiro.enc]   Lê stdin, cifra com AEGIS_BACKUP_KEY, escreve ficheiro
  decrypt [-i ficheiro.enc]   Decifra ficheiro, escreve plaintext para stdout
  gen-key                     Gera uma AEGIS_BACKUP_KEY (base64) para o .env

Exemplo:
  docker compose exec -T db pg_dump -U aegis aegis | go run ./cmd/backup encrypt -o backups/dump.enc
`)
}

func runEncrypt(args []string) {
	fs := flag.NewFlagSet("encrypt", flag.ExitOnError)
	out := fs.String("o", "", "ficheiro de saída (.enc)")
	fs.Parse(args)

	key, err := backup.KeyFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	plain, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	enc, err := backup.EncryptDump(key, plain)
	if err != nil {
		log.Fatal(err)
	}
	if *out == "" {
		log.Fatal("especifica -o ficheiro.enc")
	}
	if err := os.WriteFile(*out, enc, 0o600); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "backup cifrado: %s (%d bytes)\n", *out, len(enc))
}

func runDecrypt(args []string) {
	fs := flag.NewFlagSet("decrypt", flag.ExitOnError)
	in := fs.String("i", "", "ficheiro .enc")
	fs.Parse(args)

	key, err := backup.KeyFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	if *in == "" {
		log.Fatal("especifica -i ficheiro.enc")
	}
	enc, err := os.ReadFile(*in)
	if err != nil {
		log.Fatal(err)
	}
	plain, err := backup.DecryptDump(key, enc)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(plain); err != nil {
		log.Fatal(err)
	}
}

func runGenKey() {
	k, err := backup.GenerateKeyBase64()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(k)
	fmt.Fprintln(os.Stderr, "Adiciona ao .env: AEGIS_BACKUP_KEY="+k)
}
