# ISSUE-016 — Relocation: roher Quell-Basename erzeugt falsches Ordnerziel

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

Root-Cause [S]: Zielnamensbildung verwendet die rohe Quelle, während der anschließende Auftrag die Quelle normalisiert.
Gesprächsreferenz: `E12`; der Review bestätigt die unterschiedliche Pfadidentität in derselben Operation.

## Reach-and-Impact

Reach [S]: `cp` und `mv` in Ordnerziele mit akzeptierten normalisierbaren Quellargumenten wie `/Reports/.`.
Impact [O]: Die Offline-Vorschau verwendete `/Archive` statt des beabsichtigten `/Archive/Reports`.
Impact [S]: Tatsächliche Aufträge können am falschen Ziel kollidieren oder bei Autorename außerhalb des Zielordners landen.
Eine Remote-Verschiebung wurde nicht ausgeführt.

## Evidence

- [S] `cmd/dropbox_path.go:33-35` berechnet den Zielnamen mit `path.Base` des rohen Quellarguments.
- [S] `cmd/relocation.go:90-92` ruft erst nach dieser Berechnung `makeRelocationArg` auf.
- [S] `cmd/root.go:159-180` normalisiert die Quelle innerhalb dieses anschließenden Aufrufs.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierte Source-Metadaten.
Aufruf: `mv /Reports/. /Archive/ --if-exists=autorename --dry-run --output=json`.
Beobachtet [O]: `command_error=<nil>`, `from_path=/Reports`, `to_path=/Archive`.
Erwartet: `from_path=/Reports`, `to_path=/Archive/Reports`.
Der Slash am Ziel erzwingt Ordnerinterpretation; kein schwankender Zielzustand ist erforderlich.
Der temporäre Harness wurde entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Normalisierungsursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Gemeinsame Quellnormalisierung vor der vorhandenen Zielbildung reicht als begrenzte Korrektur.

## Proposed-Change

Gewöhnliche Quellpfade vor der Basename-Ableitung normalisieren.
Dieselbe normalisierte Quelle für Zielbildung und Relocation-Auftrag verwenden.

## Scope-and-Constraints

- Preserve: Normale Dateinamen, explizite Ordnernotation und Mehrfachquellen erhalten.
- Preserve: Die vorhandene Dropbox-Pfadkonvention mit `path` statt lokalen Dateisystemregeln erhalten.
- Exclude: Keine neue Zusage für bisher nicht unterstützte opake Relocation-Referenzen ableiten.
- Related: `ISSUE-012` behandelt unabhängig davon fehlgeschlagene Zieltypabfragen.

## Verification

Abgeschlossen [O]: Der Dry-run-Smoke reproduziert den falschen Zielpfad ohne Remote-Mutation.
Erforderlich: Den geplanten Auftrag für `/Reports/.` nach `/Archive/` prüfen.
Erwartet: `FromPath=/Reports` und `ToPath=/Archive/Reports`.
Test decision: none; keine Source- oder Teständerung in diesem Auftrag.

## Publication-Blockers

- Reproduktionsauftrag und Normalisierungsvertrag dauerhaft dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Pfadnormalisierung gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zur Quellnormalisierung vor Relocation-Zielbildung recherchieren.
Done-When: Suchabdeckung, relevante Treffer und Abgrenzung zur Zieltypprüfung sind dokumentiert.
