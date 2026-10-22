# ISSUE-012 — Zielauflösung: fehlgeschlagene Typprüfung erzeugt falschen Skip

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

Root-Cause [S]: `isRemoteFolder` übersetzt jeden Metadatenfehler in einen vermeintlich festgestellten Nicht-Ordner.
Gesprächsreferenz: `E7`; der Review korrigiert die Fehlerpriorität von High auf Medium.

## Reach-and-Impact

Reach [S]: `cp` und `mv` mit einem einzelnen Ziel ohne abschließenden Slash und vorübergehend fehlgeschlagener Typabfrage.
Impact [O]: Der Offline-Smoke meldete Erfolg durch Skip des Ordners, obwohl dessen beabsichtigte Zieldatei fehlte.
Die Quelle blieb erhalten; eine irreversible Wirkung oder laufende Remote-Zustandsänderung war nicht erforderlich.

## Evidence

- [S] `cmd/dropbox_path.go:40-49` unterscheidet Metadatenfehler nicht von einem bestätigten Nicht-Ordner.
- [S] `cmd/relocation.go:88-92` bestimmt daraus den endgültigen Zielpfad.
- [S] `cmd/relocation.go:111-120` und `cmd/relocation_if_exists.go:77-86` überspringen anhand dieses falschen Ziels.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierter HTTP-Gegenpart.
Aufruf: `cp /source/report.txt /Archive --if-exists=skip --output=json`.
Unveränderter Gegenpartzustand: `/Archive` ist ein Ordner, `/Archive/report.txt` fehlt.
Die erste Typabfrage scheitert transient; die zweite Abfrage liefert dieselben gültigen Ordner-Metadaten.
Beobachtet [O]: `command_error=<nil>` und `status=skipped` mit `to_path=/Archive`.
Kein Copy-Aufruf erfolgte; der temporäre Harness wurde entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Typabfrageursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Eine lokale Fehlerweitergabe kann die Zielidentität vor jeder Mutation stabil halten.

## Proposed-Change

Ordnererkennung mit explizitem Fehlerergebnis ausführen.
Nur bestätigtes Nichtvorhandensein als fehlendes Ordnerziel behandeln; andere Fehler vor Zielbildung zurückgeben.

## Scope-and-Constraints

- Preserve: Bestehende Ordner ohne abschließenden Slash weiterhin als Ordnerziele erkennen.
- Preserve: Explizite Ordnernotation, Mehrfachquellen und beabsichtigte Skip-Semantik erhalten.
- Exclude: Kein allgemeiner TOCTOU- oder transaktionaler Pfadvertrag wird behauptet.
- Review: Die Prioritätskorrektur berücksichtigt erforderlichen transienten Fehler und erhaltene Quelldaten.

## Verification

Abgeschlossen [O]: Der beschriebene Fehlerablauf wurde ohne Dropbox-Zugriff reproduziert.
Erforderlich: Erste Zieltypabfrage fehlschlagen lassen und nachfolgend gültige Ordner-Metadaten anbieten.
Erwartet: Der erste Fehler bleibt sichtbar; weder falscher Skip noch Mutation erfolgt.
Test decision: none; keine Source- oder Teständerung in diesem Auftrag.

## Publication-Blockers

- Reproduktionsantworten und konkrete Fehlerklasse dauerhaft nachvollziehbar dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Zielauflösung gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zu verschluckten Zieltypabfragen bei Copy und Move recherchieren.
Done-When: Suchabdeckung, Treffer und die im Review begrenzte Auswirkung sind dokumentiert.
