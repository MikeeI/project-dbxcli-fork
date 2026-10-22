# dbxcli JSON schema v1

These schemas describe command envelopes from `dbxcli --output=json` and JSON help from `dbxcli --help --output=json`.

- `success.schema.json` validates successful command responses.
- `error.schema.json` validates command error responses.
- `manifest.schema.json` validates each command manifest emitted by JSON help.
- `commands.json` lists command-specific payload definitions, result statuses, result kinds, and warning codes.
- `commands.schema.json` validates command names, input/result fields, primitive types, statuses, kinds, and warnings.

## How to validate output

Discover a command before invoking it:

```sh
dbxcli put --help --output=json > put-help.json
```

Each manifest's `input_schema` describes command arguments and flags.
Validate planned structured input, then construct CLI arguments using `x-cli-name` metadata.
Run commands requiring machine-readable results with `--output=json`:

```sh
dbxcli put --if-exists fail --output=json report.md /Reports/report.md > result.json
```

Select the response validator by its envelope:

- Validate `ok: true` with `commands.schema.json`.
- Validate `ok: false` with `error.schema.json`.
- Validate each JSON help manifest with `manifest.schema.json`.

Use schemas bundled with the installed fork version or pinned to the same fork release or commit.
Do not validate private extensions against unmodified upstream releases.
A pinned schema URL has this form, with `<fork-commit>` replaced by the installed source revision:

```text
https://raw.githubusercontent.com/MikeeI/project-dbxcli-fork/<fork-commit>/docs/json-schema/v1/commands.schema.json
```

Use mutable branch URLs only when latest development behavior is intentional.
Schema IDs are identifiers and do not replace selecting the correct versioned schema document.

## Stable envelopes

Schema v1 does not guarantee exact error-message text, exhaustive SDK enums, or the absence of later additive fields.
Stable machine contracts include `error.code`, result `status`, result `kind`, and documented field meanings.
Existing top-level fields, stable error codes, and result status meanings are not removed or renamed within v1.
New fields, commands, warning codes, and error details may be added in minor releases.
JSON responses must not include access tokens, refresh tokens, authorization codes, or app secrets.

Successful responses include:

- `ok: true`.
- `schema_version: "1"`.
- `command`, without the binary name, such as `ls` or `share-link create`.
- `input`, containing command-specific request fields.
- `results`, whose objects each include `status`, `kind`, `input`, and `result`.
- `warnings`, containing machine-actionable warnings or an empty array.

Error responses include:

- `ok: false`.
- `schema_version: "1"`.
- `command` when available, or `dbxcli` for root/pre-parse errors.
- `error.message`, containing human-readable text.
- `error.code`, containing a stable machine-readable error code.
- Optional `error.details` containing reliable structured context.
- `warnings`, containing machine-actionable warnings or an empty array.

In JSON mode, success and error envelopes go to stdout.
Human warnings, progress, diagnostics, and verbose logs go to stderr.
Error responses retain a nonzero exit status.

## Search extensions

New search producers always add a `search` summary to successful responses.
Its presence remains optional in the evolved v1 schemas so older producer responses still validate.
When the summary is present, its six fields are required:

- `limit` is the local uint64 result limit; zero means unlimited.
- `stop_reason` is `api_exhausted`, `limit`, or `error`.
- `server_has_more` is the last successful API page flag, or null before any successful page.
- `local_filters_applied` reports local metadata filtering, not server filtering.
- `api_matches_seen` counts received API matches, including duplicates and locally rejected matches.
- `api_match_limit` is 10000, the documented provider limit rather than a total-count guarantee.

`api_exhausted` means received matches were processed and the API reported no further page.
`limit` means the CLI omitted matching entries in a received page or did not fetch further pages after reaching its limit.
The last API flag can be false while local truncation occurred within that page.
API exhaustion is not proof of account-wide completeness, and there is no reliable cap-reached field.
Dropbox documents index delays, possible omitted matches, duplicates, and the 10,000-match ceiling: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/search-v-2

Search uses a dedicated metadata definition adding optional `name`, `match_type`, and `highlight_spans`.
The shared metadata payload used by other commands is unchanged.
`match_type` is an open string, not a closed list of future Dropbox enum variants.
Each highlight span requires `highlight_str` and `is_highlighted`; these are title segments, not content snippets.
Requesting highlights requires JSON output.

Search failures retain accepted entries under `error.details.partial_results` and the state under `error.details.search`.
Partial entries use the same wrapper and search metadata projection as successful results.
The original error code and exit class remain failures; partial results do not turn the whole search into success.
The error schema describes these fields through local definitions without an external schema dependency.

Older v1 replies without a summary have unknown completion state.
Do not infer summary support from `schema_version: "1"` alone.
Older closed schema files reject added fields and must be updated; tolerant readers can continue using known fields.

## Structured error details

Reusable keys have these meanings:

