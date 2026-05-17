# Repository Guidelines

## Project Structure & Module Organization

This repository contains a Go CLI tool for reading public Hatena Bookmark RSS
feeds.

- `main.go`: CLI entry point. Keep it limited to wiring dependencies and exit
  handling.
- `internal/cli/`: argument parsing, output formatting, terminal-width handling,
  and CLI-level tests.
- `internal/hateb/`: Hatena RSS URL construction, HTTP fetching, RSS parsing,
  filtering, and package-level tests.
- `.plans/`: implementation plans used before making repository changes.
- `README.md`: user-facing usage instructions.

## Build, Test, and Development Commands

- `go run . --user <hatena-id>`: run the CLI and print the first RSS page.
- `go run . --user <hatena-id> --since YYYY-MM-DD`: run with date filtering.
- `go build -o hateb`: build the local executable.
- `go test ./...`: run all package tests.
- `gofmt -w <files>`: format edited Go files before committing.

## Coding Style & Naming Conventions

Use standard Go style. Format all Go files with `gofmt`; use tabs for
indentation as produced by the formatter. Keep packages focused by role:
CLI concerns belong in `internal/cli`, and RSS/domain concerns belong in
`internal/hateb`.

Prefer small exported APIs only when another package needs them. Use clear Go
names such as `FetchBookmarks`, `ParseRSS`, and `FilterBookmarksSince`.

## Testing Guidelines

Tests use Go's standard `testing` package. Place tests next to the package under
test and name files `*_test.go`. Test functions should use the
`TestNameBehavior` pattern, for example `TestRunOutputsBookmarksWithoutSince`.

Run `go test ./...` before submitting changes. Add focused tests for argument
validation, RSS parsing, HTTP behavior, filtering, and output formatting when
those areas change.

## Commit & Pull Request Guidelines

Recent history uses short Japanese commit titles, for example:

```text
はてなブックマークRSS取得CLIを実装
```

Keep commit titles concise and descriptive. If a body is needed, explain why the
change was made and list the main changes.

Pull requests should include a short summary, the commands used for verification
such as `go test ./...`, and any behavior changes visible to CLI users. Link
related issues when available.

## Agent-Specific Instructions

Before substantial edits, create or update a plan under `.plans/` and keep it in
sync with the intended change. Do not mix unrelated refactors into feature or
bug-fix changes.
