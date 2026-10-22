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

## Private Fork Purpose and Branches

- Upstream source: [dropbox/dbxcli](https://github.com/dropbox/dbxcli).
- This checkout is the user's [personal fork](https://github.com/MikeeI/project-dbxcli-fork).
- Develop for the user's own workflows, with agents first and human usability preserved.
- Upstream integration is not a delivery goal or a prerequisite for private development.
- Do not assume upstream interest or rejection without evidence.
- `origin` is `git@github.com:MikeeI/project-dbxcli-fork.git`.
- `upstream` is `git@github.com:dropbox/dbxcli.git`.
- `personal` owns private source development, guidance, and tracking; it tracks `origin/personal`.
- `personal` remains the fork's default branch and the normal commit and push target.
- `master` remains the upstream source line; do not merge private development into it.
- Implement scoped user requests directly on `personal`; upstream-first branches and PR drafts are not required.
- Verify fixes against the current fork revision, not only the inherited upstream baseline.
- Keep changes focused and preserve existing tooling to limit future upstream merge costs.
- Initialization used `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`.
- Fetching or merging newer upstream work is a separate operation, not a prerequisite for every private change.
- Plans, goals, and findings do not authorize source implementation; require a direct scoped user instruction.
- Fork-only development does not waive existing JSON, CLI, authentication, or data-integrity contracts.

## Maintainer Role & Search Development Goal

- The user maintains this private fork, not upstream `dropbox/dbxcli`.
- Develop a better `search` for the user's agents and human workflows.
- LLM agents are the primary consumers; human usability remains required.
- Make limits and truncation explicit, preserve machine-readable fields, and never silently discard collected results.
- Prefer precise server-side filters and the existing search endpoints with the pinned SDK.
- A reverse-engineering project for the Dropbox Android app exists: `~/projects/REVERSE/project-reverse-android-dropbox`.
- Its comparison of dbxcli, the public API, and Android search is `context/findings-search-comparison.md` in that project.
- It is static analysis plus public API documentation; no live Dropbox requests were made.
- Android-internal endpoints and fields (for example DCS, Ask, date filters, snippets) are research evidence only.
- Never assume a public CLI token or app key can use them, and never copy Android behavior into `GOAL.md` as a promise.
- Concrete prioritized search goals live in `GOAL.md`.
- Search implementation follows the private-fork workflow and the scope recorded in `GOAL.md`.
- A goal is not authorization to implement source changes, alter compatibility guarantees, or publish upstream.

## Finding and Research Ledger

- Read root `ISSUES.md` at the start of every session before repository work.
- `ISSUES.md` owns the global `Next finding ID` allocator and compact cross-finding overview.
- Each `issues/ISSUE-NNN.md` owns one finding's state, evidence, drafts, and next action.
- `FORMAT.md` owns the existing record schema, evidence vocabulary, and optional upstream contribution lifecycle.
- Its upstream implementation and publication gates apply only to explicitly requested upstream work.
- Private implementation requires a scoped user instruction, not `Pull-Request-Implementation` or an upstream target.
- Preserve existing records and field meanings; do not migrate the ledger without a separate user request.
- Search the index and relevant issue records for the same symptom or root cause before adding a finding.
- Allocate new findings from `Next finding ID`; create the file, index row, and allocator increment together.
- IDs use permanent `ISSUE-NNN` values starting at `ISSUE-001`; never reuse, renumber, or partition their sequence.
- Update the issue file and index together after state, authorization, target, priority, next-action, or reference changes.
- Every issue record must follow the field and section contract in `FORMAT.md`.
- New findings start with `State: Investigating`, `Authorized-Work: Not-Selected`, and `Publication-Target: Not-Selected`.
- New findings use `External-Reference: Not published.` until an external reference exists.
- Private fixes do not imply `PR-Ready` or `Submitted`; record local scope, commits, and checks without inventing states.
- Publication candidates remain Investigating until currentness, prior art, impact, and value are evidence-backed.
- Clone detectors, AST matches, text similarity, shared names, and TODOs produce candidates only.
- A duplication finding requires shared change pressure, realistic drift, and simpler consolidation.
- Apply `skill-semantic-compression-3` when authoring or restructuring tracking content.
- Apply `skill-fork-contribution-tracking` to existing ledger mechanics and separately requested upstream work.
- This private-fork contract takes precedence over that skill's default upstream-delivery workflow.
- Run the skill's read-only validator after every ledger mutation.

### Optional upstream contributions

- Do not initiate upstream issues, comments, or pull requests unless the user explicitly requests upstream work.
- Apply `skill-maintainer-communication` before external issues, comments, reviews, discussions, or pull requests.
- Search existing upstream work and follow current contribution templates, disclosure rules, and CLA requirements.
- The user selects `Authorized-Work` and `Publication-Target`; never select them on the user's behalf.
- `Research-and-Reporting` permits research, issues, and comments but no source implementation.
- `Pull-Request-Implementation` permits only the implementation scope recorded for that finding.
- Verify upstream source claims and reproduce claimed bugs against current `upstream/master`.
- Use a separate clean branch or worktree from current `upstream/master` for each upstream implementation.
- Keep private guidance, tracking, configuration, and unrelated fork commits out of upstream contribution diffs.
- Publish one coherent root cause per issue, comment, or pull request.
- Pull-request work must verify behavior, commit, push, and reach `PR-Ready` before publication.
- Record the final external URL in `External-Reference` immediately after publication.

### External publication approval

External issues, comments, reviews, discussions, and pull requests require approval of the exact current draft and target.
Show the complete draft and target again after any change; earlier approval does not cover revisions.
Before publication, read current contribution guidance and explain applicable project policy.
The human must be able to review and own every submission statement.
Fork commits, pushes, and authorized private implementation do not require upstream publication approval.

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
- `.github/workflows/` contains inherited CI, release, and packaging automation.

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

- Run `go test ./...`, `go vet ./...`, and `go build ./...` before completing source changes.
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
- Keep changes focused; upstream contributions must also follow Dropbox's Contributor License Agreement requirements.

## Important Files

- `go.mod` owns the module path, Go language version, and dependencies.
- `CONTRIBUTING.md` owns upstream development and submission requirements.
- `docs/automation.md` owns output, error, exit-code, and authentication contracts.
- `docs/json-schema/v1/README.md` owns JSON schema stability policy.
- `.github/workflows/ci.yml` owns cross-platform validation and generated-output drift checks.
- `README.md`, `CHANGELOG.md`, `SECURITY.md`, and `LICENSE` are inherited upstream documents.

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
- Use focused existing tests while investigating; the three validation commands above gate completed source changes.
- Preserve CLI help, command metadata, generated references, and JSON schemas together when changing public commands.
- Context-only documentation changes require review, not unrelated source-test runs.
- Ledger mutations require the bundled ledger validator.
- Existing upstream quality commands retain their output behavior; do not rewrite them solely for private-fork setup.
