.PHONY: lint lint-backend migrate-status migrate-up migrate-down migrate-reset migrate-create generate-regions

MIGRATE = docker compose --profile tools run --rm migrate
GOLANGCI_LINT = docker compose --profile tools run --rm golangci-lint

lint: lint-backend

lint-backend:
	$(GOLANGCI_LINT) run

migrate-status:
	$(MIGRATE) status

migrate-up:
	$(MIGRATE) up

migrate-down:
	$(MIGRATE) down

migrate-reset:
	$(MIGRATE) reset

migrate-create:
	@if [ -z "$(name)" ]; then echo "usage: make migrate-create name=add_players"; exit 1; fi
	$(MIGRATE) -s create $(name) sql

generate-regions:
	node scripts/generate-regions.mjs
