# ISSUE-031 — Sammeloperationen: Gesamtfehler verwirft bestätigte Einzelresultate

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

Root-Cause [S]: Einzelresultate werden bei einem späteren Gesamtfehler nicht in die strukturierte Fehlerantwort übernommen.
Gesprächsreferenz: Swarm-Explore-Review `E2`; Quelle: Pain Points.

## Reach-and-Impact

Reach [S]: Rekursive Transfers und Sammeloperationen können bereits Änderungen abgeschlossen haben.
Impact [S]: Nach einem Gesamtfehler müssen Aufrufer den Ausführungsstand durch zusätzliche Zustandsabfragen rekonstruieren.
Verschobene oder gelöschte Quellen machen vollständiges Wiederholen nicht gleichwertig.
Häufigkeit und tatsächlicher Wiederherstellungsaufwand sind unbekannt.

## Evidence

- [S] `cmd/get.go:333-358` verwirft gesammelte Downloadresultate bei Fehlern.
- [S] `cmd/put.go:948-954,1023-1024` sammelt erfolgreiche Uploads, gibt sie bei Gesamtfehlern aber nicht zurück.
- [S] `cmd/relocation.go:140-155` und `cmd/rm.go:168-192` besitzen dieselbe Ergebnisgrenze.
- [S] `cmd/json_output.go:9-20` besitzt Fehlerdetails, aber keine Teilergebnisliste.
- [S] `cmd/output.go:109-116` unterstützt zusätzliche strukturierte Fehlerdetails.
- [S] `docs/json-schema/v1/error.schema.json:57-60` schließt unbekannte Detailfelder aus.

## Prior-Art

Coverage: Lokaler Index und vorhandene Akten bis `ISSUE-022` ohne gleiche Ergebnisursache geprüft.
`ISSUE-020` betrifft verlorene SDK-Fehleridentität, nicht verlorene bestätigte Einzelresultate.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Zunächst einen Produktionspfad begrenzen; andere API-Fehler nicht zusammenfassen.

## Proposed-Change

Zunächst nur bei rekursivem JSON-`get` bestätigte Resultate als `error.details.completed_results` erhalten.
Fehlerdetails, öffentliches Schema und Dokumentation gemeinsam erweitern.
`ok:false`, Fehlercode und Exitstatus beibehalten.

## Scope-and-Constraints

- Preserve: Nur bestätigte Erfolge ausweisen; der Ausgang einer fehlgeschlagenen Einzeloperation kann unbekannt bleiben.
- Preserve: Keine Teiltransaktion als Gesamterfolg darstellen.
- Cost: Ältere geschlossene Schemas akzeptieren neue Detailfelder nicht automatisch.
- Exclude: Keine allgemeine Transaktions-, Wiederaufnahme- oder Exactly-once-Architektur.

## Verification

Quellprüfung abgeschlossen; kein Transfer, Test oder Teilergebnis-Smoke ausgeführt.
Vorgeschlagen: Rekursiven Lauf mit erfolgreicher Datei und späterem Fehler prüfen.
Bestätigte Resultate müssen enthalten sein und die Antwort muss das aktualisierte Fehlerschema erfüllen.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Öffentlichen Erweiterungsvertrag und fokussierten Teilergebnisablauf bestätigen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Teilergebnisvertrag offline eingrenzen
Action: Einen rekursiven Download mit bestätigtem Erfolg und späterem Fehler ohne Live-Kontozugriff nachvollziehen.
Done-When: Erhaltene Änderungen, verlorene Ergebnisdaten und der kleinste additive Fehlervertrag sind dokumentiert.
