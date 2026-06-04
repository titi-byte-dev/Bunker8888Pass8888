// Command aegis é a CLI do AegisPass (VAULT-017).
//
// Injeta segredos do cofre em scripts via mTLS + decifragem local (Zero-Knowledge).
// A Master Password nunca sai do processo excepto para derivar a chave em memória.
package main

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"

	"github.com/titi-byte-dev/Bunker8888Pass8888/cli/internal/api"
	"github.com/titi-byte-dev/Bunker8888Pass8888/cli/internal/auth"
	clicfg "github.com/titi-byte-dev/Bunker8888Pass8888/cli/internal/config"
	"github.com/titi-byte-dev/Bunker8888Pass8888/cli/internal/device"
	"github.com/titi-byte-dev/Bunker8888Pass8888/cli/internal/inject"
	"github.com/titi-byte-dev/Bunker8888Pass8888/cli/internal/vault"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	switch cmd {
	case "version", "--version", "-version":
		fmt.Printf("aegis %s\n", version)
	case "device":
		if len(args) == 0 || args[0] != "register" {
			fatal("uso: aegis device register --name <nome> --email <email>")
		}
		runDeviceRegister(args[1:])
	case "list":
		runList(args)
	case "run":
		runInject(args)
	case "help", "-h", "--help":
		printUsage()
	default:
		fatal("comando desconhecido: %s (usa 'aegis help')", cmd)
	}
}

func printUsage() {
	fmt.Println(`AegisPass CLI — injecção de segredos via mTLS (Zero-Knowledge)

Comandos:
  aegis device register --name <nome> --email <email>
      Regista este computador (login + certificado mTLS).

  aegis list [--type login]
      Lista itens do cofre (requer dispositivo registado).

  aegis run --item <id> --field password --env SECRET -- <comando...>
      Decifra o campo e executa o comando com a variável de ambiente injectada.

  aegis version
      Mostra a versão.

Variáveis de ambiente:
  AEGIS_API_URL     API HTTP (default http://localhost:8080)
  AEGIS_MTLS_URL    API mTLS (default https://localhost:8443)
  AEGIS_CONFIG_DIR  Pasta de certificados (default ~/.aegis)
  AEGIS_MTLS_INSECURE=1  Ignora verificação TLS do servidor (só dev)`)
}

func runDeviceRegister(args []string) {
	name := flagValue(args, "--name")
	email := flagValue(args, "--email")
	if name == "" || email == "" {
		fatal("--name e --email são obrigatórios")
	}
	password := promptPassword("Master Password: ")

	apiBase := clicfg.APIBase()
	res, err := auth.Login(apiBase, email, password)
	if err != nil {
		fatal("login: %v", err)
	}
	inject.ZeroBytes([]byte(password))

	csr, keyPEM, err := device.GenerateCSR(name)
	if err != nil {
		fatal("CSR: %v", err)
	}
	reg, err := api.RegisterDevice(apiBase, res.Token, name, csr)
	if err != nil {
		fatal("registo: %v", err)
	}

	store, err := clicfg.LoadStore()
	if err != nil {
		fatal("config: %v", err)
	}
	if err := device.SaveCredentials(store, []byte(reg.CertPEM), []byte(reg.CAPEM), keyPEM, email); err != nil {
		fatal("guardar certificados: %v", err)
	}
	fmt.Printf("Dispositivo %q registado. Certificados em %s\n", name, store.Dir)
}

func runList(args []string) {
	store, err := clicfg.LoadStore()
	if err != nil || !store.HasDevice() {
		fatal("dispositivo não registado — corre: aegis device register ...")
	}
	client, err := api.NewMTLS(clicfg.MTLSBase(), store.ClientCert, store.ClientKey, store.CACert)
	if err != nil {
		fatal("mTLS: %v", err)
	}
	itemType := flagValue(args, "--type")
	items, err := client.ListVaultItems(itemType)
	if err != nil {
		fatal("list: %v", err)
	}
	for _, it := range items {
		fmt.Printf("%s\t%s\t%s\n", it.ID, it.Type, it.UpdatedAt)
	}
}

func runInject(args []string) {
	itemID := flagValue(args, "--item")
	field := flagValue(args, "--field")
	if field == "" {
		field = "password"
	}
	envKey := flagValue(args, "--env")
	if envKey == "" {
		envKey = "AEGIS_SECRET"
	}
	cmdArgs := argsAfterDoubleDash(args)
	if itemID == "" || len(cmdArgs) == 0 {
		fatal("uso: aegis run --item <id> [--field password] [--env VAR] -- <comando...>")
	}

	store, err := clicfg.LoadStore()
	if err != nil || !store.HasDevice() {
		fatal("dispositivo não registado")
	}
	email, err := device.ReadEmail(store)
	if err != nil {
		fatal("email em falta no config — re-regista o dispositivo")
	}

	password := os.Getenv("AEGIS_MASTER_PASSWORD")
	if password == "" {
		password = promptPassword("Master Password: ")
	}
	mk, err := auth.MasterKeyFromPassword(email, password, clicfg.APIBase())
	inject.ZeroBytes([]byte(password))
	if err != nil {
		fatal("master key: %v", err)
	}
	defer inject.ZeroBytes(mk)

	client, err := api.NewMTLS(clicfg.MTLSBase(), store.ClientCert, store.ClientKey, store.CACert)
	if err != nil {
		fatal("mTLS: %v", err)
	}
	item, err := client.GetVaultItem(itemID)
	if err != nil {
		fatal("item: %v", err)
	}
	login, err := vault.OpenLogin(mk, []byte(item.Blob))
	if err != nil {
		fatal("decifrar: %v", err)
	}
	secret, err := vault.FieldValue(login, field)
	if err != nil {
		fatal("%v", err)
	}
	if err := inject.RunWithEnv(envKey, secret, cmdArgs); err != nil {
		os.Exit(1)
	}
}

func flagValue(args []string, name string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == name && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(args[i], name+"=") {
			return strings.TrimPrefix(args[i], name+"=")
		}
	}
	return ""
}

func argsAfterDoubleDash(args []string) []string {
	for i, a := range args {
		if a == "--" {
			return args[i+1:]
		}
	}
	return nil
}

func promptPassword(label string) string {
	fmt.Fprint(os.Stderr, label)
	b, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		fatal("password: %v", err)
	}
	return string(b)
}

func fatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "erro: "+format+"\n", args...)
	os.Exit(1)
}
