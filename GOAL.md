# Ziel und Implementierungsplan: Agent-First Search für dbxcli

## Zweck und Freigabegrenze

Der Nutzer entwickelt diesen persönlichen Fork, nicht Upstream `dropbox/dbxcli`.
Ziel ist eine bessere Suche für LLM-Agenten; menschliche Nutzbarkeit bleibt erforderlich.
Agenten dürfen begrenzte oder fehlgeschlagene Suchläufe nicht als vollständige Ergebnislisten interpretieren.

Dieses Dokument enthält den vollständigen Plan und den jeweils belegten Implementierungsstand.
Die anschließende ausdrückliche Nutzeranweisung autorisiert die vollständige private Umsetzung dieses Scopes auf `personal`.
Die Implementierungsfreigabe allein autorisiert weder Upstream-Veröffentlichung noch globale Installation oder Kontozugriffe.
Die später gesondert freigegebenen lesenden Live-Checks sind im folgenden Nachweis dokumentiert.

Die Entwicklung erfolgt auf `personal`, ohne verpflichtende Upstream-Branches oder Pull-Request-Entwürfe.
Upstream-Integration ist kein Lieferziel, unabhängig von möglichem Interesse oder möglicher Ablehnung.
Externe Beiträge benötigen weiterhin die Freigabe des genauen aktuellen Entwurfs und Ziels.
Bestehende Strukturen und Konventionen bleiben erhalten, um spätere Upstream-Übernahmen nicht unnötig zu verteuern.

## Implementierungsstand und Prüfnachweis

Alle fünf Ziele sind im privaten Code umgesetzt und an den vorhandenen Command-/Client-Grenzen geprüft.
Die bestehende Struktur bleibt erhalten; `cmd/search-filters.go` ist die einzige neue Go-Quelldatei.
Die folgenden Abschnitte bleiben die vollständige Spezifikation dieser Umsetzung.

1. Suchzustand unterscheidet API-Ende, lokales Limit und Fehler und erhält akzeptierte Teilergebnisse.
2. `--category` und `--ext` werden als native Serverfilter übertragen.
3. Vollständige Matches erhalten Name, offenen Matchtyp und optionale JSON-only-Titelspans.
4. Nur lokal akzeptierte Datei-/Ordner-IDs werden dedupliziert und gegen das Ergebnislimit gezählt.
5. Namens-, Größen-, Zeit- und Typfilter wirken vor Deduplizierung und Auswahl mit präsenzkorrekter Discovery.

Der gemeinsame Reverse-Gleichstandsfix ist separat in `d36e7f6` auf `origin/personal` enthalten.
`ISSUE-036` bis `ISSUE-039` dokumentieren Serverfilter, Matchinformationen, Deduplizierung und lokale Filter.
`ISSUE-003`, `ISSUE-028`, `ISSUE-029`, `ISSUE-031` und `ISSUE-034` grenzen erledigte und verbleibende Scopes ab.
Der globale Allocator steht nach dieser Vergabe auf `ISSUE-040`; die Ledger-Validierung besteht.

Beobachtete Verifikation am 2026-10-22:

- Vor dem Fix scheitern der bestehende Reverse-Test und die Search-Fehlerprüfung wie erwartet.
- Nach der Umsetzung bestehen die gezielten Search-, Sortier-, Discovery-, Golden- und Live-Schemavertragsprüfungen.
- `go test ./...`, `go vet ./...` und `go build ./...` bestehen.
- Beide vorhandenen Generatoren wurden ausgeführt; nur die Search-Befehlsreferenz benötigt eine Command-Doc-Änderung.
- Kanonische Skills wurden synchronisiert und paketiert; die vorhandene Plugin-Prüfung besteht mit 29 Checks.
- Die tatsächlich gebaute CLI besteht die Skill-Vertragsprüfung mit 14 Assertions.
- Texthilfe, JSON-Discovery und native Kategorie-/Typ-Completion sind ohne Authentifizierung geprüft.
- Die vorhandenen Handlerprüfungen bestätigen stderr-Diagnose, Einzeldokumentfehler und erhaltene SDK-Exitklassen.

Es wurden keine neuen Testfunktionen, Testdateien, Abhängigkeiten, Packages oder Endpoints eingeführt.
Bis zum Implementierungsabschluss am 2026-10-22 wurden keine Live-Kontoabfragen ausgeführt.
Die gesondert freigegebenen Live-Checks vom 2026-10-23 ergänzen diesen ursprünglichen Prüfumfang.
Globale Installation, Transfers und Upstream-Veröffentlichung wurden nicht ausgeführt.
Die vollständige private Search-Implementierung ist in `547752f` auf `origin/personal` enthalten.
`ISSUE-003` und `ISSUE-036` bis `ISSUE-039` sind als privater Abschluss mit `Archive-Reason: Other` archiviert.
Das behauptet keine Upstream-Behebung; offene Listing-, Transfer- und Alt-Completion-Scopes bleiben getrennt erhalten.

## Live-Nachweis am 2026-10-23

Der Nutzer autorisierte eine Kontoidentitätsprüfung und einige lesende CLI-Aufrufe.
Ausgeführt wurden eine Kontoabfrage, zwei begrenzte Listings und fünf Suchläufe, jeweils mit `--timeout 30s`.
Alle acht Live-Aufrufe endeten erfolgreich mit `ok: true`.
Remote-Dateien und Quellcode wurden dabei nicht verändert; Dateiinhalte wurden nicht heruntergeladen.

Verwendet wurde ein temporärer nativer Build aus `48f3420e3074b8d18915e524fcc3bb023950132a`.
Die gemeldete CLI-Version war `3.7.5-0.20261022142919-48f3420e3074`, mit SDK `v6.6.2` und Spec `6d985df`.
Die folgenden Aufrufe verwenden den normalisierten Executable-Namen `dbxcli`, keine globale Installation.
Der temporäre Build wurde nach der Demonstration entfernt; der Repository-Arbeitsbaum blieb sauber.

### Konto und begrenzte Übersicht

```sh
dbxcli --timeout 30s account --output=json
dbxcli --timeout 30s ls / --limit 10 --sort name --output=json
dbxcli --timeout 30s ls /SORT_Scan --limit 6 --output=json
```

- Die Kontoabfrage bestätigt `Gerold M` mit `dropbox@in9.at`.
- Die Auth-Metadaten lauten `source: saved`, `refreshable: true` und `auth_file: default`.
- Das Root-Listing enthält zehn Ordner, darunter `/GitlabWorkspace`, `/Sharing` und `/SORT_Scan`.
- Das zweite Listing zeigt sechs Einträge aus `/SORT_Scan` als Grundlage der Suchbegriffe.
- Diese begrenzten Listings sind keine vollständige Konto- oder Ordnerinventur.
- Die erfolgreiche Kontoabfrage beweist nicht, dass im beobachteten Lauf tatsächlich ein Token erneuert wurde.

