# ISSUE-011 — Teamaufnahme: individuelles Fehlerergebnis wird als Erfolg bestätigt

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

Root-Cause [S]: `addMember` prüft den äußeren SDK-Fehler, nicht den Fehlerstatus eines abgeschlossenen individuellen Ergebnisses.
Gesprächsreferenz: `E6`; der Review bestätigt falsche Text-, Envelope- und Abschlusssemantik.

## Reach-and-Impact

Reach [S]: `team add-member` mit synchron abgeschlossenem Ergebnis wie `team_license_limit` oder `user_creation_failed`.
Impact [O]: Der isolierte Command-Smoke bestätigte die nicht erfolgte Aufnahme ausdrücklich als Erfolg.
JSON enthält den individuellen Fehlertag weiterhin; dieser mildert den Informationsverlust, nicht den falschen Gesamterfolg.

## Evidence

- [S] `cmd/add-member.go:38-49` rendert jede Antwort ohne äußeren Go-Fehler als reguläres Ergebnis.
- [S] `cmd/add-member.go:52-55` gibt bei `complete` unabhängig von den individuellen Ergebnissen Erfolg aus.
- [S] `cmd/team_json.go:286-298` bildet individuelle Fehler auf `completed` innerhalb eines erfolgreichen Envelopes ab.
- [S] SDK v6.6.2, `dropbox/team/types.go:2930-2932`: Nur `success` bedeutet tatsächlich hinzugefügt.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierte Team-API-Antwort ohne Kontooperation.
Aufruf: `team add-member user@example.com First Last`, zusätzlich separat mit `--output=json`.
Injizierte Antwort: `complete` mit genau einem individuellen `team_license_limit`-Ergebnis.
Beobachtet [O]: `command_error=<nil>` und `User successfully added to the team.`.
Beobachtet [O]: JSON enthält `ok=true`, `status=completed` und den individuellen Fehler unter `results`.
Der temporäre Harness wurde entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Ergebnisursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Ergebnisprüfung an der Command-Grenze kann Text und JSON gemeinsam korrigieren.

## Proposed-Change

Synchron abgeschlossene individuelle Ergebnisse vor dem Rendern auf `success` prüfen.
Andere Ergebnisse als Operationsfehler mit Email und ursprünglichem Ergebnistag zurückgeben.

## Scope-and-Constraints

- Preserve: `async_job_id` bleibt eine gestartete und nicht fälschlich abgeschlossene Operation.
- Preserve: Erfolgreiche Aufnahme und individuelle Ergebnisinformationen erhalten.
- Exclude: Keine automatische Job-Polling-Architektur oder Änderung des Team-Entfernens einführen.

## Verification

Abgeschlossen [O]: Beide Ausgabemodi wurden mit dem beschriebenen individuellen Fehler offline ausgeführt.
Erforderlich: `complete` mit `team_license_limit` über dieselbe öffentliche Command-Grenze prüfen.
Erwartet: Operationsfehler ohne ausdrückliche Text-Erfolgsmeldung und ohne `ok=true`.
Test decision: none; kein Source- oder Testeingriff durch diesen Tracking-Auftrag.

## Publication-Blockers

- Dauerhaft nachvollziehbare Antwortfolge und Ausgabe beider Modi dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Teamaufnahme gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zur Auswertung individueller MembersAdd-Ergebnisse recherchieren.
Done-When: Suchabdeckung, relevante Treffer und Abgrenzung zu asynchroner Jobverfolgung sind dokumentiert.
