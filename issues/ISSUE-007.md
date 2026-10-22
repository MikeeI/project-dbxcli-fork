# ISSUE-007 — Uploads: Workerfehler blockieren den Abschluss dauerhaft

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: High
Root-Cause-Confidence: High
Finding-Category: Reliability
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Worker senden blockierend auf einen Fehlerkanal, dessen Eigentümer erst nach `wg.Wait()` liest.
Gesprächsreferenz: `E2`; der Review bestätigt Deadlock und unvollständige Worker-Lebensdauer.

## Reach-and-Impact

Reach [S]: Chunk-Uploads oberhalb 32 MiB mit mehreren Append-Fehlern nach Abschluss der Jobproduktion.
Impact [O]: Ein isolierter 33-MiB-Upload kehrte trotz `--timeout=1s` nach vier Sekunden nicht zurück.
Der Netzwerk-Timeout löst den nackten Channel-Send nicht auf; daraus wird keine allgemeine Gesamtlaufzeitgarantie abgeleitet.

## Evidence

- [S] `cmd/put.go:127` begrenzt `errCh` auf genau einen Eintrag.
- [S] `cmd/put.go:130-139` lässt jeden Fehler blockierend vor dem abschließenden `wg.Done()` senden.
- [S] `cmd/put.go:177-181` wartet vor dem ersten abschließenden Fehlerempfang auf sämtliche Worker.
- [S] `cmd/put.go:150-171` besitzt zusätzlich frühe Rückkehrpfade ohne geschlossenen Arbeitskanal und Worker-Abschluss.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, Dropbox-SDK v6.6.2; ausschließlich temporäre lokale Dateien.
Aufruf: `put <33-MiB-datei> /target --output=json --timeout=1s` über `RootCmd.ExecuteC`.
Ein isolierter HTTP-Gegenpart hält drei Append-Anfragen bis zur vollständigen Einreihung an und liefert danach Fehler.
Beobachtet [O]: `command_returned=false`, `append_calls=3`, `blocked_channel_sends=2`, `waitgroup_present=true`.
Der temporäre Harness wurde entfernt; weder Dropbox noch gespeicherte Credentials wurden verwendet.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Ursache.
`ISSUE-001` betrifft Chunk-Puffer; `ISSUE-004` betrifft Traversierung, nicht Worker-Abschluss.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Eine Lifecycle-Korrektur muss sämtliche Worker-Ausgänge abdecken, nicht nur den vollen Fehlerkanal.

## Proposed-Change

Den ersten Workerfehler nichtblockierend veröffentlichen und weitere Produktion abbrechen.
Arbeitskanal und gestartete Worker auf jedem Ausgang zuverlässig beenden und abwarten.

## Scope-and-Constraints

- Preserve: Die ursprüngliche erste Fehlerursache sowie unveränderliche Retry-Payloads erhalten.
- Preserve: Begrenzte Parallelität, Chunkgrenzen und erfolgreiche Commit-Semantik erhalten.
- Exclude: Eine bloße Kanalvergrößerung behebt frühe Rückkehrpfade nicht.
- Related: `ISSUE-008` behandelt unabhängig davon die Reihenfolge des schließenden Appends.

## Verification

Abgeschlossen [O]: Der beschriebene Offline-Smoke reproduziert die zyklische Wartesituation.
Erforderlich: Mehrere Append-Fehler nach Jobproduktion sowie frühe Quell-Lesefehler prüfen.
Erwartet: Begrenzte Rückkehr mit ursprünglicher Ursache und ohne verbleibende Worker.
Test decision: none; keine Source- oder Teständerung in diesem Auftrag.

## Publication-Blockers

- Vollständig ausführbare Offline-Reproduktion und Lifecycle-Nachweis dauerhaft dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe sind nicht gewählt.

## Next-Action

Summary: Upload-Deadlock gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zum Fehlerkanal und Worker-Abschluss in `uploadChunked` recherchieren.
Done-When: Suchabdeckung, relevante Treffer und Abgrenzung zur Pufferoptimierung sind dokumentiert.
