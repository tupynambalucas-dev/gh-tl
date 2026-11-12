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

## Architectural Plans & Artifact Dual-Storage Standard

Architectural blueprints, technical analyses, and master plans follow strict enterprise storage standards:
- **Storage Location**: Blueprints, analyses, and reports reside in `plans/<artifact-slug>/` (or `plans/<artifact-slug>.md`), with `.agents/plans/<artifact-slug>/` fully recognized for backward compatibility.
- **Topologies**: Topology 1 (Single-File Artifact) for atomic refactors, analyses, and reports; Topology 2 (Multi-Artifact Suite) for multi-domain frameworks with master `README.md` and numbered domain documents.
- **Dual-Storage Invariant**: Whenever an agent creates an architectural plan, analysis, or report artifact in the Antigravity conversation directory (`<appDataDir>/brain/<conversation-id>/`), an identical synchronized copy MUST ALWAYS be mirrored in `plans/<artifact-slug>/` (or `plans/<artifact-slug>.md`).
- **Ignore & Index Triad**: All repositories ignore `plans/` and `.agents/plans/` in `.gitignore`, while explicitly un-ignoring them in `.ignore` and `.rgignore` (`!plans/`, `!plans/**`, `!.agents/plans/`, `!.agents/plans/**`) and setting `"search.useIgnoreFiles": false` in `.vscode/settings.json` so files can always be attached with `@`.

## Required Skills

When authoring or modifying code, commands, or documentation in this extension suite, agents MUST autonomously activate:

- **Atelier Plugin (`@tupynambalucas-dev/plugin-atelier`)**:
  - **`plan-architect`**: Designing, structuring, authoring, and auditing architectural plans, technical analyses, and governance reports following the dual-storage standard.
  - **`router-architect`**: Authoring, auditing, or updating hierarchical `AGENTS.md` context routers.
- **Repository Governance Plugin (`@tupynambalucas-dev/plugin-repo`)**:
  - **`readme-architect`**: Authoring or auditing `README.md` files and maintaining recursive navigation.
  - **`community-health-architect`**: Auditing and aligning community governance and health files.
- **GitHub Flow Plugin (`@tupynambalucas-dev/plugin-github-flow`)**:
  - **`github-flow-operator`**: Executing Git branching, SSH-signed Conventional Commits, and PR lifecycle workflows.