### Namenssuche und sichtbares Ergebnislimit

```sh
dbxcli --timeout 30s search 2018 /SORT_Scan --limit 3 --output=json
```

Die Ausgabe enthält drei Treffer bei drei empfangenen API-Matches.
Der Zustand ist `stop_reason: limit`, `server_has_more: true` und `local_filters_applied: false`.
Der Ausschnitt wird damit nicht als vollständiges API-Ergebnis dargestellt.

### PDF-Serverfilter, Titelspans und lokale Sortierung

```sh
dbxcli --timeout 30s search 2018 /SORT_Scan \
  --ext pdf --category pdf --highlights --limit 3 \
  --sort size --reverse --output=json
```

Die kombinierte Serverfilteranfrage liefert drei PDFs bei drei empfangenen API-Matches.
Ihre Größen sind lokal absteigend als `1805491`, `1551495` und `1517175` Byte sortiert.
Die Titelspans markieren jeweils `2018` durch `is_highlighted: true`.
Der Zustand bleibt `limit` mit `server_has_more: true`; die Sortierung beweist kein global größtes Top-N-Ergebnis.

### Kombinierte lokale Filter und API-Ende

```sh
dbxcli --timeout 30s search 2018 /SORT_Scan \
  --ext jpg --category image --type file \
  --name-glob '*.jpg' --name-regex '^2018' --exclude-name-contains _002 \
  --min-size 1 --max-size 1000000 \
  --modified-from 2018-01-01T00:00:00Z --modified-before 2019-01-01T00:00:00Z \
  --time client --limit 3 --sort time --reverse --output=json
```

Vier API-Matches ergeben einen akzeptierten Treffer: `/SORT_Scan/20180201173703_001.jpg`.
Die Datei ist `637479` Byte groß und besitzt `client_modified: 2018-02-01T16:45:53Z`.
Ihre Serveränderungszeit lautet dagegen `2025-02-18T03:40:10Z`.
Der Zustand ist `api_exhausted`, `server_has_more: false` und `local_filters_applied: true`.
Dieser Lauf belegt Filterkombination und API-Ende, aber keine gezielt herbeigeführte Mehrseiten-Pagination.

### Inhaltssuche ohne Dateidownload

```sh
dbxcli --timeout 30s search Rechnung /SORT_Scan \
  --content --ext pdf --category pdf --highlights --limit 3 --output=json
```

Die Inhaltssuche liefert drei PDFs, obwohl `Rechnung` in keinem ihrer Dateinamen steht.
Ein Treffer ist `/SORT_Scan/2018-06-20 10.10.37.pdf`.
Der Zustand ist `limit` mit `server_has_more: true`; Titelspans sind vorhanden, aber nicht hervorgehoben.
Es werden keine Inhaltssnippets oder Dateibytes ausgegeben.
Dropbox liefert auch bei der vorherigen Namenssuche `match_type: filename_and_content`.
Dieser unverändert übernommene Providerwert ist deshalb kein Beweis für die gewählte `--content`-Requestoption.

### Explizite Null-Byte-Grenze

```sh
dbxcli --timeout 30s search '2018-06-20 10.05.54.pdf' /SORT_Scan \
  --ext pdf --max-size 0 --limit 3 --output=json
```

Die zuvor gefundene PDF ist `1805491` Byte groß.
Ein API-Match wird empfangen, aber kein Ergebnis ausgegeben.
`input.max_size` bleibt ausdrücklich `0`; der Zustand ist `api_exhausted` mit lokalen Filtern.
Die Null wird damit nicht als fehlende Obergrenze behandelt.

### Aussagegrenzen des Live-Nachweises

- Bestätigt sind die genannten Querys und Filterkombinationen auf diesem Konto, keine beliebigen Suchräume.
- Die Inhaltssuche bestätigt keine garantierte OCR-Abdeckung sämtlicher Scans.
- Deduplizierung unter tatsächlich gelieferten Duplikaten wurde nicht gezielt provoziert.
- Fortsetzungsfehler und erhaltene Teilergebnisse wurden nicht am Live-Konto herbeigeführt.
- Eine ausgelöste Token-Erneuerung wurde nicht separat beobachtet.
- API-Ende bleibt wegen Indexverzug, möglichen Auslassungen und Obergrenze kein Vollständigkeitsbeweis.

## Grenzen und Nichtziele

- Verwende nur `files/search_v2` und `files/search/continue_v2`.
- Behalte das gepinnte SDK `github.com/dropbox/dropbox-sdk-go-unofficial/v6@v6.6.2` bei.
- Erhalte `main.go`, `cmd/`, `internal/output/`, `internal/jsonschema/` und die vorhandene Cobra-Integration.
- Führe keine neuen Packages, Abhängigkeiten, Endpoints oder zusätzlichen Metadatenabfragen pro Treffer ein.
- Baue keine Filter-Sprache, Retry-Engine, Worker-Pipeline oder öffentliche Cursor-Schnittstelle.
- Erhalte JSON v1 und bestehende Status-, Fehlercode- und Exitcode-Bedeutungen.
- Behandle neue Felder als Erweiterung, nicht als Umbenennung vorhandener Laufzeitfelder.
- Erhalte das bisherige Textformat der Treffer; zusätzliche Diagnosen gehören auf stderr.
- Ändere weder Login-Lifecycle noch globale Installation, Transfers oder Kontodaten.
- Ersetze serverseitige Suche nicht durch Mounts, Synchronisation oder rekursive Listings.
- Versprich keine serverseitigen Datumsbereiche, umgekehrte Serverreihenfolge oder Inhaltssnippets.
- Versprich keine Regex- oder Boolean-Suche über Dokumentinhalte.
- Füge weder Property Search noch Android-interne Suchendpoints hinzu.

Die beschriebenen neuen Flags und JSON-Felder werden gemeinsam mit ihren Schemas und Verbraucherhinweisen implementiert.
Lokale Filter beziehen sich ausschließlich auf die Treffer der angegebenen Dropbox-Query.
Sie ersetzen keine vollständige kontoweite Inventarisierung nach Name, Größe oder Datum.

## Belegter Ausgangsstand

Die folgenden Quellanker dokumentieren den historischen Bestand vor dieser Umsetzung.
Die damalige Untersuchung nutzte keine Live-Kontoabfragen; spätere Beobachtungen stehen im Live-Nachweis.

