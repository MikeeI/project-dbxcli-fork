# ISSUE-003 — Rückwärtssortierung: instabile Reihenfolge bei gleichen Schlüsseln

State: Investigating
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

Root-Cause [S]: `sortEntries` negiert bei Rückwärtssortierung das Ergebnis von `compareLess`.
Bei gleichen Schlüsseln werden dadurch beide Vergleichsrichtungen wahr.
Der Komparator verletzt die strikte schwache Ordnung des verwendeten `sort.SliceStable`.
Gesprächsreferenz: `CW2`; der anschließende Review bestätigt Ursache und lokale Korrektur.

## Reach-and-Impact

Reach [S]: `ls` und `search` verwenden dieselbe lokale Sortierung vor ihrer Ausgabe.
Gleichstände entstehen regulär bei mehreren Dateien oder Ordnern unter `--sort type --reverse`.
Impact [S]: Die bisherige Eingangsreihenfolge gleicher Sortierschlüssel ist nicht mehr geschützt.
Text- und JSON-Ausgaben verwenden die betroffene Ergebnisfolge.
Stabilität bezieht sich auf eine gegebene Eingangsfolge, nicht auf unterschiedliche Serverantworten.
Eine reale Nutzungshäufigkeit oder Produktionsstörung wurde nicht erhoben.

## Evidence

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

Im Reverse-Zweig von `sortEntries` die Operanden des bestehenden Vergleichs vertauschen, statt ihn zu negieren.
Den stabilen Sortierer und die normale Vergleichsrichtung beibehalten.
Keinen sekundären Namensvergleich und keine neue Sortieroption einführen.
Die Änderung ist ein Vorschlag und nicht zur Implementierung autorisiert.

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

Abgeschlossen [O]: Der Review hat die Gleichstands- und Ordnungsanforderungen der installierten Go-Bibliothek geprüft.
Nicht ausgeführt: Dynamische Fehlerreproduktion, Implementierung, Tests und Dropbox-Abfragen.
Nach autorisierter Implementierung ist folgende Verifikation erforderlich:

- Den bestehenden `TestSortEntriesReverse` auf zwei typgleiche Dateien und einen Ordner ausrichten.
- Vor der Korrektur den Verlust der Dateireihenfolge unter Rückwärtssortierung nach Typ reproduzieren.
- Nach der Korrektur müssen die Dateien vor dem Ordner stehen und ihre Eingangsreihenfolge behalten.
- Den fokussierten Test mit `go test ./cmd -run '^TestSortEntriesReverse$'` ausführen.

Test decision: update — `TestSortEntriesReverse` schützt die korrigierte Gleichstandsbehandlung.
Keine neue Testdatei ist vorgesehen.
Die übrigen vorgeschriebenen Source-Submission-Gates bleiben vor einer späteren Einreichung erforderlich.

## Publication-Blockers

- Eine dynamische Offline-Reproduktion des Gleichstandsfehlers fehlt.
- Prior Art und aktuelle Upstream-Beitragsvorgaben sind vor einer Veröffentlichung zu prüfen.
- `Authorized-Work` und `Publication-Target` sind vom Nutzer noch nicht ausgewählt.
- Ein exakter Veröffentlichungsentwurf und dessen Freigabe liegen nicht vor.

## Next-Action

Summary: Gleichstandsfehler offline reproduzieren
Action: Den bestehenden Reverse-Komparator mit zwei typgleichen Dateien und einem Ordner ohne Produktivänderung ausführen.
Done-When: Eingabe, tatsächliche Reihenfolge und erwartete stabile Reihenfolge sind reproduzierbar dokumentiert.
