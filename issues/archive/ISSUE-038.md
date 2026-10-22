# ISSUE-038 — Search: doppelte IDs verbrauchen das Ergebnislimit

State: Archived
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

[S] Der Ausgangsstand zählt rohe Search-Matches ohne seitenübergreifende ID-Deduplizierung gegen das lokale Limit.
Dropbox dokumentiert mögliche Duplikate: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/search-v-2
Die direkte vollständige private Implementierungsfreigabe von `GOAL.md` umfasst diese Korrektur.

## Reach-and-Impact

[S] Wiederholte Datei- oder Ordner-IDs können Ergebnisplätze verbrauchen und unterschiedliche spätere Treffer verdrängen.
[A] Die reale Häufigkeit der Serverduplikate wurde nicht erhoben.
Deduplizierung erzeugt keinen konsistenten Account-Snapshot bei sich ändernden Metadaten.

## Evidence

- [S] `cmd/search.go:collectSearchMatches` besitzt eine ausschließlich auf den Suchlauf begrenzte Accepted-ID-Menge.
- [S] Lokale Filter laufen vor der Aufnahme einer ID; erst akzeptierte unterschiedliche Matches zählen zum Limit.
- [S] Das erste passende vollständige Match bleibt erhalten, ohne Felder unterschiedlicher Versionen zu mischen.
- [O] `TestSearchUsesSearchV2AndCommandOutput` prüft Duplikate und eine später passende zuvor verworfene ID über Seiten.
- [O] `TestSearchLimitCapsResultsAndStopsPagination` unterscheidet Limitabbruch und erschöpfte finale Duplikate.
- [O] `go test ./...`, `go vet ./...` und `go build ./...` bestehen.

## Prior-Art

Lokaler Index und relevante Akten wurden nach Deduplizierung, SearchMatch und Ergebnislimits durchsucht.
`ISSUE-034` besitzt die fehlende Kürzungsdiagnose, nicht die Identitätszählung.
Upstream-Issues und PRs sind für eine mögliche spätere Veröffentlichung noch zu prüfen.

## Proposed-Change

Privat umgesetzt: nach Filterung nach stabiler Datei- oder Ordner-ID deduplizieren und danach das Limit anwenden.
Fehlende IDs und nicht unterstützte Metadaten erzeugen einen Fehler mit bisher akzeptierten Teilergebnissen.
Es gibt keinen Pfad-Fallback, globalen Cache oder zusätzlichen Metadatenabruf.

## Scope-and-Constraints

- Der Speicherbedarf folgt den akzeptierten unterschiedlichen Treffern, nicht allen verworfenen IDs.
- Native Matchinformationen des ersten qualifizierenden Objekts bleiben zusammen erhalten.
- Die API-Reihenfolge vor der lokalen Darstellung entscheidet über die behaltene Version.
- Exakt passende letzte Seiten dürfen nicht allein wegen der Ergebniszahl als gekürzt gelten.

## Verification

Test decision: update — `TestSearchUsesSearchV2AndCommandOutput` schützt Filter-/Deduplizierungsreihenfolge.
`TestSearchLimitCapsResultsAndStopsPagination` schützt den Zusammenhang mit Stopgründen und Limit.
Die Beobachtungen stammen aus der vorhandenen kontextfähigen Client-Testgrenze, nicht aus Live-Kontoabfragen.

## Publication-Blockers

Aktuelle Upstream-Prior-Art, Arbeitsmodus, Ziel, exakter Entwurf und Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: —
Action: None.
Done-When: None.

## Archive

Archive-Reason: Other
Detail: Privater Scope vollständig abgeschlossen; kein aktueller Auftrag für Upstream-Arbeit.
Evidence: Deduplizierung nach lokaler Filterung und ihre Pagination-Nachweise sind in `547752f` gepusht.
Checked: 2026-10-22
