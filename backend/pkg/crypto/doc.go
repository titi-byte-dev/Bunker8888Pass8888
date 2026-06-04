// Package crypto contém as primitivas criptográficas do AegisPass usadas pelo
// servidor e pela CLI (Go).
//
// # Fronteira Zero-Knowledge
//
// IMPORTANTE: no modelo Zero-Knowledge, a derivação da chave a partir da
// Master Password e a cifragem dos itens do cofre acontecem NO CLIENTE
// (frontend, ver frontend/src/lib/crypto.ts). O servidor nunca conhece a
// Master Key nem os dados em claro.
//
// Este pacote Go existe para:
//   - a CLI (Go), que precisa de decifrar o cofre localmente (VAULT-017);
//   - operações server-side legítimas que NÃO envolvem a Master Key do
//     utilizador (ex: cifrar metadados próprios do servidor, o "Guardião");
//   - servir de implementação de referência testável (test vectors).
//
// Está colocado em `pkg/` (e não em `internal/`) precisamente para poder ser
// importado por outros módulos do monorepo, como a CLI.
package crypto
