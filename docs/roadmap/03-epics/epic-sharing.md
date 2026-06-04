# Epic: Colaboração & Partilha Blindada — `SHARE`

> **Fase:** 2 · **Prioridade:** 🟠 Média

## Objetivo

Eliminar o envio de passwords por apps de chat/mensagens ou e-mail, oferecendo
partilha cifrada de cofres e links secretos temporários.

## Valor de negócio

Substitui o "anti-padrão" de partilhar segredos em apps de chat. Canal seguro
para informação crítica.

## Funcionalidades

- **Cofres partilhados (Shared Vaults):** pastas partilhadas por departamento;
  o RH/admin controla quem entra e sai
- **Secret Links:** link temporário para password/documento; expira após 1
  clique ou X minutos; servido via RAM, nunca persistido em disco após expirar
- **Anexo de ficheiros cifrados:** contratos, chaves `.pem`, cifrados por
  ficheiro antes do upload
- **Canais de notas cifradas temporárias:** mensagens que se apagam após leitura

## Critérios de aceitação

- [ ] Secret link expira corretamente (clique único e/ou tempo) e fica
      inacessível depois.
- [ ] Após expirar, não há rasto do segredo em disco (só existiu em RAM).
- [ ] Partilha de cofre respeita permissões e revogação imediata.
- [ ] Ficheiros anexados são cifrados com chave própria por ficheiro.

## Conceitos didáticos

> 💡 **Partilha em Zero-Knowledge:** como o servidor não tem a chave, partilhar
> um segredo significa **re-cifrar a chave do item para a chave pública do
> destinatário** (criptografia assimétrica). Assim o servidor encaminha sem
> nunca conseguir ler.

> ⚠️ **Segurança — "servido via RAM":** garantir que dados de links efémeros
> nunca tocam o disco exige cuidado: desativar swap para esse processo, não os
> escrever em logs, e ter um TTL agressivo em memória.

## Dependências

- Vault (`VAULT`) e modelo de chaves assimétricas por utilizador.

## Tasks

Ver prefixo `SHARE-*` em [`../05-tasks/backlog.md`](../05-tasks/backlog.md).
