# theta

Print the first and last lines of text files with line numbers and optional ANSI color.

Inspired by [heads-and-tails](https://github.com/managedkaos/heads-and-tails).

```bash
theta [flags] FILE [FILE...]
theta -l 5 notes.txt report.txt
theta --head 20 --tail 5 app.log
theta -T plain app.log
```

Each section begins with its filename. A 20-dash line separates the head and tail;
multiple files are separated by a blank line. Head and tail sections may overlap,
matching the legacy application. Empty files still show their name and separator.
Unreadable files are reported to stderr; other files are processed and the command
exits with status 1. Usage errors also exit with status 1. As in the legacy
application, files are read into memory and individual lines must be smaller
than 1 MiB.

| Long option | Short | Default | Purpose |
| --- | --- | --- | --- |
| `--lines N` | `-l` | `10` | Lines in both head and tail |
| `--head N` | `-h` | `--lines` | Override head count |
| `--tail N` | `-t` | `--lines` | Override tail count |
| `--color MODE` | `-c` | `auto` | `auto`, `always`, or `never` |
| `--theme NAME` | `-T` | `default` | `default` or `plain` |
| `--version` | `-v` | | Version and build metadata |
| `--help` | | | Help |

Counts must be non-negative; zero suppresses that section's lines. Explicit head
and tail counts override `--lines` regardless of argument order. Color is automatic
when stdout is a terminal; `plain` disables color even with `--color always`.
Flags can appear before or after filenames. Use `--` before names starting with `-`.
No arguments prints help. **`-h` selects head lines; use `--help` for help.**
Legacy single-dash long options such as `-lines` become `--lines` in theta.

## Development

Requires Go 1.26.4 or later.

```bash
make build       # bin/theta
make test        # tests with race detection
make verify      # formatting, vet, tests, build
make build-all   # Linux/macOS amd64 + arm64; Windows amd64
make run ARGS="--lines 3 README.md"
make man         # man/theta.1
make install     # ~/.local/bin and ~/.local/share/man (override PREFIX)
make snapshot    # GoReleaser v2: archives under dist/, no publishing
```

## CI and releases

Every push and pull request runs tests with race detection, vet, and builds on
Linux, macOS, and Windows. Linux also checks formatting and the Make build.

Push a semantic version tag such as `v1.0.0`, or dispatch **Build and Release**
with an existing version tag. The workflow tests source, builds snapshot archives,
and smoke tests all five supported platform binaries before publishing through
GoReleaser. Releases include `.tar.gz` archives for Linux/macOS, `.zip` for Windows,
README, license, man page, SHA-256 checksums, and a generated changelog.
The release job uses GitHub's built-in `GITHUB_TOKEN` with contents write permission.
