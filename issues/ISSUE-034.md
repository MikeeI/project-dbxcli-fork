# ISSUE-034 — Ergebnislisten: begrenzte Ausschnitte bleiben unmarkiert

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Medium
Root-Cause-Confidence: High
Finding-Category: UI
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Ergebnisbegrenzung verwirft Einträge oder beendet Pagination, ohne diesen Zustand auszugeben.
Gesprächsreferenz: Swarm-Explore-Review `E6`; Quelle: Surface.

## Reach-and-Impact

Reach [S]: `ls --limit N` und `search --limit N` liefern begrenzte Ergebnisse.
Impact [S]: Die Ergebniszahl allein unterscheidet vollständige Ergebnisse nicht zuverlässig von einem Ausschnitt.
Tatsächliche Fehlentscheidungen oder Nutzungshäufigkeiten wurden nicht beobachtet.

## Evidence

- [S] `cmd/ls.go:190-230` begrenzt allgemeine und gefilterte Auflistungen.
- [S] `cmd/search.go:79-93,165-179` begrenzt Suchergebnisse ohne entsprechende Ausgabekennzeichnung.
- [S] `cmd/json_output.go:23-27,35-42` besitzt bereits strukturierte Warnungen.
- [S] `docs/json-schema/v1/README.md:67-70` erlaubt additive Warncodes mit konsistentem Vertrag.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Ausgabekennzeichnung geprüft.
Die Sortierhilfe aus Swarm-Explore-Review `E13` betrifft einen anderen Informationsmangel.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Vorhandenen Seitenzustand nutzen; keine neue Pagination-API benötigen.

## Proposed-Change

Tatsächlich ausgelassene Einträge und verbleibende Seiten bis zur Ausgabe erhalten.
Textdiagnose auf stderr und dokumentierten JSON-Warncode mit entsprechend präziser Bedeutung verwenden.

## Scope-and-Constraints

- Preserve: Ergebniszeilen, Exitcodes und bestehende Statuswerte erhalten.
- Risk: Eine Ergebniszahl gleich dem Limit beweist keine Kürzung.
- Risk: Weitere gefilterte Seiten beweisen keine weiteren passenden Treffer.
- Exclude: Keine zusätzlichen Dropbox-Abfragen, erfundenen Gesamttrefferzahlen oder neuen Cursor-Verträge.
- Cost: Warnmetadaten und generierte Schemas müssen gemeinsam aktualisiert werden.

## Verification

Quellprüfung abgeschlossen; keine Laufzeitprüfung ausgeführt.
Vorgeschlagen: Vollständige, tatsächlich abgeschnittene und genau passende Ergebnisse vergleichen.
Bei `--only-deleted` gefilterte Seiten und lediglich mögliche weitere Treffer unterscheiden.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Verlässliche Kürzungs- und Seitenzustände am Ausgabeweg prüfen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Kürzungssignale am Ausgabeweg prüfen
Action: Bekannte Kürzung, mögliche Folgeseiten und vollständige Ergebnisse in beiden Befehlen abgrenzen.
Done-When: Der vorhandene Zustand trägt eindeutige Warnungen ohne zusätzliche API-Abfragen oder falsche Trefferzusagen.
