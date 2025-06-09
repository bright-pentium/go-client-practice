GOBIN ?= $$(go env GOBIN)
PG_ENV_PATH ?= ./configs/pg.env
APP_ENV_PATH ?= ./configs/app.env
CONTAINER_NAME ?= my-echo-app
PORT := $(shell grep -E '^PORT=' $(APP_ENV_PATH) | cut -d '=' -f2)

.PHONY: db-up
db-up:
	docker compose -f deployments/docker-compose.yml --env-file $(PG_ENV_PATH) up

.PHONY: db-down
db-down:
	docker compose -f deployments/docker-compose.yml down --volumes

.PHONY: migrate-create
migrate-create:
	@if [ -z "$(MSG)" ]; then echo "Error: MSG is required."; exit 1; fi
	$(GOBIN)/migrate create -ext sql -dir internal/migrations -seq "$(MSG)"

.PHONY: swag-gen
swag-gen:
	$(GOBIN)/swag init -g server.go --dir internal/delivery/echo/ --output docs/swaggo --parseDependency

.PHONY: migrate-up
migrate-up:
	@DATABASE_URL=$$(go run cmd/herlper/pgUrlparser.go $(PG_ENV_PATH)); \
	if [ -z "$$DATABASE_URL" ]; then \
		echo "Error: DATABASE_URL is required."; exit 1; \
	fi; \
	$(GOBIN)/migrate -database "$$DATABASE_URL" -path internal/migrations up

.PHONY: app-build
app-build:
	docker build -t $(CONTAINER_NAME) -f deployments/Dockerfile .

.PHONY: app-up
app-up:
	docker run --name $(CONTAINER_NAME) \
		-p $(PORT):$(PORT) \
		-v $(APP_ENV_PATH):/app/configs/app.env \
		-v $(PG_ENV_PATH):/app/configs/pg.env \
		--rm

.PHONY: mock-gen
mock-gen:
	$(GOBIN)/mockery

.PHONY: setup-gen
setup-gen: mock-gen swag-gen

.PHONY: setup-db
setup-db: db-up migrate-up

.PHONY: run
run: docker-build app-up
