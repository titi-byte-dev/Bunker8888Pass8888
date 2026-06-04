# /crypto-check

Valida que código de criptografia novo ou alterado segue as regras do projeto.

## Checklist a aplicar

1. **Biblioteca:** usa `crypto/*` (Go) ou WebCrypto (browser)? 🚫 Nada caseiro.
2. **Algoritmo:** AES-GCM-256 para simétrica; Argon2id/PBKDF2 (alto custo) para
   derivação de chave.
3. **Nonce/IV:** único por operação, de `crypto/rand` / `crypto.getRandomValues`;
   nunca reutilizado com a mesma chave; armazenado junto do ciphertext.
4. **Chaves:** Master Key só no cliente; servidor só recebe o Auth Hash;
   `extractable: false` no WebCrypto quando aplicável.
5. **Autenticação:** modo autenticado (GCM) para detetar adulteração.
6. **Testes:** existem test vectors / table-driven tests para o novo código?
   (ver [`docs/roadmap/06-testing/security-testing.md`](../../docs/roadmap/06-testing/security-testing.md))

## Resultado

Lista cada item como ✅ / ⚠️ / 🔴 com a localização (ficheiro:linha) e a correção
recomendada. Não alterar código sem confirmação.
