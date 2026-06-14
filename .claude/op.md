# OpenProject CLI (`op`)

`op` is a CLI for interacting with OpenProject project management instances. Use it proactively whenever the user discusses work packages, projects, time entries, notifications, or anything related to OpenProject — just as naturally as you'd use `git` for version control or `ls` to explore the filesystem.

## Global flags

These flags apply to every command:

```
--format string    Output format: text (default) or json
--profile string   Profile name (overrides OP_CLI_PROFILE env var, default "default")
--verbose          Print verbose output
```

## Orientation

```bash
op whoami                          # Show current user and server (verify connection)
op whoami --profile <name>         # Same for a specific profile
op project list                    # Discover available projects (get IDs for other commands)
```

## Work packages

IDs accept either a numeric ID (`42`) or a project-based identifier (`PROJ-123`) everywhere.

```bash
op work-package list                              # All visible work packages
op work-package list -p <project-id-or-slug>      # Filter by project (numeric ID or identifier/slug)
op work-package list -s open                      # Filter: open / closed / <id> / comma-separated IDs
op work-package list -s '!<id>'                   # Exclude a status (prefix with !)
op work-package list -a me                        # Filter by assignee (me or user ID)
op work-package list -t <type-id>                 # Filter by type (comma-separated IDs, ! prefix to exclude)
op work-package list -v <version-id>              # Filter by version
op work-package list --not-version <id>           # Exclude a version
op work-package list --parent-id <wp-id>          # Direct children of a work package
op work-package list --include-sub-projects       # Include sub-project work packages
op work-package list --sub-project <id>           # Limit sub-projects to include (with --include-sub-projects)
op work-package list --not-sub-project <id>       # Exclude sub-projects (with --include-sub-projects)
op work-package list --timestamp 2025-01-01       # Work packages as of a date (YYYY-MM-DD or YYYY-MM-DDTHH:MM:SSZ)
op work-package list --total                      # Show only the total count

op work-package search <query>...                 # Search by subject, type, status, project name, or identifier
op work-package search <query>... -p <project>    # Limit search to a project; multiple words are ANDed; up to 100 results

op work-package inspect <id>                      # Full details of one work package
op work-package inspect <id> --types              # Also list available types on the work package
op work-package inspect <id> --open               # Open in default browser

op work-package create "Subject" -p <id>          # Create (project accepts numeric ID or slug)
op work-package create "Subject" -p <id> \
  --type <type> --assignee <uid> \
  --description "markdown text" --open            # Full create with all options

op work-package update <id> --subject "…"         # Update subject
op work-package update <id> --type <type>         # Change type
op work-package update <id> --assignee <uid>      # Change assignee
op work-package update <id> --description "…"     # Update description (markdown)
op work-package update <id> --action <action>     # Execute a custom action (e.g. status transition)
op work-package update <id> --attach <filepath>   # Attach a file
```

## Time entries

```bash
op time-entry list                                # Time entries for current user
op time-entry list -u <user-id>                   # Time entries for a specific user

op time-entry create \
  -w <wp-id> \
  --hours 1.5 \
  --activity "Development" \
  --spent-on 2025-01-15 \
  --comment "Fixed the bug" \
  -u <user-id>                                    # All flags; --spent-on defaults to today
```

## Budgets

```bash
op budget list -p <project-id-or-slug>            # Budgets for a project
op budget inspect <id>                            # Full details of a budget
```

## Projects

```bash
op project list                                   # All visible projects
op project inspect <id-or-slug>                   # Project details (numeric ID or identifier/slug)
op project inspect <id-or-slug> --open            # Open in default browser
```

## Other resources

```bash
op notification list                              # Unread notifications
op notification list -r <reason>                  # Filter by reason

op user search <query>                            # Find a user

op status list                                    # List work package statuses
op type list                                      # List work package types

op activity list                                  # All activities
op activity list --work-package <wp-id>           # Activities for a specific work package
```

## Git integration

```bash
op git start workpackage <id>   # Create a branch from WP id/subject/type and switch to it
```

## Profiles (multiple OpenProject instances)

The `--profile` flag selects which instance to talk to (stored in `~/.config/openproject/config`). The default profile is used when omitted. The env var `OP_CLI_PROFILE` also selects a profile.

```bash
op login                           # Authenticate (prompts for host, token, profile name)
op login --profile work            # Authenticate a named profile
op logout                          # Remove the default profile
op logout --profile work           # Remove a named profile
op whoami                          # Show all profiles with their server and user
```

## JSON output

Use `--format json` to get machine-readable output or raw markdown descriptions:

```bash
op work-package inspect <id> --format json
op work-package list --format json -p <project>
```

## Workflow for use-case analysis

1. `op whoami` — confirm which instance you're on
2. `op project list` — find relevant project IDs
3. `op work-package list -p <id> -s open` — get the work items
4. `op work-package inspect --format json <id>` — drill into specifics (raw markdown)
5. `op work-package search <terms> -p <id>` — find by subject or identifier fragment
