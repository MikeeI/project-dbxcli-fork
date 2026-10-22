---
name: dbxcli
description: Safely operate Dropbox through a locally installed dbxcli command, using its JSON manifest and schema-backed machine contract. Use for Dropbox file, shared-link, team, or account work; do not call the Dropbox API directly.
---

# dbxcli

Use the local `dbxcli` executable as the only Dropbox integration.
Do not reimplement Dropbox API calls, scrape text help, or maintain a command catalog in this skill.
The CLI's JSON help manifest is authoritative for the installed version.

## Start safely

1. Check availability with `command -v dbxcli`, then run `dbxcli version --output=json`.
   If unavailable, report that and give operating-system-specific [installation guidance](https://github.com/dropbox/dbxcli/releases); do not install without authorization.
2. Before a command not already discovered in the current task, run `dbxcli [command path] --help --output=json`.
   Begin with `dbxcli --help --output=json` when the command path is unknown.
3. Inspect `supports_structured_output`, `input_schema`, `stdin_stdout`, `destructive_level`, `flags`, and `args`.
   Do not infer commands or flags from memory.
4. Validate intended JSON-shaped arguments and flags against `input_schema` before constructing an invocation.
   Map fields through `x-cli-name`; preserve arrays and omit unset optional filters instead of filling technical defaults.

Read [tool-integration.md](references/tool-integration.md) for discovery, validation, result, and error handling.
Read [automation.md](references/automation.md) for writes and confirmation behavior.
Read [safety.md](references/safety.md) before handling credentials, transfers, deletion, replacement, or sharing.

## Execution contract

For normal execution, always pass `--output=json` and parse stdout as one JSON envelope.
Treat stderr as diagnostics and check both the process exit code and `.ok`.

- For `.ok: true`, use documented `results[].status`, `results[].kind`, and `warnings`, not prose or undocumented fields.
- For `.ok: false`, branch on stable `.error.code`, not `.error.message`, and use relevant structured `.error.details`.
- Do not blindly retry writes or authentication failures; surface a concise explanation without sensitive values.
- When `supports_structured_output` is false, explain the limitation or choose a safe supported alternative.
- JSON help remains available even when normal structured command output is unsupported.

For destructive or externally visible actions, discover and validate first, then prefer an available `--dry-run`.
Use an explicit available `--if-exists` policy rather than assuming overwrite or conflict behavior.
A successful dry-run does not authorize execution; obtain explicit user confirmation before destructive execution.
Use `--yes` only after that confirmation when the discovered command exposes it.

## Large listings, search, and multi-step work

The CLI follows Dropbox pagination internally; agents must not invent or pass cursors.
For broad listings or search, discover the command and bound output with the available `--limit`.
Start with a narrow folder or search path instead of recursively enumerating an entire Dropbox.
A limited result is a selection, not proof that no additional matches exist.

For search requests about text inside files, use discovered `--content` only when content search was requested.
Otherwise search remains filename-oriented.
Server filters narrow the query; local name, size, time, and type filters only inspect its returned metadata.
Local filters do not turn a query into an exhaustive account inventory and may require many API pages.

For successful search, inspect `.search.stop_reason` when present:

- `api_exhausted` means the API cursor ended, not that every matching account file was found.
- `limit` means the local result limit stopped collection; unvisited pages may contain no further local matches.
- `.search.server_has_more` is the last API page state, not a standalone completeness or resume signal.
- `.search.api_matches_seen` includes rejected and duplicate API matches; `.search.api_match_limit` reports the API cap.
- If `.search` is absent, completeness is unknown, even when `schema_version` is `"1"`.

On a search failure, retain `.error.details.partial_results` when present and inspect `.error.details.search`.
These entries are usable partial results, not a successful complete search; the error code and exit status still matter.
Deduplication retains the first locally qualifying match per file or folder ID, not a consistent account snapshot.
Title highlight spans are not content snippets, and requesting them requires JSON output.
Omitting a size filter differs from explicitly setting zero; do not invent an upper bound of zero.

For search → get → process, discover search, select an exact metadata path, and then discover get.
Download to a named local file, never stdout, process it, and report the output path or concise result.
A later upload, share-link, or replacement needs separate discovery, safety policy, and authorization.

## Boundaries

- Never expose tokens, auth codes, app secrets, environment dumps, or auth-file contents in prompts or artifacts.
- Do not put these values in commands, logs, JSON fixtures, tool results, or commits; refer to names and paths only.
- Do not use `DBXCLI_ACCESS_TOKEN=value` inline; use an already-provisioned execution environment.
- Keep `DBXCLI_AUTH_FILE` outside the repository and do not read, upload, or commit it.
- Never send binary file data through a tool result; download to a local file and report its path and metadata.
- The local operand `-` is a byte stream and cannot be combined with JSON output for downloads.
- Do not expose links or alter permissions, overwrite, move, restore, or delete without user-authorized scope.
- Validate against schemas distributed with the installed fork version, not upstream schemas missing private extensions.
