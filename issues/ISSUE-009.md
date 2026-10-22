# ISSUE-009 — Credentials: Schreibfehler zerstört die gültige Auth-Datei

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: High
Root-Cause-Confidence: High
Finding-Category: Reliability
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: `writeTokens` überschreibt die autoritative Auth-Datei direkt statt sie erst vollständig zu ersetzen.
Gesprächsreferenz: `E4`; der Review bestätigt den Verlust des vorherigen gültigen Zustands bei Schreibfehlern.

## Reach-and-Impact

Reach [S]: Login und automatischer Refresh gespeicherter Credentials verwenden denselben Writer.
Impact [O]: Ein fehlgeschlagener Refresh beschädigte im isolierten Smoke auch die nicht aktualisierten Token-Typen.
Betroffen waren ausschließlich künstliche Credentials in einer temporären Datei.

## Evidence

- [S] `cmd/auth.go:217-228` serialisiert die gesamte Token-Map und schreibt sie mit `os.WriteFile` direkt zurück.
- [S] `cmd/auth.go:259,274` verwendet diesen Writer für Login und automatischen Refresh.
- [S] Der Go-Vertrag von `os.WriteFile` erlaubt nach vorheriger Trunkierung einen teilweise geschriebenen Fehlerzustand.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; künstliche OAuth-Antwort und temporäre Auth-Datei.
Der echte `account --output=json`-Command refresht abgelaufene künstliche Credentials.
Ein prozesslokales `RLIMIT_FSIZE=64` erzwingt den Schreibfehler nach zuvor gültiger 216-Byte-Datei.
Beobachtet [O]: `command_error=write ...: file too large`, `after_bytes=64`, `after_json_valid=false`.
Beobachtet [O]: Weder persönlicher noch Team-Refresh-Token blieb erhalten.
Limits wurden vor der Auswertung wiederhergestellt; temporärer Harness und Datei wurden entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Persistenzursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Atomare Veröffentlichung ist unabhängig von der Prozesskoordination in `ISSUE-010`.

## Proposed-Change

Geschützte Nachbardatei vollständig schreiben und schließen.
Die bestehende Auth-Datei erst nach erfolgreicher Fertigstellung atomar ersetzen.

## Scope-and-Constraints

- Preserve: Token-Map, Dateiauswahl, JSON-Kompatibilität und beabsichtigte Dateiberechtigungen erhalten.
- Preserve: Vor einer erfolgreichen Veröffentlichung müssen die vorherigen Credentials unverändert lesbar bleiben.
- Exclude: Ohne verifizierte Sync- und Storage-Garantien keine Stromausfall-Durabilität behaupten.
- Related: Atomarer Austausch allein verhindert die veralteten Snapshots aus `ISSUE-010` nicht.

## Verification

Abgeschlossen [O]: Der vorstehende Schreibfehler-Smoke reproduziert den Verlust des gültigen Altzustands.
Erforderlich: Schreibfehler während des Refreshs vor erfolgreicher Veröffentlichung injizieren.
Erwartet: Fehler wird zurückgegeben; ursprüngliche Datei bleibt unverändert und vollständig lesbar.
Test decision: none; keine Source- oder Teständerung in diesem Tracking-Auftrag.

## Publication-Blockers

- Ausführbare Fehlerinjektion und relevante Dateiberechtigungen für den Korrekturvorschlag dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Credential-Persistenz gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zu atomarem Speichern und Refresh-Schreibfehlern recherchieren.
Done-When: Suchabdeckung, relevante Treffer und Abgrenzung zum Snapshot-Rennen sind dokumentiert.
