# Machine-contract integration

The installed CLI, not this reference, owns command discovery.
JSON help works without Dropbox authentication:

```sh
dbxcli --help --output=json
dbxcli put --help --output=json
dbxcli share-link create --help --output=json
```

Inspect `results[].result.supports_structured_output` before normal execution.
Its `input_schema` describes positional arguments and flags with JSON-friendly names such as `if_exists`.
Use each property's `x-cli-name` when constructing the invocation.
Validate input first, preserve string-array values, and omit optional filters that the user did not select.
A missing size bound is not zero; an explicitly supplied maximum of zero requests empty files.

Run supported operations with `--output=json`.
Stdout contains exactly one success or error envelope; stderr may contain diagnostics, progress, and warnings.
Never parse text output as a fallback.

```sh
dbxcli ls --output=json /
dbxcli put --if-exists fail --output=json report.md /Reports/report.md
```

Success includes `ok: true`, `schema_version`, `command`, `input`, `results`, and `warnings`.
Use documented result statuses and kinds as stable outcomes.
Errors contain `ok: false` and `error`; check the exit status and stable `error.code`, not human-facing message text.
Known error details provide relevant context, not permission to expose sensitive values.

Exit-code classes include auth (2), permission (3), not found (4), conflict (5), rate limit (6), and validation (7).
Partial stdout transfer uses exit code 8.
A rate-limit error may provide `error.details.retry_after_seconds`.
Wait or retry only while the task remains safe; never automatically retry non-idempotent or destructive operations.

## Search results and partial failures

The CLI owns pagination; do not invent cursor loops in an agent.
The discovered `--limit` bounds returned entries, not the number of API matches inspected or account completeness.
Local filters and accepted-ID deduplication run before that limit; local sorting affects only the selected results.
A local date or regex filter only narrows matches from the supplied Dropbox query, not the entire account.
Use discovered `--content` only for requested content searches and prefer a narrow optional path scope.

New successful search responses include a `search` object:

- `stop_reason: api_exhausted` means the API returned no further page after processing the received matches.
- `stop_reason: limit` means a matching page remainder or unvisited pages were excluded by the local limit.
- `server_has_more` is the last successful API page flag, which can be false despite local truncation within that page.
- `local_filters_applied` identifies local metadata filtering, not server filters.
- `api_matches_seen` counts raw API matches, including duplicates and locally rejected entries.
- `api_match_limit` records the documented provider cap, not a reliable total or a server cap-reached signal.

If the summary is absent in an older v1 response, treat the completion state as unknown.
The same `schema_version` does not establish support for every later optional field.
API exhaustion does not prove complete account coverage because indexing can lag and matches can be omitted.

On failure, `error.details.search` contains `stop_reason: error` and the last known server state.
Before any successful page, `server_has_more` is null.
`error.details.partial_results` retains accepted entries in the same wrapper and metadata form as success results.
Use those entries without interpreting the overall operation as successful or ignoring its exit code.

Search results may add `name`, an open-string `match_type`, and native title `highlight_spans`.
Highlights contain `highlight_str` and `is_highlighted`; they are not content snippets or character offsets.
The highlight request requires `--output=json`.
Duplicate IDs retain the first locally qualifying full match rather than synthesizing metadata from different revisions.

## Schema selection

Use `manifest.schema.json`, `commands.schema.json`, and `error.schema.json` from the installed fork's version.
Prefer a pinned fork release or commit, or its bundled local schema files, over mutable remote `master` URLs.
Unmodified upstream schemas do not describe private search extensions.

The evolved v1 schemas accept older replies without optional search fields.
Older closed schemas reject newly added fields and must be updated alongside the consuming integration.
Error schemas remain self-contained and do not require fetching a separate command schema.
