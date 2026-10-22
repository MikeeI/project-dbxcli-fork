# ISSUE-018 — Freigabelinks: Raw-Änderung liefert veraltete Metadaten zurück

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

Root-Cause [S]: Der Raw-Modify-Pfad verwirft die aktualisierten Linkmetadaten und gibt das vorherige Linkobjekt weiter.
Gesprächsreferenz: `E14`; der Review bestätigt die falsche Zustandsdarstellung nach erfolgreicher Änderung.

## Reach-and-Impact

Reach [S]: `share-link create` verwendet einen bestehenden Link und ändert Einstellungen mit `--disallow-download`.
Impact [O]: Nach modellierter erfolgreicher Änderung enthielt das JSON-Ergebnis weiterhin alte Berechtigungen und Ablaufdatum.
Die tatsächliche Remote-Änderung wurde nicht gegen Dropbox ausgeführt.

## Evidence

- [S] `cmd/share_link_raw.go:83-106` verwirft den erfolgreichen Antwortinhalt von `modify_shared_link_settings`.
- [S] `cmd/share_link_create.go:270-274` gibt danach unverändert das ursprüngliche Linkobjekt zurück.
- [S] `cmd/share_link_json.go:61-69,111-123` serialisiert Ablaufdatum und Berechtigungen direkt aus diesem Objekt.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierter Freigabelink-Gegenpart.
Aufruf: `share-link create /report.pdf --disallow-download --remove-expiration --output=json`.
Create liefert einen bestehenden Link mit aktivem Download und Ablaufdatum `2026-12-01T00:00:00Z`.
Modify liefert erfolgreich neue Metadaten ohne Ablaufdatum und mit `allow_download=false`.
Beobachtet [O]: `command_error=<nil>`, `remote_modify_completed=true`, Ausgabe weiterhin mit altem Datum und `allow_download=true`.
Der Status `existing` allein ist kein Fehler; betroffen sind die veralteten Ergebnisdaten.
Der temporäre Harness wurde entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Link-Ergebnisursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Rückgabe der vorhandenen Modify-Antwort genügt; keine zusätzliche Zustandsabfrage ist erforderlich.

## Proposed-Change

Den Raw-Modify-Client aktualisierte Linkmetadaten zurückgeben lassen und dessen Antwort mit dem vorhandenen Parser dekodieren.
Das aktualisierte Objekt im Existing-Link-Create-Zweig weiterreichen.

## Scope-and-Constraints

- Preserve: Explizite `false`-Werte im Raw-Settings-Payload erhalten.
- Preserve: Wiederverwendete URL, individuelle Fehlerdetails und vorhandene Link-Typen erhalten.
- Exclude: Keine zusätzliche JSON-only-Leseoperation als Ersatz für bereits vorhandene Modify-Metadaten einführen.
- Preserve: Aufrufer des geänderten internen Clientvertrags vollständig berücksichtigen.

## Verification

Abgeschlossen [O]: Der isolierte Command-Smoke reproduziert alte Ergebnisdaten nach erfolgreicher Modify-Antwort.
Erforderlich: Downloads deaktivieren und Ablaufdatum entfernen; Create-Ergebnis mit Modify-Antwort vergleichen.
Erwartet: Beide Antworten zeigen deaktivierten Download und kein Ablaufdatum.
Test decision: none; kein Source- oder Testeingriff in diesem Auftrag.

## Publication-Blockers

- Vollständige Create-Fehler- und Modify-Erfolgsantworten dauerhaft dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Linkmetadaten gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zu veralteten Rückgabedaten nach Raw-Linkänderungen recherchieren.
Done-When: Suchabdeckung, relevante Treffer und betroffene Verbraucher des internen Clientvertrags sind dokumentiert.
