# ISSUE-010 — Credentials: veralteter Snapshot verliert zwischenzeitliche Anmeldungen

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: High
Root-Cause-Confidence: High
Finding-Category: Correctness
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Die unkoordinierte Read-modify-write-Sequenz veröffentlicht nach OAuth-I/O den gesamten alten Token-Snapshot.
Gesprächsreferenz: `E5`; der Review bestätigt verlorene Änderungen zwischenzeitlicher Credential-Schreibvorgänge.

## Reach-and-Impact

Reach [S]: Überlappender Login beziehungsweise Refresh auf derselben gespeicherten Token-Map.
Impact [O]: Der isolierte Interleaving-Smoke verlor neue Team-Credentials bei erfolgreichem persönlichem Refresh.
Ein tatsächlich parallel abgeschlossener Live-Login wurde nicht ausgeführt.

## Evidence

- [S] `cmd/auth.go:237-248` liest die vollständige Token-Map vor der OAuth-Operation.
- [S] `cmd/auth.go:254-259,264-274` speichert diesen Snapshot nach Login beziehungsweise Refresh vollständig zurück.
- [S] Die Schreibpfade besitzen keine gemeinsame prozessübergreifende Read-modify-write-Sperre.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierter OAuth- und Account-Gegenpart.
Während des OAuth-Refreshs schreibt der Harness einen neuen `teamManage`-Eintrag in dieselbe temporäre Auth-Datei.
Danach führt der echte `account --output=json`-Command seinen begonnenen Refresh erfolgreich zu Ende.
Beobachtet [O]: `command_error=<nil>`, `interleaved_write_completed=true`, `after_json_valid=true`.
Beobachtet [O]: `new_personal_preserved=true`, `new_team_preserved=false`, `old_team_restored=true`.
Diese kontrollierte Verschachtelung bildet einen zulässigen zweiten Writer ab; Harness und Datei wurden entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Lost-update-Ursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Getrennt von atomarem Dateiaustausch behandeln; beide Eigenschaften werden für robuste Persistenz benötigt.

## Proposed-Change

Lesen, Ändern und Ersetzen mit einer stabilen prozessübergreifenden Sperre koordinieren.
Alle mutierenden Credential-Pfade einschließlich Logout auf dieselbe Koordinationsgrenze abstimmen.

## Scope-and-Constraints

- Preserve: Änderungen anderer Domains und Token-Typen dürfen nicht durch veraltete Snapshots verschwinden.
- Preserve: Sperrfreigabe auf Fehlerpfaden und responsive Wartebehandlung explizit besitzen.
- Exclude: Eine Sperre ausschließlich beim Schreiben oder eine reine Prozess-Mutex reicht nicht.
- Related: `ISSUE-009` besitzt die unabhängige atomare Veröffentlichung bei Schreibfehlern.
- Cost: Eine Sperre über OAuth-I/O kann andere Prozesse verzögern und benötigt eine bewusste Lifecycle-Grenze.

## Verification

Abgeschlossen [O]: Der deterministische Interleaving-Smoke reproduziert den Snapshot-Verlust.
Erforderlich: Refresh und Speicherung eines anderen Token-Typs auf dieselbe Datei verschachteln.
Erwartet: Beide bestätigten Änderungen bleiben erhalten; Fehlerpfade geben jede Sperre zuverlässig frei.
Test decision: none; keine Source- oder Teständerung in diesem Auftrag.

## Publication-Blockers

- Prozessübergreifenden Ablauf und Lifecycle der vorgeschlagenen Sperre belastbar dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Credential-Rennen gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zu Credential-Lost-updates und Prozesskoordination recherchieren.
Done-When: Suchabdeckung, relevante Treffer und Abgrenzung zu atomarem Dateiaustausch sind dokumentiert.