- `go.mod:6` pinnt das benötigte SDK bereits auf `v6.6.2`.
- `cmd/search.go:79-93` sammelt Metadaten, begrenzt sie und sortiert erst danach lokal.
- `cmd/search.go:82-86` verlässt den Lauf bei einem Fortsetzungsfehler vor der Ergebnisausgabe.
- `cmd/search.go:149-159` verwirft beim Entpacken Matchtyp und Highlight-Spans.
- `cmd/files_client.go:24-25` bietet die beiden benötigten kontextfähigen Suchmethoden bereits an.
- `cmd/json_metadata.go:22-73` besitzt den gemeinsamen Metadatentranslator für die weitere Verwendung.
- `cmd/output.go:65-116` kann Fehlerdetails ergänzen und die ursprüngliche Fehleridentität erhalten.
- `cmd/root.go:327-331` übernimmt die endgültige Fehlerausgabe und den Exitcode.
- `cmd/help_input_schema.go:18-38` besitzt ein Array-Muster für variadische Argumente.
- `cmd/help_input_schema.go:49-81,129-184` behandelt Listenflags und präsenzabhängige Defaults noch nicht passend.
- `internal/jsonschema/commands.go:133-163` erzeugt geschlossene Erfolgs- und Ergebnisobjekte.
- `tools/gen-json-schemas/main.go:24-64` liest `commands.json` und erzeugt nur `commands.schema.json`.

Die SDK-Quellanker beziehen sich auf `dropbox/files/` innerhalb der gepinnten Modulversion.
`types.go:4043-4143,4179-4190` enthält Filter, Matchdaten und die moderne Highlight-Option bereits.
`types.go:1357-1375,1470-1474` definiert Datei- und Ordner-IDs sowie dateispezifische Größe und Zeiten.
`client.go:758-770` dokumentiert Begrenzung, Indexverzug, Duplikate und möglicherweise ausgelassene Treffer.

Die aktuelle API-Dokumentation bestätigt diese Grenzen und die Einschränkung der Serverfilter auf aktive Suche.
Quelle: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/search-v-2
Die Fortsetzung besitzt keinen eigenen Total-Count oder ein verlässliches `cap_reached`-Feld.
Quelle: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/search-continue-v-2

## Die fünf Ziele in Prioritätsreihenfolge

1. Mach Ergebnisbegrenzung und Fehlerabbrüche sichtbar und erhalte bereits akzeptierte Teilergebnisse.
2. Ergänze Kategorie und Dateierweiterung als echte serverseitige Filter.
3. Erhalte Matchtyp und optional angeforderte Titel-Highlights in der Search-Ausgabe.
4. Entferne seitenübergreifende Duplikate anhand der Datei- oder Ordner-ID.
5. Filtere empfangene Metadaten lokal, bevor unterschiedliche passende Treffer gegen das Limit zählen.

## Datei- und Ordnerstruktur

`M` kennzeichnet vorhandene zu ändernde Dateien; `N` kennzeichnet die einzige neue Quelldatei.
`G` kennzeichnet generierte Ausgaben; `S` kennzeichnet synchronisierte Skill-Kopien.

```text
cmd/
├── search.go                                      M
├── search-filters.go                              N
├── format.go                                      M
├── help_manifest.go                               M
├── help_json.go                                   M
├── help_input_schema.go                           M
├── search_test.go                                 M
├── format_test.go                                 M
├── help_json_test.go                              M
├── json_contract_test.go                          M
├── json_command_schema_test.go                    M
└── testdata/json_contract/
    ├── success_schemas.json                       M
    ├── success_outputs.json                       M
    └── error_outputs.json                         M

internal/
└── jsonschema/
    ├── commands.go                                M
    └── definition_schemas.go                      M

docs/
├── automation.md                                  M
├── commands/
│   └── dbxcli_search.md                            G
└── json-schema/v1/
    ├── README.md                                  M
    ├── commands.json                              M
    ├── commands.schema.json                       G
    ├── success.schema.json                        M
    └── error.schema.json                          M

skills/
└── dbxcli/
    ├── SKILL.md                                   M
    └── references/
        └── tool-integration.md                    M

plugin/
└── claude/skills/dbxcli/
    ├── SKILL.md                                   S
    └── references/
        └── tool-integration.md                    S

scripts/
└── test-skill-contract.sh                         M
```

`main.go`, `cmd/root.go`, `cmd/files_client.go` und `internal/output/` bleiben unverändert.
Auch `cmd/json_metadata.go` und `cmd/json_output.go` benötigen keine Erweiterung ihrer gemeinsamen Datentypen.
Die vorhandenen Generatoren und das Sync-Skript werden verwendet, nicht durch neue Werkzeuge ersetzt.
Die zahlreichen Vertragsdateien sind bestehende Darstellungen derselben CLI, keine neuen Architekturschichten.

## Gemeinsamer Such-Lifecycle

### Datenhaltung und Zuständigkeiten

`cmd/search.go` bleibt Owner für Command, Request, Pagination, Deduplizierung und Ausgabe.
`cmd/search-filters.go` besitzt ausschließlich lokale Filterkonfiguration, Validierung und Metadatenprüfung.
Die Filterdatei erhält weder Netzwerkzugriff noch Pagination- oder Ausgabezuständigkeit.

Behalte `[]*files.SearchMatchV2` bis zur Ausgabe, statt früh auf `[]files.IsMetadata` zu reduzieren.
Ein zusätzlicher `searchHit`-Wrapper wäre ohne weitere Verantwortung überflüssig.
Die bestehende `filesClient`-Grenze genügt; ein zweites Search-Client-Interface ist nicht erforderlich.

Geplante interne Typen in `cmd/search.go`:

```diff
+type searchResults struct {
+    matches        []*files.SearchMatchV2
+    serverHasMore  *bool
+    stopReason     string
+    apiMatchesSeen uint64
+}
+
+type searchSummary struct {
+    Limit               uint64 `json:"limit"`
+    StopReason          string `json:"stop_reason"`
+    ServerHasMore       *bool  `json:"server_has_more"`
+    LocalFiltersApplied bool   `json:"local_filters_applied"`
+    APIMatchesSeen      uint64 `json:"api_matches_seen"`
+    APIMatchLimit       uint64 `json:"api_match_limit"`
+}
+
+type jsonSearchMetadata struct {
+    jsonMetadata
+    Name           string                 `json:"name,omitempty"`
+    MatchType      string                 `json:"match_type,omitempty"`
+    HighlightSpans []*files.HighlightSpan `json:"highlight_spans,omitempty"`
+}
+
+type jsonSearchOutput struct {
+    jsonOperationOutput
+    Search searchSummary `json:"search"`
+}
```

`searchResults` enthält keinen zusätzlichen Fehler; der Fehler wird separat als `error` zurückgegeben.
`collectSearchMatches()` erhält Kontext, bestehenden Client, Request, lokale Filter und Ergebnislimit explizit.
Sein Vertrag erklärt, dass die zurückgegebenen Ergebnisse auch bei einem Fehler gültige Teilergebnisse enthalten.
Kontext und veränderlicher Suchzustand bleiben auf den einzelnen Command-Aufruf begrenzt.

