# project-dbxcli-fork

<essential-rule>
AGENTS.md is the sole authoritative project context file.
Read and edit AGENTS.md directly.
</essential-rule>

## Development Rules

Before launching agents, apply skill-xray, skill-expert, and skill-brutal to the task.
Surface expert-level issues, non-obvious issues, blindspots, stale assumptions, and hidden dependencies.
Also surface missed constraints, edge cases, false positives, verification gaps, overclaims, and weak assumptions.
Identify improvement potential, inefficiencies, and what is wrong without softening.
Use these findings to design safe slices, sequencing, checks, and boundaries for complete agent results.

Every agent prompt must require skill-xray, skill-expert, and skill-brutal for the assigned scope before acting.
It must surface non-obvious issues, blindspots, stale assumptions, hidden dependencies, and edge cases.
It must also surface verification gaps, overclaims, failure modes, weak assumptions, and what is wrong.
The agent must adjust its approach, challenge its assumptions, and flag misleading or incomplete output risks.

Implementation assignments must cover existing patterns, callers, exported-symbol consumers, and failure modes.
They must also cover concurrency safety and lifecycle cleanup.
Each assignment must state `Test decision: none` or `Test decision: update`.
`update` must name the exact existing test that follows an intentional contract change.
Never request new tests.
Prohibit broad edits, unrelated cleanup, and unassigned files.

No vague agents.
Each assignment needs exact targets, non-goals, evidence anchors, acceptance criteria, and an output contract.

Commit completed units continuously.
Before each commit, use skill-git-commit-format to determine whether staged effects are one coherent unit.
The skill owns commit-message format and evidence.
After the boundary is valid, run the repository-owned commit and push workflow.
Do not commit every trivial edit immediately or defer unrelated work into one end-of-session commit.

Every project-level quality command is quiet by default and verbose on demand.
This policy applies regardless of language or toolchain.
It covers Make targets, package scripts, Python CLIs, shell quality gates, and test runners.
Successful checks print only compact status such as `format: ok`, `lint: ok`, `test: ok`, or `check: ok`.
On failure, exit non-zero and print the failing step, exit code, and enough output to act without rerunning.
Full raw output must remain available through `--verbose`, `VERBOSE=1`, or the underlying tool's verbose mode.
New quality commands and future language setup must follow this policy instead of inventing another logging contract.

# Repository Guidelines

## Project Overview

`dbxcli` is Dropbox's Go command-line client for files, server-side search, shared links, and team administration.
It supports interactive use, scripts, CI, and agent workflows through text output and versioned JSON contracts.
Dropbox engineers maintain upstream, but it is not a formally supported Dropbox product.

## Fork & Upstream Contribution Intent

