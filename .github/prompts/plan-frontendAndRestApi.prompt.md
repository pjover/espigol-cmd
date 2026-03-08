## Plan: REST API Architecture

This is a comprehensive plan to add a REST API server to `espigol`, adhering to your Hexagonal Architecture and testing requirements.

**Phases Overview**
- **Phase 1**: Expand the CLI with server lifecycle commands (`start`, `stop`, `status`).
- **Phase 2**: Build the REST API adapter (HTTP layer), applying Google Auth and Swagger documentation for `Partner` and `ExpenseForecast` entities.

***

### Phase 1: Add the Server Commands to Espigol
*We will use PID files (e.g., in `/tmp` or `.espigol/`) to track the server daemon's lifecycle across different CLI executions.*

**Steps**
1. **Create the Base Server Command**
   - Create `internal/adapters/cli/server/server.go` to define the parent `server` command.
   - Write unit tests mocking `ports.CommandManager` to ensure the command registers properly.
   - *Commit:* `feat(cli): add server command`
2. **Implement Server Start Command**
   - Create `internal/adapters/cli/server/start.go` for `espigol server start`.
   - Implement the logic to write a PID file and block/run the server (the actual HTTP server will be dummy in this step).
   - Write unit tests mocking standard OS operations (if needed) and checking CLI flags.
   - *Commit:* `feat(cli): implement server start command and PID tracking`
3. **Implement Server Status & Stop Commands**
   - Create `internal/adapters/cli/server/status.go` to read the PID file and check if the process is alive.
   - Create `internal/adapters/cli/server/stop.go` to send a termination signal (`SIGTERM`) to the PID.
   - Write tests simulating active and missing PID files.
   - *Commit:* `feat(cli): implement server stop and status commands`
4. **Wire Dependencies**
   - Update `internal/dependency_injection.go` to inject the new server commands into the root `CommandManager`.
   - *Commit:* `refactor(di): inject server commands into CLI manager`

### Phase 2: Implement the REST API Server (Go)
*We will use Go 1.22's enhanced `net/http` standard multiplexer to keep dependencies lean, and `swaggo/swag` for OpenAPI generation.*

**Steps**
1. **Define HTTP Server Adapter and Port**
   - Configure the server's port at `configs/espigol.yaml` file and load it dinamically
   - Define a `Server` port in `internal/domain/ports/server.go`.
   - Create the HTTP adapter backbone in `internal/adapters/http/server.go` that initializes a `net/http` `ServerMux`.
   - Wire this into the `start` command from Phase 1.
   - *Commit:* `feat(http): initialize net/http multiplexer adapter`
2. **Implement CRUD Handlers for Partners**
   - Create `internal/adapters/http/partner_handler.go` defining `GET`, `POST`, `PUT`, `DELETE` routes for `Partners`.
   - Use `testify/mock` to map requests to existing domain layer services. Provide exhaustive handler anonymous unit tests (checking correct status codes like 200, 201, 400, 404).
   - Add Swaggo `@Summary` comments to all methods.
   - *Commit:* `feat(api): add REST endpoints for Partner entity`
3. **Implement CRUD Handlers for Expense Forecasts**
   - Create `internal/adapters/http/expense_forecast_handler.go`.
   - Implement handlers, mock testing setup, and Swaggo annotations just like Partners.
   - *Commit:* `feat(api): add REST endpoints for ExpenseForecast entity`
4. **Generate API Documentation**
   - Run `swag init` to generate the `docs/` folder. Apply an HTTP handler (`http-swagger`) to serve the UI at `/swagger/index.html`.
   - *Commit:* `docs: generate and serve OpenAPI swagger docs`

