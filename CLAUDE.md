# OpenProject CLI — Developer Guide

## Build & Run

```bash
# Build
go build -o op .

# Run tests
go test ./...

# Run a specific test package
go test ./components/printer/...

# Install locally
go install .
```

## Architecture

The codebase is organized in strict layers. Each layer only depends on layers below it.

```
cmd/                        # Cobra commands — CLI entry points only, no business logic
components/
  resources/                # Business logic — API calls, option handling
  paths/                    # API path constants (paths.go — single file)
  requests/                 # HTTP client (GET, POST, PATCH)
  parser/                   # JSON response parsing
  printer/                  # Terminal output formatting
  routes/                   # Browser URL generation
  common/                   # Shared utilities (string, slice, math)
  configuration/            # Config file reading, CLI version
  launch/                   # Browser launcher
models/                     # Domain models (plain structs, no logic)
dtos/                       # JSON DTOs with Convert() to models
```

### Data flow

API response → `parser.Parse[SomethingDto]()` → `dto.Convert()` → `models.Something` → `printer.Something()`

### DTO conventions

- DTOs live in `dtos/`, named `<Resource>Dto`
- DTOs mirror the OpenProject API v3 HAL JSON structure
- Links use `*LinkDto` with `Href` and `Title` fields, serialized as `_links`
- Every DTO implements `Convert() *models.Something` to produce a domain model
- Collections follow the pattern: `<Resource>CollectionDto` with `Embedded.<Resource>Elements`
- `omitempty` on all JSON fields in DTOs used for POST/PATCH bodies

### Command conventions

- Commands follow `op <verb> <resource>` (verb-first)
- Each verb has its own package under `cmd/` (`create`, `list`, `update`, `inspect`, `search`)
- Each verb package exposes a `RootCmd` registered in `cmd/root.go`
- One file per resource within each verb package (e.g. `cmd/list/work_packages.go`)
- Flags: always define long flag names; add short flags (`-p`, `-o`) for frequently used ones
- Flags that resolve a resource (e.g. `--type`, `--assignee`) perform an API lookup and store the resolved link in the DTO

### Resource conventions

- Each resource has its own package under `components/resources/`
- Operation types use `iota` enums: `CreateOption`, `UpdateOption`, `FilterOption`
- Operations are dispatched via a `map[Option]func(...)` — add new options by extending the map
- Public API: `Create(...)`, `Lookup(id)`, `All(filters, query, ...)`, `Update(id, options)`

### Paths conventions

- All API paths are defined in `components/paths/paths.go`
- Functions are named after the resource: `WorkPackage(id)`, `WorkPackages()`, `ProjectWorkPackages(projectId)`
- All paths are relative (no host), starting with `/api/v3`

### Printer conventions

- All terminal output goes through `printer/` — never `fmt.Println` directly in commands
- Color scheme: Red = ID, Green = type, Cyan = subject/name, Yellow = status
- `printer.Error(err)` for errors, `printer.ErrorText(msg)` for plain error strings
- `printer.Info(msg)` for progress messages, `printer.Done()` after successful mutations

## Testing conventions

- Test files use external test packages: `package printer_test`, `package requests_test`
- `TestMain` in `printer_test` initializes shared state (routes, printer) for the package
- Tests use plain `t.Errorf` — no test framework, no assertions library
- Tests only exist for `printer`, `requests`, and `common` — no tests on `cmd/` or `resources/`
- When adding a new printer function, add a corresponding test in `components/printer/`

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/fatih/color` | Terminal colors (via printer) |
| `github.com/briandowns/spinner` | Progress spinner |
| `github.com/go-git/go-git/v5` | Git integration (`op git` commands) |
| `github.com/sosodev/duration` | ISO 8601 duration parsing |
