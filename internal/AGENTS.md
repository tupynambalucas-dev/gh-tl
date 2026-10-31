# Internal Packages Agent Context

This document defines the agent bounded context for shared internal packages within `gh-tl`.

## Packages

- [README.md](./README.md): Overview of internal modules.
- `config`: Handles Git configuration files and `hasconfig:remote.*.url` includes.
- `git`: Interacts with system `git` CLI, parses semver versions.
- `github`: REST API client powered by `go-gh/v2`.
- `ssh`: Ed25519 keypair management and OpenSSH `allowed_signers`.
- `ui`: Terminal output, ANSI colors, and compliance gate reporting.

## Guidelines

- Code inside `internal/` must remain domain-agnostic and reusable across all `cmd/` packages.
- Functions must return structured errors; CLI output is reserved for `cmd/` or `ui`.
- All unit tests for internal packages must be colocated (`*_test.go`).
