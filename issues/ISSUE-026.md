# ISSUE-026 — Listing: Löschzustand wird über Textmarkierungen transportiert

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Low
Root-Cause-Confidence: High
Finding-Category: Maintainability
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: Listing-Aufbereitung dekoriert `PathDisplay`; JSON liest den Löschzustand aus dieser Darstellung zurück.
Gesprächsreferenz: Swarm-Explore-Review `E11`; Quelle: Seam.

## Reach-and-Impact

Reach [S]: Gelöschte Einträge werden für Text und strukturiertes `ls` gemeinsam aufbereitet.
Impact [S]: Änderungen an Textmarkierungen betreffen Datenaufbereitung und JSON-Zustandsrekonstruktion gemeinsam.
Eine tatsächlich falsche Löschklassifikation oder Produktionsstörung ist nicht nachgewiesen.

## Evidence

- [S] `cmd/ls.go:96-104` verändert `PathDisplay` mit der Löschmarkierung.
- [S] `cmd/ls.go:239-274` verbindet Revisionsanreicherung, Filter und Dekoration.
- [S] `cmd/ls.go:324-340` entfernt `<<…>>` und setzt daraus `Deleted=true`.
- [S] JSON besitzt bereits ein eigenes Löschzustandsfeld.

## Prior-Art

Coverage: Index und vorhandene Akten bis `ISSUE-022` ohne gleiche Listing-Zustandsursache geprüft.
`ISSUE-021` besitzt die unabhängige Sortierung nach Revisionsanreicherung, nicht die Text-/Zustandskopplung.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind nicht recherchiert.
Contribution fit: Lokale Darstellungskopplung entfernen, keinen unbelegten Korrektheitsfehler behaupten.

## Proposed-Change

Metadaten und ursprünglichen Löschzustand in einem lokalen Listing-Ergebnis explizit erhalten.
Text ausschließlich beim Rendern dekorieren; JSON unmittelbar aus dem gespeicherten Zustand erzeugen.
`setPathDisplayAsDeleted` und die Rückübersetzung aus der Textmarkierung entfernen.

## Scope-and-Constraints

- Preserve: Textmarkierung, JSON-v1-Felder, Revisionsanreicherung und Filterlimits erhalten.
- Preserve: Gültige Sortiersemantik erhalten; die unabhängige Reihenfolgekorrektur bleibt in `ISSUE-021`.
- Exclude: Keine gemeinsame SDK-Metadatenarchitektur und keine Änderungen anderer Befehle.
- Cost: Interne Listing-Verbraucher geschlossen migrieren; keine öffentliche Go-Signatur ändern.

## Shared-Change-Pressure

Textdarstellung und JSON-Datenzustand hängen aktuell an derselben Zeichenkonvention.
Der Vorschlag entfernt diese Rückübersetzung statt nur Code zu verschieben.

## Verification

Quellprüfung abgeschlossen; kein Test oder Laufzeitvergleich ausgeführt.
Bestehende Verträge: `TestLsJSONDeletedEntryIsStructured` und `TestLsOnlyDeletedLimitCountsFilteredEntries`.
Zusätzlich Textdarstellung und Sortierverträge vergleichen; keine bestehende Fehlreihenfolge als Invariante festschreiben.
Test decision: none; nur Tracking, keine Source- oder Teständerung.

## Publication-Blockers

- Verhaltensgleichheit der Listing-Pipeline und enger Migrationsumfang bestätigen.
- Aktuellen kanonischen Upstream und Prior Art prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Veröffentlichungsfreigabe fehlen.

## Next-Action

Summary: Listing-Zustandsgrenze gezielt prüfen
Action: Löschzustand, Revisionsanreicherung und beide Renderer bis zu ihren internen Verbrauchern abgrenzen.
Done-When: Eine lokale Umstellung erhält Filter und Text-/JSON-Verträge ohne Rückübersetzung oder neue Sortierpolitik.
