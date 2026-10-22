# ISSUE-022 — Listings: erforderlicher Recoveryfehler wird leerer Erfolg

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Low
Root-Cause-Confidence: High
Finding-Category: Correctness
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Die erforderliche GetMetadata-Recovery nach `not_folder` verwirft ihren Fehler und liefert einen nil-Eintrag.
Gesprächsreferenz: `E18`; der Review begrenzt den Befund auf erforderliche Recovery, nicht optionale Revisionsanreicherung.

## Reach-and-Impact

Reach [S]: Ein zunächst als Ordner erkannter Pfad wird vor ListFolder zur Datei; die Ersatzmetadatenabfrage scheitert.
Impact [O]: Der Offline-Smoke meldete danach ein erfolgreiches leeres JSON-Listing ohne Warnung.
Die geringe Priorität berücksichtigt den zusätzlich erforderlichen Zustandswechsel und Recoveryfehler.
Eine Häufigkeit realer solcher Abläufe wurde nicht erhoben.

## Evidence

- [S] `cmd/ls.go:133-149` fragt bei expliziten Pfaden zunächst Metadaten und anschließend den Ordnerinhalt ab.
- [S] `cmd/ls.go:155-160` beschreibt die notwendige Ersatzantwort, verwirft aber deren Fehler mit `metaRes, _`.
- [S] `cmd/ls.go:264-271` übernimmt nur bekannte Metadatentypen; der resultierende nil-Eintrag verschwindet lautlos.
- [S] `cmd/ls.go:250-253` prüft denselben Metadatenfehler in der analogen Not-file-Recovery korrekt.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierter Metadaten- und ListFolder-Gegenpart.
Aufruf: `ls /changing --output=json` über `RootCmd.ExecuteC`.
Injizierte Folge: FolderMetadata, ListFolder-Fehler `path/not_folder`, anschließend `context.DeadlineExceeded` bei GetMetadata.
Beobachtet [O]: `command_error=<nil>`, `ok=true`, `results=[]`, `warnings=[]`.
Wegen Filterung von nil entsteht in diesem Pfad kein späterer JSON-Typfehler.
Der temporäre Harness wurde entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche erforderliche Recovery-Ursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Die Fehlerweitergabe ist lokal begrenzt und benötigt keinen neuen Listing- oder Timeoutvertrag.

## Proposed-Change

Den Fehler der zweiten Metadatenabfrage erfassen und vor Aufbau der Ergebnisliste zurückgeben.
Die erforderliche Ersatzantwort darf nicht durch einen erfolgreichen leeren Zustand ersetzt werden.

## Scope-and-Constraints

- Preserve: Erfolgreiche Not-folder-Recovery und ursprüngliche SDK-/Timeout-Fehleridentität erhalten.
- Exclude: Fehler optionaler ListRevisions-Anreicherung sind ohne zusätzlichen Pflichtvertrag kein Teil dieses Befunds.
- Exclude: Keine transaktionale Snapshot-Garantie für mehrere Remote-Abfragen erfinden.
- Review: Der frühere breitere Anreicherungsfehler wurde ausdrücklich nicht als bestätigter Befund übernommen.

## Verification

Abgeschlossen [O]: Die vorstehende lokale Antwortfolge reproduziert den falschen leeren Erfolg.
Erforderlich: Ordner-Datei-Wechsel mit anschließend fehlschlagender Ersatzmetadatenabfrage ausführen.
Erwartet: Der letzte Fehler bleibt erkennbar; weder Text noch JSON meldet erfolgreich ein leeres Listing.
Test decision: none; keine Source- oder Teständerung in diesem Auftrag.

## Publication-Blockers

- Antwortfolge und realistische Zustandswechselvoraussetzung dauerhaft nachvollziehbar dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Listing-Recovery gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zu verlorenen Ersatzmetadatenfehlern nach Not-folder-Recovery recherchieren.
Done-When: Suchabdeckung, relevante Treffer und die Abgrenzung zur optionalen Anreicherung sind dokumentiert.
