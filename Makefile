MODULE=github.com/pjover/espigol

.PHONY: build run tidy import-partners up down test

build:
	mkdir -p bin
	go build -o bin/espigol .

run: build
	go run . $(ARGS)

tidy:
	go mod tidy

import-partners:
	$(eval CSVPATH=$(if $(CSV),$(CSV),private/CSV/partners.csv))
	go run . import partners --csv=$(CSVPATH)

up:
	docker-compose --env-file .env -f docker-compose.yaml up -d

down:
	docker-compose -f docker-compose.yaml down

test:
	go test ./...