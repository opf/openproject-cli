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
as `op VERB RESOURCE [additional information]`. You will see plenty of examples within this section.

### Discoverability

Discoverability is key. As we won't document every single command within this README, it is important for the CLI tool,
to provide helpful information at any time and any level of the command structure. To achieve that, the OpenProject CLI
has a help page on every command. Adding the `--help/-h` flag to any chain of commands will tell the user of all
possibilities he can use from here.

```shell
# General help
op -h

# Help about what things I can list
op list -h

# Help about how can I narrow down my list of work
# packages to get exactly what I was looking for
op list workpackages -h
```

The second basement is the autocompletion. To set it up correctly, follow the steps explained in the `Completion`
section above.

Once it is working, you can use it to discover possible commands and flags while typing.

```shell
# Chaining commands: hitting completion key after
op list
# returns
activities     -- Lists activities for work package
notifications  -- Lists notificationswork package
projects       -- Lists projects
workpackages   -- Lists work packages

# Discover flags: hitting completion key after
op update workpackge 42 -
# returns
--action    -a  -- Executes a custom action on a work package
--assignee      -- Assign a user to the work package
--attach        -- Attach a file to the work package
--help      -h  -- help for workpackage
--subject       -- Change the subject of the work package
--type      -t  -- Change the work package type
```

### Prominent examples

There are a couple of use cases, you might want to execute from the command line. In this section we provide a handful
of examples, that might be useful for a great number of people.

#### Creation

```shell
# Creating a work package in a project only by subject.
# Work package is created with many default values (as for type and status),
# very similar to how a work package is created inline in a work package table.
op create workpackge --project 11 'Document new CLI tool'

# Same command with shorthands and directly open it in a browser to continue working on it.
op create workpackge -p11 'Document new CLI tool' -o

# Validating the creation of a child work package without persisting it.
# The parent determines the project automatically.
op create workpackage --parent 74316 --type Implementation 'Build reusable skill' --dry-run --json
```

#### Listing

```shell
# Get a list of unread notifications and filter them by reason
op list notifications --reason mentioned

# Get a list of all work packages assigned to me
op list workpackages --assignee me
```

#### Updating

```shell
# Executing a custom action on a work package
op update workpackage 42 --action Claim

# Batch updating some properties of a work package
# Valid input will get processed, while invalid (e.g. wrongly typed) input will get omitted
op update workpackage 42 --subject 'The new subject' --status 'In Progress' --type Implementation

# Uploading an attachment to a work package
op update workpackage 42 --attach ./Downloads/Report.pdf

# Resolving and validating field updates as JSON before applying them
# This first slice supports schema-resolved custom fields via --set.
op update workpackage 74316 --set 'Votes=3' --dry-run --json
```

#### Inspecting

```shell
# Inspecting a work package with more details,
# then in the work package list command
op inspect workpackage 42

# Inspecting a work package and its direct children as machine-readable JSON
op inspect workpackage 74316 --children --json
```

#### JSON error codes

The `--json` variants of `inspect`, `create`, and `update` return a stable error envelope:

```json
{"error": {"code": "<code>", "message": "<message>"}}
```

| Code | Trigger | Mutation persisted? |
|---|---|---|
| `invalid_argument` | Missing positional arg, invalid id, local validation failure, malformed `--set`, or no update options under `--json`/`--dry-run` | No |
| `conflicting_arguments` | Incompatible flag combinations, e.g. `--description` + `--set`, `--open` + `--json`, `--dry-run` without `--json` | No |
| `api_error` | Underlying OpenProject API call failed while resolving, creating, updating, or dry-running a work package | No |
| `post_apply_inspect_failed` | CREATE or PATCH succeeded but the follow-up inspect to build the JSON response failed — the mutation **has persisted**, only the response payload could not be assembled | **Yes** |
| `ambiguous_field` | A `--set` label matched more than one schema field | No |
| `duplicate_field` | A `--set` assignment names the same API field twice | No |
| `unknown_field` | A `--set` label or API field is not present in the resolved schema | No |
| `unsupported_field_type` | A `--set` field's schema type is not yet coercible | No |
| `invalid_field_value` | A `--set` value fails type coercion (e.g. non-integer for an `Integer` field) | No |
| `non_writable_field` | A `--set` field exists in the schema but is not writable | No |

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
