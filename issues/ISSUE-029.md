# ISSUE-029 — Shell-Completion: bekannte Konfliktstrategien werden nicht angeboten

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Low
Root-Cause-Confidence: High
Finding-Category: UI
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Bestehende Shell-Completion ist nicht mit den vorhandenen Enum-Metadaten für Flagwerte verbunden.
Gesprächsreferenz: Swarm-Explore-Review `E15`; Quellen: Friction, Surface.

## Reach-and-Impact

Reach [S]: Interaktive Nutzer wählen Konfliktstrategien für `put`, `cp` und `mv`.
Impact [S]: Werte müssen erinnert oder aus Hilfe übernommen werden statt sie direkt zu vervollständigen.
Die tatsächlich ersparte Zeit und Nutzungshäufigkeit sind nicht erhoben.

## Evidence

- [S] `cmd/completion.go:42-94` erzeugt die bestehenden Shellskripte.
- [S] `cmd/help_manifest.go:116,185,202` besitzt die befehlsspezifischen `--if-exists`-Enums.
- [S] Im geprüften `cmd/` fehlen Registrierungen über `RegisterFlagCompletionFunc`.
- [S] `cmd/root.go:121-128` nimmt Completion von der Authentifizierung aus.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Completion-Anbindung geprüft.
`ISSUE-005` betrifft Manifest-Eigentümerschaft, nicht fehlende interaktive Werteauswahl.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Vorhandene Enums nutzen, kein zweites Register und keine Remote-Abfragen.

## Proposed-Change

Zunächst nur `--if-exists` für `put`, `cp` und `mv` über Cobras Flag-Completion vervollständigen.
Werte aus vorhandenen Metadaten beziehen und irrelevante Dateivorschläge unterdrücken.

## Scope-and-Constraints

- Preserve: Nutzer entscheidet die Strategie selbst; Completion führt keine Dateioperation aus.
- Preserve: `overwrite` ausschließlich bei `put` anbieten.
- Exclude: Keine zweite Werteliste, neue Abhängigkeit oder allgemeine Remote-Pfadvervollständigung.
- Cost: Statische Altskripte gegebenenfalls regenerieren; andere Completion unverändert erhalten.
- Later: Die breitere Surface-Linse nennt auch `--sort`, `--time` und `--order-by`; sie gehören nicht zum ersten Slice.

## Verification

Quellprüfung abgeschlossen; kein Completion-Smoke ausgeführt.
Vorgeschlagen: `dbxcli __complete put --if-exists a` muss `autorename` ohne Dateivorschläge liefern.
Danach befehlsspezifische Werte bei `cp` und `mv` prüfen; dort darf `overwrite` nicht erscheinen.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Native Completion-Protokollantwort und Metadatenanbindung offline prüfen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Enum-Completion offline abgrenzen
Action: Den vorhandenen Completion-Weg und die drei Konfliktstrategie-Enums auf eine gemeinsame Anbindung prüfen.
Done-When: Befehlsspezifische Werte und unterdrückte Dateivorschläge sind ohne Anmeldung reproduzierbar spezifiziert.
