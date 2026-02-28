MODULE=github.com/pjover/espigol

.PHONY: format build run tidy import-partners import-expense-forecasts-common import-expense-forecasts-partners up down test init-db

format:
	go fmt ./...

build: format
	mkdir -p bin
	go build -o bin/espigol ./cmd/espigol

run: build
	go run ./cmd/espigol $(ARGS)

tidy:
	go mod tidy

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

init-db:
	@docker exec espigol-mongo_server-1 mongosh --quiet --eval ' \
		if (db.getMongo().getDBNames().indexOf("espigol") < 0) { \
			db.getSiblingDB("espigol").createCollection("_init"); \
			print("Database espigol created"); \
		} else { \
			print("Database espigol already exists"); \
		} \
	'

test:
	go test ./...