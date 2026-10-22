# ISSUE-033 — Uploads: revisionsgebundenes Zurückschreiben fehlt im Bearbeitungszyklus

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Low
Root-Cause-Confidence: High
Finding-Category: API
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Vorhandene Upload-Policies bieten Add oder Overwrite, aber keine serverseitige Revisionsvorbedingung.
Gesprächsreferenz: Swarm-Explore-Review `E5`; Quelle: Adjacent Capabilities.
Der Datensatz beschreibt eine optionale Fähigkeit, keinen nachgewiesenen Defekt.

## Reach-and-Impact

Reach [S]: Der dokumentierte Agent-Ablauf verbindet Download, lokale Bearbeitung und separat autorisierte Ersetzung.
Impact [A]: Revisionsgebundenes Schreiben könnte zwischenzeitliche Fremdänderungen vor blindem Ersatz schützen.
Ein tatsächlicher Bedarf oder eine Nutzungshäufigkeit sind nicht belegt.

## Evidence

- [S] `skills/dbxcli/SKILL.md:72-77` beschreibt den Bearbeitungszyklus.
- [S] `cmd/get.go:397-409` und `cmd/json_metadata.go:29-38` liefern Downloadrevisionen.
- [S] `cmd/put.go:590-635` verwendet gemeinsame Commit-Informationen ohne Update-Modus.
- [S] `cmd/put.go:318-330` kann ein Dateiziel automatisch auf ein Ordnerziel umleiten.
- [S] SDK v6.6.2, `files/types.go:263-269,5426-5444`, besitzt striktes revisionsgebundenes Update.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Fähigkeitslücke geprüft.
`ISSUE-013` betrifft dagegen bereits vorhandene mehrdeutige Autorename-Retries.
Gaps: Tatsächlicher Bedarf sowie externe Issues, PRs, Diskussionen und Releases sind ungeklärt.
Contribution fit: Erst Bedarf bestätigen; kein allgemeines Synchronisationsprodukt ableiten.

## Proposed-Change

Bei bestätigtem Bedarf optionales `put --if-rev <revision>` für eine lokale Datei mit explizitem Dateiziel anbieten.
`WriteModeUpdate`, übergebene Revision, `StrictConflict=true` und `Autorename=false` verwenden.
Automatische Zielumleitung und gleichzeitig explizite `--if-exists`-Strategien in diesem Modus ausschließen.

## Scope-and-Constraints

- Preserve: Menschliche Schreibautorisierung bleibt unabhängig von der technischen Revisionsbedingung erforderlich.
- Preserve: Bestehende Upload-Policies außerhalb des neuen Modus erhalten.
- Cost: Command-Metadaten, Dry-run-Eingabe und generierte Verträge müssen die Bedingung ausdrücken.
- Risk: Verlorene Erfolgsantworten können einen späteren Revisionskonflikt erzeugen; keine Exactly-once-Garantie.
- Exclude: Keine Rekursion, automatische Zusammenführung oder Übernahme einer neueren Revision.

## Verification

Quellprüfung und SDK-Vertrag geprüft; kein Bedarfsgespräch, Transfer oder Test ausgeführt.
Vorgeschlagen: Passende Revision muss schreiben; abweichende Revision muss ohne Ersatz oder Autorename scheitern.
Mehrdeutige Erfolgsantworten separat vom Revisionskonflikt bewerten.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Bedarf und eindeutigen Eingabe-, Dry-run- und Konfliktvertrag bestätigen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Bedarf für Revisionsbedingung bestätigen
Action: Einen tatsächlichen Bearbeitungsablauf auf parallele Änderungen und erforderliche Konfliktentscheidungen prüfen.
Done-When: Bedarf, vorhandene Alternativen und die kleinste nützliche Revisionsbedingung sind dokumentiert.
