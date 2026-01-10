ifneq (,$(wildcard .env))
include .env
export
endif

COMPOSE_FILE=deploy/local/docker-compose.yml
ENV_FILE=.env
GOOSE_MIGRATION_DIR=./internal/workstream-tasks/migrations 

.PHONY: run stop logs build rebuild rpcgen goose_up goose_down

run:
	docker compose \
		-f $(COMPOSE_FILE) \
		--env-file $(ENV_FILE) \
		up -d

stop:
	docker compose \
		-f $(COMPOSE_FILE) \
		--env-file $(ENV_FILE) \
		down

logs:
	docker compose \
		-f $(COMPOSE_FILE) \
		--env-file $(ENV_FILE) \
		logs -f

build:
	docker compose \
		-f $(COMPOSE_FILE) \
		--env-file $(ENV_FILE) \
		build --no-cache

rebuild: stop
	@COMPOSE_BAKE=true
	docker compose \
		-f $(COMPOSE_FILE) \
		--env-file $(ENV_FILE) \
		up -d --build --force-recreate

rpcgen:
	protoc \
		-I api/proto \
		--go_out=internal/gen \
		--go_opt=paths=source_relative \
		--go-grpc_out=internal/gen \
		--go-grpc_opt=paths=source_relative \
		api/proto/tasks/v1/tasks.proto

goose_up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(POSTGRES_DSN) \
	goose -dir=$(GOOSE_MIGRATION_DIR) up

goose_down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(POSTGRES_DSN) \
	goose -dir=$(GOOSE_MIGRATION_DIR) down
