# ISSUE-004 — Uploads: doppelte Traversierung für Planung und Ausführung

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Low
Root-Cause-Confidence: High
Finding-Category: Maintainability
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Rekursive Upload-Planung und Ausführung besitzen getrennte Implementierungen derselben Traversierungsregeln.
Der Kommentar bei `cmd/put.go:842` verlangt ausdrücklich übereinstimmende Semantik.
Gesprächsreferenz: `S1`; der anschließende Review bestätigt den begrenzten Wartbarkeitsvorschlag.

## Reach-and-Impact

Reach [S]: `put --recursive` verwendet je nach Dry-run und Ausgabeformat einen der beiden Abläufe.
Impact [S]: Änderungen an Pfadbildung, Filtern oder Traversierungsregeln müssen beide Implementierungen berücksichtigen.
Ein tatsächlicher Drift, Laufzeitdefekt oder Geschwindigkeitsgewinn wurde nicht nachgewiesen.
Die geschätzte Nettoeinsparung beträgt 40–50 Produktionszeilen, nicht eine gemessene Patchbilanz.

## Evidence

- [S] `cmd/put.go:333-362` verteilt auf Dry-run, normalen Textmodus und normalen JSON-Modus.
- [S] `cmd/put.go:844-911` implementiert die separate Planung mit zwei `filepath.WalkDir`-Durchläufen.
- [S] `cmd/put.go:913-1027` wiederholt Traversierung, Symlink-Filter und relative Pfadbildung für echte Uploads.
- [S] `cmd/put.go:207` besitzt mit `putOptions.dryRun` bereits die erforderliche Modusinformation.
- [S] LSP bestätigt genau einen Aufrufer von `plannedPutRecursiveResults`, bei `cmd/put.go:335`.
- [S] `cmd/put.go:569` und `cmd/put.go:1039` zeigen reale Seiteneffekte unterhalb der Traversierung.
- [S] `cmd/put_test.go:1112-1167` verbietet Upload und Ordnererzeugung während rekursivem Dry-run.
- [S] `cmd/put_test.go:1258-1279` verbietet zusätzliche Ordner-Metadatenabfragen im normalen Textmodus.
- [O] Der lokale Vergleich mit dem angegebenen `upstream/master` zeigt keine Differenz in `cmd/put.go`.

## Prior-Art

Coverage: Der lokale Index und alle vorhandenen Finding-Akten wurden auf dieselbe Ursache geprüft.
`ISSUE-001` betrifft Chunk-Puffer und ist trotz gemeinsamer Datei eine andere Ursache.
Gaps: Upstream-Issues, PRs, Diskussionen und Releases wurden für diesen Vorschlag noch nicht recherchiert.
Contribution fit: Offen; der begrenzte Wartbarkeitsgewinn muss gegen bestehende Arbeiten abgewogen werden.
Der Source-Abgleich verwendet den vorhandenen Upstream-Ref, keinen erneuten Netzwerkabruf.

## Proposed-Change

Den separaten Planer entfernen und seinen Aufrufer auf `putRecursiveInternal(src, dst, opts, true)` umstellen.
An Datei-, Wurzel- und weiteren Ordneraktionen explizit zwischen Planung und bestehender Ausführung unterscheiden.
Planungszweige müssen vor Client-Erzeugung, Statusausgabe und mutierenden Aufrufen zurückkehren.
`renderPlannedPutResults`, Ergebnisgeneratoren und `collectResults` bleiben erhalten.
Die weiterhin gültige Synchronisierungsrationale beim gemeinsamen Ablauf erhalten.
Die Änderung ist ein Vorschlag und nicht zur Implementierung autorisiert.

## Scope-and-Constraints

- Preserve: Zwei Durchläufe und die Reihenfolge Dateien → Wurzel → weitere Verzeichnisse erhalten.
- Preserve: `dirsWithFiles` markiert weiterhin nur direkte Eltern regulärer Dateien.
- Preserve: Symlink-Warnungen, Spezialdateifilter, Pfadnormalisierung und Fehlerdetails erhalten.
- Preserve: Walk-Fehler brechen unmittelbar ab; echte Uploadfehler werden weiterhin gesammelt.
- Preserve: Normaler Textmodus sammelt keine Ergebnislisten und verursacht keine zusätzlichen Metadatenabfragen.
- Preserve: Transfer-, Retry-, Auth-, Timeout- und Cleanup-Funktionen bleiben unverändert.
- Exclude: Kein generischer Visitor, Aktionsdispatcher oder vorab materialisierter Gesamtplan.
- Cost: Die bisher getrennte Schreibfreiheitsgrenze des Planers wird durch drei sichtbare Kontrollflussgrenzen ersetzt.
- Cost: Geschätzt 68–71 Produktionszeilen entfernen und 20–30 ergänzen; Tests und generierte Inhalte zählen nicht.

## Shared-Change-Pressure

Gemeinsam sind beide Walk-Phasen, lokale Pfadbildung, Filter, Symlink-Warnungen und Walk-Fehlerübersetzung.
Der bestehende Kommentar dokumentiert die erforderliche Synchronisierung bereits ausdrücklich.
Die Zusammenführung lohnt nur ohne zusätzliche Abstraktionsschicht und ohne geänderte Fehlerzeitpunkte.

## Verification

Abgeschlossen [O]: Quelltext, Aufrufer und vorhandene Testverträge wurden im Audit und Review geprüft.
Nicht ausgeführt: Implementierung, Tests, Transfer-Szenarien und Messung einer tatsächlichen Löschbilanz.
Nach autorisierter Implementierung ist folgende Verifikation erforderlich:

- `go test ./cmd -run '^TestPut(Recursive|JSONRecursive|JSONDryRunRecursive)'` muss weiterhin bestehen.
- Im Diff Dry-run-Schreibfreiheit, Aktionsreihenfolge und unveränderte Fehlerzeitpunkte prüfen.
- Die rekursiven Ergebnistests vergleichen teilweise Mengen und beweisen allein keine vollständige Reihenfolge.
- Die vorgeschriebenen Source-Submission-Gates bleiben vor einer späteren Einreichung erforderlich.

## Publication-Blockers

- Aktuellen kanonischen Upstream-Stand und Prior Art vor einer Beitragsempfehlung prüfen.
- Verhaltensgleichheit und tatsächliche Nettoeinsparung sind ohne Implementierung nicht dynamisch bestätigt.
- `Authorized-Work` und `Publication-Target` sind vom Nutzer noch nicht ausgewählt.
- Ein exakter Veröffentlichungsentwurf und dessen Freigabe liegen nicht vor.

## Next-Action

Summary: Traversierungsrefaktorierung gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zur rekursiven Upload-Planung und Traversierungsvereinheitlichung recherchieren.
Done-When: Relevante Treffer, Dublettenbewertung und Nutzen des begrenzten Vorschlags sind dokumentiert.