Definiere Stopgründe, Kategorien, Schema-Detailkeys und die API-Obergrenze beim jeweiligen Owner als benannte Werte.
Verwende den SDK-Konstruktor als Owner der Standard-Seitengröße, statt den Default im CLI erneut zu definieren.
Neue Flag-Defaults werden einmal bei ihrer Registrierung festgelegt.

### Verarbeitungsreihenfolge

1. Lies und validiere die Optionen; kompiliere Regex und parse Zeitgrenzen einmal.
2. Baue den Request einschließlich Serverfilter und optionaler Highlight-Anforderung.
3. Rufe die nächste Seite seriell mit dem bestehenden Command-Kontext ab.
4. Zähle empfangene API-Matches einschließlich Duplikaten und lokal abgelehnten Matches.
5. Prüfe jeden Match auf verarbeitbare Metadaten und notwendige Identität.
6. Wende lokale Filter auf den originalen SDK-Metadatensatz an.
7. Überspringe bereits akzeptierte IDs.
8. Übernimm neue passende Matches, solange das lokale Ergebnislimit nicht erreicht ist.
9. Prüfe auch den Rest der bereits empfangenen Seite auf weitere passende unterschiedliche Treffer.
10. Entscheide anhand dieses Zustands und `HasMore`, ob eine weitere Seite benötigt wird.
11. Sortiere die akzeptierten vollständigen Search-Matches lokal.
12. Erzeuge genau eine zur Ausgabeart passende Ergebnis- oder Fehlerdarstellung.

Eine leere oder vollständig herausgefilterte Seite beendet den Lauf bei `HasMore=true` nicht.
Es gibt keine zusätzliche Lookahead-Abfrage nur zur Bestätigung einer Kürzung.
Die API begrenzt die abrufbare Suchmenge; die CLI erfindet keinen separaten Scan-Cutoff oder Gesamtzähler.

Filterung vor der Deduplizierung verhindert, dass ein früher nicht passender Datensatz eine spätere ID-Darstellung sperrt.
Die ID-Map erfasst die akzeptierten Treffer, nicht sämtliche lokal abgelehnten Datensätze.

## Ziel 1: Abschlusszustand und Teilergebnisse

### Suchzusammenfassung

Erfolgreiche neue Search-Antworten erhalten zusätzlich das Objekt `search`.
Vorhandene Felder, Ergebnis-Wrapper, Statuswerte und Metadatentypen werden weder entfernt noch umgedeutet.
`limit` enthält auch den bisherigen Wert `0` für ein unbegrenztes lokales Ergebnislimit.
`local_filters_applied` beschreibt nur lokale Filter, nicht die Serverfilter.
`api_matches_seen` zählt empfangene Rohmatches, nicht nur ausgegebene Dateien.
`api_match_limit` nennt die dokumentierte Obergrenze von 10.000 API-Matches.

Die Stopgründe haben folgende Bedeutung:

- `api_exhausted` bedeutet, dass alle empfangenen Matches verarbeitet sind und die API keine weitere Seite meldet.
- `limit` bedeutet, dass passende zusätzliche Matches verworfen oder weitere Seiten wegen des Limits nicht abgerufen wurden.
- `error` bedeutet, dass Abruf oder Verarbeitung durch einen Fehler unterbrochen wurden.

`server_has_more` bezeichnet ausschließlich den Wert der letzten erfolgreichen API-Antwort.
Vor der ersten erfolgreichen Antwort ist dieser Wert `null`, nicht ein erfundenes `false`.
Der Name präzisiert das frühere Planfeld `has_more`, das nicht als öffentlicher Laufzeitvertrag eingeführt wurde.
`stop_reason` bleibt das maßgebliche CLI-Abschlusssignal.

Verwende weder `complete` noch ein selbst erfundenes `cap_reached` als verlässliche Serveraussage.
`api_exhausted` garantiert keine vollständige Liste aller passenden Dateien im Konto.
Auch bei 10.000 empfangenen Matches bleibt unklar, ob natürliches Ende und API-Obergrenze zusammenfallen.

### Entscheidende Grenzfälle

- Genau zehn passende Treffer bei Limit zehn und ohne Rest oder Folgeseite ergeben `api_exhausted`.
- Zwanzig passende Treffer auf der letzten Seite bei Limit zehn ergeben `limit` und `server_has_more: false`.
- Weitere API-Seiten bei erreichtem Limit ergeben `limit`, ohne weitere lokal passende Treffer zu versprechen.
- Seiten ohne lokal akzeptierte Treffer werden weiter abgefragt, solange die API weitere Seiten meldet.

### Fehlerausgabe

Bei Fehlern enthalten `error.details.search` den Suchzustand und `error.details.partial_results` die Teilergebnisse.
Teilergebnisse verwenden denselben Ergebnis-Wrapper und dieselbe Search-Metadatenprojektion wie Erfolgsantworten.
Ein initialer Suchfehler kann eine leere Teilergebnisliste und `server_has_more: null` enthalten.
Erhalte `ok: false`, ursprüngliche Fehlerursache, Fehlercode und die bestehende Exitcode-Zuordnung.
Verwende dafür `withJSONErrorDetails()` statt eines neuen Fehlerframeworks.

Im JSON-Modus schreibt `search()` vor der Fehlerrückgabe kein Erfolgsdokument.
Der vorhandene Root-Renderer erzeugt anschließend genau eine Fehlerhülle.
Erfolgs-JSON gefolgt von einer zweiten Fehlerantwort ist ausdrücklich ausgeschlossen.

Im Textmodus bleiben akzeptierte Treffer auf stdout verfügbar.
Suchzustand, Filterhinweis und Abbruchdiagnose gehen auf stderr; der Fehler-Exitcode bleibt erhalten.
Schreibfehler werden weitergegeben und dürfen weder Erfolg vortäuschen noch die ursprüngliche API-Ursache verbergen.

## Ziel 2: Serverseitige Kategorien und Erweiterungen

Vorgesehene Aufrufe:

```sh
dbxcli search report --category pdf
dbxcli search report --category document,spreadsheet
dbxcli search report --ext pdf,docx
```

Registriere `--category` und `--ext` als native StringSlice-Flags.
Übertrage ihre Werte direkt auf `SearchOptions.FileCategories` und `SearchOptions.FileExtensions`.
`files.NewSearchOptions()` setzt bereits den für diese Filter benötigten aktiven Dateistatus.
Füge keinen neuen Status-Flag und keine lokale Ersatzimplementierung dieser Serverfilter hinzu.

Die auswählbaren Kategorien sind:

```text
image, document, pdf, spreadsheet, presentation,
audio, video, folder, paper, others
```

`other` ist ein SDK-Union-Fallback und wird nicht als auswählbare Kategorie angeboten.
Validierung, Manifest und die Completion der neuen Kategorie- und Typwerte verwenden dieselben Wertelisten.

Für Erweiterungen gilt:

