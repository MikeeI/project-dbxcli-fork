# ISSUE-005 — CLI-Metadaten: getrennte Register und redundanter Auditstatus

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

Root-Cause [S]: Command-Metadaten verteilen sich auf zwei über Command-Schlüssel verbundene Register.
Abgeleitete Ergebnisfelder und `Known` repräsentieren bereits vorhandene Informationen zusätzlich.
Gesprächsreferenz: `S2`; der anschließende Review bestätigt die begrenzte Eigentümerschaftskorrektur.

## Reach-and-Impact

Reach [S]: JSON-Hilfe und der Dokumentationsgenerator konsumieren die daraus gebildeten Command-Manifeste.
Impact [S]: 28 Command-Schlüssel und 34 explizite `Known: true`-Zuweisungen müssen mitgeführt werden.
Ein falsches Manifest, Laufzeitdefekt oder Geschwindigkeitsgewinn wurde nicht nachgewiesen.
Die geschätzte Nettoeinsparung beträgt 35–45 Produktionszeilen, nicht eine gemessene Patchbilanz.

## Evidence

- [S] `cmd/help_manifest.go:18-28` besitzt drei abgeleitete Ergebnislisten und das Feld `Known`.
- [S] `cmd/help_manifest.go:83-352` definiert 34 Manifest-Einträge mit `Known: true`.
- [S] `cmd/help_manifest.go:354-383` wiederholt 28 Command-Schlüssel für Status-, Kind- und Warning-Metadaten.
- [S] `cmd/help_manifest.go:409-417` verbindet die Register und kopiert beziehungsweise sortiert die Vertragslisten.
- [S] `cmd/help_json.go:401-403` kopiert und sortiert dieselben Listen erneut an der Ausgabegrenze.
- [S] `cmd/help_manifest.go:431-449` behandelt Vertragsvorhandensein und Auditstatus getrennt.
- [S] LSP findet neben der Definition zwei Produktionszugriffe auf `commandContractRegistry` und einen Testzugriff.
- [S] `cmd/help_json_test.go:437,465` und `cmd/dry_run_contract_test.go:73` konsumieren die internen Register direkt.
- [S] `docs/json-schema/v1/README.md:159-190` beschreibt öffentliche Manifestfelder und deren Bedeutung.
- [O] Der lokale Upstream-Vergleich zeigt keine Differenz in `cmd/help_manifest.go` und `cmd/help_json.go`.

## Prior-Art

Coverage: Der lokale Index und alle vorhandenen Finding-Akten wurden ohne gleiche Metadatenursache geprüft.
Gaps: Upstream-Issues, PRs, Diskussionen und Releases wurden für diesen Vorschlag noch nicht recherchiert.
Contribution fit: Offen; ein kleiner Wartbarkeitsgewinn rechtfertigt keinen breiten Manifest-Umbau.
Der Source-Abgleich verwendet den vorhandenen Upstream-Ref, keinen erneuten Netzwerkabruf.

## Proposed-Change

Den vorhandenen Vertragstyp als optionales Feld `Contract *jsonCommandContractMetadata` im Manifest-Eintrag besitzen.
Alle Vertragswerte unverändert übernehmen und `commandContractRegistry` vollständig entfernen.
Die drei abgeleiteten Ergebnisfelder und den Register-Join entfernen.
Kopieren, Sortieren und Normalisieren an der bestehenden finalen Ausgabegrenze erhalten.
`Known` entfernen und Auditstatus ausschließlich aus der Manifest-Registermitgliedschaft ableiten.
Schema-Referenzen weiterhin nur bei vorhandenem Vertrag und erfüllter Structured-Output-/Help-Bedingung anbieten.
Bestehende interne Testzugriffe migrieren, ohne Assertions oder Abdeckung zu entfernen.
Die Änderung ist ein Vorschlag und nicht zur Implementierung autorisiert.

## Scope-and-Constraints

- Preserve: Alle Status-, Kind-, Warning-, Flag- und Schema-Metadaten unverändert erhalten.
- Preserve: Auditstatus und Ausführungsvertrag bleiben unabhängig; insbesondere bei `login` und `completion`.
- Preserve: Ein vorhandener optionaler Vertrag bleibt von einem fehlenden Vertrag unterscheidbar.
- Preserve: Unbekannte Metadatenpfade erhalten keine Auditklassifikation oder Vertragsreferenzen.
- Preserve: Finale defensive Kopien, deterministische Sortierung und `[]` statt `null` erhalten.
- Preserve: Exportierte Manifestfunktionen, JSON-Typen, Feldnamen und Schema-Validierungen bleiben bestehen.
- Exclude: Kein Aliasregister, neues Generatorframework oder geänderter CLI-Vertrag.
- Cost: Bestehende Assertions müssen dieselben Bedingungen über die neue interne Darstellung prüfen.
- Cost: Ohne versetzte Vertragsliterale geschätzt 45–50 Zeilen entfernen und 5–10 ergänzen.
- Cost: Tests, generierte Inhalte und bloßes Verschieben von Vertragswerten zählen nicht als Einsparung.

## Shared-Change-Pressure

Beide Register beschreiben dieselben Command-Identitäten und werden beim Manifestaufbau zusammengeführt.
Ein gemeinsamer Eintrag entfernt den zweiten Index und dessen Join, nicht die unabhängige Vertragsanwesenheit.
Ein erheblicher Teil der kleinen Nettoeinsparung entfällt auf den redundanten Auditstatus.

## Verification

Abgeschlossen [O]: Quelltext, symbolische Verbraucher und vorhandene Manifest-/Dry-run-Audits wurden geprüft.
Nicht ausgeführt: Implementierung, Tests, dynamischer Manifestvergleich und Messung einer tatsächlichen Löschbilanz.
Nach autorisierter Implementierung ist folgende Verifikation erforderlich:

- `TestJSONHelpManifestRegistryAudit` muss Registerabdeckung und Flag-Metadaten weiterhin prüfen.
- `TestJSONHelpIsDeterministic` muss die deterministische Ausgabe weiterhin bestätigen.
- `TestJSONHelpManifestsValidateAgainstPublicManifestSchema` muss die öffentlichen Manifeste weiterhin validieren.
- `TestDryRunCommandsDeclarePlannedStatus` muss Vertragsanwesenheit und `planned` weiterhin prüfen.
- Diese bestehenden Tests mit `go test ./cmd -run '^(TestJSONHelpManifestRegistryAudit|TestJSONHelpIsDeterministic|TestJSONHelpManifestsValidateAgainstPublicManifestSchema|TestDryRunCommandsDeclarePlannedStatus)$'` fokussiert ausführen.
- JSON-Manifeste vor und nach dem Umbau vergleichen, besonders auditierte Commands ohne Ausführungsvertrag.
- Die vorgeschriebenen Source-Submission-Gates bleiben vor einer späteren Einreichung erforderlich.

## Publication-Blockers

- Aktuellen kanonischen Upstream-Stand und Prior Art vor einer Beitragsempfehlung prüfen.
- Verhaltensgleichheit und tatsächliche Nettoeinsparung sind ohne Implementierung nicht dynamisch bestätigt.
- `Authorized-Work` und `Publication-Target` sind vom Nutzer noch nicht ausgewählt.
- Ein exakter Veröffentlichungsentwurf und dessen Freigabe liegen nicht vor.

## Next-Action

Summary: Registervereinfachung gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zur Eigentümerschaft von Command-Manifesten und Vertragsmetadaten recherchieren.
Done-When: Relevante Treffer, Dublettenbewertung und Nutzen des begrenzten Vorschlags sind dokumentiert.
