# load env file
ENVFILE ?= .env
include $(ENVFILE)
export

MIGRATE_DSN ?= "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable"
MIGRATE := docker run --rm -v $(shell pwd)/migrations:/migrations --user "$(shell id -u):$(shell id -g)" --network host migrate/migrate -path=/migrations -database "$(MIGRATE_DSN)"

.PHONY: default
default: help

# generate help info from comments: thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## help information about make commands
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: test-server
test-server: ## run server tests
	go test ./internal/server -v

.PHONY: sync-sqlboiler-conf
sync-sqlboiler-conf: ## sync sqlboiler config file with passed envs
	@echo "Syncing sqlboiler config file with passed envs..."
	@sed -r -i 's/[[:space:]]*dbname[[:space:]]*=[[:space:]]*(.+)/dbname = "$(POSTGRES_DB)"/' sqlboiler.toml
	@sed -r -i 's/[[:space:]]*port[[:space:]]*=[[:space:]]*(.+)/port = $(POSTGRES_PORT)/' sqlboiler.toml
	@sed -r -i 's/[[:space:]]*user[[:space:]]*=[[:space:]]*(.+)/user = "$(POSTGRES_USER)"/' sqlboiler.toml
	@sed -r -i 's/[[:space:]]*pass[[:space:]]*=[[:space:]]*(.+)/pass = "$(POSTGRES_PASSWORD)"/' sqlboiler.toml

.PHONY: generate-models
generate-models: sync-sqlboiler-conf ## run sqlboiler to generate models
	sqlboiler --config sqlboiler.toml --output models --no-auto-timestamps --wipe psql

.PHONY: migrate
migrate: ## run all new database migrations
	@echo "Running all new database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-down
migrate-down: ## revert database to the last migration step
	@echo "Reverting database to the last migrate step..."
	@$(MIGRATE) down 1

.PHONY: migrate-drop
migrate-drop: ## drop all database migrations
	@echo "dropping database..."
	@$(MIGRATE) drop -f

.PHONY: migrate-new
migrate-new: ## create a new database migration
	@read -p "Enter the name of the new migration: " name; \
	$(MIGRATE) create -ext sql -dir /migrations/ $${name}

.PHONY: migrate-reset
migrate-reset: ## reset database and re-run all migrations
	@echo "Resetting database..."
	@$(MIGRATE) drop -f
	@echo "Running all database migrations..."
	@$(MIGRATE) up

.PHONY: migrate-arg
migrate-arg: ## run migration command with argument ARG
	@echo "Running migration command with argument: $(ARG)"
	@$(MIGRATE) $(ARG)
