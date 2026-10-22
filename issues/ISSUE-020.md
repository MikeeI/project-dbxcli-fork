# ISSUE-020 — Fehleraggregation: einzelne SDK-Ursache wird generisch umcodiert

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Medium
Root-Cause-Confidence: High
Finding-Category: API
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Aggregationsgrenzen ersetzen auch genau eine konkrete SDK-Ursache durch einen generischen Anzahlfehler.
Gesprächsreferenz: `E16`; der Review begrenzt den Befund auf eindeutig klassifizierbare Einzelursachen.
Betroffene Instanzen sind Relocation, rekursiver Download und rekursiver Freigabelink-Download.

## Reach-and-Impact

Reach [S]: Genau eine fehlgeschlagene Teiloperation mit spezifischer SDK-Fehlerklasse.
Impact [S]: JSON-Fehlerklasse, daraus abgeleiteter Exitstatus und vorhandene Rate-Limit-Recoverydaten gehen verloren.
Impact [O]: Der Relocation-Smoke meldete für `from_lookup/not_found` ausschließlich `command_failed`.
Keine neue allgemeine Politik zur Darstellung mehrerer unterschiedlicher Fehler wird vorausgesetzt.

## Evidence

- [S] `cmd/relocation.go:136-155` verliert Ursachen; `cmd/relocation_if_exists.go:131-136` erzwingt `command_failed`.
- [S] `cmd/get.go:354-358` und `cmd/share_link_download.go:282-286` reduzieren Fehler ebenfalls auf Anzahltexte.
- [S] `cmd/share_link_download.go:297-310` trennt SDK-Fehler zusätzlich durch `%v` von ihrer typisierten Identität.
- [S] `cmd/output.go:354-382,417-426,491-501` leitet Exitstatus und Recoverydaten aus Fehlerklasse und Typidentität ab.
- [S] `docs/automation.md:332-351` definiert spezifische Klassen und Exitcodes für Automatisierung.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierter HTTP-Gegenpart.
Aufruf: `cp /missing.txt /new.txt --output=json` mit genau einem Copy-Fehler `from_lookup/not_found`.
Beobachtet [O]: `command_error=cp: 1 error(s)`, `error_code=command_failed`; die Ursache blieb nur in stderr.
Ein weiterer isolierter rekursiver Shared-Link-Ablauf lieferte beim einzigen verschachtelten Listing ein `429`.
Beobachtet [O]: Zurückgegeben wurde ausschließlich `share-link download: 1 error(s)`.
Der spätere generische Exitstatus folgt aus dem Source; der Harness maß ihn nicht separat am Prozess.
Die temporären Harness-Dateien wurden entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Fehleridentitätsursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Eine erhaltene Einzelursache korrigiert den dokumentierten Fehlervertrag ohne neues Mehrfachfehlerprotokoll.

## Proposed-Change

Bei genau einem Fehler die ursprüngliche Ursache mit `%w` und ohne generischen Code-Override weitergeben.
Betroffene Listing-Wrapper ebenfalls auf erhaltende Weitergabe umstellen.

## Scope-and-Constraints

- Preserve: Pfadkontext, vorhandene SDK-Identität und bereitgestellte Retry-Details erhalten.
- Preserve: Erfolgsabläufe und bestehende Behandlung tatsächlich gemischter Mehrfachfehler nicht stillschweigend umdefinieren.
- Review: `errors.Join` allein genügt nicht; `cmd/output.go:617-623` traversiert nur einfache Unwrap-Ketten.
- Exclude: Kein allgemeiner Fehler-Orchestrator oder neues Envelope-Schema ohne gesonderten Vertrag.

## Verification

Abgeschlossen [O]: Einzelne Not-found- und verschachtelte Rate-Limit-Abläufe wurden offline ausgeführt.
Erforderlich: Einzelursachen durch Relocation, rekursiven Download und Shared-Link-Listing führen.
Erwartet: Spezifische Fehlerklasse, zugehöriger Exitstatus und vorhandenes `retry_after_seconds` bleiben erhalten.
Test decision: none; keine Source- oder Teständerung in diesem Auftrag.

## Publication-Blockers

- Einzelursachenreproduktionen einschließlich endgültiger JSON-/Exit-Ausgabe dauerhaft dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Einzelursachen gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zum Erhalt einzelner SDK-Fehler an Aggregationsgrenzen recherchieren.
Done-When: Suchabdeckung, relevante Treffer und die Begrenzung auf Einzelursachen sind dokumentiert.
