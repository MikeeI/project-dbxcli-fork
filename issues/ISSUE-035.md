# ISSUE-035 — Uploadstatus: Lesefortschritt wirkt wie bestätigte Übertragung

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

Root-Cause [S]: Der als Upload dargestellte Zähler misst Quelldatenlesen vor der parallelen Chunk-Verarbeitung.
Gesprächsreferenz: Swarm-Explore-Review `E7`; Quelle: Operator.
Der Review begrenzt die Korrektur auf ehrliche Statusdarstellung statt zusätzlicher Worker-Telemetrie.

## Reach-and-Impact

Reach [S]: Chunk-Uploads oberhalb 32 MiB lesen Chunks vor deren Bestätigung durch den Server.
Impact [S]: Vollständig gelesene Quelldaten bedeuten weder abgeschlossene Appends noch erfolgreichen Finish.
Reale Wartezeiten und Nutzungshäufigkeit sind unbekannt.

## Evidence

- [S] `cmd/put.go:79-90` benennt den Reader-Fortschritt als `Uploading`.
- [S] `cmd/put.go:149-178,188-195` trennt Lesen, Append-Abschluss und Finish.
- [S] `cmd/put.go:598` setzt den Fortschrittsreader vor die Chunk-Pipeline.
- [S] Gepinntes `ioprogress/reader.go:29-57` zählt ausdrücklich gelesene Bytes.
- [S] LSP bestätigt als Produktionsverbraucher `cmd/put.go:100` und `cmd/put.go:598`.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Statusursache geprüft.
`ISSUE-007` besitzt Worker-Cleanup, `ISSUE-008` die Close-Barriere; beide sind unabhängige Fehler.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Sichtbare Messgröße präzisieren, keinen neuen Fortschrittszustandsautomaten einführen.

## Proposed-Change

Chunk-Fortschritt ausdrücklich als Lesen der Quelldaten beschriften.
Den noch ausstehenden Uploadabschluss getrennt kennzeichnen.

## Scope-and-Constraints

- Preserve: Keine Änderung an Transferparallelität, Retry-Payloads, JSON-Ergebnissen oder Exitcodes.
- Exclude: Keine synchronisierten Bestätigungszähler in Workern und keine neuen Nutzerentscheidungen.
- Preserve: Single-shot-Verhalten außerhalb der begrenzten Statusänderung erhalten.

## Verification

Quellprüfung abgeschlossen; kein verzögerter Transfer-Smoke ausgeführt.
Vorgeschlagen: Appends und Finish verzögern und die sichtbaren Phasen vergleichen.
Vollständig gelesene Quelldaten dürfen keinen abgeschlossenen Upload suggerieren.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Messgrößen und Statusübergänge im verzögerten Ablauf bestätigen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Uploadstatus bei Verzögerung prüfen
Action: Lesen, ausstehende Appends und verzögerten Finish in einem isolierten Statusablauf unterscheiden.
Done-When: Die kleinste Anzeigeänderung benennt jede Phase ohne vorzeitige Erfolgsbehauptung.
