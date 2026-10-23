# ISSUE-030 — Gelöschte Dateien: neun ungenutzte Revisionen werden mitabgerufen

State: Archived
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Low
Root-Cause-Confidence: High
Finding-Category: Performance
Created: 2026-10-22
Updated: 2026-10-23
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Revisionsanreicherung verwendet SDK-Standard `Limit=10`, obwohl nur das erste Ergebnis konsumiert wird.
Gesprächsreferenz: Swarm-Explore-Review `E16`; Quelle: Cheap Wins.

## Reach-and-Impact

Reach [S]: `ls --include-deleted` und `ls --only-deleted` reichern gelöschte Dateien mit Revisionsmetadaten an.
Impact [S]: Bei ausreichend langer Historie werden bis zu neun ungenutzte Metadatenobjekte übertragen und decodiert.
API-Aufrufzahl bleibt gleich; Laufzeit- und Byteersparnis wurden nicht gemessen.

## Evidence

- [S] `cmd/ls.go:244-260` erzeugt das Revisionsargument und konsumiert ausschließlich `res.Entries[0]`.
- [S] `go.mod:6` pinnt Dropbox-SDK v6.6.2.
- [S] SDK v6.6.2, `files/types.go:2168-2175`, setzt im Konstruktor `Limit=10`.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Revisionsmengenursache geprüft.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Eine Request-Einstellung beim vorhandenen Owner begrenzen, keinen Cache einführen.

## Proposed-Change

Vor `ListRevisionsContext` ausdrücklich `revisionArg.Limit = 1` setzen.
Anreicherung und Auswahl des ersten Ergebnisses unverändert erhalten.
Der Nutzer autorisierte die private Umsetzung dieses Review-Befunds `P4` am 2026-10-23 mit „fix all issues“.
Diese Freigabe wählt weder einen Upstream-Arbeitsmodus noch ein Veröffentlichungsziel.

## Scope-and-Constraints

- Preserve: Fehlerbehandlung, leere Revisionslisten und `not_file`-Fallback erhalten.
- Preserve: JSON- und Langformat-Metadaten weiterhin anreichern.
- Exclude: Keine neuen Flags, zusätzliche API-Abfragen oder allgemeine Listing-Refaktorierung.
- Limit: Bis zu neun Objekte sind eine Obergrenze je geeignetem Eintrag, kein gemessener Geschwindigkeitsgewinn.

## Verification

### Historischer Untersuchungsstand

Quell- und SDK-Prüfung abgeschlossen; kein Request-Smoke oder Benchmark ausgeführt.
Vorgeschlagen: `limit=1` im Request und unveränderte erste Revisionsmetadaten nachweisen.
Leerantwort und `not_file`-Fallback müssen unverändert behandelt werden.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

### Privater Abschluss am 2026-10-23

[O] Commit `e4a0229` auf `personal` wurde nach `origin/personal` gepusht.
[S] `cmd/ls.go:32,245-247` begrenzt ausschließlich die Revisionsanreicherung auf einen Eintrag.
[O] `TestLsJSONDeletedEntryIsStructured` scheiterte vor der Korrektur mit `ListRevisions limit = 10` und besteht danach.
[O] Temporäre Smokes bestätigen unveränderte Leerantworten und den `not_file`-Fallback mit notwendiger Metadatenabfrage.
[O] `go test ./...`, `go vet ./...`, `go build ./...` und `go test -race ./...` bestehen.
[O] `golangci-lint run --fix ./...` meldet ausschließlich den unveränderten ungenutzten Helfer `cmd/json_metadata.go:64`.
Test decision: update; der bestehende JSON-Listingtest schützt zusätzlich das Requestlimit des externen API-Vertrags.
Keine permanente Testfunktion ergänzt oder entfernt; temporäre Smokes wurden entfernt.
Keine Live-Kontoabfragen, Byte-/Decodierzeitmessung, globale Installation oder Upstream-Veröffentlichung.

## Publication-Blockers

- Reale Antwortbytes und Decodierzeit sind nicht gemessen; der Offline-Nachweis bestätigt Requestlimit und bestehende Ergebnis-/Fallback-Verträge.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: —
Action: None.
Done-When: None.

## Archive

Archive-Reason: Other
Detail: Privat autorisiertes Revisionslimit implementiert, geprüft und auf origin/personal gepusht; keine Upstream-Erledigung behauptet.
Evidence: `e4a0229`; bestehender JSON-Listingtest vor/nach Korrektur, Fallback-Smokes und vorgeschriebene Go-Prüfungen.
Checked: 2026-10-23