- Official upstream: [dropbox/dbxcli](https://github.com/dropbox/dbxcli).
- This checkout is the [MikeeI/project-dbxcli-fork](https://github.com/MikeeI/project-dbxcli-fork) fork, not an independently owned product.
- The goal is to support upstream with evidence-backed issues, comments, and pull requests.
- `ISSUES.md` provides the compact finding overview and global ID allocator.
- `issues/ISSUE-NNN.md` owns the complete durable record for one root cause.
- `FORMAT.md` owns research, drafting, implementation authorization, approval, and publication rules.
- Prefer small, well-scoped corrections with outsized maintainer or user value.
- Recommend a pull request when a bounded verified fix is ready and no active implementation owns it.
- Otherwise recommend a comment, a new issue, or continued investigation according to the evidence.
- Apply `skill-fork-contribution-tracking` for ledger, lifecycle, personal-branch, and upstream handoff work.
- Apply `skill-maintainer-communication` before external issues, pull requests, reviews, comments, or discussions.
- Search existing work first and follow current upstream templates and disclosure rules.
- Apply `skill-semantic-compression-3` when authoring or restructuring tracking content.
- Apply `skill-git-commit-format` while respecting explicit upstream contribution conventions.
- Never choose `Authorized-Work` or `Publication-Target` on the user's behalf.
- `Research-and-Reporting` permits research, issues, and comments but no source implementation.
- `Pull-Request-Implementation` authorizes only the scoped implementation recorded for that finding.
- Base upstream contribution branches on current `upstream/master`.
- Keep fork-only context, ledgers, configuration, and personal commits out of upstream contribution diffs.
- Reproduce claimed bugs against current upstream and run the narrowest conclusive verification.
- Publish one coherent root cause per issue, comment, or pull request.
- `origin` is `git@github.com:MikeeI/project-dbxcli-fork.git`; `upstream` is `git@github.com:dropbox/dbxcli.git`.
- `personal` owns fork-only guidance and tracking, tracks `origin/personal`, and is the fork's default branch.
- `master` remains the upstream source line; never merge personal tracking commits into contribution branches.
- Use a separate clean branch or worktree from `upstream/master` for each authorized upstream implementation.
- Initialization used `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`; refresh upstream before contributions.

## Maintainer Role & Search Development Goal

- The user is not an upstream `dropbox/dbxcli` developer or maintainer; they are an external contributor working through this fork.
- The user wants to keep developing `dbxcli` and help implement a better `search` for both agents and humans.
- Primary consumers are LLM agents that integrate the CLI and iterate faster and more precisely than humans; human usability remains required.
- Design search output and options for agents first: explicit limits and truncation, stable machine-readable fields, precise server-side filters, and no silently lost results.
- Prefer improvements that use the existing `files/search_v2` and `files/search/continue_v2` endpoints and the pinned SDK before proposing new endpoints.
- Android-internal endpoints and fields are research evidence only; never assume a public CLI token or app key can use them.
- Related evidence lives in the reverse-engineering project: `~/projects/REVERSE/project-reverse-android-dropbox/context/findings-search-comparison.md`.
- Concrete prioritized search goals live in `GOAL.md`.
- Search work follows the existing contribution contract: findings, `Authorized-Work`, upstream-first branches, and approval-gated publication.
- Never treat this goal as authorization to implement a finding or publish to upstream.

## Finding and Contribution Ledger

- At the start of every agent session, agents MUST read root `ISSUES.md` before repository work.
- `ISSUES.md` owns the global `Next finding ID` allocator and compact cross-finding overview.
- Each `issues/ISSUE-NNN.md` owns one finding's state, evidence, drafts, and next action.
- `FORMAT.md` is authoritative for research, drafting, implementation boundaries, and publication format.
- Before adding a finding, search the index and every relevant issue record for the same symptom or root cause.
- New findings MUST use `Next finding ID`.
- Create the issue file, add its index row, and increment the allocator together.
- Finding IDs use `ISSUE-NNN`, start at `ISSUE-001`, and remain permanent.
- Never reuse, renumber, or scope IDs by subsystem, status, session, or contribution type.
- Update the issue file and `ISSUES.md` together after state, authorization, target, priority, next action, or reference changes.
- Every issue record MUST use the field and section contract in `FORMAT.md`.
- New findings start with `State: Investigating`, `Authorized-Work: Not-Selected`, and `Publication-Target: Not-Selected`.
- New findings use `External-Reference: Not published.` until an external reference exists.
- Keep findings Investigating until currentness, prior art, impact, and correction value are evidence-backed.
- Clone detectors, AST matches, text similarity, shared names, and TODOs produce candidates only.
- A duplication finding requires shared change pressure, realistic drift, and simpler consolidation.
- The user selects `Authorized-Work` for each finding.
- `Research-and-Reporting` MUST NOT implement the finding.
- `Pull-Request-Implementation` MAY implement only the recorded scope after research resolves callers and failure modes.
- Pull-request work MUST verify behavior, commit, push, and reach `PR-Ready` before publication.
- Show the exact draft and target before publishing an issue, comment, or pull request to official upstream.
- Publish to official upstream only after the user approves the exact current draft and target.
- Any draft or target change requires showing the complete current draft and target again before publication.
- Run the read-only validator bundled with `skill-fork-contribution-tracking` after every ledger mutation.
- Record the final external URL in `External-Reference` immediately after publication.
- Keep `FORMAT.md`, `ISSUES.md`, `issues/`, and fork-only `AGENTS.md` changes out of upstream contribution diffs.

### External publication approval

Only an external issue, comment, review, discussion, or pull request write is approval-gated.
Before publication, read current contribution guidance and explain applicable project policy.
The human must be able to review and own every submission statement.
Fork commits, pushes, tracking updates, and source implementation follow the active repository contract.

## Architecture & Data Flow

- `main.go` resolves build-version metadata, registers the version command, and delegates execution to `cmd`.
- `cmd/` owns Cobra commands, OAuth credentials, Dropbox SDK calls, validation, and command output.
- `internal/output/` owns output-format selection and rendering.
- `internal/jsonschema/` owns command manifests and generated JSON schema definitions.
- `cmd/search.go` sends queries to Dropbox `files/search_v2` and pages with `files/search/continue_v2`.
- Search is server-side; do not replace it with mounts, local synchronization, or recursive listings and local filtering.
- Search-result presentation may sort returned metadata locally; distinguish that from server-side search ordering.

## Key Directories

- `docs/commands/` contains generated Cobra command references.
- `docs/json-schema/v1/` contains generated automation contracts and schema documentation.
- `tools/gen-docs/` and `tools/gen-json-schemas/` own their respective generators.
- `packaging/` and `build.sh` own release archives and platform packaging.
- `skills/dbxcli/` owns the canonical CLI skill; `plugin/claude/skills/dbxcli/` is its synchronized copy.
- `.github/workflows/` owns upstream CI, release, and packaging automation.

## Development Commands

```sh
go run . --help
go run . search --help
go build .
go test ./...
go vet ./...
go build ./...
go test -race ./...
go run ./tools/gen-docs
go run ./tools/gen-json-schemas
```

- Run the three upstream validation commands `go test ./...`, `go vet ./...`, and `go build ./...` before submission.
- Run `go test -race ./...` for concurrency-sensitive changes.
- Regenerate command documentation and JSON schemas when command structure, help, or JSON metadata changes.
- Run `./scripts/sync-claude-skill.sh` after canonical CLI skill changes.
- Run `VERSION=<version> ./build.sh` only when cross-platform release binaries are needed.

## Code Conventions

- Preserve the upstream module path `github.com/dropbox/dbxcli/v3` and root `main.go` plus `cmd/` layout.
- This upstream fork deliberately does not use the house `src/` layout or generated Makefile workflow.
- Follow existing Cobra, Dropbox SDK, OAuth, output-envelope, and schema conventions.
- Preserve stable JSON v1 success/error envelopes and machine-readable help discovery.
- Treat generated documentation and schemas as outputs; edit their owners and regenerate them.
- Keep contributions focused and follow upstream's Dropbox Contributor License Agreement requirements.

## Important Files

- `go.mod` owns the module path, Go language version, and dependencies.
- `CONTRIBUTING.md` owns upstream development and submission requirements.
- `docs/automation.md` owns output, error, exit-code, and authentication contracts.
- `docs/json-schema/v1/README.md` owns JSON schema stability policy.
- `.github/workflows/ci.yml` owns cross-platform validation and generated-output drift checks.
- `README.md`, `CHANGELOG.md`, `SECURITY.md`, and `LICENSE` remain upstream-owned project documents.

## Runtime/Tooling Preferences

- `go.mod` currently requires Go 1.26.0; upstream CI pins Go 1.26.6.
- Preserve upstream's tooling instead of running setup-project or installing house templates over it.
- No repository Makefile or JavaScript package workflow exists in the initial upstream checkout.
- Saved OAuth credentials default to `~/.config/dbxcli/auth.json`; `DBXCLI_AUTH_FILE` selects another file.
- `DBXCLI_ACCESS_TOKEN` supplies a direct token that is neither saved nor refreshed.
- Saved login credentials support refresh; direct-token automation has a separate lifecycle.
- Go's `HTTPS_PROXY`, `HTTP_PROXY`, and `NO_PROXY` settings apply to API and OAuth network requests.
- Repository initialization does not authorize Dropbox login, account queries, global installation, or file transfers.

## Testing & QA

- Tests live alongside their owning Go packages; search contracts are covered in `cmd/search_test.go`.
- Use focused existing tests while investigating; upstream validation remains the submission gate for source changes.
- Preserve CLI help, command metadata, generated references, and JSON schemas together when changing public commands.
- Fork-only context and empty-ledger initialization require ledger validation, not unrelated source-test runs.
- Existing upstream quality commands retain upstream output behavior; do not rewrite them for this fork setup.
