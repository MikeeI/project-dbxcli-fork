# Automation and JSON output

`dbxcli` provides text output, structured JSON execution output, and JSON help for machine-readable discovery.
In JSON mode, result and error envelopes go to stdout, while human diagnostics and progress go to stderr.
The `warnings` array contains machine-actionable warnings.

## Which interface should I use?

For ordinary scripts and CI jobs, use `--output=json` and check the exit code, `ok`, and stable result/error fields.
Basic invocation does not require schema validation.
Wrappers and agents needing discovery or pre-invocation validation should use JSON manifests and `input_schema`.

## JSON command output

Text is the default; supported commands accept the global `--output` flag:

```sh
dbxcli ls --output=json /
dbxcli account --output=json
dbxcli logout --output=json
```

Discover structured-output support before relying on it:

```sh
dbxcli put --help --output=json
```

Successful responses use a stable envelope:

```json
{
  "ok": true,
  "schema_version": "1",
  "command": "ls",
  "input": {
    "path": "/Reports"
  },
  "results": [
    {
      "status": "listed",
      "kind": "file",
      "input": {},
      "result": {
        "type": "file",
        "path_display": "/Reports/q1.pdf",
        "path_lower": "/reports/q1.pdf"
      }
    }
  ],
  "warnings": []
}
```

Errors go to stdout as a single envelope and retain a nonzero exit status:

```json
{
  "ok": false,
  "schema_version": "1",
  "command": "rm",
  "error": {
    "message": "path exists and is not a folder: /old-file.txt",
    "code": "path_conflict"
  },
  "warnings": []
}
```

## Search selection, state, and partial results

Search uses Dropbox `files/search_v2` and `files/search/continue_v2`, not a recursive listing or local index.
The query always goes to Dropbox; optional local filters only narrow the metadata returned for that query.
Use a path scope, server filters, a result limit, and an explicit timeout when practical:

```sh
dbxcli --timeout 30s search report /Reports --ext pdf,docx --category document --limit 20 --output=json
dbxcli search report /Reports --name-glob '*.pdf' --exclude-name-contains Draft --min-size 0 --output=json
dbxcli search report /Reports --modified-from 2026-01-01T00:00:00Z --modified-before 2027-01-01T00:00:00Z --output=json
```

`--category` and `--ext` are native comma-separated, repeatable string-slice flags sent to Dropbox.
Discover supported categories through JSON help or completion; the catch-all category is `others`, not `other`.
Extensions are trimmed, lowercased, and stripped of one optional leading dot; blanks, globs, and paths are invalid.
Do not assume special multi-extension behavior beyond the server's extension filter.
`--content` changes filename-only search to filename-and-content search without disabling server filters.
Search retains the SDK's active-file status default.

Local predicates combine with AND:

- `--name-glob` matches the entire basename with Go `path.Match`; `**` has no recursive traversal meaning.
- `--name-regex` uses Go regexp search semantics; use anchors for full-name matching or `(?i)` for case-insensitivity.
- `--name-contains` performs a case-sensitive literal substring match.
- Repeated `--exclude-name-contains` values reject any containing name; commas are literal and empty values are invalid.
- `--min-size` and `--max-size` are inclusive byte bounds for files, including an explicitly supplied zero.
- `--modified-from` is inclusive and `--modified-before` exclusive; both require RFC3339 timestamps with a timezone.
- `--time server|client` chooses the modified timestamp used for local date filtering and existing time presentation.
- `--type file|folder` selects a metadata kind; omitting it allows both.

Folders do not satisfy file size or modification-time bounds.
Combining `--type folder` with those bounds is invalid instead of fabricating zero sizes or timestamps.
Local filtering uses SDK-sized pages even with a small result limit, so rejected entries do not cause tiny requests.
It can inspect many matches before finding enough accepted results; `--limit` is not an API-work bound.

The CLI filters each received match, retains the first qualifying full match per metadata ID, and then applies its limit.
It sorts only that accepted selection, not all possible server matches.
Duplicates do not consume the result limit, and a rejected earlier revision does not hide a later qualifying match.
Missing IDs or malformed metadata fail with accepted partial results rather than silently discarding entries.

New successful search responses add this summary alongside the normal envelope:

```json
{
  "search": {
    "limit": 20,
    "stop_reason": "limit",
    "server_has_more": true,
    "local_filters_applied": true,
    "api_matches_seen": 100,
    "api_match_limit": 10000
  }
}
```

This fragment illustrates the summary, not a complete response.
Interpret the fields together:

