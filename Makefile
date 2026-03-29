MODULE=github.com/pjover/espigol

.PHONY: format build run tidy swag-init import-partners import-expense-forecasts-common import-expense-forecasts-partners up down test init-db server-start server-stop server-status generate-expense-forecast-report

format:
	go fmt ./...

build: format
	mkdir -p bin
	go build -o bin/espigol ./cmd/espigol

run: build
	go run ./cmd/espigol $(ARGS)

tidy:
	go mod tidy

swag-init: ## Regenerate OpenAPI docs from Swaggo annotations (output: docs/)
	swag init -g cmd/espigol/main.go --parseDependency --output docs

import-partners:
	$(eval CSVPATH=$(if $(CSV),$(CSV),private/CSV/partners.csv))
	go run ./cmd/espigol import partners --file=$(CSVPATH)

import-expense-forecasts-common:
	$(eval CSVPATH=$(if $(CSV),$(CSV),private/CSV/expense-forecasts-common.csv))
	go run ./cmd/espigol import expense-forecasts --file=$(CSVPATH)

import-expense-forecasts-partners:
	$(eval CSVPATH=$(if $(CSV),$(CSV),private/CSV/expense-forecasts-partners.csv))
	go run ./cmd/espigol import expense-forecasts --file=$(CSVPATH)

up:
	docker-compose --env-file .env -f docker-compose.yaml up -d

down:
	docker-compose -f docker-compose.yaml down

start: up
	go run ./cmd/espigol server start

stop:
	go run ./cmd/espigol server stop

status:
	go run ./cmd/espigol server status

init-db:
	@docker exec espigol-mongo_server-1 mongosh --quiet --eval ' \
		if (db.getMongo().getDBNames().indexOf("espigol") < 0) { \
			db.getSiblingDB("espigol").createCollection("_init"); \
			print("Database espigol created"); \
		} else { \
			print("Database espigol already exists"); \
		} \
		const espigolDb = db.getSiblingDB("espigol"); \
		espigolDb.partner.createIndex({ email: 1 }, { unique: true, sparse: true }); \
		print("Index on partner.email ensured"); \
	'

test:
	go test ./...

generate-expense-forecast-report: build ## Generate expense categories PDF report for YEAR (default: current year)
	./bin/espigol generate expense-forecast-report $(if $(YEAR),--year=$(YEAR),)