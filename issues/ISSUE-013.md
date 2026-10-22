# ISSUE-013 — Uploads: Autorename-Retry erzeugt zusätzliche Remote-Dateien

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Medium
Root-Cause-Confidence: High
Finding-Category: Correctness
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Ein nicht-idempotenter Single-shot-Commit wird nach mehrdeutigem Transportfehler vollständig wiederholt.
Gesprächsreferenz: `E8`; der Review bestätigt die falsche Wiederholungsidentität bei Autorename.

## Reach-and-Impact

Reach [S]: Kleine Uploads bis einschließlich 32 MiB mit `--if-exists=autorename` und verlorener Commit-Bestätigung.
Impact [S]: Ein bereits gespeicherter erster Versuch kann beim Retry eine zweite Datei erzeugen.
Impact [O]: Der isolierte Gegenpart beobachtete zwei Commit-Versuche; die CLI berichtete ausschließlich den letzten.
Kein realer Remote-Commit wurde ausgeführt.

## Evidence

- [S] `cmd/put.go:93-102` setzt den Reader zurück und wiederholt `UploadContext` vollständig.
- [S] `cmd/put.go:590-610` setzt für Autorename Add-Modus, `StrictConflict=true` und `Autorename=true`.
- [S] SDK v6.6.2, `dropbox/files/types.go:247-249,263-268`: Konflikte führen zu Autorename, auch bei identischem Inhalt.
- [S] `cmd/retry.go:52-62,116-124` klassifiziert mehrdeutige Lese- und Netzwerkfehler als wiederholbar.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; temporäre kleine Quelldatei.
Aufruf: `put <kleine-datei> /target --if-exists=autorename --output=json` über `RootCmd.ExecuteC`.
Der isolierte Gegenpart modelliert einen gespeicherten ersten Versuch, liefert jedoch `io.ErrUnexpectedEOF` statt Bestätigung.
Er beantwortet den zweiten Upload mit Metadaten für `/target (1)`.
Beobachtet [O]: `command_error=<nil>`, `commits=2` und genau ein zurückgegebenes `autorenamed`-Ergebnis.
`commits=2` zählt simulierte Commit-Verarbeitung, keine beobachteten Live-Dropbox-Dateien.
Der temporäre Harness wurde entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Commit-Identitätsursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Retry-Politik am schreibenden Aufrufer korrigieren, nicht alle Netzwerk-Retries pauschal entfernen.

## Proposed-Change

Nicht-idempotente Commits nach mehrdeutigen Transportfehlern nicht blind wiederholen.
Den unklaren Commit-Ausgang ausdrücklich behandeln oder anhand belastbarer Operationsidentität rekonstruieren.

## Scope-and-Constraints

- Preserve: Eindeutig abgewiesene, sicher wiederholbare Anfragen dürfen retryfähig bleiben.
- Preserve: Reguläre Konfliktpolitik und Ergebnisidentität erhalten.
- Exclude: Keine Inhaltsgleichheits-Deduplizierung unterstellen; `StrictConflict` hebt diese Gegenhypothese auf.
- Exclude: Kein genereller Retry-Framework-Umbau ohne weitere nachgewiesene Anforderungen.

## Verification

Abgeschlossen [O]: Der Offline-Smoke reproduziert die erneute nicht-idempotente Anfrage nach verlorener Bestätigung.
Erforderlich: Ersten Commit akzeptieren, dessen Antwort verlieren und den resultierenden zweiten Seiteneffekt beobachten.
Erwartet: Kein zweiter Remote-Dateieffekt für denselben logischen Upload.
Test decision: none; keine Source- oder Teständerung in diesem Auftrag.

## Publication-Blockers

- Protokollablauf und Abgrenzung zwischen Anfrageversuch und tatsächlich bestätigtem Commit dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Autorename-Retries gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zu mehrdeutigen Single-shot-Commits mit Autorename recherchieren.
Done-When: Suchabdeckung, relevante Treffer und sichere gegenüber mehrdeutigen Retry-Fällen sind dokumentiert.