- Entferne äußere Leerzeichen und einen optionalen führenden Punkt.
- Normalisiere Großbuchstaben zu Kleinbuchstaben.
- Lehne leere Einträge und Glob-Ausdrücke wie `*.pdf` ab.
- Erfinde keine zusätzlichen Aliasse oder garantierte Behandlung mehrteiliger Erweiterungen.

Die genaue kombinierte Semantik der beiden Filterlisten wird nicht stärker zugesagt als im API-Vertrag.
Die dokumentierte Active-Einschränkung verbietet ihre Kombination mit `--content` nicht.

## Ziel 3: Matchtyp und Titel-Highlights

`match_type` wird übernommen, wenn die API ihn liefert; fehlende Werte werden nicht geraten.
Bewahre ihn als String und begrenze ihn nicht durch ein geschlossenes CLI-Enum.
Bekannte SDK-Tags umfassen auch `image_content` und `metadata`, nicht nur Name, Inhalt und beides.

`--highlights` fordert Titelspans über `SearchV2Arg.MatchFieldOptions.IncludeHighlights` an.
Das veraltete `SearchV2Arg.IncludeHighlights` bleibt unbenutzt.
Die Spans behalten ihre native Form mit `highlight_str` und `is_highlighted`.
Es sind Titelbestandteile, keine Zeichenoffsets und keine Inhaltssnippets.
Quelle: https://dropbox.tech/developers/search-files-using-the-dropbox-api

`jsonSearchMetadata` bettet den vorhandenen `jsonMetadata`-Wert ein.
`Name` stammt aus den Originalmetadaten und bleibt auch ohne `path_display` als Filterziel nachvollziehbar.
Der gemeinsame Metadatentyp anderer Commands erhält keine Search-Felder.

Highlights sind zunächst ausdrücklich JSON-only.
`--highlights` zusammen mit Textausgabe erzeugt einen Validierungsfehler mit Hinweis auf `--output=json`.
Die Prüfung erfolgt in `parseSearchOptions()` vor dem Suchrequest, nicht durch einen Umbau des Auth-Lifecycles.
Ohne diesen Flag bleibt das Textformat unverändert; eine neue Farb- oder Spaltenansicht gehört nicht zu diesem Scope.

## Ziel 4: Deduplizierung

Dateien und Ordner werden anhand ihrer SDK-ID dedupliziert, nicht anhand ihres Pfades.
Der erste lokal passende vollständige Search-Match pro ID gewinnt.
Duplikate zählen nicht gegen das Ergebnislimit.
Die ID-Map lebt nur während eines Suchaufrufs und wird nicht persistent oder global gespeichert.

Führe widersprüchliche Revisionen, Titelspans und Matchtypen nicht zu einem erfundenen Match zusammen.
Die Matchinformationen beschreiben den behaltenen Datensatz, nicht eine vollständige Historie aller Duplikate.
Fehlende notwendige IDs oder nicht verarbeitbare Metadaten werden als Fehler mit Teilergebnissen gemeldet.
Verwende weder einen erfundenen Pfad-Fallback noch stilles Überspringen als angeblich vollständigen Erfolg.
Vorhandene v1-Ergebniskinds werden wegen dieser internen Validierung nicht aus den Schemas entfernt.

## Ziel 5: Lokale Metadatenfilter

### Namen und Typ

- `--name-glob PATTERN` verwendet `path.Match()` auf `Metadata.Name`.
- `--name-regex REGEX` verwendet einmalig kompiliertes Go-`regexp` auf demselben Namen.
- `--name-contains TEXT` prüft eine wörtliche case-sensitive Teilzeichenfolge.
- `--exclude-name-contains TEXT` ist wiederholbar und verwirft Matches bei einem passenden Ausschluss.
- `--type file|folder` filtert den Metadatentyp; ohne Flag sind Dateien und Ordner zugelassen.

Der Glob bezieht sich auf den vollständigen Namen, nicht auf den Pfad.
Es gibt keine rekursive `**`-Pfadsuche.
Regex verwendet die normale Go-Suchsemantik; vollständiges Matching benötigt `^` und `$`.
Case-insensitive Regex kann ausdrücklich mit `(?i)` angefordert werden.
Ausschlüsse verwenden StringArray statt StringSlice, damit Kommas innerhalb eines Wertes erhalten bleiben.
Leere Ausschlusswerte werden abgelehnt, statt unbemerkt sämtliche Namen auszuschließen.

### Größe und Änderungszeit

- `--min-size BYTES` setzt eine inklusive Untergrenze als native `uint64`-Bytezahl.
- `--max-size BYTES` setzt eine inklusive Obergrenze; ausdrücklich gesetztes `0` bedeutet null Byte.
- `--modified-from RFC3339` setzt eine inklusive Zeituntergrenze.
- `--modified-before RFC3339` setzt eine exklusive Zeitobergrenze.

Die Zeitgrenzen bilden ein halb offenes Intervall mit explizitem RFC3339-Zeitzonenbezug.
Das bestehende `--time server|client` wählt auch das für die lokalen Zeitfilter verwendete Metadatenfeld.
Parse Zeiten einmal und bewahre ihre Präzision auch im Eingabe-Echo.
Lehne eine minimale Größe oberhalb der maximalen Größe und leere oder umgekehrte Zeitintervalle ab.

Größen- und Zeitfilter gelten nur für Dateien.
Ordner werden nicht als null Byte oder mit einem erfundenen Änderungsdatum behandelt.
Lehne `--type folder` zusammen mit Größen- oder Zeitfiltern als widersprüchlich ab.
Fehlende notwendige Dateidaten werden nicht durch erfundene Werte ersetzt.

### Datei-, Server- und Aufnahmezeit

`client_modified` ist der vom Client beziehungsweise Uploader gesetzte Dateiänderungszeitpunkt.
Dropbox übernimmt diesen Wert ungeprüft; ein API-Uploader kann ihn ausdrücklich angeben.
Er ist weder ein verlässlicher Erstellungszeitpunkt noch ein Synchronisationsnachweis.
Für Änderungs- und Konflikterkennung dient stattdessen die Revision `rev`.
Quellen: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/get-metadata und https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/upload

`server_modified` beschreibt die letzte Änderung der Datei auf Dropbox.
Der Wert ist nicht automatisch das Datum des ersten Uploads und kein Fotoaufnahme- oder Dokumenterstellungsdatum.
Quelle: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/get-metadata

Eine mögliche Foto-/Videoaufnahmezeit ist das eigene optionale Feld `media_info.metadata.time_taken`.
Dropbox extrahiert es aus eingebetteten Medieninformationen, beispielsweise EXIF `DateTimeOriginal`.
Es ist kein universeller Dateisystem-Erstellungszeitpunkt und kann fehlen oder noch nicht verfügbar sein.
Quellen: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/get-metadata und https://community.dropbox.com/en/discussion/253446/fully-understanding-filesphotometadata

