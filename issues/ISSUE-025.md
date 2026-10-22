# ISSUE-025 — Pipe-Uploads: vollständige Eingabe benötigt temporären Speicher

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Medium
Root-Cause-Confidence: High
Finding-Category: Other
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Stdin wird vor dem Upload vollständig in eine temporäre Datei kopiert.
Gesprächsreferenz: Swarm-Explore-Review `E10`; Quelle: Pain Points.
Dies ist eine belegte Betriebsbelastung, kein nachgewiesener Implementierungsdefekt.
Die ID ist unabhängig vom zurückgezogenen `E10` des semantischen Audits.

## Reach-and-Impact

Reach [S]: Der dokumentierte Archiv-Workflow `tar … | dbxcli put - …` nutzt vollständiges Spooling.
Impact [S]: Temporärer Speicher für die gesamte Eingabe ist erforderlich; Erzeugung und Upload überlappen nicht.
Konkrete Platzprobleme, Eingabegrößen und Nutzungshäufigkeit wurden nicht beobachtet.

## Evidence

- [S] `README.md:78-82` dokumentiert den Pipe-Workflow.
- [S] `cmd/put.go:442-447` wartet vor dem Upload auf vollständiges Spooling.
- [S] `cmd/stdin.go:13,26-37` verwendet `os.CreateTemp` mit leerem Verzeichnisargument.
- [S] `docs/automation.md:320-322` dokumentiert bereits die vollständige Speicheranforderung.
- [O] `go doc os.CreateTemp` bestätigt die Verwendung von `os.TempDir` bei leerem Verzeichnisargument.
- [O] `go doc os.TempDir` bestätigt unter Unix `$TMPDIR`, sonst `/tmp`.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-017` ohne gleiche Betriebsbelastung geprüft.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Bestehenden Betriebsausweg erklären; eine vermeidbare Spooling-Anforderung ist nicht bewiesen.

## Proposed-Change

In der vorhandenen Pipe-Dokumentation die Unix-Auswahl über `TMPDIR` erläutern.
Ein bereits vorhandenes, schreibbares Verzeichnis mit ausreichend freiem Speicher voraussetzen.
Vollständiges Spooling und wiederholbares Lesen beibehalten.

## Scope-and-Constraints

- Preserve: Retry-Fähigkeit, Tempfile-Cleanup und bestehende Stdin-Semantik erhalten.
- Exclude: Kein unbelegter Streaming-Umbau und keine neue CLI-Einstellung.
- Limit: `TMPDIR` verlagert den Platzbedarf, beseitigt weder dessen Umfang noch die Vorlaufphase.
- Scope: Die genannte Umgebungsvariable ist der Unix-Vertrag; andere Plattformen haben eigene Temp-Auswahlregeln.

## Verification

Standardbibliotheksvertrag geprüft; kein Spool-Smoke oder Transfer ausgeführt.
Vorgeschlagen: Lokalen Spool-Smoke mit gesetztem `TMPDIR` auf Ablageort, Vollständigkeit und Cleanup prüfen.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Ablageort und Cleanup mit vorhandenem `TMPDIR` offline prüfen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: TMPDIR-Ausweg offline bestätigen
Action: Einen lokalen Spooling-Ablauf mit eigens gewähltem temporärem Verzeichnis ohne Dropbox-Transfer prüfen.
Done-When: Ablageort, vollständige Eingabe und Entfernung der Spool-Datei sind dokumentiert.
