# ISSUE-037 — Search: Matchtyp und Titel-Highlights gehen verloren

State: Archived
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Medium
Root-Cause-Confidence: High
Finding-Category: API
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

[S] Der Ausgangsstand reduziert `files.SearchMatchV2` vor der Darstellung auf gemeinsame Metadaten.
Matchtyp und native Highlight-Spans erreichen deshalb die Ausgabe nicht.
Die vollständige private Umsetzung von `GOAL.md` ist direkt autorisiert, eine Upstream-Veröffentlichung nicht.

## Reach-and-Impact

[S] JSON-Verbraucher verlieren servergelieferte Informationen darüber, weshalb ein Treffer gefunden wurde.
Titelspans ermöglichen Markierung, sind jedoch weder Inhaltssnippets noch verlässliche Offsets.
Das Vorhandensein dieser Informationen für konkrete Live-Suchläufe ist nicht beobachtet.

## Evidence

- [S] Das gepinnte SDK enthält `SearchMatchV2.MatchType`, `HighlightSpans` und `SearchMatchFieldOptions`.
- [S] `cmd/search.go:collectSearchMatches` erhält vollständige Matchobjekte bis nach der lokalen Auswahl.
- [S] `searchOperationResults` ergänzt nur Search-Metadaten um Name, offenen Matchtyp und native Titelspans.
- [S] `newSearchV2Arg` setzt die moderne `MatchFieldOptions.IncludeHighlights`, nicht das veraltete Feld.
- [O] `TestSearchJSONOutputsInputAndResults` prüft Request, Matchtyp und beide nativen Span-Felder.
- [O] Die bestehenden Live-Output-Schematests validieren tatsächliche Handler-Ausgaben mit Matchinformationen.
- [O] `go test ./...`, `go vet ./...` und `go build ./...` bestehen.

## Prior-Art

Lokaler Index und passende Akten wurden auf Matchdatenverlust und Highlights durchsucht.
Die bisherige Kürzungsakte `ISSUE-034` besitzt einen anderen Informationsverlust.
Eine externe Prior-Art-Recherche oder Bewertung eines Upstream-Ziels ist nicht Teil der privaten Lieferung.

## Proposed-Change

Privat umgesetzt: vollständige Matches, Search-spezifische Metadaten und JSON-only-Flag `--highlights`.
Fehlende Matchinformationen werden nicht geraten, und unbekannte zukünftige Matchtypen bleiben gültige Strings.
Gemeinsame Metadaten anderer Commands werden nicht erweitert.

## Scope-and-Constraints

- Titelspans bleiben SDK-native Objekte mit `highlight_str` und `is_highlighted`.
- Text-Trefferzeilen erhalten keine Highlight- oder Snippetdarstellung.
- `--highlights` ohne JSON-Ausgabe ist ein Validierungsfehler.
- Neue optionale Felder erfordern aktualisierte geschlossene Schemas beim Verbraucher.

## Verification

Test decision: update — `TestSearchJSONOutputsInputAndResults` schützt den öffentlichen Search-Vertrag.
Die bestehenden Erfolgs-/Fehler-Schematests und Golden-Fixtures prüfen dieselbe Metadatenform.
Texthilfe und Skill-Referenzen benennen JSON-only-Highlights ausdrücklich als Titelspans.
Keine Live-Dropbox-Anfrage wurde ausgeführt.

## Publication-Blockers

Upstream-Arbeitsmodus, Veröffentlichungsziel, aktuelle Prior Art und exakte Freigabe fehlen.

## Next-Action

Summary: —
Action: None.
Done-When: None.

## Archive

Archive-Reason: Other
Detail: Privater Scope vollständig abgeschlossen; kein aktueller Auftrag für Upstream-Arbeit.
Evidence: Matchtyp, Titelspans und geprüfte Ausgabe-/Schema-Verträge sind in `547752f` auf `origin/personal` enthalten.
Checked: 2026-10-22
