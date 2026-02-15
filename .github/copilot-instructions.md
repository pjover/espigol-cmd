# AI Code Assistant Guide: espigol

This document provides guidance for the AI Code Assistant on how to effectively interact with the `espigol` repository.

## About the Project

Espigol is a Go project designed to manage partners and expense forecasts. It follows a Hexagonal Architecture (Ports and Adapters) pattern, using MongoDB for persistence and Cobra/Viper for configuration and CLI management.

## Key Technologies

- **Programming Language:** Go 1.23+
- **Dependency Management:** Go Modules
- **Database:** MongoDB (via `go.mongodb.org/mongo-driver`)
- **CLI Framework:** Cobra (`github.com/spf13/cobra`)
- **Configuration:** Viper (`github.com/spf13/viper`)
- **Containerization:** Docker, Docker Compose
- **Testing:** Native Go testing (`testing` package) with `testify` (assert/mock) if available.
- **Linting & Formatting:** `gofmt`

## Directory Structure

- `domain/`: Core business logic.
  - `domain/model/`: Pure domain entities (e.g., `Partner`, `ExpenseForecast`).
  - `domain/ports/`: Interfaces defining services and repositories (e.g., `ConfigService`, `PartnerRepository`).
  - `domain/services/`: Domain service implementations.
- `adapters/`: Infrastructure implementations.
  - `adapters/cfg/`: Configuration adapter (Viper).
  - `adapters/cli/`: CLI command manager and commands.
  - `adapters/mongodb/`: MongoDB repository implementations.
- `private/`: Private data, logs, and documentation.
- `Makefile`: Build and automation tasks.
- `docker-compose.yaml`: Infrastructure services (MongoDB, Mongo Express).
- `dependency_injection.go`: Application wiring.
- `main.go`: Application entry point.

## Commands

### Getting started

To set up the project for local development:

1.  **Install dependencies:**
    ```bash
    go mod download
    ```
2.  **Start infrastructure:**
    ```bash
    docker-compose up -d
    ```

### Building and Running

-   **Build the binary:**
    ```bash
    make build
    ```
-   **Run the application:**
    ```bash
    make run
    ```
    *(Or `./bin/espigol` after building)*

### Testing

-   **Run all tests:**
    Inside Visual studio code is better to simply run `Test: Run All Tests` 
    if not, the run `make test`
    or to test a specific package run `go test ./domain/model/...`
    or run tests with verbose output with `go test -v ./...`

### Formatting

-   **Format code:**
    ```bash
    make format
    ```

## Core Principles

### 1. Always Verify File State First
**Before starting any task**, check the current file system state. Don't trust cached file contents.

```bash
# Check recent changes
git status
git diff

# Verify file contents
cat path/to/file.go
```

**Why**: The user may have made manual edits, run formatters, or updated files since your last interaction.

### 2. Never Commit Changes
**CRITICAL**: Do NOT commit any changes to git. The user handles all commits.

```bash
# ❌ NEVER do this
git add .
git commit -m "..."
git push

# ✅ Only verify changes
git status
git diff
```

### 3. Add a commit message suggestion in the final summary

After finishing the task, add to the task summary a commit message suggestion following the **Conventional Commits** guidelines.

#### Commit Message Format
```
<type>[optional scope]: <description>

[optional body]
```

#### Types
- `fix`: Patches a bug.
- `feat`: Introduces a new feature.
- `build`: Changes to build system or dependencies.
- `chore`: Maintenance tasks (refactoring, cleanup).
- `ci`: CI/CD changes.
- `docs`: Documentation updates.
- `test`: Adding or missing tests.
- `refactor`: Code change that neither fixes a bug nor adds a feature.


### 4. Always Run Quality Checks

- After making changes, **always** ensure code is formatted and passes tests.
- Fix all issues before considering the task complete.

## Project-Specific Patterns

### Error Handling
- Use `fmt.Errorf` with `%w` to wrap errors: `fmt.Errorf("context: %w", err)`.
- Return errors explicitly; do not use `panic` unless absolutely necessary (e.g., startup failure in `main`).

### Configuration
- Use `adapters/cfg` via the `ports.ConfigService` interface.
- Do not access `os.Getenv` directly in domain logic.

### Dependency Injection
- Use `dependency_injection.go` to wire up components.
- Avoid global state; inject dependencies via constructors.

### Testing Strategy
- **Unit Tests**: Place `_test.go` files next to the source files.
- **Naming**: Test functions should start with `Test` (e.g., `TestPartnerCreation`).
- **Mocks**: 
    - Use interfaces (`domain/ports`) to mock dependencies in tests.
    - Use [testify](github.com/stretchr/testify) to mock adapters

## Development Workflow

1.  **Understand the Request**: Read carefully, check context.
2.  **Gather Context**: `git status`, `ls -R`, `cat` relevant files.
3.  **Plan**: Think through the Hexagonal Architecture implications.
4.  **Implement**: Make minimal, focused changes.
5.  **Verify**: `make format`, `go test ./...`.
6.  **Suggest Commit**: Propose a conventional commit message.

## Temporary Files
If you need to create temporary files, use `/tmp`.
**DO NOT CREATE** temporary files in the source directories.

## Summary Checklist
Before completing any task, verify:
- File state checked.
- Minimal changes made.
- Code formatted (`make format`).
- `make format` passes.
- No git commits made.
- Commit message proposed.