- `argument` identifies a positional argument related to validation or remediation.
- `arguments` lists positional arguments related to the error.
- `flag` identifies a CLI flag without leading dashes.
- `flags` lists related flags without leading dashes.
- `value` records a relevant invalid string, number, or boolean value.
- `path` identifies a directly relevant Dropbox or local path.
- `revision` identifies a related Dropbox revision.
- `email` identifies the related address.
- `member_id` identifies a selected Dropbox team member, such as `--as-member`.
- `from_path` and `to_path` identify relocation source and destination paths.
- `url` identifies a related shared-link or API URL.
- `operation` identifies an operation such as `upload`, `download`, `delete`, `restore`, or `share_link_create`.
- `token_type` identifies the related token kind.
- `login_command` gives the suggested authentication remediation command.
- `env_var` names the relevant environment variable.
- `api_summary` preserves an available Dropbox error summary.
- `api_endpoint` names a Dropbox endpoint extracted from an SDK error message.
- `bytes_written` records output confirmed before a failed stdout transfer.
- `retry_after_seconds` records the suggested delay for a rate-limited request.
- `search` and `partial_results` preserve failed-search state and accepted results as described above.

Prefer existing path keys over synonyms: `path` for one relevant path and `from_path`/`to_path` for relocations.
Add keys such as `target` or `local_path` only when the existing names would be ambiguous.

## Generation and supported commands

`login` and `completion` intentionally lack normal structured command-result JSON.
Their help remains available as JSON manifests, including root, version, and `help version` discovery.
`dbxcli --output=json` without `--help` still prints text help; shell-completion protocols remain text-only.

The current JSON-enabled command paths are listed in `commands.json`.
The catalog is maintained input, not generated output.
`commands.schema.json` is generated from it and typed metadata in `internal/jsonschema/`:

```sh
go run ./tools/gen-json-schemas
```

The generator does not rewrite `success.schema.json`, `error.schema.json`, or `commands.json`.
Command-specific schemas constrain allowed fields, stable required fields, primitive types, nested objects and enums.
They avoid over-modeling Dropbox SDK enum values not owned by dbxcli.
Golden fixtures under `cmd/testdata/json_contract/` lock command input/result contracts.

## Command Manifest v1

JSON help is the canonical command manifest surface:

```sh
dbxcli --help --output=json
dbxcli put --help --output=json
dbxcli --output=json help put
```

Each help result has `status: "described"`, `kind: "command"`, and a result validating against `manifest.schema.json`.
Manifest v1 preserves original help fields and adds:

- `manifest_version: "1"` and structured `args`.
- Enriched `flags` with enums, conflicts, prompt behavior, and sensitive-input metadata.
- Generated `input_schema` with bounded-value enums.
- `examples` and `schema_refs`.
- Best-effort audited `dropbox_scopes` and `stdin_stdout`.
- `result_statuses`, `result_kinds`, and `warning_codes`.

`input_schema` supports tools, MCP-style integrations, and planners validating inputs before invocation.
It excludes `--help` and `--output`, maps hyphens to underscores, and retains original names in `x-cli-name`.
No-auth commands may omit globally accepted flags that do not affect them.
Sensitive properties use `writeOnly` and `x-sensitive`; conflicts use `x-conflicts`.
StringSlice/StringArray flags are arrays with string items; enum constraints apply to items.
Optional presence-sensitive size bounds have no schema default, even though Cobra's technical zero is visible in flags.
Omit an unselected bound instead of materializing zero; explicitly supplied zero is preserved in execution input.

For structured commands, `schema_refs.command_success_schema` identifies a definition in `commands.schema.json`.
`schema_refs.command_contract` identifies the source catalog entry in `commands.json`.
`scope_accuracy: audited_best_effort` describes audited planning metadata, not a substitute for runtime API errors.

Account results contain auth context under `results[0].result.auth`:

- `source` is `saved` or `env`.
- `refreshable` is a boolean.
- `auth_file` is `default`, `custom`, or `none`; the full auth-file path is not included by default.

Warnings have a stable `code` and human-readable `message` and may contain optional command-specific context.
Current codes include `deprecated_command` for deprecated paths and `skipped_symlink` for recursive uploads.
`logout` may return `token_revoke_failed` after local removal when remote revocation failed.

## Stable error codes

- `invalid_arguments`: Invalid command arguments or flags.
- `path_conflict`: A local or Dropbox path conflicts with the requested operation.
- `auth_required`: No usable saved credentials, or Dropbox rejected the token.
- `auth_refresh_failed`: Saved refreshable credentials could not be refreshed.
- `app_key_required`: Login or refresh requires a Dropbox app key.
- `auth_exchange_failed`: OAuth exchange failed or returned unusable tokens.
- `not_found`: Dropbox reported that the requested object was not found.
- `partial_transfer`: A download-to-stdout stream failed after partial output.
- `permission_denied`: Permissions, scope, member selection, or account state denied access.
- `rate_limited`: Dropbox rate limited the request.
- `dropbox_api_error`: An API failure without a more specific mapped code.
- `env_token_still_active`: `DBXCLI_ACCESS_TOKEN` must be unset before logout can complete.
- `structured_output_unsupported`: The command does not support normal `--output=json` execution.
- `unknown_command`: Cobra could not resolve the command.
- `unknown_flag`: Cobra could not resolve a flag.
- `command_failed`: Fallback when no more specific stable code applies.
