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
  configuration/            # Multi-profile config (INI), CLI version
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
- `WorkPackageDto` includes a `displayId` field (camelCase, from the API) mapped to `DisplayId` on the model; always present — holds the semantic identifier (e.g. `PROJ-123`) when project-based identifiers are enabled, or the numeric id as a string otherwise

### Command conventions

- Commands follow `op <noun> <verb>` (noun-first), e.g. `op work-package list`, `op time-entry create`
- Each noun has its own package under `cmd/` (`workpackage`, `timeentry`, `project`, `user`, …) — package names use no hyphens even when the command name does (e.g. `work-package` → `package workpackage`)
- Each noun package exposes a `RootCmd` registered in `cmd/root.go`
- One file per verb within each noun package (e.g. `cmd/workpackage/list.go`, `cmd/workpackage/create.go`)
- Flags: always define long flag names; add short flags (`-p`, `-o`) for frequently used ones
- Flags that resolve a resource (e.g. `--type`, `--assignee`) perform an API lookup and store the resolved link in the DTO

### Resource conventions

- Each resource has its own package under `components/resources/`
- Operation types use `iota` enums: `CreateOption`, `UpdateOption`, `FilterOption`
- Operations are dispatched via a `map[Option]func(...)` — add new options by extending the map
- Public API: `Create(...)`, `Lookup(id)`, `All(filters, query, ...)`, `Update(id, options)`

### Configuration conventions

- Config stored as INI at `~/.config/openproject/config` (or `$XDG_CONFIG_HOME/openproject/config`)
- Each profile is an INI section: `[name]` with `host` and `token` keys
- Profile names: letters, digits, `-`, `_` only; no leading/trailing hyphens; validated by `ValidateProfileName`, sanitized by `SanitizeProfileName`
- Key constants: `DefaultProfile = "default"`, `EnvProfile = "OP_CLI_PROFILE"`
- Key functions: `ReadConfig(profile)`, `WriteConfigForProfile(profile, host, token)`, `DeleteProfile(profile)`, `AllProfiles()`
- `OP_CLI_HOST` / `OP_CLI_TOKEN` env vars override all profiles; `OP_CLI_PROFILE` selects a profile (overridden by `--profile` flag)
- Old single-line format (`host token`) is auto-migrated to `[default]` on first read

### Paths conventions

- All API paths are defined in `components/paths/paths.go`
- Functions are named after the resource: `WorkPackage(id)`, `WorkPackages()`, `ProjectWorkPackages(projectId)`
- All paths are relative (no host), starting with `/api/v3`
- Project path functions (`Project`, `ProjectWorkPackages`, `ProjectVersions`, `ProjectBudgets`) take a `string` that may be either a numeric ID (`"42"`) or a human-readable identifier (`"my-project"`); the OpenProject API accepts both forms at the same endpoints
- `WorkPackage(id)` and `WorkPackageActivities(id)` take a `string` that may be either a numeric ID (`"12345"`) or a project-based semantic identifier (`"PROJ-123"`); validated by `work_packages.ValidateIdentifier`

### Printer conventions

- All terminal output goes through `printer/` — never `fmt.Println` directly in commands
- Color scheme: Red = ID, Green = type, Cyan = subject/name, Yellow = status
- `printer.Error(err)` for errors, `printer.ErrorText(msg)` for plain error strings
- `printer.Info(msg)` for progress messages, `printer.Done()` after successful mutations
- Work packages display `DisplayId` from the API: semantic form (e.g. `PROJ-123`) when project-based identifiers are enabled, `#N` for numeric-only systems (where `displayId` equals the numeric id)
- Work package browser URLs use the short `wp/<displayId>` form (e.g. `wp/PROJ-123` or `wp/42`)

## Testing conventions

- Test files use external test packages: `package printer_test`, `package requests_test`
- `TestMain` in `printer_test` initializes shared state (routes, printer) for the package
- Tests use plain `t.Errorf` — no test framework, no assertions library
- Tests only exist for `printer`, `requests`, `common`, and `configuration` — no tests on `cmd/` or `resources/`
- When adding a new printer function, add a corresponding test in `components/printer/`
- When adding a new configuration function, add a corresponding test in `components/configuration/`

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `github.com/fatih/color` | Terminal colors (via printer) |
| `github.com/briandowns/spinner` | Progress spinner |
| `github.com/go-git/go-git/v5` | Git integration (`op git` commands) |
| `github.com/sosodev/duration` | ISO 8601 duration parsing |
