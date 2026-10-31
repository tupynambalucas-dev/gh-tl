# Internal Shared Packages (`internal/`)

The `internal/` directory houses private, reusable core domain modules consumed exclusively by `gh-tl` commands.

---

## Package Directory

| Package | Purpose |
| :--- | :--- |
| [`internal/config`](./config/config.go) | Manages Git global configurations, corporate identity templates, and conditional includes. |
| [`internal/git`](./git/git.go) | Low-level Git binary execution, semantic version parsing, and global configuration getters/setters. |
| [`internal/github`](./github/keys.go) | Integrates with GitHub API v3 using `go-gh/v2` for user profile queries and SSH key registration. |
| [`internal/ssh`](./ssh/ssh.go) | Ed25519 cryptographic key generation, file existence checks, and `allowed_signers` management. |
| [`internal/ui`](./ui/ui.go) | Cross-platform terminal typography, ANSI colors, headers, and gate status printers. |
