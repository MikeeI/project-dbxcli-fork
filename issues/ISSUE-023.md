# ISSUE-023 — Rekursiver Upload: Textmodus verschweigt übersprungene Symlinks

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Medium
Root-Cause-Confidence: High
Finding-Category: UI
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Symlink-Auslassungen werden ausschließlich bei aktivierter Ergebnissammlung erfasst.
Gesprächsreferenz: Swarm-Explore-Review `E8`; Quelle: Operator.

## Reach-and-Impact

Reach [S]: Reguläres `put -r` verwendet im Textmodus keine Ergebnissammlung.
Impact [S]: Ausgelassene Symlinks sind dort nicht sichtbar, während JSON bereits `skipped_symlink` ausweist.
Die tatsächliche Nutzungshäufigkeit und betroffene Datenmengen sind nicht erhoben.

## Evidence

- [S] `cmd/put.go:348-350` wählt den Textpfad ohne Ergebnissammlung.
- [S] `cmd/put.go:927-935` bindet die Symlink-Warnung an `collectResults`.
- [S] `cmd/json_output.go:29-33` besitzt den bestehenden Warncode `skipped_symlink`.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Auslassungsdiagnostik geprüft.
`ISSUE-004` betrifft gemeinsame Traversierung, nicht die unterschiedliche Sichtbarkeit der Auslassungen.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Beobachtbare Ausgabe lokal verbessern, Symlinks weiterhin bewusst nicht verfolgen.

## Proposed-Change

Tatsächlich übersprungene Symlinks ohne vollständige Ergebnissammlung zählen.
Im Textmodus einmalig auf stderr die Anzahl und nicht übertragenen Inhalte klar benennen.

## Scope-and-Constraints

- Preserve: Symlink-Ziele nicht hochladen und vorhandene JSON-Pfadwarnungen erhalten.
- Preserve: Keine neuen Flags oder Exitcodes.
- Exclude: Keine Warnung pro Symlink und keine zusätzliche Ergebnisliste im Textmodus.
- Related: Bei einer späteren Traversierungsvereinheitlichung aus `ISSUE-004` die Diagnose mitführen.

## Verification

Quellprüfung abgeschlossen; kein Transfer oder Test ausgeführt.
Vorgeschlagen: Reguläre Datei und Symlink in Text und JSON vergleichen.
Ohne Symlink darf keine Zusatzwarnung erscheinen; Zielinhalte bleiben unverändert ausgeschlossen.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Auslassung und einmalige Zusammenfassung über denselben Traversierungspfad prüfen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Symlink-Diagnose offline vergleichen
Action: Einen identischen Verzeichnisbaum mit regulärer Datei und Symlink in beiden Ausgabemodi nachvollziehen.
Done-When: Unterschiedliche Warnungsausgabe und eine ergebnislistenfreie Textzusammenfassung sind abgegrenzt.
