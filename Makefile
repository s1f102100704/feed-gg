.PHONY: migrate-status migrate-up migrate-down migrate-reset migrate-create

MIGRATE = docker compose --profile tools run --rm migrate

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
