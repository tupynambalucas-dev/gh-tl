# gh-tl

Official enterprise GitHub CLI extension suite for the `tupynambalucas-dev` organization, delivering modular engineering tooling, automated onboarding, cryptographic compliance, and multi-repository workflows.

---

## Overview

`gh-tl` serves as the centralized command-line interface for Tupynambá Lucas engineering teams. Built on top of the GitHub CLI (`gh`) and the official `go-gh` SDK, it provides an extensible, modular command suite designed to eliminate operational friction and enforce organizational governance.

### Architectural Principles

The repository adheres to core engineering principles defined in the organizational design system:

- **Domain-Driven Design (DDD)**: Tooling capabilities are divided into isolated, domain-centric packages within `cmd/` (such as onboarding, verification, and future governance tools).
- **Single Responsibility (SRP)**: Each command module encapsulates a single domain workflow and maintains its own dedicated documentation.
- **Open/Closed (OCP)**: The CLI suite is open to extension via new Cobra subcommands without requiring modifications to existing domain implementations.
- **Single Source of Truth (SSOT)**: Domain specifications and compliance rules reside in their respective command packages rather than being duplicated at the root layer.

---

## Installation

### Prerequisites

1. **GitHub CLI (`gh`)**:
   - Windows: `winget install GitHub.cli` or `scoop install gh`
   - macOS: `brew install gh`
   - Linux: `sudo apt install gh`
2. **Authenticate with GitHub**:
   ```bash
   gh auth login --git-protocol ssh --web
   ```

### Install Extension

Install directly from GitHub:

```bash
gh extension install tupynambalucas-dev/gh-tl
```

To update to the latest release:

```bash
gh extension upgrade tl
```

---

## Command Catalog

The suite organizes commands into domain-specific modules:

| Command | Domain | Description | Documentation |
| :--- | :--- | :--- | :--- |
| `gh tl onboard` | Onboarding | Automated developer bootstrapping, Ed25519 key generation, GitHub API registration, and conditional Git inclusion. | [cmd/onboard/README.md](./cmd/onboard/README.md) |
| `gh tl verify` | Compliance | Read-only pre-flight audit validating workstation compliance across 7 strict readiness gates. | [cmd/onboard/README.md](./cmd/onboard/README.md#gh-tl-verify) |
| `gh tl audit` | Governance | Organization security auditing, branch protection rulesets, and secret scanning. | *Planned* |
| `gh tl workspace` | Multi-Repo | Workspace synchronization and multi-repository batch coordination with Mani. | *Planned* |

To explore command-specific flags and options:

```bash
gh tl --help
gh tl <command> --help
```

---

## Repository Structure

```text
gh-tl/
├── cmd/                          # Command layer (Cobra CLI subcommands)
│   ├── root.go                   # Suite entrypoint and command registrar
│   ├── onboard/                  # Onboarding & verification domain
│   │   ├── onboard.go            # 'gh tl onboard' implementation
│   │   ├── verify.go             # 'gh tl verify' implementation
│   │   ├── README.md             # Domain specification and 7-gate reference
│   │   └── AGENTS.md             # Onboarding domain agent router
│   ├── README.md                 # Command layer documentation
│   └── AGENTS.md                 # Command layer agent router
├── internal/                     # Private, domain-agnostic shared modules
│   ├── config/                   # Git configuration and include templates
│   ├── git/                      # Git execution and version parsing
│   ├── github/                   # GitHub REST client (go-gh/v2)
│   ├── ssh/                      # Ed25519 key management and allowed_signers
│   ├── ui/                       # Terminal formatting and typography
│   ├── README.md                 # Internal modules documentation
│   └── AGENTS.md                 # Internal modules agent router
├── .github/                      # Workflows and extension metadata
│   └── workflows/
│       ├── ci.yml                # Multi-platform build and unit test matrix
│       └── release.yml           # Automated cross-compilation pipeline
├── AGENTS.md                     # Root context router
├── go.mod                        # Go module definition
├── go.sum                        # Go dependency checksums
├── LICENSE                       # Apache License 2.0
├── main.go                       # Binary entrypoint
└── README.md                     # Root suite documentation
```

---

## Development and Contribution

### Build Locally

```bash
git clone git@github.com:tupynambalucas-dev/gh-tl.git
cd gh-tl

# Execute unit tests
go test -v ./...

# Build local binary
go build -v .

# Install locally into gh for testing
gh extension install .
```

### Release Pipeline

Releases are automatically cross-compiled for Windows, macOS, and Linux using `cli/gh-extension-precompile@v2` when a signed Git tag is pushed:

```bash
git tag -s v0.1.0 -m "release: v0.1.0"
git push origin v0.1.0
```

---

## License

Copyright (c) 2026 tupynambalucas.dev. Licensed under the [Apache License 2.0](./LICENSE).
