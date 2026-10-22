# ISSUE-027 — Rekursives get: Textpfad verwirft den gewählten Fehlerwriter

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Low
Root-Cause-Confidence: High
Finding-Category: Correctness
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Ein Textmodus-Wrapper ersetzt die vom Kommando gewählten Optionen durch `getOptions{}`.
Gesprächsreferenz: Swarm-Explore-Review `E12`; Quelle: Seam.

## Reach-and-Impact

Reach [S]: Rekursive Textdownloads mit eigens gesetztem Cobra-Fehlerwriter.
Impact [S]: Status und Einzelfehler fallen auf `os.Stderr` zurück statt das gewählte Ziel zu verwenden.
Gewöhnliches CLI-Stderr ist davon nicht unterschiedlich betroffen; die Beitragspriorität bleibt Low.

## Evidence

- [S] `cmd/get.go:140-144` verwendet nur im JSON-Pfad die bereits ausgewählten Optionen.
- [S] `cmd/get.go:248-250` setzt im Textwrapper die Optionen zurück.
- [S] `cmd/get.go:332-358` verwendet den daraus abgeleiteten Writer für Status und Fehler.
- [S] LSP bestätigt einen Produktionsaufrufer des zurücksetzenden Wrappers.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Writer-Ursache geprüft.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Einen bestehenden Produktionsaufruf und seinen privaten Wrapper begrenzen.

## Proposed-Change

Textrekursion mit den vorhandenen `opts` direkt an `getRecursiveInternal` anschließen.
`collectResults=false` erhalten und `getRecursiveWithRootMetadata` nach Migration entfernen.

## Scope-and-Constraints

- Preserve: JSON-Ergebnisse, ressourcensparenden Textpfad und gewöhnliches CLI-Stderr erhalten.
- Preserve: Bestehende direkte interne Verbraucher bei der Wrapperentfernung berücksichtigen.
- Exclude: Keine neue Rendering-Abstraktion oder allgemeine Downloadumstrukturierung.
- Related: `ISSUE-024` benötigt denselben tatsächlichen Diagnosewriter für die Terminalpolitik.

## Verification

Quell- und Verbraucherprüfung abgeschlossen; kein Download-Smoke oder Test ausgeführt.
Vorgeschlagen: Rekursiven Textdownload mit gesetztem Fehlerwriter auf Status und Einzelfehler prüfen.
Keine Ausgabe darf versehentlich auf das globale stderr wechseln.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Writer-Fluss bis zur Fortschrittsanzeige und Fehleraggregation offline prüfen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Rekursiven Fehlerwriter offline nachweisen
Action: Einen rekursiven Textdownload mit eigens gewähltem Fehlerwriter und kontrolliertem Einzelfehler nachvollziehen.
Done-When: Writer-Verlust und unveränderte Ergebnislistenpolitik sind an der Command-Grenze dokumentiert.
