# Command Layer Agent Context

This document defines the agent bounded context for the `cmd/` package in `gh-tl`.

## Subsystem Navigation

- [README.md](./README.md): Command layer architecture and command catalog.
- [root.go](./root.go): Root command initialization and subcommand registration.
- [onboard/](./onboard/AGENTS.md): Onboarding and pre-flight readiness gate domain.

## Architecture Guidelines

- Commands must follow Cobra idiomatic patterns with decoupled execution functions.
- Domain commands must reside in separate sub-packages under `cmd/`.
- Every command package must include both `README.md` and `AGENTS.md`.
