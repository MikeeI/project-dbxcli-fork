# ISSUE-006 — Downloads: Exportnamen überschreiben andere rekursive Quellen

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: High
Root-Cause-Confidence: High
Finding-Category: Correctness
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Rekursive Downloads bilden unterschiedliche Remote-Quellen ungeprüft auf denselben endgültigen Exportnamen ab.
Gesprächsreferenz: `E1`; der Review bestätigt die Kollision, nicht eine allgemeine Überschreibsperre.

## Reach-and-Impact

Reach [S]: `get --recursive` mit einem Exportdokument und einer gewöhnlichen Datei gleichen endgültigen Namens.
Impact [O]: Der isolierte Command-Smoke meldete beide Dateien als heruntergeladen; lokal blieb nur der Exportinhalt.
Nicht belegt sind eine Live-Dropbox-Kollision oder deren Häufigkeit.

## Evidence

- [S] `cmd/get.go:327-345` leitet zunächst Ziele aus den unterschiedlichen Remote-Pfaden ab.
- [S] `cmd/get.go:524-525` ersetzt bei `dstExplicit=false` den Namen durch `ExportMetadata.Name`.
- [S] `cmd/get.go:545-555` veröffentlicht den Export ohne Besitzprüfung des endgültigen Zielpfads.
- [O] Der aktuelle Source-Abgleich nach `git fetch -q upstream master` zeigt keine Produktionsdifferenz zum angegebenen Upstream.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, Dropbox-SDK v6.6.2; isolierter HTTP-Gegenpart ohne Kontozugriff.
Aufruf: `get /docs <temporäres-ziel> --recursive --output=json` über `RootCmd.ExecuteC`.
Antwortfolge: gewöhnliche `/docs/Report.md`, danach `/docs/Report.paper` mit Exportname `Report.md`.
Beobachtet [O]: `command_error=<nil>` und zwei `downloaded`-Ergebnisse mit demselben lokalen Ziel.
Beobachtet [O]: Der erhaltene Dateiinhalt war ausschließlich `exported-paper`.
Das konkrete Exportnamenpaar wurde injiziert, nicht als Live-Serverantwort aufgezeichnet.
Der temporäre Harness wurde nach dem Audit entfernt.

## Prior-Art

Coverage: `ISSUES.md` und die vorhandenen Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Ursache.
Gaps: Upstream-Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Eine einzelne Zielkollisionskorrektur ist abgrenzbar; Beitragstyp und Ziel bleiben offen.

## Proposed-Change

Endgültige lokale Ziele unterschiedlicher Quellen vor einem überschreibenden Commit auf Kollision prüfen.
Bei einer Kollision klar abbrechen oder eine ausdrücklich festgelegte deterministische Ausweichbenennung verwenden.

## Scope-and-Constraints

- Preserve: Gewolltes Überschreiben vorhandener lokaler Ziele außerhalb einer Mehrquellenkollision erhalten.
- Preserve: Sichere Tempfile-Commits, Exportformat und Remote-Identitäten erhalten.
- Exclude: Keine vollständige Transaktion für rekursive Downloads oder allgemeine Pfadarchitektur einführen.
- Review: `E10` ist zurückgezogen; die injizierte falsche Exportlänge beweist keinen zugesagten Validierungsvertrag.
- Review: Für `E10` wird deshalb kein weiterer bestätigter Befund und keine Implementierungsfreigabe angelegt.

## Verification

Abgeschlossen [O]: Der vorstehende Offline-Smoke reproduziert den stillen Inhaltsverlust im damaligen Audit.
Erforderlich: Beide Auflistungsreihenfolgen mit kollidierenden Export- und gewöhnlichen Dateinamen prüfen.
Erwartet: Kein erfolgreicher Abschluss, der zwei erhaltene Inhalte für denselben endgültigen Zielpfad behauptet.
Test decision: none; dieser Tracking-Auftrag ändert weder Source noch Tests.

## Publication-Blockers

- Reproduktionsharness und vollständige Ausgabe dauerhaft nachvollziehbar dokumentieren; keine Live-Wirkung behaupten.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- `Authorized-Work`, `Publication-Target` und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Exportkollision gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zu Exportnamenkollisionen beim rekursiven Download recherchieren.
Done-When: Suchabdeckung, relevante Treffer und Abgrenzung der tatsächlichen Zielkollision sind dokumentiert.
