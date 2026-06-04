# Makefile — atalhos de desenvolvimento do AegisPass.
#
# Didático: cada "alvo" (target) é um comando nomeado. Corre com `make <alvo>`.
# .PHONY diz ao make que estes alvos NÃO são ficheiros (evita conflitos com
# ficheiros/pastas com o mesmo nome).
.PHONY: help backend-run backend-test backend-build frontend-install frontend-dev frontend-test test

help:
	@echo "Alvos disponiveis:"
	@echo "  backend-run       - arranca a API Go"
	@echo "  backend-test      - corre os testes Go (inclui cripto)"
	@echo "  backend-build     - compila o servidor e a CLI"
	@echo "  frontend-install  - instala dependencias do frontend"
	@echo "  frontend-dev      - arranca o dev server do frontend"
	@echo "  frontend-test     - corre os testes do frontend (vitest)"
	@echo "  test              - corre todos os testes"

backend-run:
	cd backend && go run ./cmd/server

backend-test:
	cd backend && go test ./...

backend-build:
	cd backend && go build ./...
	cd cli && go build ./...

frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-test:
	cd frontend && npm run test

test: backend-test frontend-test