Der öffentliche Metadatenweg verwendet `get_metadata` mit `include_media_info: true`.
`list_folder` und seine Fortsetzung liefern das Medienfeld seit dem 2019-12-02 nicht mehr.
Das alte SDK-Listenflag `IncludeMediaInfo` ist entsprechend wirkungslos und keine Lösung für Aufnahmezeitfilter.
Quelle: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/get-metadata
SDK-Belege: `files/types.go:122-124,1376-1379,1896-1899,2521-2528` in der gepinnten Version.

Die aktuelle CLI gibt in `cmd/json_metadata.go` nur Client- und Serveränderungszeit aus.
`cmd/format.go:getTime` und die lokalen Filter wählen ausschließlich zwischen `--time server` und `--time client`.
Es gibt keinen `--time taken`-Modus und keinen Filter nach Medienaufnahmezeit.
Diese Untersuchung erweitert weder die CLI noch den vereinbarten Suchscope um zusätzliche Metadatenabfragen.

Ein konkret beobachtetes Gegenbeispiel ist `/SORT_Scan/20180119172250_001.jpg`.
Das datumsähnliche Namenspräfix lautet `20180119172250`, aber `client_modified` ist `2018-02-05T10:32:20Z`.
`server_modified` derselben Datei ist `2025-02-18T03:39:47Z`.
Die Werte belegen unterschiedliche Zeitangaben, aber weder einen konkreten Importvorgang noch eine echte Aufnahmezeit.
Die eingebetteten Medien-/EXIF-Daten dieser Datei wurden nicht abgerufen.

### Kombination, Limit und Aufwand

Unterschiedliche positive Filter werden mit AND kombiniert.
Ausschlüsse werden mit OR kombiniert und anschließend negiert.
Nur passende unterschiedliche IDs zählen gegen das Ergebnislimit.
Lokale Filter speichern, laden oder durchsuchen keine Dateiinhalte.
Sie reduzieren keine Serverarbeit und können trotz kleinem Ergebnislimit viele Seiten erfordern.

Bei aktiven lokalen Filtern bleibt die SDK-Standard-Seitengröße erhalten.
Ein lokales Ergebnislimit von eins soll nicht automatisch Suchseiten mit nur einem Rohmatch erzwingen.
Ohne lokale Filter bleibt die bisherige Verkleinerung kleiner Requests erhalten.
Es wird kein zusätzlicher `--page-size`-Flag eingeführt.
Der vorhandene Kontext einschließlich `--timeout` gilt weiterhin für alle seriellen Seitenabrufe.

## Änderungen je vorhandenem Code-Owner

### `cmd/search.go`

- Erweitere `searchCommandOptions` und `searchInput` um die beschriebenen Optionen.
- Halte die Flagregistrierung für Produktionscommand und bestehendes Testkommando in einer gemeinsamen Funktion.
- Erweitere `newSearchV2Arg()` um Serverfilter, Highlights und die Seitengrößenregel.
- Ersetze die verlustbehaftete Metadatensammlung durch `collectSearchMatches()`.
- Entferne danach obsolete `appendSearchMatches()`- und `searchLimitReached()`-Pfade ohne Kompatibilitätsshims.
- Sortiere vollständige Search-Matches über den gemeinsamen Metadatenkomparator.
- Verwende für die JSON-Projektion weiterhin `jsonMetadataFromDropbox()` und `newJSONOperationResult()`.
- Erweitere `renderSearchOutput()` und `newSearchInput()` um die neuen Verträge.
- Verwende für unveränderte Textzeilen weiterhin die vorhandenen Metadatenformatierer.

### `cmd/format.go`

Im Ausgangsstand negiert der Reverse-Komparator `less` und verletzt bei Gleichständen die Sortierordnung.
Die begrenzte Korrektur gehört zum bereits vorhandenen `ISSUE-003`, nicht zu einer neuen Sortierarchitektur.
Belege: `cmd/format.go:compareLess` und `issues/archive/ISSUE-003.md`.

Ordne Reverse zentral in `compareLess()` durch vertauschte Operanden.
Entferne die zusätzliche Negation aus `sortEntries()`.
Die neue Search-Match-Sortierung verwendet denselben Komparator.
Behalte stabile Sortierung ohne neue sekundäre Namensschlüssel bei.
Die Korrektur betrifft auch vorhandene Listing-Aufrufer und benötigt den bestehenden Reverse-Regressionsnachweis.

### `cmd/help_manifest.go`, `cmd/help_json.go` und `cmd/help_input_schema.go`

Ergänze die neuen Flag-Metadaten, Enums, Wertarten und passende Beispiele.
Erkenne StringSlice und StringArray als Arrays mit String-Items.
Setze Kategorie-Enums auf `items.enum`, nicht auf das gesamte Array.
Verwende das vorhandene Array-Muster variadischer Argumente und die bereits vorhandene `Items`-Struktur.

Präsenzabhängige Filter erhalten im `input_schema` keinen Default.
Insbesondere darf ein technischer Cobra-Default `0` nicht als semantischer Default von `max_size` erscheinen.
Sonst könnte ein Agent durch Übernahme von Defaults unbeabsichtigt nur leere Dateien suchen.
Unterscheide in der CLI fehlende Werte mit `Flags().Changed()` und im Eingabe-Echo mit optionalen Pointer-Feldern.
Ein ausdrücklich gesetztes `0` darf nicht durch `omitempty` verschwinden.
Erfinde dafür weder `-1` noch eine künstliche maximale Dateigröße.

Auch ein String `"[]"` ist kein gültiger Array-Default im Eingabeschema.
Die neuen optionalen Listenfilter benötigen keinen publizierten Schema-Default.
Andere bestehende Flags behalten ihre bisherige Default-Bedeutung.
Die interne Default-Präsenzregel wird durch `jsonCommandFlag` transportiert, ohne das öffentliche Manifest zu erweitern.

### Schema-Katalog und Generator

Ergänze `search_summary`, Search-Metadaten und Highlight-Spans in den vorhandenen Definitions-Ownern.
Search verweist künftig auf die search-spezifische Metadatendefinition statt auf eine verbreiterte gemeinsame Definition.
Erweitere `CommandContract` gezielt um eine optionale Referenz `SearchSummary` mit JSON-Key `search_summary`.
Der Generator erlaubt damit das zusätzliche Feld `search` nur in den passenden Command-Verträgen.
Führe kein generisches `ExtraFields`-System ein.

`commands.json` und die manuell gepflegten Basisschemas werden an ihren vorhandenen Ownern geändert.
Nur `commands.schema.json` wird durch den vorhandenen JSON-Schema-Generator neu erzeugt.
Der generische Erfolgsvertrag erlaubt zusätzlich ein optionales `search`-Objekt.
Das Fehlerschema erlaubt die neuen Detailkeys und beschreibt die Search-Daten mit lokalen Definitionen.

