MODULE=github.com/pjover/espigol

.PHONY: build run tidy importar-socis

build:
	mkdir -p bin
	go build -o bin/espigol .

run: build
	go run . $(ARGS)

tidy:
	go mod tidy

importar-socis:
	$(eval CSVPATH=$(if $(CSV),$(CSV),private/CSV/socis.csv))
	go run . importar socis --csv=$(CSVPATH)
