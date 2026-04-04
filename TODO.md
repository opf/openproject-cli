# TODO

Issues identified during code review — non-blocking, to address in future iterations.

## Code quality

- [ ] **Extract uint64 argument validation helper** — `inspect` and `update` commands all repeat the same pattern (check `len(args) == 1`, parse uint64). Extract to `cmd/common.go` to reduce duplication.

- [ ] **Replace `fmt.Println/Printf` with `printer.*` in `root.go` and `login.go`** — Both files use `fmt` directly for terminal output, violating the convention that all output goes through `printer/`. `root.go` uses it for the version string, `login.go` for the token prompt and error messages.

## Known limitations

- [ ] **`op time-entry create --activity` cannot list available activities** — `GET /api/v3/time_entries/activities` returns 404 in some OpenProject instances (version or permission issue). Investigate using the form endpoint `GET /api/v3/time_entries/form` which may include allowed activities in its schema.

## CLI UX

- [ ] **Add `-w` short flag to `activities --work-package`** — `--work-package` has no short option, inconsistent with similar flags elsewhere. Add `-w` or document why it's intentionally omitted.

- [ ] **`--format` flag has no short option** — inconsistent with the convention of adding short flags for frequently used ones. Consider `-f` if it doesn't conflict.
