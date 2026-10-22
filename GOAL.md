# Goal: Agent-First Search for dbxcli

## Purpose

The maintainer of this fork is an external contributor, not an upstream `dropbox/dbxcli` developer.
The goal is a better `dbxcli search` for LLM agents first and humans second.
Agents integrate the CLI and iterate faster and more precisely than humans, so they must not misread incomplete results.

This file records goals and priorities only.
It does not authorize implementation or upstream publication.
Each goal still needs a finding in `ISSUES.md` and an explicit `Authorized-Work` selection.

## Scope Boundary

- Use only `files/search_v2` and `files/search/continue_v2` plus the pinned Go SDK (`dropbox-sdk-go-unofficial/v6@v6.6.2`).
- Do not add new endpoints.
- Do not use Android-internal endpoints or fields.
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
- These are title highlights, not content snippets like the Android app's.

### 4. Merge duplicates across pages (small)

- Merge matches across pages by file ID.
- The official documentation explicitly warns that duplicates may appear across pages.

## Acceptance Direction

- An agent can distinguish complete, limit-truncated, and error-truncated results from the output alone.
- No result collected before a failure is discarded.
- Filters apply server-side, before the result limit.
- Existing text output and JSON v1 consumers keep working.

## Non-Goals

- Android-only server filters such as date range, `reverse_order`, and snippets; they are unverified for a public app key.
- DCS, Ask, suggestions, and connector search.
- Regex or boolean query languages over document content; no API contract supports them.
- Property search; it would add a new endpoint family.

## Known Limits

- Official docs cap one search plus continuations at 10,000 matches.
- Indexing may lag, and some results may not be returned.
- The CLI cannot promise complete result lists.

## Evidence

- Reverse-engineering research lives in `~/projects/REVERSE/project-reverse-android-dropbox`.
- Main comparison: `context/findings-search-comparison.md` in that project.
- That research is static Android analysis plus public API documentation; no live Dropbox requests were made.
- Treat Android behavior as evidence for ideas, never as a promise of backend acceptance.
