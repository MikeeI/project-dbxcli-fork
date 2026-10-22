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

Root-Cause [S]: Im Ausgangsstand erklärt die Hilfe nicht, dass Auswahl vor lokaler Sortierung endet.
Gesprächsreferenz: Swarm-Explore-Review `E13`; Quellen: Surface, Taste.

## Reach-and-Impact

Reach [S]: Nutzer kombinieren bei `search` und `ls` ein Limit mit lokaler Sortierung.
Impact [S]: `search report --sort size --reverse --limit 1` liefert nicht zuverlässig den global größten Treffer.
Die damalige Hilfe legt die Einschränkung nicht offen; tatsächliche Fehlentscheidungen wurden nicht beobachtet.

## Evidence

- [S] `cmd/search.go:79-93` begrenzt vor `sortEntries`.
- [S] `cmd/ls.go:190-211` begrenzt vor der nachfolgenden Sortierung.
- [S] `cmd/search.go:229-233` und `cmd/ls.go:416-418` erklären die Reihenfolge nicht.
- [S] `--order-by` beschreibt serverseitige Suchreihenfolge, `--sort` lokale Darstellung.

Die historischen Search-Anker wurden durch den privaten `GOAL.md`-Auftrag korrigiert.
`registerSearchFlags` erklärt Auswahl vor Sortierung und schließt eine globale Top-N-Zusage ausdrücklich aus.
Die Listing-Hilfe gehört nicht zur implementierten Search-Erweiterung.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Erklärungslücke geprüft.
`ISSUE-003` besitzt den unabhängigen Komparatorfehler, nicht das Bedienmodell.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Hilfebesitzer präzisieren und vorhandene Generatoren verwenden.

## Proposed-Change

Serverreihenfolge, Begrenzung und anschließende lokale Sortierung ausdrücklich erklären.
`--reverse` auf die gewählte lokale Sortierung beziehen.
Search-Hilfe und generierte Referenz sind privat korrigiert; die entsprechende Listing-Erklärung bleibt offen.

## Scope-and-Constraints

- Preserve: Bestehende Such-, Limit- und Sortiersemantik unverändert erhalten.
- Exclude: Keine globale Top-N-Auswahl, zusätzlichen API-Abfragen oder neuen Sortiermodi.
- Preserve: JSON-Hilfe übernimmt dieselben präzisierten Flagbeschreibungen.
- Related: `ISSUE-034` behandelt die sichtbare Ergebnisbegrenzung, nicht deren Sortierreihenfolge.

## Verification

- [O] Die tatsächliche `search --help`-Ausgabe nennt `applied before local sorting`.
- [O] `--sort` erklärt `not a global top-N selection`, und `--reverse` bezieht sich auf den lokalen Sortiermodus.
- [O] `go run ./tools/gen-docs` regeneriert die bestehende Search-Referenz aus dem Code-Owner.

Test decision: none — die reine Hilfesemantik wird an der tatsächlichen CLI geprüft.
Die neuen Discovery-Filterverträge werden separat durch bestehende Tests geschützt.

## Publication-Blockers

- Die verbleibende Listing-Hilfefassung mit ihrem vorhandenen Auswahlpfad abgleichen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Verbleibende Listing-Hilfe offline abgleichen
Action: Die unveränderte Listing-Hilfe mit Begrenzung und anschließender lokaler Sortierung vergleichen.
Done-When: Die verbleibende Listing-Erklärungslücke ist ohne globale Top-N-Zusage präzise dokumentiert.
