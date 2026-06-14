# TODO

Issues identified during code review — non-blocking, to address in future iterations.

## Known limitations

- [ ] **`op time-entry create --activity` cannot list available activities** — `GET /api/v3/time_entries/activities` returns 404 in some OpenProject instances (version or permission issue). Investigate using the form endpoint `GET /api/v3/time_entries/form` which may include allowed activities in its schema.

## Project identifier refactoring

- [ ] **`validatedVersionId()` in `cmd/workpackage/list.go` swallows errors** — after `printer.Error(err)` the function continues with a nil `project`/`versions` instead of returning early. Should `return ""` immediately after the error print to avoid fragile control flow.

## CLI UX

- [ ] **Add `-w` short flag to `activities --work-package`** — `--work-package` has no short option, inconsistent with similar flags elsewhere. Add `-w` or document why it's intentionally omitted.

- [ ] **`--format` flag has no short option** — inconsistent with the convention of adding short flags for frequently used ones. Consider `-f` if it doesn't conflict.
