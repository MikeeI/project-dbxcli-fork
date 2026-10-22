# ISSUE-036 — Search: öffentliche Kategorie- und Erweiterungsfilter fehlen

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

[S] Der Ausgangsstand übergibt keine Kategorien oder Dateierweiterungen an die bereits vorhandenen SDK-Suchoptionen.
Die direkte Nutzeranweisung zur vollständigen privaten Umsetzung von `GOAL.md` autorisiert diesen lokalen Scope.
Die Upstream-Felder bleiben ohne ausgewählten Arbeitsmodus oder Veröffentlichungsauftrag unverändert.

## Reach-and-Impact

[S] Agenten und Menschen konnten den Suchraum nicht über diese nativen öffentlichen Serverfilter einschränken.
[A] Der Nutzen hängt vom jeweiligen Suchraum und der tatsächlichen Serverauswertung kombinierter Filter ab.
Eine Produktionshäufigkeit oder gemessene Beschleunigung wird nicht behauptet.

## Evidence

- [S] Das gepinnte SDK `v6.6.2` stellt `files.SearchOptions.FileCategories` und `FileExtensions` bereit.
- [S] `cmd/search.go:parseSearchOptions` validiert Kategorien und normalisiert Erweiterungen.
- [S] `cmd/search.go:newSearchV2Arg` überträgt beide Filter ohne neue Endpoints oder zusätzliche Einzelabfragen.
- [O] `TestSearchJSONOutputsInputAndResults` prüft den SDK-Request und das normalisierte Eingabe-Echo.
- [O] `go test ./...`, `go vet ./...` und `go build ./...` sind im privaten Implementierungsstand bestanden.
- [O] Die tatsächliche CLI besteht `scripts/test-skill-contract.sh` mit 14 Assertions ohne Kontoabfragen.

## Prior-Art

Der vollständige lokale Index und symptomatisch passende Akten wurden vor der Vergabe auf gleiche Ursachen geprüft.
`ISSUE-029` betrifft alte Konfliktstrategie-Completions, nicht diese fehlenden Suchfilter.
Externe Issues, PRs und aktuelle Upstream-Änderungen wurden für eine Veröffentlichung nicht recherchiert.

## Proposed-Change

Privat umgesetzt: native StringSlice-Flags `--category` und `--ext` mit SDK-Request, Eingabe-Echo und Discovery.
Kategorien stammen aus einer gemeinsamen Werteliste; der unterstützte Catch-all heißt `others`.
Erweiterungen erlauben weder leere Werte noch Glob- oder Pfadsyntax.

## Scope-and-Constraints

- Die beiden bestehenden Suchendpoints und das gepinnte SDK bleiben die Laufzeitgrenze.
- Der aktive SDK-Dateistatus und die Kombination mit `--content` bleiben erhalten.
- Die neue Kategorie-Completion nutzt native Cobra-APIs ohne eigenes Completion-Framework.
- Eine Live-Bestätigung der Filterkombination und zusätzliche Serverfunktionen sind nicht beauftragt.

## Verification

Test decision: update — `TestSearchJSONOutputsInputAndResults` schützt Request und Eingabe-Echo.
`TestJSONHelpManifestV1SelectedCommandMetadata` schützt Array- und Item-Enum-Verträge.
Die tatsächliche Texthilfe zeigt die neuen Flags; `__complete search --category doc` liefert native Kategorien.
Die Serverwirkung wurde nicht an einem Dropbox-Konto geprüft.

## Publication-Blockers

Upstream-Arbeitsmodus, Ziel, Prior Art, exakter Entwurf und dessen Freigabe fehlen.
Die private Implementierungsfreigabe ersetzt diese Voraussetzungen nicht.

## Next-Action

Summary: Privaten Abschluss revisionsfest erfassen
Action: Den geprüften privaten Search-Vertrag committen und auf `origin/personal` pushen.
Done-When: Commit und Push des vollständigen privaten Vertrags sind mit konkreter Revision belegt.
