# Espígol

Espígol és una eina per gestionar els socis i les previsions de despeses de la cooperativa. Ofereix una interfície de línia de comandes (CLI) per importar dades i una API REST per consultar i modificar els recursos.

---

## Resum de la implementació

### Arquitectura

El projecte segueix l'**Arquitectura Hexagonal** (Ports i Adaptadors):

- `internal/domain/model/` — Entitats de domini (`Partner`, `ExpenseForecast`).
- `internal/domain/ports/` — Interfícies (`DbService`, `ConfigService`, `Server`, `CommandManager`).
- `internal/domain/services/` — Serveis de domini (importadors CSV).
- `internal/adapters/cfg/` — Adaptador de configuració (Viper).
- `internal/adapters/cli/` — Adaptador CLI (Cobra): comandes `import` i `server`.
- `internal/adapters/http/` — Adaptador HTTP (`net/http`): handlers REST per a `Partner` i `ExpenseForecast`.
- `internal/adapters/mongodb/` — Adaptador de persistència (MongoDB).
- `internal/dependency_injection.go` — Cablejat de dependències.

### Funcionalitats implementades

| Àmbit                          | Descripció                                                                             |
| ------------------------------ | -------------------------------------------------------------------------------------- |
| **Importació CSV**             | Importa socis i previsions de despeses des de CSV via CLI o Makefile                   |
| **API REST – Partners**        | `GET`, `POST`, `PUT`, `DELETE` a `/partners` i `/partners/{id}`                        |
| **API REST – Previsions**      | `GET`, `POST`, `PUT`, `DELETE` a `/expense-forecasts` i `/expense-forecasts/{id}`      |
| **Health check**               | `GET /health` — retorna `200 OK`                                                       |
| **Swagger UI**                 | Documentació interactiva a `/swagger/index.html`                                       |
| **Servidor HTTP**              | Inici i aturada graciosa via senyals UNIX (`SIGTERM`/`SIGINT`)                         |
| **Cicle de vida del servidor** | Comandes CLI `server start`, `server stop`, `server status` (seguiment via fitxer PID) |
| **Persistència**               | MongoDB via `go.mongodb.org/mongo-driver`                                              |
| **Categoria de despesa**       | Classificació automàtica de les previsions en despesa corrent o d'inversió             |
| **Límits de subvenció**        | Límits anuals de subvenció llegits des de `configs/espigol.yaml`                       |

---

## Requisits previs

- Go 1.23+
- Docker i Docker Compose
- `swag` CLI (per regenerar docs): `go install github.com/swaggo/swag/cmd/swag@latest`

---

## Posada en marxa

```bash
# 1. Descarregar dependències
go mod download

# 2. Aixecar la infraestructura (MongoDB + Mongo Express)
make up

# 3. Inicialitzar la base de dades (índexs i col·leccions)
make init-db

# 4. Compilar
make build
```

---

## Comandes CLI

### Estructura

```
espigol
├── import                  Importar dades des de fitxers CSV
│   ├── partners            Importar socis
│   └── expense-forecasts   Importar previsions de despeses
└── server                  Gestionar el servidor REST
    ├── start               Iniciar el servidor
    ├── stop                Aturar el servidor
    └── status              Consultar l'estat del servidor
```

### Importació de dades

```bash
# Importar socis (per defecte: private/CSV/partners.csv)
make import-partners

# Importar socis des d'un fitxer concret
make import-partners CSV=~/Downloads/partners.csv

# Importar previsions de despeses comunes
make import-expense-forecasts-common

# Importar previsions de despeses per soci
make import-expense-forecasts-partners

# O directament amb la CLI
./bin/espigol import partners --file=ruta/al/fitxer.csv
./bin/espigol import expense-forecasts --file=ruta/al/fitxer.csv
```

### Gestió del servidor REST

```bash
# Iniciar el servidor (en segon pla)
make server-start
# O: ./bin/espigol server start

# Comprovar si el servidor és actiu
make server-status
# O: ./bin/espigol server status

# Aturar el servidor
make server-stop
# O: ./bin/espigol server stop
```

---

## URLs del servei

| Recurs            | URL                                      | Descripció                        |
| ----------------- | ---------------------------------------- | --------------------------------- |
| **Health**        | http://localhost:8080/health             | Comprova que el servidor és actiu |
| **Swagger UI**    | http://localhost:8080/swagger/index.html | Documentació interactiva de l'API |
| **Swagger JSON**  | http://localhost:8080/swagger/doc.json   | Especificació OpenAPI 2.0         |
| **Mongo Express** | http://localhost:8081/db/espigol         | Interfície web de MongoDB         |

> Els ports es configuren a `configs/espigol.yaml`.

---

## API REST

### Partners

| Mètode   | Ruta             | Descripció                  |
| -------- | ---------------- | --------------------------- |
| `GET`    | `/partners`      | Llista tots els socis       |
| `GET`    | `/partners/{id}` | Obté un soci per ID         |
| `POST`   | `/partners`      | Crea un nou soci            |
| `PUT`    | `/partners/{id}` | Actualitza un soci existent |
| `DELETE` | `/partners/{id}` | Elimina un soci             |

### Previsions de despeses

| Mètode   | Ruta                      | Descripció                       |
| -------- | ------------------------- | -------------------------------- |
| `GET`    | `/expense-forecasts`      | Llista totes les previsions      |
| `GET`    | `/expense-forecasts/{id}` | Obté una previsió per ID         |
| `POST`   | `/expense-forecasts`      | Crea una nova previsió           |
| `PUT`    | `/expense-forecasts/{id}` | Actualitza una previsió existent |
| `DELETE` | `/expense-forecasts/{id}` | Elimina una previsió             |

Cada previsió exposa els camps addicionals `year` (any derivat de la data prevista) i `expenseCategory` (`Despesa corrent` o `Despesa d'inversió`, derivat del subtipus de despesa).

---

## Targets del Makefile

```bash
make build                          # Compila el binari a bin/espigol
make run ARGS="..."                 # Compila i executa amb arguments
make test                           # Executa tots els tests
make format                         # Formata el codi (gofmt)
make tidy                           # Neteja les dependències (go mod tidy)
make swag-init                      # Regenera la documentació OpenAPI (docs/)
make up                             # Aixeca la infraestructura Docker
make down                           # Atura la infraestructura Docker
make init-db                        # Inicialitza la BD (índexs i col·leccions)
make import-partners [CSV=path]     # Importa socis des de CSV
make import-expense-forecasts-common [CSV=path]   # Importa despeses comunes
make import-expense-forecasts-partners [CSV=path] # Importa despeses per soci
make server-start                   # Inicia el servidor REST
make server-stop                    # Atura el servidor REST
make server-status                  # Comprova l'estat del servidor REST
```

---

## Configuració

El fitxer de configuració principal és `configs/espigol.yaml`:

```yaml
db:
  name: espigol
  server: mongodb://localhost:27017
expenses:
  limits:
    "2026":
      current: 30000
      investment: 70000
server:
  port: 8080
urls:
  mongoexpress: http://localhost:8081/db/espigol
```

Els límits de subvenció (`expenses.limits`) s'organitzen per any. Per a cada any es defineixen el màxim de **despesa corrent** (`current`) i de **despesa d'inversió** (`investment`). El sistema els llegeix dinàmicament amb `LimitsForYear(year, config)`.

---

## Regenerar la documentació de l'API

Quan s'afegeixen o modifiquen anotacions Swaggo als handlers, cal regenerar els fitxers `docs/`:

```bash
make swag-init
```

Això actualitza `docs/swagger.json`, `docs/swagger.yaml` i `docs/docs.go`.

