# Makefile — atalhos de desenvolvimento do AegisPass.
#
# Didático: cada "alvo" (target) é um comando nomeado. Corre com `make <alvo>`.
# .PHONY diz ao make que estes alvos NÃO são ficheiros (evita conflitos com
# ficheiros/pastas com o mesmo nome).
.PHONY: help backend-run backend-test backend-build frontend-install frontend-dev frontend-test test db-backup backup-gen-key

help:
	@echo "Alvos disponiveis:"
	@echo "  backend-run       - arranca a API Go"
	@echo "  backend-test      - corre os testes Go (inclui cripto)"
	@echo "  backend-build     - compila o servidor, backup e a CLI"
	@echo "  frontend-install  - instala dependencias do frontend"
	@echo "  frontend-dev      - arranca o dev server do frontend"
	@echo "  frontend-test     - corre os testes do frontend (vitest)"
	@echo "  test              - corre todos os testes"
	@echo "  backup-gen-key    - gera AEGIS_BACKUP_KEY para o .env"
	@echo "  db-backup         - pg_dump cifrado para backups/ (docker compose + .env)"

backend-run:
	cd backend && go run ./cmd/server

backend-test:
	cd backend && go test ./...

backend-build:
	cd backend && go build ./...
	cd backend && go build -o ../bin/backup ./cmd/backup
	cd cli && go build ./...

backup-gen-key:
	cd backend && go run ./cmd/backup gen-key

# Requer: docker compose up (db healthy), AEGIS_BACKUP_KEY no ambiente
db-backup:
	mkdir -p backups
	docker compose exec -T db pg_dump -U aegis aegis | (cd backend && go run ./cmd/backup encrypt -o ../backups/aegis-$$(date +%Y%m%d-%H%M%S).sql.enc)

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-test:
	cd frontend && npm run test

test: backend-test frontend-test
