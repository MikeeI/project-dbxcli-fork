# Goal: Agent-First Search for dbxcli

## Purpose

The maintainer of this fork is an external contributor, not an upstream `dropbox/dbxcli` developer.
The goal is a better `dbxcli search` for LLM agents first and humans second.
Agents integrate the CLI and iterate faster and more precisely than humans, so they must not misread incomplete results.

This file records goals and priorities only.
It does not authorize implementation or upstream publication.
Each goal still needs a finding in `ISSUES.md` and an explicit `Authorized-Work` selection.

## Contribution Intent

These goals are for this fork only.
Upstream `dropbox/dbxcli` is an official Dropbox project, and upstream interest in these changes is unlikely.
Goals are implemented on the fork's `personal` branch and are not prepared as upstream pull requests.
Upstream publication happens only if the user later decides so and approves the exact draft and target.
Staying close to upstream code and conventions remains useful so that merging future upstream releases stays cheap.

## Scope Boundary

- Use only `files/search_v2` and `files/search/continue_v2` plus the pinned Go SDK (`dropbox-sdk-go-unofficial/v6@v6.6.2`).
- Do not add new endpoints.
- Keep the JSON v1 success and error envelopes stable; add fields compatibly.
- Regenerate command docs and JSON schemas when help or JSON metadata changes.
- Flag names below are proposals, not final contracts.

## Goals in Priority Order

### 1. Make truncated searches detectable (small)

- Output `has_more`, the local limit, and the reason the search stopped.
- When a continuation request fails, still output the results collected so far and report the error with them.
- Today an agent cannot tell "no more matches" from "limit reached".

### 2. Add category and extension as server-side filters (small)

- Proposed flags: `--category pdf` and `--ext pdf,docx`.
- They set `file_categories` and `file_extensions` in `SearchOptions`.
- Both fields already exist in the SDK and are currently never set.
- Official docs support these filters only for active-file search.

### 3. Output match type and title highlights (small)

- Keep `match_type` when converting matches to output metadata.
- Optionally request title highlights through `match_field_options.include_highlights`.
- The agent can then see whether a match came from the name, the content, or both.
- These are title highlights only, not content snippets.

### 4. Merge duplicates across pages (small)

- Merge matches across pages by file ID.
- The official documentation explicitly warns that duplicates may appear across pages.

### 5. Filter returned metadata locally (small to medium)

- Filter results the server cannot filter, using only fields already present in the returned metadata.
- Name: glob, regex, contains, and exclusions, for example everything except names containing "Draft".
- Size: minimum and maximum.
- Modified time: range on `server_modified` or `client_modified`.
- Entry type: files only or folders only, from the `.tag` field.
- The result limit counts only entries that pass the filters.
- Collect further pages until enough matching entries exist or the search ends or hits the 10,000-match cap.
- Never filter an already limited result list; that silently loses matches on later pages.
- Local filters save no server work, and a local date filter is not a server-side date filter.
- Output must say that local filters were applied and whether the search ended before the limit was met.

## Acceptance Direction

- An agent can distinguish complete, limit-truncated, and error-truncated results from the output alone.
- No result collected before a failure is discarded.
- Server-side filters apply to the request; local filters apply before the result limit counts entries.
- Existing text output and JSON v1 consumers keep working.

## Non-Goals

- Server-side date range, reverse ordering, and content snippets; the pinned public `SearchOptions` has no such fields.
- Other search surfaces beyond `files/search_v2`, such as suggestions or connector search.
- Regex or boolean query languages over document content; no API contract supports them.
- Property search; it would add a new endpoint family.

## Known Limits

- Official docs cap one search plus continuations at 10,000 matches.
- Indexing may lag, and some results may not be returned.
- The CLI cannot promise complete result lists.
