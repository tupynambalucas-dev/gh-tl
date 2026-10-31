# Command Layer (`cmd/`)

The `cmd/` package defines the entrypoint and command hierarchy for the `gh-tl` extension suite using the Cobra CLI framework.

---

## Architectural Model

Following the Single Responsibility Principle (SRP) and Open/Closed Principle (OCP):
- Each command domain is encapsulated inside its own dedicated sub-package (e.g. `cmd/onboard/`).
- The root command (`cmd/root.go`) is responsible solely for registering domain subcommands and providing the root CLI entrypoint.
- New commands can be added by introducing a new sub-package without altering existing command logic.

---

## Command Catalog

| Command | Sub-Package | Description | Status |
| :--- | :--- | :--- | :--- |
| `gh tl onboard` | [`cmd/onboard`](./onboard/README.md) | Automated workstation bootstrapper, SSH key provisioning, and identity isolation. | Active |
| `gh tl verify` | [`cmd/onboard`](./onboard/README.md) | Audit 7 pre-flight readiness compliance gates in read-only mode. | Active |
| `gh tl audit` | `cmd/audit` | Organizational security, rulesets, and secret compliance auditor. | Planned |
| `gh tl workspace` | `cmd/workspace` | Multi-repository synchronization and Mani workspace orchestrator. | Planned |

---

## Adding a New Command

1. Create a dedicated package directory: `cmd/<domain>/`.
2. Implement a constructor function exporting `New<Domain>Cmd() *cobra.Command`.
3. Provide a dedicated `README.md` and `AGENTS.md` documenting flags, workflows, and specifications.
4. Register the new command in `cmd/root.go`:
   ```go
   rootCmd.AddCommand(<domain>.New<Domain>Cmd())
   ```
