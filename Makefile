MODULE=github.com/pjover/espigol

.PHONY: build run tidy importar-socis

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