Halte `error.schema.json` eigenständig und ohne externen `$ref` auf das generierte Commands-Schema.
Damit entstehen weder zusätzliche Netzwerkauflösung noch eine versteckte Mischung gepinnter und aktueller Schemas.
Die begrenzte doppelte Schemadarstellung wird durch die vorhandenen Vertragsprüfungen gegen echte Ausgabeformen geprüft.
Der bestehende Testcompiler registriert bislang nur ein Schema pro Aufruf.
Beleg: `cmd/json_command_schema_test.go:343-366`.

### Dokumentation, Skills und vorhandene Fixtures

Beschreibe Serverfilter, lokale Filter, Stopgründe, Teilergebnisse und die Auswahl vor lokaler Sortierung.
Passe bestehende Erfolgs- und Fehlerbeispiele sowie die Feldinventare im vorhandenen Vertragsaudit an.
Setze für eingebettete Search-Metadaten die reflektierten Basis- und Zusatzfelder gezielt zusammen.
Der bisherige `jsonFieldNames()`-Helper erfasst unexportierte eingebettete Structs nicht automatisch.
Beleg: `cmd/json_contract_test.go:1154-1179`.

Kanonische Skill-Dateien erklären den Verbrauch neuer Suchzustände und Teilergebnisse.
Sie verweisen auf die zur verwendeten Fork-Version gehörenden Schemas, nicht auf unveränderte Upstream-Schemas.
Die Claude-Kopien werden ausschließlich über das bestehende Sync-Skript aktualisiert.

## JSON-v1-Kompatibilität

Die neue CLI liefert die Suchzusammenfassung bei erfolgreichen Suchen immer.
Im weiterentwickelten v1-Schema bleibt das neue Top-Level-Feld `search` optional.
Seine enthaltenen Felder sind bei vorhandenem Objekt vollständig typisiert und erforderlich.
Vorhandene Search-Tests schützen die Anwesenheit des Objekts in der neuen Implementierung.

Ein fehlendes `search` in einer älteren Antwort bedeutet unbekannten Abschlusszustand, nicht Vollständigkeit.
`schema_version: "1"` allein beweist nicht, dass die neue Suchzusammenfassung unterstützt wird.
Bestehende Erfolgshüllen, Fehlerhüllen, Ergebnis-Wrapper und ihre bisherigen Feldnamen bleiben erhalten.
Die neuen Matchfelder und Search-Fehlerdetails werden optional ergänzt.

Die Kompatibilitätsrichtungen sind ausdrücklich zu unterscheiden:

- Alte Leser bekannter Felder können neue Antworten weiterverarbeiten, wenn sie unbekannte Felder tolerieren.
- Alte geschlossene Schemas lehnen neue Felder ab und müssen aktualisiert werden.
- Aktualisierte v1-Schemas akzeptieren ältere Antworten ohne die neuen optionalen Search-Felder.
- Neue Agenten dürfen bei älteren Antworten keine nicht vorhandenen Abschlussinformationen voraussetzen.

Eine vollständig transparente Kompatibilität mit alten gepinnten geschlossenen Schema-Dateien wird nicht behauptet.
Veröffentliche Schemas und Verbraucherhinweise zusammen mit der betreffenden Fork-Version.

## Risiken und verbleibende Unsicherheiten

- Die API kann Treffer auslassen und Änderungen verzögert indizieren; die CLI kann keinen vollständigen Suchraum beweisen.
- Ohne Total-Count oder Cap-Signal bleibt die Ursache eines API-Endes an der Obergrenze teilweise unbekannt.
- Lokale Filter benötigen möglicherweise viele Seiten; es gibt keine Zusage gleicher Latenz bei gleichem Ergebnislimit.
- Doppelte IDs können unterschiedliche Metadaten enthalten; die erste passende Darstellung ist keine Snapshot-Garantie.
- Die Serverwirkung außerhalb der im Live-Nachweis erfassten Querys und Filterkombinationen bleibt ungeprüft.
- Neue Array-Flags und fehlende Filterwerte benötigen korrekte Discovery-Schemas, nicht nur korrekte Cobra-Parser.
- Ausgabe- und Schemaänderungen müssen gemeinsam geprüft werden, nicht erst in einer abschließenden Aufräumphase.
- Die Reverse-Korrektur beeinflusst auch Listings; ihre begrenzte Regression muss separat nachgewiesen werden.
- Der bestehende Auth-PreRun kann vor der Search-Validierung stattfinden; `ISSUE-032` bleibt ein separater Befund.
- Die CI beschränkt Push- und PR-Trigger auf `master`; ein Push auf `personal` ist kein automatischer Prüfnachweis.

Der CI-Anker ist `.github/workflows/ci.yml:3-7`.
Eine Änderung der CI-Trigger oder eine allgemeine Ledger-Migration gehört nicht zu diesem Search-Plan.

## Tests und Verifikation

Test decision: update.
Die Änderungen schützen materielle CLI-Verträge, eine bestehende Sortierregression und laufzeitabhängige Pagination.
Es werden vorhandene Tests erweitert, keine neuen Testdateien oder Testfunktionen angelegt.

### Gezielte bestehende Tests

- `TestSearchUsesSearchV2AndCommandOutput` schützt Filterung und Deduplizierung über mehrere Seiten vor dem Limit.
- `TestSearchJSONOutputsInputAndResults` schützt Eingabe-Echo, Suchzustand und erhaltene Matchinformationen.
- `TestSearchJSONOmitsPathWithoutScope` erhält die bestehende optionale Pfadsemantik.
- `TestSearchLimitCapsResultsAndStopsPagination` unterscheidet echte Begrenzung und exakt ausgeschöpfte Ergebnisse.
- `TestSearchJSONErrorWritesNoOutput` erhält den initialen Fehlerfall und ergänzt einen Fortsetzungsfehler.
- `TestSortEntriesReverse` reproduziert und schützt die Gleichstandsbehandlung des gemeinsamen Komparators.
- `TestJSONHelpManifestV1SelectedCommandMetadata` prüft Arrays, Item-Enums und fehlende Filter-Defaults im Discovery-Vertrag.

Im Fortsetzungsfehlerfall schreibt `search()` im JSON-Modus weiterhin noch kein Dokument.
Der Test ruft anschließend den regulären Fehlerrenderer auf und prüft genau eine Hülle mit Teilergebnissen.
Er prüft auch den erhaltenen ursprünglichen Fehlercode und dessen Exitcode-Zuordnung.
Der initiale Fehlerfall wird nicht als bereits vorhandene Teilergebnisabdeckung ausgegeben.

