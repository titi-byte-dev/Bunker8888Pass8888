# Notas pessoais (`_private/`)

Esta pasta é para as **tuas notas pessoais e didáticas** — rascunhos, dúvidas,
resumos de estudo de Go/Svelte/cripto, links, etc.

## Como funciona o "esconder" no git

O git é transparente: quem clona um repositório vê todos os ficheiros versionados.
Por isso há duas formas de manter estas notas "escondidas".

### Opção 1 — Local apenas (ativa por defeito) ✅

A pasta `_private/` está no `.gitignore` (a regra `**/_private/`), por isso
**nada aqui dentro é enviado para o git** — exceto este `README.md`, que é a
exceção explícita (`!**/_private/README.md`) para a equipa perceber a convenção.

- ✅ Privado de verdade.
- ⚠️ Sem backup no git: se perderes a máquina, perdes as notas. Faz backup à parte.

### Opção 2 — No git, mas cifrado (`git-crypt`) 🔐

Se quiseres as notas **dentro do git mas com conteúdo ilegível** para quem não
tem a chave (encaixa no tema Zero-Knowledge do AegisPass):

> 💡 **Conceito:** `git-crypt` cifra ficheiros de forma transparente ao fazer
> commit e decifra-os ao fazer checkout. No remoto ficam como bytes cifrados.

1. Instalar: `git-crypt` + `gnupg` (GPG).
2. No repositório: `git-crypt init`
3. Criar/editar `.gitattributes` na raiz com, por exemplo:
   ```gitattributes
   docs/roadmap/_private/** filter=git-crypt diff=git-crypt
   ```
4. Adicionar a tua chave: `git-crypt add-gpg-user <o-teu-email-GPG>`
5. (E **remover** a regra `**/_private/` do `.gitignore`, senão o git ignora-os.)

A partir daí, os ficheiros são versionados e sincronizados, mas o conteúdo
viaja cifrado. Só quem tiver a chave GPG (tu) os consegue ler.

---

Por defeito ficamos na **Opção 1**. Diz quando quiseres ativar a Opção 2 e eu
configuro o `git-crypt` e o `.gitattributes`.
