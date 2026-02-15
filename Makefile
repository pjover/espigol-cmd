MODULE=github.com/pjover/espigol

.PHONY: format build run tidy import-partners import-expense-forecasts-common import-expense-forecasts-partners up down test

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

test:
	go test ./...