- `limit: 0` means no local result limit.
- `stop_reason: api_exhausted` means received matches were processed and the last API page reported no more pages.
- `stop_reason: limit` means matching page leftovers or unvisited API pages were excluded after reaching the limit.
- `server_has_more` records the last successful page flag, not whether the CLI returned every qualifying page entry.
- `local_filters_applied` reports local predicates only; server category/extension filters do not set it.
- `api_matches_seen` counts raw matches from received pages, including rejected and duplicate entries.
- `api_match_limit` reports the provider's 10,000-match ceiling, not a total count or a cap-reached signal.

Reaching exactly the limit on a fully consumed final page is `api_exhausted`, not `limit`.
Omitting a qualifying match from that final page is `limit` even when `server_has_more` is false.
Unvisited pages may contain no further locally qualifying entries; the CLI does not perform a proof-only lookahead.
Neither API exhaustion nor a result count below 10,000 proves account-wide completeness.
Dropbox documents indexing delays, omitted matches, duplicates, and the search ceiling: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/search-v-2

Search result metadata adds `name` and, when supplied, `match_type` and `highlight_spans`.
`match_type` is an open string; do not assume only filename or file-content variants exist.
Request native title spans with `--highlights --output=json`; highlights in text mode return `invalid_arguments`.
Each span has `highlight_str` and `is_highlighted`; these are not content snippets, offsets, or HTML.
Missing server match information is omitted rather than inferred.

If search fails, stdout contains one error envelope, not a preceding success envelope or JSON fragment.
`error.details.search` records `stop_reason: error`; its server flag is null before any successful page.
`error.details.partial_results` contains accepted results in exactly the same per-result shape as success.
The original error classification and exit class remain authoritative even when partial results are useful.
Text mode prints accepted rows on stdout and a deterministic `search key=value ...` summary on stderr before the error.

The evolved v1 schemas accept older replies without `search`; absence means unknown search state.
Older closed schema files reject the additive fields and must be updated with their consumer.
Use schemas from the matching fork version rather than assuming upstream schemas describe private extensions.

## Schema-first tool integration

Use the CLI and matching schema files as the stable discovery and validation interface:

- Run `dbxcli --help --output=json` for command discovery.
- Validate planned arguments and flags against each manifest's `input_schema`.
- Validate successes with `commands.schema.json` and failures with `error.schema.json`.
- Pin schema files to the installed fork release or commit for reproducible integration.

The command catalog, error codes, and schemas are documented in [JSON schema v1](json-schema/v1/README.md).
Command-specific success validation covers input, results, primitive types, statuses, kinds, and warning codes.
There is no separate machine protocol; invoke the CLI and parse stdout JSON rather than scraping text output.

## JSON help manifest

JSON help is separate from execution output, requires no Dropbox authentication, and works for unsupported execution modes:

```sh
dbxcli --help --output=json
dbxcli put --help --output=json
dbxcli share-link create --help --output=json
dbxcli --output=json help share-link create
```

It discovers command paths, args, flags, input schemas, aliases, auth modes, destructive levels, and stdin/stdout behavior.
It also exposes schema references and whether normal structured execution is supported.
Check `results[].result.supports_structured_output`; false may yield `structured_output_unsupported` during execution.

Each `input_schema` uses JSON-friendly names such as `if_exists` and retains CLI spellings in `x-cli-name`.
Bounded values have enum constraints; string-slice and string-array flags use arrays with string items.
Optional size bounds have no semantic default: omit unset bounds instead of filling a technical flag default of zero.
An explicitly supplied zero remains present in execution input.

## Safe scripting patterns

Prefer explicit automation policies:

```sh
dbxcli put --if-exists fail report.md /Reports/report.md
dbxcli put --if-exists skip report.md /Reports/report.md --output=json
dbxcli share-link create /Reports/report.md --output=json
```

Use JSON when stable statuses, result kinds, warnings, or error codes matter.
Use text for human terminal workflows or commands intentionally writing file bytes to stdout.
Set an explicit timeout to avoid indefinite network operations in CI:

```sh
dbxcli --timeout 2m ls --output=json /
```

`--timeout` uses Go durations such as `30s`, `2m`, or `1h`; its default of zero disables the command deadline.
For a corporate or local proxy, set the standard Go HTTP environment variables:

```sh
export HTTPS_PROXY=http://proxy.company.example:8080
export NO_PROXY=localhost,127.0.0.1,.company.example
dbxcli ls --output=json /
```

`HTTPS_PROXY`, `HTTP_PROXY`, and `NO_PROXY` apply to API and OAuth exchange/refresh requests.
Go also supports lowercase forms such as `https_proxy` and `no_proxy`.
The login browser authorization step is outside dbxcli and may need separate browser or operating-system proxy settings.
Check auth and identity before running a job:

```sh
dbxcli account --output=json
```

Use direct-token authentication for short-lived CI jobs:

```sh
DBXCLI_ACCESS_TOKEN="$DROPBOX_ACCESS_TOKEN" dbxcli ls --output=json /
```

