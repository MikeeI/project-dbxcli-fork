# ISSUE-021 — Listings: Löschmetadaten werden vor ihrer Anreicherung sortiert

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

Root-Cause [S]: `finalizeLsEntries` sortiert rohe Löschmetadaten vor deren Ersatz durch die tatsächlich ausgegebenen Revisionen.
Gesprächsreferenz: `E17`; der Review bestätigt die falsche Reihenfolge bei unterschiedlichen ausgegebenen Sortierwerten.

## Reach-and-Impact

Reach [S]: Listings mit gelöschten Einträgen und `--sort=size` oder `--sort=time` im normalen Finalisierungspfad.
Impact [O]: Ein unbegrenztes `--only-deleted --sort=size` lieferte 100 MiB vor 1 KiB.
Der limitierte Only-deleted-Sonderpfad bereitet zuerst auf und ist deshalb nicht identisch betroffen.

## Evidence

- [S] `cmd/ls.go:185-187` ruft `sortEntries` vor `prepareLsEntries` auf.
- [S] `cmd/ls.go:239-260` ersetzt rohe DeletedMetadata nachträglich durch gelieferte FileMetadata der Revision.
- [S] `cmd/format.go:86-97` liefert für rohe DeletedMetadata Größe null und Nullzeit.
- [S] `cmd/ls.go:190-197,214-223` bereitet im limitierten Only-deleted-Pfad dagegen bereits vor dem Sortieren auf.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierter Listen- und Revisionsgegenpart.
Aufruf: `ls --only-deleted --sort=size --output=json` ohne Limit.
ListFolder liefert `/large`, danach `/small` als DeletedMetadata; Revisionen ergeben 100 MiB und 1 KiB.
Beobachtet [O]: `command_error=<nil>` mit der ausgegebenen Größenfolge `104857600`, danach `1024`.
Der temporäre Harness wurde entfernt; kein realer Papierkorb wurde abgefragt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` wurden auf die gleiche Ursache geprüft.
`ISSUE-003` betrifft den negierten Reverse-Komparator bei Gleichständen und ist eine andere Ursache.
Hier bestehen unterschiedliche echte Sortierwerte; weder Reverse noch eine zugesagte Gleichstandsordnung sind erforderlich.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Eine lokale Reihenfolgekorrektur genügt, ohne globale Top-N- oder Suchauswahlsemantik zu ändern.

## Proposed-Change

In `finalizeLsEntries` zuerst `prepareLsEntries` ausführen und dessen Fehler weitergeben.
Anschließend die vorbereiteten, tatsächlich ausgegebenen Metadaten sortieren.

## Scope-and-Constraints

- Preserve: Bestehende Ergebnisauswahl, Limit-Politik, Pagination und serverseitige Suchreihenfolge erhalten.
- Preserve: Gelöscht-Markierung und zulässige optionale Revisionsanreicherung erhalten.
- Exclude: Keine neue sekundäre Sortierung oder Änderung des Reverse-Komparators in diesem Befund.
- Related: `ISSUE-003` bleibt eigenständig und wird durch diese Korrektur nicht automatisch gelöst.

## Verification

Abgeschlossen [O]: Der unbegrenzte Größen-Smoke reproduziert die falsche Reihenfolge.
Erforderlich: Zwei DeletedMetadata mit erfolgreichen Revisionen unterschiedlicher Größe ohne Limit auflisten.
Erwartet: Aufsteigende tatsächliche Größenfolge; ein ausreichend großes Limit darf die Sortierbedeutung nicht verändern.
Test decision: none; kein Source- oder Testeingriff in diesem Auftrag.

## Publication-Blockers

- Reproduktionsantworten und Text-/JSON-Reihenfolge dauerhaft dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Löschmetadaten-Sortierung gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zur Sortierung nach Auflösung historischer Dateimetadaten recherchieren.
Done-When: Suchabdeckung, relevante Treffer und die Abgrenzung zum Reverse-Gleichstandsfehler sind dokumentiert.
