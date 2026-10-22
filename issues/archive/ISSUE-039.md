# ISSUE-039 — Search: lokale Metadatenfilter fehlen vor der Ergebnisauswahl

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

[S] Der Ausgangsstand bietet keine lokalen Namens-, Größen-, Zeit- oder Typprädikate vor der Ergebnisauswahl.
Ein nachträglicher Verbraucherfilter kann bereits durch das CLI-Limit ausgeschlossene passende Treffer nicht zurückholen.
Die direkte Nutzeranweisung zur privaten Umsetzung von `GOAL.md` autorisiert diesen vollständigen lokalen Vertrag.

## Reach-and-Impact

[S] Präzise Auswahl aus einer Dropbox-Query benötigt Filterung vor Accepted-ID-Deduplizierung und Ergebnislimit.
Lokale Filter können ausschließlich empfangene Suchmatches einschränken, keine vollständige Kontoinventur ersetzen.
[A] Der zusätzliche Seitenbedarf hängt von der Selektivität der Query und den lokalen Prädikaten ab.

## Evidence

- [S] `cmd/search-filters.go` besitzt Parsing, vorbereitete Prädikate und deren gemeinsame Auswertung.
- [S] `parseSearchFilters` unterscheidet fehlende Größenbounds durch `Changed()` von explizit gesetztem null Byte.
- [S] Datumspointer erhalten RFC3339-Zeitzone und Subsekunden; das Intervall ist unten inklusiv und oben exklusiv.
- [S] Ordner erfüllen keine dateispezifischen Größen-/Datumsbounds; widersprüchliche Ordnerfilter werden zurückgewiesen.
- [O] `TestSearchUsesSearchV2AndCommandOutput` prüft alle lokalen Prädikate an einem mehrseitigen Suchlauf.
- [O] Die Live-Erfolgsschemavalidierung schützt ein ausdrücklich gesetztes `max_size: 0` bei einer leeren Datei.
- [O] `TestJSONHelpManifestV1SelectedCommandMetadata` und 14 tatsächliche Skill-CLI-Assertions bestehen.
- [O] `go test ./...`, `go vet ./...` und `go build ./...` bestehen.

## Prior-Art

Lokaler Index und passende Akten wurden vor der Vergabe auf Suchfilter und dieselbe Auswahlursache durchsucht.
`ISSUE-028` besitzt die Erklärung lokaler Sortierung, `ISSUE-032` die bestehende globale Auth-/Validierungsreihenfolge.
Beide werden durch diesen Filter-Scope nicht allgemein behoben.
Externe Prior Art und Upstream-Ziele wurden nicht untersucht.

## Proposed-Change

Privat umgesetzt: Glob, Go-Regexp, Literalinklusion/-exklusion, inklusive Bytegrenzen, Zeitintervall und Dateityp.
Positive Filter wirken gemeinsam; jede passende Literalexklusion verwirft den Namen.
Ein lokaler Filter hält die SDK-Seitengröße bei, statt das Suchlimit zur winzigen API-Seitengröße zu machen.
Discovery liefert echte String-Arrays, Item-Enums und keine semantischen Defaults für optionale Bytegrenzen.

## Scope-and-Constraints

- Bestehende `cmd/`- und Cobra-Struktur bleiben ohne neue Abhängigkeit, Filtersprache oder Worker-Pipeline erhalten.
- Regex wird einmal kompiliert; Glob bezieht sich auf den Basename statt rekursiver Pfadtraversierung.
- Exklusionsflags sind native StringArray-Werte; Kommata bleiben literal und leere Exklusionen sind ungültig.
- `--time server|client` wählt auch das Datum für lokale Filter.
- Fehlende Größenbounds bleiben im Echo abwesend; ein ausdrücklich gesetztes null Byte bleibt erhalten.
- Filterung bedeutet keine Zusage vollständiger API-Ergebnisse oder konstanter Latenz.

## Verification

Test decision: update — `TestSearchUsesSearchV2AndCommandOutput` schützt den mehrseitigen Filtervertrag.
`TestJSONHelpManifestV1SelectedCommandMetadata` schützt Typen und präsenzabhängige Schema-Defaults.
Die vorhandene Live-Output-Schemavalidierung schützt Größe null, neue Suchmetadaten und ältere v1-Antworten ohne Summary.
Die tatsächliche CLI-Texthilfe beschreibt alle Filter; Kontoabfragen und globale Installation wurden nicht ausgeführt.

## Publication-Blockers

Upstream-Arbeitsmodus, Veröffentlichungsziel, aktuelle Prior Art, exakter Entwurf und dessen Freigabe fehlen.

## Next-Action

Summary: —
Action: None.
Done-When: None.

## Archive

Archive-Reason: Other
Detail: Privater Scope vollständig abgeschlossen; kein aktueller Auftrag für Upstream-Arbeit.
Evidence: Lokale Filter, präsenzkorrekte Discovery und alle vorhandenen Vertragsprüfungen sind in `547752f` gepusht.
Checked: 2026-10-22
