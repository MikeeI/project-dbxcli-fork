# ISSUE-024 — Transferdiagnostik: dynamischer Fortschritt wird ohne Terminal ausgegeben

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Low
Root-Cause-Confidence: High
Finding-Category: UI
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Drei Fortschrittspfade prüfen die Terminalfähigkeit ihres tatsächlichen Diagnosewriters nicht.
Gesprächsreferenz: Swarm-Explore-Review `E9`; Quelle: Default.

## Reach-and-Impact

Reach [S]: Uploads, gewöhnliche Downloads und Shared-Link-Downloads zeichnen Fortschritt auf stderr.
Impact [S]: Umgeleitete Diagnosen enthalten dynamische Zeilen mit Carriage Returns.
Verworfenes stderr oder daraus folgende übersehene Fehler wurden nicht beobachtet.

## Evidence

- [S] `cmd/put.go:79-90`, `cmd/get.go:475-498` und `cmd/share_link_download.go:432-462` zeichnen ungeprüft.
- [S] Gepinntes `ioprogress/draw.go:33-54` schreibt `\r` ohne Terminalprüfung.
- [S] `ioprogress/reader.go:101-105` fällt bei `DrawFunc=nil` auf seinen Standardcallback zurück.
- [S] `ioprogress/draw.go:17-20` bindet diesen Standardcallback an stdout.
- [S] `go.mod:14` enthält bereits `golang.org/x/term`.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Terminalauswahl geprüft.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Diagnosepolitik besitzen; keine Änderung der Transfer- oder Ergebnisverträge.

## Proposed-Change

Standard `auto` am tatsächlichen Diagnoseziel auswerten.
Bei deaktiviertem Fortschritt ausdrücklich einen stummen Zeichencallback verwenden, nicht `nil`.
`always` und `never` ausschließlich auf dynamischen Fortschritt anwenden.

## Scope-and-Constraints

- Preserve: Warnungen, Fehler und stdout unabhängig von der Fortschrittseinstellung erhalten.
- Compatibility: Bisherige Fortschrittsauswertung über `always` ermöglichen und den Standardwechsel dokumentieren.
- Cost: Flag-Metadaten und generierte Referenzen gemeinsam nachführen.
- Prerequisite: Den gewählten Writer auch im rekursiven Textdownload erhalten; siehe `ISSUE-027`.
- Exclude: Keine Terminalprüfung an stdout statt am tatsächlichen Diagnosewriter.

## Verification

Quell- und Bibliotheksprüfung abgeschlossen; kein Terminal- oder Transfer-Smoke ausgeführt.
Vorgeschlagen: Terminal, umgeleitetes stderr sowie explizite Modi vergleichen.
Keine Fortschrittsbytes dürfen auf stdout landen; Warnungen und Fehler müssen erhalten bleiben.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Writer-Erkennung und explizit stumme Bibliotheksanbindung im Ablauf prüfen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Terminalpolitik am Diagnosewriter prüfen
Action: Für die drei Transferpfade terminalfähige und umgeleitete Diagnosewriter samt Bibliotheksfallback abgrenzen.
Done-When: Auto-Auswahl, alternative Modi und stdout-Sauberkeit sind ohne versteckte Writer-Wechsel spezifiziert.
