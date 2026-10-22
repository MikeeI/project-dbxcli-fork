# ISSUE-032 — CLI-Validierung: fehlende Operanden werden erst nach Authentifizierung erkannt

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

Root-Cause [S]: Die Argumentanzahlprüfung liegt in Handlern hinter dem gemeinsamen Authentifizierungshook.
Gesprächsreferenz: Swarm-Explore-Review `E4`; Quelle: Friction.

## Reach-and-Impact

Reach [S]: Ungültige Aufrufe von `get`, `put` und `search` durchlaufen zuerst die Authentifizierung.
Impact [S]: Ohne Credentials erscheint zuerst der Loginfehler statt der bekannten Operandenanforderung.
Mit Credentials können Refresh und Namespace-RPC vor dem unvermeidlichen Eingabefehler erfolgen.
Reale Nutzungshäufigkeit und zusätzliche Wartezeit sind nicht erhoben.

## Evidence

- [S] `cmd/root.go:247-249,307` authentifiziert über `PersistentPreRunE`.
- [S] `cmd/get.go:68-70`, `cmd/put.go:291-294` und `cmd/search.go:52-57` prüfen die Argumentanzahl im Handler.
- [S] `go.mod:11` pinnt Cobra v1.10.2.
- [S] Cobra v1.10.2, `command.go:968-986`, führt `ValidateArgs` vor `PersistentPreRunE` aus.
- [S] `cmd/output.go:123-128` besitzt typisierte `invalid_arguments`-Fehler.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Validierungsursache geprüft.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Drei bekannte Befehle mit bestehenden Fehlerverträgen begrenzen.

## Proposed-Change

Dieselben typisierten Argumentanzahlfehler über Cobras `Args` vor der Authentifizierung erzeugen.
Gemeinsame Validierungslogik für weiterhin direkt konsumierte Handler berücksichtigen.

## Scope-and-Constraints

- Preserve: Bestehende JSON-Fehlercodes und Operandenangaben erhalten.
- Compatibility: Bei ungültiger Eingabe und fehlenden Credentials erhält bewusst der Eingabefehler Vorrang.
- Preserve: Gültige Aufrufe behalten ihre Authentifizierung.
- Exclude: Keine allgemeine Neuordnung sämtlicher Flag- oder Pfadprüfungen.

## Verification

Quellprüfung abgeschlossen; keine Command-Smokes oder Tests ausgeführt.
Vorgeschlagen: Operandlose Aufrufe ohne Credentials müssen ohne Authentifizierungsstart `invalid_arguments` liefern.
Gültige Aufrufe müssen ohne Credentials weiterhin `auth_required` liefern.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Fehlerreihenfolge offline über die öffentliche Command-Grenze prüfen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Validierungsreihenfolge offline nachweisen
Action: Ungültige und gültige Aufrufe der drei Befehle mit isolierter fehlender Authentifizierung vergleichen.
Done-When: Fehlercodes, JSON-Details und der Authentifizierungsstart sind für beide Eingabefälle dokumentiert.
