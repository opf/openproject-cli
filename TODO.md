# TODO

Issues identified during code review — non-blocking, to address in future iterations.

## Known limitations

- [ ] **`op time-entry create --activity` cannot list available activities** — `GET /api/v3/time_entries/activities` returns 404 in some OpenProject instances (version or permission issue). Investigate using the form endpoint `GET /api/v3/time_entries/form` which may include allowed activities in its schema.

## CLI UX

- [ ] **`--format` flag has no short option** — inconsistent with the convention of adding short flags for frequently used ones. Consider `-f` if it doesn't conflict.