Use a dedicated saved-auth file when a job needs refreshable credentials.
Store it in a private temporary or secret-backed path rather than the repository:

```sh
export DBXCLI_AUTH_FILE="${RUNNER_TEMP:-/tmp}/dbxcli-auth.json"
dbxcli login
dbxcli ls /
```

Do not commit, cache, or upload this file as an artifact because it can contain refresh tokens.
Capture JSON and handle failures by exit status:

```sh
if ! dbxcli put --if-exists fail --output=json report.md /Reports/report.md >result.json; then
  jq -r '.error.code + ": " + .error.message' result.json >&2
  exit 1
fi

jq -r '.results[] | [.status, .kind] | @tsv' result.json
```

This GitHub Actions example assumes dbxcli is already installed:

```yaml
jobs:
  publish-report:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v7
      - name: Upload report to Dropbox
        env:
          DBXCLI_ACCESS_TOKEN: ${{ secrets.DROPBOX_ACCESS_TOKEN }}
        run: |
          dbxcli account --output=json
          dbxcli put --if-exists fail --output=json report.md /Reports/report.md
```

## Authentication for automation

Saved OAuth credentials default to:

```text
~/.config/dbxcli/auth.json
```

Save refreshable credentials with `dbxcli login`.
Check auth and identity with `dbxcli account --output=json`.
Account auth metadata is under `results[0].result.auth`; this fragment shows saved authentication:

```json
{
  "auth": {
    "source": "saved",
    "refreshable": true,
    "auth_file": "default"
  }
}
```

Stable auth fields are `source` (`saved` or `env`), `refreshable` (boolean), and `auth_file` (`default`, `custom`, or `none`).
Supply short-lived direct access tokens through `DBXCLI_ACCESS_TOKEN`:

```sh
DBXCLI_ACCESS_TOKEN=sl.xxxxxx dbxcli ls --output=json /
```

Direct tokens are not saved or refreshed; expired tokens must be replaced by the caller.
The relevant result fragment is:

```json
{
  "auth": {
    "source": "env",
    "refreshable": false,
    "auth_file": "none"
  }
}
```

Select another saved credentials file with `DBXCLI_AUTH_FILE`:

```sh
DBXCLI_AUTH_FILE=/path/to/auth.json dbxcli login
DBXCLI_AUTH_FILE=/path/to/auth.json dbxcli ls /
```

The relevant result fragment is:

```json
{
  "auth": {
    "source": "saved",
    "refreshable": true,
    "auth_file": "custom"
  }
}
```

`dbxcli logout` revokes saved tokens and removes saved credentials.
It cannot remove environment-provided tokens; unset `DBXCLI_ACCESS_TOKEN` before logout when it is set.

## Stdin and stdout

Use `-` as the local operand for byte streams:

```sh
printf 'hello' | dbxcli put - /hello.txt
tar cz ./src | dbxcli put - /backups/src.tgz
dbxcli get /backups/src.tgz - | tar tz
dbxcli get /file.txt - > local-copy.txt
```

Stdin uploads spool the full input to a temporary file and require corresponding disk space.
Stdout downloads keep bytes clean by sending diagnostics, progress, warnings, and verbose logs to stderr.
Byte-output commands cannot simultaneously put result JSON on stdout.
`dbxcli get <path> - --output=json` and `dbxcli share-link download <url> - --output=json` return `invalid_arguments`.
A bare `-` is a byte-stream operand only in the local position; Dropbox paths named `-` remain valid.
For example, `dbxcli put - /-` and `dbxcli get /- -` are valid, while local files named `-` use `./-`.

## Exit status

Text and JSON modes share stable exit codes.
Success is zero, including successful commands containing warnings.

| Exit code | Meaning | JSON error codes |
| --- | --- | --- |
| `0` | Success | none |
| `1` | Generic failure | `command_failed`, `dropbox_api_error` |
| `2` | Authentication failure | `auth_required`, `auth_refresh_failed`, `auth_exchange_failed`, `app_key_required`, `env_token_still_active` |
| `3` | Permission denied | `permission_denied` |
| `4` | Not found | `not_found` |
| `5` | Conflict | `path_conflict` |
| `6` | Rate limited | `rate_limited` |
| `7` | Validation or usage | `invalid_arguments`, `unknown_command`, `unknown_flag`, `structured_output_unsupported` |
| `8` | Partial stdout transfer | `partial_transfer` |

For failures, inspect both the exit code and `error.code` for the most specific remediation.

## Shell completion

Completion scripts and protocol output are text-only because shells require those formats:

```sh
dbxcli completion bash
dbxcli completion zsh
dbxcli completion fish
dbxcli completion powershell
```

Passing `--output=json` to completion commands returns `structured_output_unsupported`.
