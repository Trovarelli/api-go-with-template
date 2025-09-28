## Run tests locally (app runs outside Docker)
test:
	go test ./src/internal/config -v

## Lint using dockerized golangci-lint
lint:
	docker run --rm -v ${PWD}:/app -w /app golangci/golangci-lint:v2.3.0 golangci-lint run

## Start only the Postgres database via Docker Compose
db-up:
	docker compose up -d db

## Stop containers (DB) and remove volumes (danger: wipes data)
db-down:
	docker compose down -v

## Recreate DB from init.sql (drops data and volume)
db-reseed:
	docker compose down -v && docker compose up -d db

## Run the Go app locally
run-local:
	go run ./src

## Simple CI alias
ci: lint test