Erweitere außerdem die bestehenden Golden-, Schema- und Live-Output-Vertragsaudits im Package `cmd`.
Nutze insbesondere die vorhandenen Erfolgs- und Fehler-Schemavalidierungen für die neuen Ausgabeformen.
Die Tests prüfen auch eine ältere v1-Antwort ohne `search` gegen das aktualisierte Schema.
Bei den Grenzfällen zählen unterschiedliche beobachtbare Verträge, nicht sämtliche äquivalenten Eingabevarianten.

### Auszuführende Schritte nach Implementierungsfreigabe

Reproduziere zunächst die betroffenen Fehler mit den engsten bestehenden Testgrenzen.
Führe nach einer Teständerung zuerst nur die betroffenen bestehenden Tests aus.
Regeneriere Ausgaben erst nach Änderungen ihrer Owner und synchronisiere Skills nach kanonischen Skill-Änderungen.

```sh
go run ./tools/gen-docs
go run ./tools/gen-json-schemas
./scripts/sync-claude-skill.sh
```

Prüfe anschließend Text- und JSON-Hilfe sowie die tatsächliche CLI-Oberfläche einschließlich stderr und Exitcodes.
Ohne geeigneten Account-Zugriff wird die Serverfilterwirkung nicht als live bestätigt ausgegeben.
Verwende dafür die engste vorhandene CLI-/Client-Testgrenze und benenne die verbleibende Live-Evidenzlücke.

Vor Abschluss von Quelländerungen gelten die vorhandenen Repository-Gates:

```sh
go test ./...
go vet ./...
go build ./...
```

Dieser Entwurf führt keine neue Nebenläufigkeit ein.
Race-Prüfungen sind zusätzlich erforderlich, falls die tatsächliche Umsetzung doch nebenläufigkeitsrelevante Änderungen enthält.
Die bestehenden Plugin- und Skill-Vertragschecks werden nach betroffenen Skill-Änderungen ausgeführt.
Ausgeführte Checks werden als Beobachtung dokumentiert; vorgeschlagene Checks gelten nicht als bestanden.

## Umsetzungseinheiten des Plans

### Einheit 1: Verlustfreie interne Datenerhaltung

Behalte vollständige Search-Matches und trenne die Seitensammlung von ihrer Darstellung.
Bewahre zunächst das bestehende öffentliche Ausgabeformat und die bisherige Ergebnisauswahl.
Sichere die begrenzte Reverse-Korrektur mit dem bestehenden Regressionstest ab.
Diese Einheit liefert noch keine behauptete neue öffentliche Search-Funktion.

### Einheit 2: Abschlusszustand, Teilergebnisse und Deduplizierung

Ergänze Suchzusammenfassung, korrekte Limitentscheidung, ID-Deduplizierung und Fehlerdetails als zusammenhängenden Vertrag.
Passe die betreffenden Schemas, Fixtures, bestehenden Tests und Verbraucherhinweise in derselben Einheit an.
Prüfe zuerst akzeptierte Ergebnisse und danach auftretende Fehler einschließlich des Einzeldokumentvertrags.

### Einheit 3: Serverfilter und Matchinformationen

Ergänze Kategorie, Erweiterung, Matchtyp und den JSON-only-Highlight-Flag.
Liefere die zugehörigen Flag-Metadaten, Array-Schemas, Ausgabedefinitionen, Beispiele und Tests gemeinsam.

### Einheit 4: Lokale Filter

Ergänze die Filterdatei und ihre Integration vor Deduplizierung und Ergebnislimit.
Liefere präsenzabhängige Defaults, Eingabe-Echo, Größen-/Zeitsemantik und lokale Filterdiagnose gemeinsam.
Prüfe einen mehrseitigen Lauf mit lokal abgelehnten Matches und einem Duplikat vor dem später passenden Treffer.

Eine öffentliche Verhaltensänderung ist erst abgeschlossen, wenn ihr maschinenlesbarer Vertrag ebenfalls stimmt.
Schemas, Skill-Hinweise und Vertragsprüfungen sind daher kein nachgelagerter Abschlussblock.

## Finding-Bezüge und Dokumentationsstatus

`ISSUE-034` besitzt bereits die unmarkierte Ergebnisbegrenzung.
`ISSUE-031` untersucht den verwandten Verlust bestätigter Teilergebnisse bei Sammeloperationen.
`ISSUE-028` besitzt die fehlende Erklärung der Auswahl vor lokaler Sortierung.
`ISSUE-003` besitzt den Reverse-Gleichstandsfehler.
`ISSUE-029` betrifft die bestehende Lücke bei Enum-Completions.

Die Umsetzung ergänzt bestehende Akten für Sortierung, Suchzustand, Teilergebnisse und Auswahlhilfe.
Neue Akten für Serverfilter, Matchinformationen, ID-Deduplizierung und lokale Filter verwenden den globalen Allocator.
Die direkte private Implementierungsfreigabe wählt weder `Authorized-Work` noch `Publication-Target` für Upstream.
Eine allgemeine Ledger-Migration und Upstream-Publikation bleiben außerhalb des Scopes.

## Abnahmekriterien

- Aus der neuen Ausgabe ist erkennbar, ob die API endet, das lokale Limit beendet oder ein Fehler abbricht.
- Server-Seitenzustand und CLI-Abbruchgrund haben unterschiedliche, ausdrücklich erklärte Bedeutungen.
- Bereits akzeptierte Treffer bleiben bei späteren Fehlern erhalten; bewusste Filterung und Begrenzung bleiben ausgewiesen.
- Serverseitige Filter stehen im Request; lokale Filter wirken vor dem Zählen unterschiedlicher passender IDs.
- Ein früher lokal unpassender Datensatz sperrt keine spätere passende Darstellung derselben ID.
- Exakt ausgeschöpfte letzte Seiten werden nicht allein wegen der Ergebniszahl als gekürzt bezeichnet.
- Matchtypen werden nicht geraten, und Titelspans werden nicht als Inhaltssnippets ausgegeben.
- Fehlende Filterwerte bleiben von ausdrücklich gesetzten Nullwerten unterscheidbar.
- Neue Discovery-Schemas beschreiben Listen als Arrays und liefern keine unbeabsichtigten Filter-Defaults.
- Erfolgs- und Fehlerausgaben bilden im JSON-Modus jeweils genau ein Dokument.
- Neue Erfolgsantworten enthalten `search`, während aktualisierte v1-Schemas auch ältere Antworten akzeptieren.
- Bestehende Statuswerte, Fehlercodes, Exitcodes und Text-Trefferformate bleiben bis auf den begrenzten Sortierfix erhalten.
- Die verwendeten Schemas und Agentenhinweise gehören zur tatsächlichen Fork-Version.
- Erforderliche bestehende Tests, Generatoren und Repository-Gates sind mit beobachtetem Ergebnis dokumentiert.
- Weder API-Ende noch eine niedrige Ergebniszahl werden als Beweis vollständiger Kontosuche dargestellt.
