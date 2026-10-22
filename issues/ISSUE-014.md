# ISSUE-014 — Uploads: Close-Recovery behandelt zulässiges closed nicht

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Medium
Root-Cause-Confidence: High
Finding-Category: Reliability
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Recovery akzeptierter Appends berücksichtigt ausschließlich Offset-Antworten, nicht den eigenen Close-Übergang.
Gesprächsreferenz: `E9`; der Review begrenzt die Aussage auf den zulässigen, nicht behandelten `closed`-Zweig.

## Reach-and-Impact

Reach [S]: Ein eigener finaler Close wurde möglicherweise akzeptiert, seine Bestätigung ging jedoch verloren.
Impact [O]: Der serielle Offline-Smoke brach bei danach geliefertem `closed` vor jedem Finish-Aufruf ab.
Nicht nachgewiesen ist, welchen Fehler Dropbox bei genau diesem Retry tatsächlich bevorzugt zurückgibt.
Es wird nicht behauptet, dass jeder verlorene Close-Erfolg scheitert.

## Evidence

- [S] `cmd/put.go:58-76` wiederholt Append und akzeptiert ausschließlich `incorrect_offset` mit erwartetem Folgeoffset.
- [S] `cmd/put.go:165` markiert den finalen Bytebereich als schließend.
- [S] SDK v6.6.2, `dropbox/files/types.go:4688-4691,4864-4866`: Close beendet Appends; `closed` ist ein eigener Fehlerzustand.
- [S] SDK v6.6.2, `dropbox/files/types.go:4896-4903`: Der Session-Cursor schützt gegen Datenverlust und Duplikation.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierter Session-Gegenpart.
Aufruf: `put <33-MiB-datei> /target --workers=1 --output=json`.
Der Gegenpart akzeptiert den finalen Close, liefert zunächst `io.ErrUnexpectedEOF` und beim Retry `closed`.
Beobachtet [O]: `command_error=closed/...`, `append_calls=4`, `finish_calls=0`.
Serielle Worker schließen das unabhängige Überholproblem aus `ISSUE-008` für diese Reproduktion aus.
Der konkrete Fehlerprioritätsablauf ist injiziert; temporärer Harness und Dateien wurden entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Close-Recovery-Ursache.
Gaps: Live-Fehlerpriorität sowie externe Issues, PRs, Diskussionen und Releases sind nicht geklärt.
Contribution fit: Eine eng begrenzte Recovery darf beliebiges `closed` nicht in Erfolg umdeuten.

## Proposed-Change

Nach bestätigten vorherigen Appends einen mehrdeutigen eigenen finalen Close über `Finish` auflösen.
Den erwarteten Gesamt-Offset prüfen lassen; beliebiges `closed` ohne diese Voraussetzungen bleibt ein Fehler.

## Scope-and-Constraints

- Preserve: Ursprüngliche Sessionidentität, vollständige Bytebereiche und Fehleridentität erhalten.
- Prerequisite: Die Close-Barriere aus `ISSUE-008` muss frühere Appends zuverlässig abschließen.
- Exclude: Keine pauschale Behandlung aller Closed-Zustände als erfolgreicher Inhaltstransfer.
- Review: Eine SDK-zulässige Antwort ist nicht mit einem aufgezeichneten Live-Servertrace gleichzusetzen.

## Verification

Abgeschlossen [O]: Der serielle Protokoll-Smoke beweist den unbehandelten Closed-Zweig.
Erforderlich: Finale Bestätigung verlieren und danach `closed` liefern.
Erwartet: Vollständiger bestätigter Commit oder unverfälschter Finish-Prüfungsfehler, kein ungeprüfter Erfolg.
Test decision: none; kein Source- oder Testeingriff in diesem Auftrag.

## Publication-Blockers

- Tatsächliche Dropbox-Fehlerpriorität oder ausreichend eindeutigen API-Vertrag für diesen Ablauf klären.
- Prior Art, aktuelle Upstream-Beitragsvorgaben und Reproduktionsartefakte prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Close-Recovery am Protokoll eingrenzen
Action: Den Fehlervertrag des erneut gesendeten finalen Close-Appends anhand autoritativer Protokollevidenz klären.
Done-When: Unterstützte Antwortzweige und die Grenze zwischen zulässigem Zustand und beobachtetem Verhalten sind dokumentiert.
