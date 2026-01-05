COMPOSE_FILE=deploy/local/docker-compose.yml
ENV_FILE=.env

.PHONY: run stop logs build rebuild rpcgen

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
