# gh-tl CLI Extension Suite Agent Context

This document defines the agent bounded context router for `gh-tl`, the enterprise GitHub CLI extension suite for the `tupynambalucas-dev` organization.

## Architectural Boundaries

- **Package Namespace**: `github.com/tupynambalucas-dev/gh-tl`
- **CLI Framework**: `github.com/spf13/cobra` combined with `github.com/cli/go-gh/v2`
- **Distribution Model**: Precompiled multi-platform binaries via `cli/gh-extension-precompile@v2` on GitHub Releases.
- **Root Invocation**: `gh tl <command>`

## Domain Commands

- `gh tl onboard`: Interactive and automated developer workstation bootstrapper.
- `gh tl verify`: Pre-flight compliance verification auditing 7 strict readiness gates.

## Subsystem Navigation

- [README.md](./README.md): Primary developer documentation and usage manual.
- [cmd/root.go](./cmd/root.go): Root Cobra command definition (`gh tl`).
- [cmd/onboard.go](./cmd/onboard.go): Automated bootstrapping and key synchronization workflows.
- [cmd/verify.go](./cmd/verify.go): 7 pre-flight readiness gate compliance checks.
- [internal/git](./internal/git/git.go): Git version parser and global config wrapper.
- [internal/ssh](./internal/ssh/ssh.go): Ed25519 keypair generation and allowed_signers manager.
- [internal/config](./internal/config/config.go): Conditional include injector (`hasconfig:remote.*.url`).
- [internal/github](./internal/github/keys.go): GitHub API authentication and signing key registration.

## Governance Standards

- All commits and release tags MUST be cryptographically signed with SSH.
- Conventional Commits standard (`feat`, `fix`, `chore`, etc.).
- Strictly zero emojis in code comments, terminal output, or documentation.
- All documentation, comments, and CLI strings MUST be written in en-US.
