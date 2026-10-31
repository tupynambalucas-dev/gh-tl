# Onboarding Domain Agent Context

This document defines the agent bounded context for the `onboard` and `verify` command package within `gh-tl`.

## Bounded Context Scope

- **Package**: `github.com/tupynambalucas-dev/gh-tl/cmd/onboard`
- **Responsibilities**:
  - Provisioning Ed25519 SSH keys for corporate identity.
  - Interfacing with GitHub API to register authentication and signing keys.
  - Configuring global Git defaults and `hasconfig:remote.*.url` conditional identity files.
  - Executing the 7 pre-flight readiness gates in read-only audit mode.

## Key Files

- [README.md](./README.md): Primary technical specification and 7-gate reference table.
- [onboard.go](./onboard.go): Implementation of `gh tl onboard` setup workflow.
- [verify.go](./verify.go): Implementation of `gh tl verify` pre-flight readiness audit.

## Invariants

- Non-destructive execution: `gh tl verify` MUST never mutate files or configurations.
- SSH path normalization: Git config values on Windows must use forward slashes (`/`).
- Idempotency: Running `onboard` multiple times must preserve existing valid keys without duplicating allowed_signers entries.
