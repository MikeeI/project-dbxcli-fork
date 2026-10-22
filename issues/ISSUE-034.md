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

Root-Cause [S]: Im Ausgangsstand beendet das Limit die Ergebnissammlung ohne sichtbare Kennzeichnung.
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

Privater Search-Fortschritt [S]: `searchSummary` trennt `stop_reason` vom letzten `server_has_more`.
Die historischen Search-Quellanker beschreiben nicht mehr den aktuellen privaten Implementierungsstand.
Der entsprechende Listing-Pfad bleibt unverändert und weiter offen.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Ausgabekennzeichnung geprüft.
Die Sortierhilfe aus Swarm-Explore-Review `E13` betrifft einen anderen Informationsmangel.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Vorhandenen Seitenzustand nutzen; keine neue Pagination-API benötigen.

## Proposed-Change

Für Search privat umgesetzt: Suchzusammenfassung mit API-Ende, lokalem Limit und Fehlerzustand statt eines Warncodes.
Textdiagnosen gehen auf stderr; JSON erhält den typisierten Suchzustand.
Für Listings bleibt die minimale verlässliche Kürzungskennzeichnung getrennt zu untersuchen.

## Scope-and-Constraints

- Preserve: Ergebniszeilen, Exitcodes und bestehende Statuswerte erhalten.
- Risk: Eine Ergebniszahl gleich dem Limit beweist keine Kürzung.
- Risk: Weitere gefilterte Seiten beweisen keine weiteren passenden Treffer.
- Exclude: Keine zusätzlichen Dropbox-Abfragen, erfundenen Gesamttrefferzahlen oder neuen Cursor-Verträge.
- Cost: Suchzustand, Fehlerdetails, Schemas, vorhandene Tests und Verbraucherhinweise werden gemeinsam versioniert.

## Verification

- [O] `TestSearchLimitCapsResultsAndStopsPagination` unterscheidet genaue letzte Seiten und tatsächlich ausgelassene Treffer.
- [O] `TestSearchJSONErrorWritesNoOutput` prüft Initial-/Fortsetzungsfehler einschließlich Teilergebnissen und Exitklasse.
- [O] `TestSearchUsesSearchV2AndCommandOutput` prüft die Textdiagnose mit dem tatsächlichen Search-Handler.
- [O] Alle vorgeschriebenen Go-Gates bestehen; eine neue Listing-Kürzungsprüfung wurde nicht ergänzt.

Test decision: update — die genannten vorhandenen Search-Tests schützen den privaten Teilvertrag.
Für `ls --only-deleted` bleiben Filterseiten und nur mögliche Folgetreffer separat abzugrenzen.

## Publication-Blockers

- Verlässliche Kürzungs- und Seitenzustände verbleibender Listing-Pfade prüfen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Verbleibende Listing-Kürzung offline abgrenzen
Action: Den unveränderten Listing-Pfad auf bekannte Kürzung, mögliche Folgeseiten und vollständige Ergebnisse eingrenzen.
Done-When: Der verbleibende Listing-Vertrag ist ohne zusätzliche API-Abfragen oder falsche Trefferzusagen belegt.
