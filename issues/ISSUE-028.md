# ISSUE-028 — CLI-Hilfe: Begrenzung vor lokaler Sortierung bleibt unklar

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

Root-Cause [S]: Hilfe beschreibt Limit und Sortierung unabhängig, obwohl die Auswahl vor der Sortierung endet.
Gesprächsreferenz: Swarm-Explore-Review `E13`; Quellen: Surface, Taste.

## Reach-and-Impact

Reach [S]: Nutzer kombinieren bei `search` und `ls` ein Limit mit lokaler Sortierung.
Impact [S]: `search report --sort size --reverse --limit 1` liefert nicht zuverlässig den global größten Treffer.
Die Hilfe legt die Einschränkung nicht offen; tatsächliche Fehlentscheidungen wurden nicht beobachtet.

## Evidence

- [S] `cmd/search.go:79-93` begrenzt vor `sortEntries`.
- [S] `cmd/ls.go:190-211` begrenzt vor der nachfolgenden Sortierung.
- [S] `cmd/search.go:229-233` und `cmd/ls.go:416-418` erklären die Reihenfolge nicht.
- [S] `--order-by` beschreibt serverseitige Suchreihenfolge, `--sort` lokale Darstellung.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Erklärungslücke geprüft.
`ISSUE-003` besitzt den unabhängigen Komparatorfehler, nicht das Bedienmodell.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Hilfebesitzer präzisieren und vorhandene Generatoren verwenden.

## Proposed-Change

Serverreihenfolge, Begrenzung und anschließende lokale Sortierung ausdrücklich erklären.
`--reverse` auf die gewählte lokale Sortierung beziehen.
Hilfetexte von `search` und `ls` ändern und Befehlsreferenzen aus ihren Besitzern regenerieren.

## Scope-and-Constraints

- Preserve: Bestehende Such-, Limit- und Sortiersemantik unverändert erhalten.
- Exclude: Keine globale Top-N-Auswahl, zusätzlichen API-Abfragen oder neuen Sortiermodi.
- Preserve: JSON-Hilfe übernimmt dieselben präzisierten Flagbeschreibungen.
- Related: `ISSUE-021` behandelt die sichtbare Ergebnisbegrenzung, nicht deren Sortierreihenfolge.

## Verification

Quellprüfung abgeschlossen; keine Hilfeausführung oder Generatorausführung im Tracking-Auftrag.
Vorgeschlagen: Text- und JSON-Hilfe am Beispiel `--sort size --reverse --limit 1` gegen den Ablauf prüfen.
Die Beschreibung darf keinen global größten Treffer versprechen.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Präzise aktuelle Hilfefassung und generierte Referenzen gemeinsam prüfen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Auswahlmodell in Hilfe abgleichen
Action: Aktuelle Text- und JSON-Hilfe anhand begrenzter lokaler Sortierung mit dem Kontrollfluss vergleichen.
Done-When: Die Erklärungslücke und eine minimale semantisch korrekte Hilfekorrektur sind dokumentiert.
