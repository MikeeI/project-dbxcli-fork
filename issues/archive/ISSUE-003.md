# ISSUE-003 — Rückwärtssortierung: instabile Reihenfolge bei gleichen Schlüsseln

State: Archived
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Medium
Root-Cause-Confidence: High
Finding-Category: Correctness
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Im referenzierten Ausgangsstand negiert `sortEntries` bei Reverse das Ergebnis von `compareLess`.
Bei gleichen Schlüsseln werden dadurch beide Vergleichsrichtungen wahr.
Der Komparator verletzt die strikte schwache Ordnung des verwendeten `sort.SliceStable`.
Gesprächsreferenz: `CW2`; der anschließende Review bestätigt Ursache und lokale Korrektur.
Zusätzliche Gesprächsreferenz: Swarm-Explore-Review `E3`; Quelle: Cheap Wins.
Der Review bestätigt vertauschte Operanden ohne zusätzliche Gleichstandsregel.

## Reach-and-Impact

Reach [S]: `ls` und `search` verwenden dieselbe lokale Sortierung vor ihrer Ausgabe.
Gleichstände entstehen regulär bei mehreren Dateien oder Ordnern unter `--sort type --reverse`.
Impact [S]: Die bisherige Eingangsreihenfolge gleicher Sortierschlüssel ist nicht mehr geschützt.
Text- und JSON-Ausgaben verwenden die betroffene Ergebnisfolge.
Stabilität bezieht sich auf eine gegebene Eingangsfolge, nicht auf unterschiedliche Serverantworten.
Eine reale Nutzungshäufigkeit oder Produktionsstörung wurde nicht erhoben.

## Evidence

Die folgenden historischen Quellanker beziehen sich auf den genannten Upstream-Ausgangsstand.

- [S] `cmd/format.go:58-64` übergibt bei `reverse` den negierten Vergleich an `sort.SliceStable`.
- [S] `cmd/format.go:67-79` liefert für gleiche Namen, Größen, Zeiten oder Typklassen in beiden Richtungen `false`.
- [S] `cmd/format.go:100-107` ordnet allen Dateien beziehungsweise Ordnern jeweils denselben Typ-Schlüssel zu.
- [S] `cmd/ls.go:185-197` und `cmd/search.go:91-93` verwenden den gemeinsamen Sortierpfad.
- [S] Go 1.27.3 verlangt eine strikte schwache Ordnung: `/usr/local/go/src/sort/sort.go:24-31`.
- [S] `sort.SliceStable` erhält gleiche Elemente in Eingangsreihenfolge: `/usr/local/go/src/sort/slice.go:32-44`.
- [S] `cmd/format_test.go:139-151` prüft Reverse nur mit unterschiedlichen Namen und reproduziert keinen Gleichstand.
- [O] Der lesende Source-Vergleich mit dem bestätigten Upstream-Stand zeigte keine relevante Quelldifferenz.

## Prior-Art

Coverage: Der aktuelle Index und die bestehenden Akten `ISSUE-001` und `ISSUE-002` enthalten keine gleiche Ursache.
Gaps: Upstream-Issues, PRs, Diskussionen und Releases wurden für diesen Befund noch nicht recherchiert.
Contribution fit: Eine lokale Korrektur ist abgegrenzt; bestehende Arbeiten und Reproduzierbarkeit sind noch zu prüfen.

## Proposed-Change

Privat umgesetzt: Reverse vertauscht zentral in `compareLess` die Operanden, statt `less` in `sortEntries` zu negieren.
Der stabile Sortierer und die normale Vergleichsrichtung bleiben erhalten.
Es gibt keinen sekundären Namensvergleich und keine neue Sortieroption.
Die direkte vollständige private Implementierungsfreigabe von `GOAL.md` umfasst diesen Fix.

## Scope-and-Constraints

- Preserve: Die absteigende Reihenfolge unterschiedlicher Schlüssel bleibt erhalten.
- Preserve: Gleiche Schlüssel behalten ihre jeweilige Eingangsreihenfolge.
- Preserve: Server-Suchreihenfolge, Pagination und Ergebnisauswahl vor der lokalen Sortierung bleiben unverändert.
- Preserve: Flags, Textfelder, JSON-Schemas und Metadaten bleiben unverändert.
- Exclude: Keine neue Sortierbibliothek, zusätzliche Sortierregel oder allgemeine Refaktorierung einführen.
- Cost: Aufrufer übernehmen die Korrektur automatisch; keine Datenmigration oder Nutzerumstellung ist erforderlich.
- Compatibility: Die bisher fehlerhafte Gleichstandsreihenfolge ändert sich bewusst.
- Rollback: Ausschließlich den lokalen Reverse-Vergleich zurücknehmen; persistierter Zustand ist nicht betroffen.

## Verification

- [O] Vor dem Fix scheitert `TestSortEntriesReverse` mit `[/second.txt /first.txt /folder]` statt stabiler Dateifolge.
- [O] Die Reproduktion verwendet `go test ./cmd -run '^(TestSortEntriesReverse|TestSearchJSONErrorWritesNoOutput)$'`.
- [O] Nach dem Fix besteht `go test ./cmd -run '^(TestSearch|TestRenderSearch|TestSortEntriesReverse)'`.
- [O] `go test ./...`, `go vet ./...` und `go build ./...` bestehen.
- [O] Commit `d36e7f6` enthält ausschließlich Komparatorkorrektur und vorhandenen Regressionstest und ist gepusht.

Test decision: update — `TestSortEntriesReverse` schützt gleiche Typ-Schlüssel ohne neue Testfunktion.
Die Nachweise benötigen weder Dropbox-Abfragen noch persistierte Zustandsänderungen.

## Publication-Blockers

- Die private Korrektur und ihre Reproduktion sind abgeschlossen; eine Upstream-Veröffentlichung ist nicht beauftragt.
- Prior Art und aktuelle Upstream-Beitragsvorgaben sind vor einer Veröffentlichung zu prüfen.
- `Authorized-Work` und `Publication-Target` sind vom Nutzer noch nicht ausgewählt.
- Ein exakter Veröffentlichungsentwurf und dessen Freigabe liegen nicht vor.

## Next-Action

Summary: —
Action: None.
Done-When: None.

## Archive

Archive-Reason: Other
Detail: Privater Scope vollständig abgeschlossen; kein aktueller Auftrag für Upstream-Arbeit.
Evidence: Fix `d36e7f6` und Search-Vertrag `547752f` sind geprüft und auf `origin/personal` gepusht.
Checked: 2026-10-22
