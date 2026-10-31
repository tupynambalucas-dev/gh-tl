# Onboarding & Readiness Verification Domain (`gh tl onboard` / `gh tl verify`)

This directory houses the implementation and specification for developer workstation provisioning, cryptographic SSH key registration, conditional Git identity injection, and pre-flight readiness gate compliance.

---

## Commands

### `gh tl onboard`

Executes the automated workstation setup workflow:
- Verifies Git version compatibility (`>= 2.36.0`).
- Resolves developer identity from authenticated GitHub profile or input flags.
- Provisions an Ed25519 SSH signing keypair (`~/.ssh/github_tupynambalucas-dev_ed25519`).
- Synchronizes the public key with the GitHub API for both transport authentication and commit signing.
- Configures isolated corporate Git identity via `hasconfig:remote.*.url`.
- Updates `~/.ssh/allowed_signers` for local signature verification.
- Executes the 7 pre-flight readiness gates upon completion.

```bash
# Standard interactive setup
gh tl onboard

# Custom parameters
gh tl onboard --name "Developer Name" --email "contato@tupynambalucas.dev" --org "tupynambalucas-dev"
```

#### Flags

| Flag | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `--name` | `string` | *(GitHub Name / Login)* | Full developer name for Git commits |
| `--email` | `string` | `contato@tupynambalucas.dev` | Corporate email address |
| `--key-path` | `string` | `~/.ssh/github_tupynambalucas-dev_ed25519` | File path for the generated private key |
| `--org` | `string` | `tupynambalucas-dev` | Target GitHub organization name |
| `--skip-upload` | `bool` | `false` | Skip registering keys via GitHub API |
| `--verify` | `bool` | `false` | Run verification gates without modifying configuration |

---

### `gh tl verify`

Audits the workstation environment across the 7 compliance gates in read-only mode, without making any modifications to disk or Git configurations. Exits with code `0` on success or code `1` if any gate fails.

```bash
gh tl verify
```

---

## Pre-Flight Readiness Gates Specification

| Gate | Verification Target | Compliance Requirement | Remediation Procedure |
| :--- | :--- | :--- | :--- |
| **Gate 1: Git Engine** | `git --version` | Version `>= 2.36.0` | Upgrade Git via package manager (`winget upgrade Git.Git` or `brew upgrade git`). |
| **Gate 2: Signing Key** | File existence | Ed25519 public key in `~/.ssh` or `user.signingkey` | Execute `gh tl onboard` to generate the corporate keypair. |
| **Gate 3: Commit Signing** | Git config | `gpg.format = ssh` and `commit.gpgsign = true` | Execute: `git config --global gpg.format ssh` and `git config --global commit.gpgsign true`. |
| **Gate 4: Mani Orchestrator** | `mani --version` | Mani CLI available in `PATH` | Install Mani via Scoop (`scoop install mani`) or download from GitHub Releases. |
| **Gate 5: Developer Runtime** | `node`, `pnpm` | Node.js `>= 22` and active `pnpm` | Install Node.js LTS and enable corepack via `corepack enable`. |
| **Gate 6: SSH Transport** | Git config | `url.git@github.com:.insteadOf = https://github.com/` | Execute: `git config --global url."git@github.com:".insteadOf https://github.com/`. |
| **Gate 7: GitHub CLI Session** | `gh auth status` | Valid session for `github.com` | Execute: `gh auth login --git-protocol ssh --web`. |

---

## Technical Architecture

### Conditional Git Configuration Injection

To eliminate filesystem-dependent assumptions, `gh tl onboard` configures remote URL globbing inside the user's global `~/.gitconfig`:

```ini
[includeIf "hasconfig:remote.*.url:git@github.com:tupynambalucas-dev/**"]
    path = ~/.gitconfig-tupynambalucas-dev

[includeIf "hasconfig:remote.*.url:https://github.com/tupynambalucas-dev/**"]
    path = ~/.gitconfig-tupynambalucas-dev
```

The isolated file `~/.gitconfig-tupynambalucas-dev` defines:

```ini
[user]
    name = Developer Name
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
