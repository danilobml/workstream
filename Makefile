COMPOSE_FILE=deploy/local/docker-compose.yml
ENV_FILE=.env

.PHONY: run stop logs build rebuild

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

rebuild:
	docker compose \
		-f $(COMPOSE_FILE) \
		--env-file $(ENV_FILE) \
		up -d --build --force-recreate
