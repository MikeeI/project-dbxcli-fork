# ISSUE-040 — Uploads: unveränderte Chunkdaten werden bei Retries erneut gehasht

State: Archived
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Low
Root-Cause-Confidence: High
Finding-Category: Performance
Created: 2026-10-23
Updated: 2026-10-23
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: `uploadOneChunk` übergibt bei jedem Versuch dieselben unveränderten Chunkdaten mit leerem `ContentHash`.
Das SDK ergänzt den Hash nur in einer Argumentkopie und berechnet ihn deshalb beim nächsten Versuch erneut.
Gesprächsreferenz: `P2` des Performance-Reviews; die Einstufung wurde mangels belegter Nutzungshäufigkeit auf Low korrigiert.
Der private Ausgangsstand dieser Umsetzung ist `personal@d0e5aab1952eb5765d1774c08d4b33fa8249d140`.

## Reach-and-Impact

Reach [S]: Chunk-Uploads oberhalb 32 MiB verwenden `uploadOneChunk` einschließlich vorhandener Retry-Policy.
Impact [S]: Jeder zusätzliche Versuch hasht die gesamte Payload erneut; Standardchunks umfassen 16 MiB.
Ohne Retry entfällt kein Hashdurchlauf; Häufigkeit und reale Transferbeschleunigung sind nicht gemessen.

## Evidence

- [S] `cmd/put.go:59-76` besitzt nach der privaten Korrektur eine lokale Argumentkopie mit wiederverwendetem Hash.
- [S] SDK v6.6.2, `files/content_hash.go:131-141`, ergänzt ausschließlich die SDK-interne Argumentkopie.
- [S] SDK v6.6.2, `files/content_hash.go:170-195`, berechnet nur bei fehlendem Hash den Inhalt erneut.
- [O] Ein temporärer Offline-Smoke mit echtem SDK bestätigt drei Versuche mit identischem Hash und identischen 16-MiB-Payloads.

## Performance-Evidence

Cost model [S]: Bei `r` Retries entfallen nach der Korrektur `r × Chunkgröße` zusätzliche Hash-Eingabebytes.
[O] Vorher wurden bei drei Versuchen 50331648 Byte automatisch im SDK gehasht; nachher waren es null.
[S] Stattdessen berechnet der CLI-Owner die 16777216 Byte einmal beim ersten zugelassenen Versuch.
[O] Explizite Caller-Hashes bleiben erhalten; abgebrochene Kontexte verursachen keinen weiteren Append.
Measurement status: SDK-Bytezähler unter Go 1.27.3 auf Linux/amd64; keine Live-Transfers oder gemessene Ende-zu-Ende-Ersparnis.

## Prior-Art

Coverage: Aktueller Index und offene Issue-Akten wurden auf Hash-, Prüfsummen- und ContentHash-Ursachen geprüft.
`ISSUE-001` betrifft Lesepuffer; `ISSUE-008` betrifft die Abschlussreihenfolge und fordert lediglich Erhalt des Hashings.
Gaps: Externe Issues, PRs und Releases wurden nicht recherchiert.
Contribution fit: Privater Fix ohne Upstream-Veröffentlichungsauftrag.

## Proposed-Change

Den Hash beim ersten zugelassenen Versuch mit dem vorhandenen SDK-`contenthash.Compute` berechnen.
Die lokale Argumentkopie behält ihn ausschließlich für die Lebensdauer dieser unveränderten Chunk-Payload.
Der Nutzer autorisierte die private Umsetzung der vier Review-Befunde mit „fix all issues“.
Das wählt weder einen Upstream-Arbeitsmodus noch ein Veröffentlichungsziel.

## Scope-and-Constraints

- Preserve: Bestehende Hashwerte, unveränderte Payloads, Caller-Argumente, Kontextprüfung, Backoff und Offset-Recovery erhalten.
- Preserve: SDK-Integritätsprüfung bleibt aktiv; kein `WithoutAutoContentHash` und kein eigener Hashalgorithmus.
- Exclude: Keine neuen Abhängigkeiten, kein globaler Cache und keine Änderung des Worker-Lifecycles.
- Related: `ISSUE-007`, `ISSUE-008` und `ISSUE-014` bleiben eigenständige, nicht durch diesen Fix erledigte Ursachen.

## Verification

[O] Vorhandene `TestUploadOneChunk_*`- und `TestUploadChunked_*`-Prüfungen bestehen nach der Korrektur.
[O] `TestUploadOneChunk_RetryOnServerError` wurde um identische Retry-Payloads, Hashwerte und unveränderte Caller-Argumente erweitert.
Test decision: update; der bestehende Test schützt den Integritätsvertrag wiederholter Upload-Requests.
[O] Der temporäre SDK-Smoke sowie Prüfungen expliziter Hashwerte und abgebrochener Kontexte bestehen und wurden entfernt.
[O] Commit `5843c2f` auf `personal` wurde nach `origin/personal` gepusht.
[O] `go test ./...`, `go vet ./...`, `go build ./...` und `go test -race ./...` bestehen.
[O] `golangci-lint run --fix ./...` meldet ausschließlich den unveränderten ungenutzten Helfer `cmd/json_metadata.go:64`.

## Publication-Blockers

Externe Prior Art, aktuelle Upstream-Verifikation, Veröffentlichungsziel und genaue Veröffentlichungsfreigabe fehlen.
Die private Umsetzung autorisiert keine externe Einreichung.

## Next-Action

Summary: —
Action: None.
Done-When: None.

## Archive

Archive-Reason: Other
Detail: Privat autorisierte Retry-Hash-Optimierung implementiert, geprüft und auf origin/personal gepusht; keine Upstream-Erledigung behauptet.
Evidence: `5843c2f`; SDK-Hashzähler vor/nach Korrektur, Retry-Vertragsprüfung und vorgeschriebene Go-Prüfungen.
Checked: 2026-10-23
