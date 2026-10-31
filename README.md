# gh-onboard

Official GitHub CLI extension for automated workstation provisioning, cryptographic SSH commit signing, and multi-repository readiness compliance across the `tupynambalucas-dev` organization.

---

## Overview

`gh-onboard` streamlines the developer onboarding experience by replacing fragile manual setup steps and platform-dependent scripts with a single native GitHub CLI command. It provisions cryptographic Ed25519 signing keys, synchronizes credentials with GitHub, configures directory-agnostic conditional Git identities (`hasconfig:remote.*.url`), and validates workstation compliance across 7 pre-flight readiness gates.

### Why gh-onboard?

- **Zero-Friction Onboarding**: New developers do not need to clone repositories or manually generate SSH keys before their first contribution.
- **Cross-Platform Binary**: Written in Go and compiled natively for Windows, macOS, and Linux via GitHub CLI extensions.
- **Dual-Purpose SSH Keys**: Automates the dual registration of SSH keys for both Git transport (`authentication`) and commit validation (`signing`).
- **Directory-Agnostic Isolation**: Uses Git conditional includes (`hasconfig:remote.*.url`) to apply corporate identity strictly to organizational repositories, regardless of where they are cloned on disk.
- **Pre-Flight Readiness Gates**: Validates toolchains, orchestrators (Mani), signing parameters, and CLI authentication status in seconds.

---

## Installation

### Prerequisites

1. **GitHub CLI (`gh`)**: Install via your package manager:
   - Windows: `winget install GitHub.cli` or `scoop install gh`
   - macOS: `brew install gh`
   - Linux: `sudo apt install gh` or distribution equivalent
2. **Authenticate with GitHub**:
   ```bash
   gh auth login --git-protocol ssh --web
   ```

### Install Extension

Install directly from GitHub via the CLI extension manager:

```bash
gh extension install tupynambalucas-dev/gh-onboard
```

To update to the latest release:

```bash
gh extension upgrade onboard
```

---

## Usage

### 1. Automated Workstation Bootstrapping

Run the interactive setup workflow:

```bash
gh onboard
```

Options:
```text
Flags:
      --name string       Developer full name for Git commits (defaults to GitHub profile)
      --email string      Corporate email address (default "contato@tupynambalucas.dev")
      --key-path string   Target SSH private key path (defaults to ~/.ssh/github_tupynambalucas-dev_ed25519)
      --org string        Target GitHub organization name (default "tupynambalucas-dev")
      --skip-upload       Skip registering SSH keys via GitHub API
      --verify            Execute pre-flight readiness verification only
```

Example with custom name and corporate email:
```bash
gh onboard --name "Lucas Tupynamba" --email "contato@tupynambalucas.dev"
```

### 2. Pre-Flight Readiness Verification

To audit your current workstation without applying any configuration changes, run the `verify` subcommand:

```bash
gh onboard verify
```

Output preview:
```text
==========================================================
tupynambalucas-dev Workstation Readiness Gate Verification
==========================================================

[Gate 1/7] Inspecting Git version...  PASS git version 2.54.0 (>= 2.36.0)
[Gate 2/7] Inspecting Ed25519 signing key...  PASS Active signing key detected: ~/.ssh/github_tupynambalucas-dev_ed25519.pub
[Gate 3/7] Inspecting cryptographic commit signing config...  PASS gpg.format=ssh, commit.gpgsign=true
[Gate 4/7] Inspecting Mani multi-repository orchestrator...  PASS Version: 0.32.1
[Gate 5/7] Inspecting Node.js and pnpm runtime...  PASS Node v24.19.0, pnpm 10.33.2
[Gate 6/7] Inspecting default Git SSH transport override...  PASS SSH transport rewrite configured
[Gate 7/7] Inspecting GitHub CLI authentication...  PASS Authenticated as @tupynambalucas

----------------------------------------------------------
 PASS Workstation is fully compliant with all readiness gates.
```

---

## Readiness Gates Specification

The extension enforces 7 strict compliance gates:

| Gate | Target Component | Requirement | Purpose |
| :--- | :--- | :--- | :--- |
| **Gate 1** | Git Engine | `>= 2.36.0` | Required for `hasconfig:remote.*.url` and native SSH commit signing |
| **Gate 2** | Cryptographic Key | Ed25519 Keypair | Non-repudiation and high-performance elliptic curve signature |
| **Gate 3** | Commit Signing | `gpg.format=ssh`, `commit.gpgsign=true` | Zero-compromise commit and tag signing requirement |
| **Gate 4** | Multi-Repo Orchestrator | Mani CLI (`>= 0.30`) | Batch command execution and repository synchronization |
| **Gate 5** | Developer Runtime | Node.js (`>= 22`), pnpm (`>= 9`) | Monorepo toolchains, linting, and package dependency resolution |
| **Gate 6** | Transport Rewrite | `url.git@github.com:.insteadOf` | Enforces SSH key authentication over HTTPS password prompts |
| **Gate 7** | GitHub CLI Auth | `gh auth status` | Validates active session and API token privileges |

---

## Architecture and Git Configuration

### Conditional Identity Isolation

`gh-onboard` configures your global `~/.gitconfig` with remote URL matching:

```ini
[includeIf "hasconfig:remote.*.url:git@github.com:tupynambalucas-dev/**"]
    path = ~/.gitconfig-tupynambalucas-dev

[includeIf "hasconfig:remote.*.url:https://github.com/tupynambalucas-dev/**"]
    path = ~/.gitconfig-tupynambalucas-dev
```

This guarantees that:
- Any repository belonging to `tupynambalucas-dev` automatically commits with your corporate email and cryptographic SSH key.
- Personal or open-source repositories residing outside the organization retain your standard global Git profile.

### Isolated Corporate Configuration (`~/.gitconfig-tupynambalucas-dev`)

```ini
[user]
    name = Your Name
    email = contato@tupynambalucas.dev
    signingkey = ~/.ssh/github_tupynambalucas-dev_ed25519.pub

[commit]
    gpgsign = true

[tag]
    gpgsign = true

[gpg]
    format = ssh

[gpg "ssh"]
    allowedSignersFile = ~/.ssh/allowed_signers
```

---

## Development and Testing

### Build from Source

```bash
git clone git@github.com:tupynambalucas-dev/gh-onboard.git
cd gh-onboard

# Run tests
go test -v ./...

# Build local binary
go build -v .

# Install locally as an extension for testing
gh extension install .
```

### Automated Release Pipeline

Releases are automatically cross-compiled and published via GitHub Actions using `cli/gh-extension-precompile@v2`. Tagging a commit with `v*` triggers multi-platform binary compilation and asset attachment:

```bash
git tag -s v0.1.0 -m "release: v0.1.0"
git push origin v0.1.0
```

---

## License

Copyright (c) 2026 tupynambalucas.dev. Licensed under the [Apache License 2.0](./LICENSE).
