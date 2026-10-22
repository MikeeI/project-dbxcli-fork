# ISSUE-008 — Uploads: finaler Append schließt vor ausstehenden Chunks

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: High
Root-Cause-Confidence: High
Finding-Category: Correctness
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Der schließende Append läuft ohne Abschlussbarriere im selben Pool wie frühere Bytebereiche.
Gesprächsreferenz: `E3`; der Review bestätigt den fehlerhaften Sessionübergang.

## Reach-and-Impact

Reach [S]: Parallele Chunk-Uploads mit verzögerten oder wiederholten früheren Appends.
Impact [O]: Ein isolierter Upload erhielt nach vorgezogenem Close einen `closed`-Fehler und erreichte keinen Finish-Aufruf.
Eine Live-Dropbox-Scheduling-Häufigkeit ist nicht gemessen.

## Evidence

- [S] `cmd/put.go:130-139` überträgt `chunk.close` unmittelbar in `args.Close`.
- [S] `cmd/put.go:162-169` sendet den anhand der Dateiposition bestimmten finalen Chunk in denselben Workerpool.
- [S] `cmd/put.go:177-178` wartet erst nach dessen Übergabe auf alle Worker.
- [S] SDK v6.6.2, `dropbox/files/types.go:4688-4691`: Nach Close sind weitere Appends unzulässig.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierter HTTP-Gegenpart.
Aufruf: `put <33-MiB-datei> /target --output=json` über `RootCmd.ExecuteC`.
Frühere Appends werden verzögert; der finale Append schließt die Session vor ihrer Bestätigung.
Beobachtet [O]: `command_error=closed/...`, `append_calls=3`, `finish_calls=0`.
Der Gegenpart bildet den dokumentierten Closed-Zustand ab; kein realer Dropbox-Transfer wurde ausgeführt.
Der temporäre Harness wurde nach dem Audit entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Ursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Die Close-Barriere ist von Pufferoptimierung und allgemeinem Fehlerkanal-Lifecycle getrennt.

## Proposed-Change

Den finalen Chunk aus der parallelen Menge zurückhalten.
Erst alle vorherigen Appends erfolgreich bestätigen lassen, dann den finalen Close-Append senden und finishen.

## Scope-and-Constraints

- Preserve: Zulässige Parallelität nicht-finaler, korrekt ausgerichteter Bytebereiche erhalten.
- Preserve: SDK-Hashing, Cursor-Offsets und unveränderliche Retry-Payloads erhalten.
- Related: `ISSUE-007` besitzt Workerfehler und Cleanup; `ISSUE-014` besitzt mehrdeutige Close-Recovery.
- Exclude: Keine vollständige Serialisierung aller Chunks und kein neues Transferframework.

## Verification

Abgeschlossen [O]: Der Protokoll-Smoke reproduziert den vorgezogenen Sessionabschluss.
Erforderlich: Einen frühen Append verzögern oder wiederholen und die Request-Reihenfolge beobachten.
Erwartet: Der Close-Request folgt auf erfolgreiche Bestätigung aller früheren Bytebereiche.
Test decision: none; dieser Auftrag implementiert die Korrektur nicht.

## Publication-Blockers

- Reproduktionsablauf und SDK-Zustandsvertrag dauerhaft nachvollziehbar zusammenführen.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Sessionabschluss gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zur Abschlussbarriere paralleler Upload-Sessions recherchieren.
Done-When: Suchabdeckung, Treffer und Abgrenzung zu Worker-Lifecycle und Close-Recovery sind dokumentiert.
