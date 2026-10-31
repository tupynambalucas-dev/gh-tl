# gh-tl CLI Extension Suite Agent Context

This document defines the root agent bounded context router for `gh-tl`, the enterprise GitHub CLI extension suite for the `tupynambalucas-dev` organization.

## Architectural Boundaries

- **Package Namespace**: `github.com/tupynambalucas-dev/gh-tl`
- **CLI Framework**: `github.com/spf13/cobra` combined with `github.com/cli/go-gh/v2`
- **Distribution Model**: Precompiled multi-platform binaries via `cli/gh-extension-precompile@v2` on GitHub Releases.
- **Root Invocation**: `gh tl <command>`
- **Design Principles**: Domain-Driven Design (DDD) bounded contexts, SOLID principles (SRP, OCP), and Single Source of Truth (SSOT).

## Subsystem Navigation

- [README.md](./README.md): Primary suite documentation, installation, and command catalog.
- [cmd/](./cmd/AGENTS.md): Command layer and Cobra subcommand definitions.
  - [cmd/onboard/](./cmd/onboard/AGENTS.md): Onboarding and pre-flight readiness gate domain.
- [internal/](./internal/AGENTS.md): Shared, domain-agnostic internal modules (`config`, `git`, `github`, `ssh`, `ui`).

## Governance Standards

- All commits and release tags MUST be cryptographically signed with SSH.
- Conventional Commits standard (`feat`, `fix`, `docs`, `chore`, etc.).
- Strictly zero emojis across all code, comments, terminal output, PRs, and documentation.
- All technical documentation, comments, and CLI strings MUST be written in en-US.
- The AGENTS-README co-existence invariant is strictly enforced: every package containing an `AGENTS.md` must also contain a corresponding `README.md`.
