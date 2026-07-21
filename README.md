# OpenProject CLI

[![CI](https://github.com/opf/openproject-cli/actions/workflows/ci.yml/badge.svg)](https://github.com/opf/openproject-cli/actions/workflows/ci.yml)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/opf/openproject-cli)

OpenProject CLI is a tool for operating your OpenProject instances from the command line. Therefore, it provides a
subset of use cases for OpenProject.

⚠️ **IMPORTANT**: This tool is currently in a pre-release phase. It is not meant to be full-featured, complete, nor will
it have official technical support by the OpenProject GmbH.

The OpenProject CLI is meant to be operated in a self-explanatory and easy-to-use way. It provides meaningful commands
without abbreviations, that can be executed in a human-readable manner. For more experienced users, many commands
provide aliases for fewer keystrokes.

## Setup

There are many ways to install the OpenProject CLI:

### Download executable

Executables are provided for Linux, macOS and Windows (currently Linux x64 only).

You can find the executables in the [**Releases**](https://github.com/opf/openproject-cli/releases/) section of this
repository. The latest release is available under
[https://github.com/opf/openproject-cli/releases/latest](https://github.com/opf/openproject-cli/releases/latest).

Download the zip archive for your operating system and desired version of the OpenProject CLI:

```shell
curl -O https://github.com/opf/openproject-cli/releases/download/0.2.0/openproject-cli_linux_x64_X.Y.Z.zip
```

Extract the executable from the zip archive and move the executable to a location that is within your `PATH`:

```shell
unzip openproject-cli_linux_x64_X.Y.Z.zip
sudo mv op /usr/local/bin
```

Check if the executable is working:

```shell
op --version

OpenProject CLI: X.Y.Z
        commit: 8dc3232
        built: Thu Jun 29 20:41:41 UTC 2023
        built with: go1.20.5
```

> MacOS will place the executable in quarantine by default. If your Mac displays an alert when you run the executable, you can dequarantine the `op` executable via the following command:

```shell
xattr -d com.apple.quarantine /usr/local/bin/op
```

### Go toolchain

If you already have the Go toolchain installed on your system, you can install the OpenProject CLI with the `go install`
command:

```shell
go install github.com/opf/openproject-cli
```

This will install the most current commit development version for now. Keep in mind that the resulting command is also
not `op` but `openproject-cli`. You can though rename the binary installed in your `GOPATH` or `GOBIN` directory.

## Completion

The OpenProject CLI supports auto-completion out of the box. You can generate the completion script by running:

```shell
op completion [shell]
```

Replace `[shell]` with either `bash`, `fish`, `powershell`, `zsh` to generate the shell specific completion script.

You need to put the completion in a place where your shellcan find it. For ZSH this could look like:

```shell
op completion zsh > ~/.oh-my-zsh/completions/_op
rm ~/.zcompdump*
```

Open a new shell and completion should work by using the `TAB` key as usual.

## Usage

The OpenProject CLI commands are structured in a common, human-readable pattern. Every command is built
as `op NOUN VERB [additional information]`. You will see plenty of examples within this section.

### Breaking change: noun-first commands

Earlier releases (up to 0.5.5) used verb-first commands. These have been replaced without aliases;
scripts need to migrate:

| Old (verb-first)          | New (noun-first)                     |
|---------------------------|--------------------------------------|
| `op list workpackages`    | `op work-package list`               |
| `op create workpackage`   | `op work-package create`             |
| `op update workpackage`   | `op work-package update`             |
| `op inspect workpackage`  | `op work-package inspect`            |
| `op search workpackages`  | `op work-package search`             |
| `op list projects`        | `op project list`                    |
| `op inspect project 42`   | `op project inspect 42`              |
| `op list notifications`   | `op notification list`               |
| `op list activities 42`   | `op activity list --work-package 42` |
| `op list status`          | `op status list`                     |
| `op list types`           | `op type list`                       |
| `op list timeentries`     | `op time-entry list`                 |
| `op search user X`        | `op user search X`                   |

### Discoverability

Discoverability is key. As we won't document every single command within this README, it is important for the CLI tool,
to provide helpful information at any time and any level of the command structure. To achieve that, the OpenProject CLI
has a help page on every command. Adding the `--help/-h` flag to any chain of commands will tell the user of all
possibilities he can use from here.

```shell
# General help
op -h

# Help about work package commands
op work-package -h

# Help about how can I narrow down my list of work packages
op work-package list -h
```

The second basement is the autocompletion. To set it up correctly, follow the steps explained in the `Completion`
section above.

Once it is working, you can use it to discover possible commands and flags while typing.

```shell
# Chaining commands: hitting completion key after
op work-package
# returns
create  -- Create work package in project
inspect -- Inspect a work package
list    -- Lists work packages
search  -- Searches for work packages
update  -- Update a work package

# Discover flags: hitting completion key after
op work-package update 42 -
# returns
--action    -a  -- Executes a custom action on a work package
--assignee      -- Assign a user to the work package
--attach        -- Attach a file to the work package
--help      -h  -- help for work-package
--subject       -- Change the subject of the work package
--type      -t  -- Change the work package type
```

### Authentication

#### Single instance (default)

Run `op login` to authenticate. You will be prompted for the host URL, an API
token, and a profile name (defaults to `default`):

```shell
op login
# Profile name? [default]
# OpenProject host URL: https://community.openproject.org
# OpenProject API Token: ...
```

#### Multiple instances

Use named profiles to work with more than one OpenProject instance:

```shell
# Create a profile named "work"
op login --profile work

# Create a profile named "staging"
op login --profile staging
```

All commands accept a `--profile` flag to select which instance to use:

```shell
op work-package list --profile work
op project list --profile staging
```

The `OP_CLI_PROFILE` environment variable is an alternative to `--profile`,
useful in scripts or shell sessions where you always want the same instance:

```shell
export OP_CLI_PROFILE=work
op work-package list   # uses the "work" profile
```

When both are set, `--profile` takes precedence. The explicit credential
variables `OP_CLI_HOST` and `OP_CLI_TOKEN` override everything.

#### Inspecting and removing profiles

`op whoami` shows all configured profiles together with their server URL and
authenticated user:

```shell
op whoami
# Profile: default
# Server:  https://community.openproject.org
# User:    #42 Jane Doe
#
# Profile: work
# Server:  https://work.example.com
# User:    #7 John Smith
```

Pass `--profile` to inspect a single profile:

```shell
op whoami --profile work
```

`op logout` removes a profile (defaults to `default`):

```shell
op logout              # removes the "default" profile
op logout --profile work
```

#### Profile names

Profile names may only contain letters, digits, `-` and `_`; hyphens must not
be consecutive, leading, or trailing. When prompted interactively the CLI
suggests a sanitized version of any invalid input.

### Prominent examples

There are a couple of use cases, you might want to execute from the command line. In this section we provide a handful
of examples, that might be useful for a great number of people.

#### Creation

```shell
# Creating a work package in a project only by subject.
# Work package is created with many default values (as for type and status),
# very similar to how a work package is created inline in a work package table.
# --project accepts either a numeric ID or the project identifier from the URL.
op work-package create --project 11 'Document new CLI tool'
op work-package create --project my-project 'Document new CLI tool'

# Same command with shorthands and directly open it in a browser to continue working on it.
op work-package create -p11 'Document new CLI tool' -o
```

#### Listing

```shell
# Get a list of unread notifications and filter them by reason
op notification list --reason mentioned

# Get a list of all work packages assigned to me
op work-package list --assignee me

# Get a list of direct children of a work package
# --parent-id accepts either a numeric ID or a project-based identifier (e.g. PROJ-123)
op work-package list --parent-id 42
op work-package list --parent-id PROJ-123
```

#### Updating

```shell
# Executing a custom action on a work package
# The ID argument accepts either a numeric ID or a project-based identifier (e.g. PROJ-123)
op work-package update 42 --action Claim
op work-package update PROJ-123 --action Claim

# Batch updating some properties of a work package
op work-package update 42 --subject 'The new subject' --type Implementation

# Local inputs are validated before the first change is sent. Each API action
# remains a separate request, so an earlier change may remain if a later
# request fails.

# Uploading an attachment to a work package
op work-package update 42 --attach ./Downloads/Report.pdf
```

#### Searching

```shell
# Search work packages by subject, type, status, project name, or identifier.
# Multiple words are ANDed: all terms must match. Returns up to 100 results.
op work-package search cascade
op work-package search fix login

# Limit search to a specific project (numeric ID or identifier)
op work-package search cascade --project my-project
op work-package search cascade -p 11
```

#### Inspecting

```shell
# Inspecting a work package with more details,
# then in the work package list command
# Accepts either a numeric ID or a project-based identifier (e.g. PROJ-123)
op work-package inspect 42
op work-package inspect PROJ-123
```

## AI agent integration (Claude Code)

The repository ships a [`.claude/op.md`](.claude/op.md) reference file that teaches Claude Code how to use `op`. Once set up, the agent will use `op` proactively whenever you discuss work packages, projects, or anything OpenProject-related.

`~/.claude/CLAUDE.md` is Claude Code's global instruction file — create it if it doesn't exist yet.

### Option A — available across all your projects (recommended)

Symlink the file into your Claude configuration directory from the repo root:

```shell
ln -s $(pwd)/.claude/op.md ~/.claude/op.md
```

Then add a line to `~/.claude/CLAUDE.md`:

```markdown
Use `op.md` as a reference for all `op` CLI commands (work packages, projects, time entries, search, profiles…). Read it whenever the user asks about OpenProject or wants to use `op`.
```

The symlink keeps the file in sync automatically as you pull new versions.

### Option B — available in this project only

Add a line to the project's `CLAUDE.md` at the repo root:

```markdown
Use `.claude/op.md` as a reference for all `op` CLI commands (work packages, projects, time entries, search, profiles…). Read it whenever the user asks about OpenProject or wants to use `op`.
```

No copy or symlink needed — the path is relative to the project root.

## Creating a release

Releases are triggered by pushing a git tag. The tag name becomes the version string embedded in the binaries.

```shell
git tag vX.Y.Z
git push origin vX.Y.Z
```

This triggers the [Release workflow](https://github.com/opf/openproject-cli/actions/workflows/release.yml), which builds binaries for all platforms and publishes a GitHub release with zip archives attached.

Once the workflow completes, go to the [GitHub releases page](https://github.com/opf/openproject-cli/releases) and manually edit the release to add the release notes.

## Used open source tools and libraries

- [Go](https://github.com/golang/go) programming language and tools
- [Cobra](https://github.com/spf13/cobra) library for creating powerful modern CLI applications
