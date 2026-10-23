# Optimierungen für agentenorientierte Dropbox-Suche

## Ergebnis und Status

Der größte strukturelle Hebel ist weniger wiederholte Netzwerkarbeit: größere Suchseiten, wiederverwendbarer Root-Kontext und verlustfreie Ergebnissnapshots.
Ein suchspezifischer Custom-Decoder und vorberechnete Namenssortierschlüssel zeigen zusätzlich messbare CPU-Vorteile.
Caches, Deltaindex, unabhängige Suchpartitionen und ein möglicher Sessionbetrieb erweitern den Lösungsraum deutlich über weitere CLI-Flags hinaus.
Ein vollständiger SDK-Rewrite ist dafür weder Voraussetzung noch durch die bisherigen Messungen gerechtfertigt.

Stand der Untersuchung: 2026-10-23.
Dokumentationsbasis des Forks: `personal@beb52aadf4c0987e443036df069d3c53c13a9c88`.
Während der Untersuchung wurden parallel andere Performance-Fixes abgeschlossen; die Such- und Sortierbefunde wurden gegen diese Basis abgeglichen.
Insbesondere überspringen `account` und `du` inzwischen die unbenötigte Root-Abfrage; `search` benötigt weiterhin die hier untersuchte Initialisierung.

Dieses Dokument hält Forschung, Messungen, Vorschläge, Gegenbeweise und offene Nachweise fest.
Es ersetzt weder `AGENTS.md` noch den Finding-Ledger und vergibt keine neuen `ISSUE-NNN`-Kennungen.
Die früheren Nichtziele in `GOAL.md` begrenzen die damalige Umsetzung, nicht den hier ausdrücklich erweiterten Untersuchungsraum.
Die Dokumentation autorisiert keine Produktionsimplementierung, Abhängigkeitsmigration, Installation, Kontoinventur oder authentifizierte Live-Untersuchung.

## Ziel und Bewertungsmaßstab

Gewünscht sind möglichst viele brauchbare Ergebnisse mit möglichst wenigen Dropbox-Requests und möglichst geringer Latenz für LLM-Agenten.
Custom-SDKs, SDK-Forks, eigene Decoder, direkte API-Adapter und Caches sind ausdrücklich zulässige Entwurfsoptionen.
Entscheidend ist nachgewiesene oder belastbar abgeleitete Arbeitsersparnis bei erhaltenen Ergebnis-, Authentifizierungs- und Fehlerverträgen.

Diese Ziele sind nicht identisch:

- Wenige Requests bevorzugen große Seiten, Wiederverwendung und Zusammenlegen identischer Abrufe.
- Geringe Zeit bis zum ersten Treffer kann kleinere Antworten und progressive Ausgabe bevorzugen.
- Geringe Gesamtdauer kann begrenzte Parallelität unabhängiger Suchen rechtfertigen.
- Hoher Recall kann zusätzliche Suchpartitionen und damit mehr Requests erfordern.
- Geringer Modellkontextverbrauch erfordert kompakte Ausgabe, nicht zwingend weniger geladene Treffer.
- Hohe Aktualität begrenzt Cache-Wiederverwendung und benötigt gegebenenfalls zusätzliche Abgleiche.

API-Seitengröße, lokal gewünschte Ergebniszahl, sichtbarer Agent-Ausschnitt und Arbeitsbudget müssen getrennte Größen sein.
Keine dieser Größen darf stillschweigend als Vollständigkeitsgarantie ausgegeben werden.

## Evidenz und Untersuchungsumfang

### Evidenzklassen

- **Gemessen:** ausgeführter Offline-Probe mit angegebenem Datensatz und Laufzeitumgebung.
- **Quellbelegt:** direkt aus aktuellem Code oder API-Vertrag belegtes Verhalten.
- **Abgeleitet:** Kostenmodell unter ausdrücklich genannten Annahmen, ohne gemessenen End-to-End-Gewinn.
- **Offen:** plausible Option, deren Wirksamkeit, Backend-Akzeptanz oder semantische Gleichwertigkeit noch nicht verifiziert ist.

### Durchführung

Vier unabhängige Untersuchungen deckten Transport/Auth, SDK-Serialisierung, Caches/Resume und API-Limits/alternative Suchflächen ab.
Main untersuchte zusätzlich CLI-Ausgabe und Sortierung und führte die beschriebenen Offline-Probes aus.
Der Schwerpunkt lag auf dem vollständigen Suchpfad, nicht auf jeder Funktion sämtlicher geklonter Bibliotheken.
Cobra und pflag wurden als Referenzen für CLI-Lifecycle und Parsing bereitgestellt; ein Austausch ist durch keinen belegten Suchhotspot begründet.

Es wurden keine authentifizierten Dropbox-Requests, Downloads oder Kontoinventuren ausgeführt.
Die Probes verwendeten synthetische Daten, einen lokalen HTTP-Server beziehungsweise einen kontrollierten `RoundTripper`.
Die Messungen umfassen keine reale Dropbox-Netzwerklatenz und keinen vollständigen CLI-End-to-End-Vergleich.

### Referenzklone

Die folgenden Repositorys liegen lokal unter `references/`:

| Verzeichnis | Version oder Revision | Commit |
| --- | --- | --- |
| `dropbox-sdk-go/` | `v6.6.2` | `fad5ca444338899f8f4fa5c8c8cafd8521101f58` |
| `oauth2/` | `v0.37.0` | `c624b89dadc3221560b7345c090bbe69e90808ee` |
| `cobra/` | `v1.10.2` | `88b30ab89da2d0d0abb153818746c5a2d30eccec` |
| `pflag/` | `v1.0.10` | `0491e5702ad2bb108bc519a5221bcc0f52aa9564` |
| `dropbox-api-spec/` | untersuchter aktueller Checkout | `b63b5d5be0d1473b6bf237843fea977574c18fad` |

Das SDK-Modul liegt im Clone unter `references/dropbox-sdk-go/v6/`.
Die untersuchten Dateien `dropbox/sdk.go` und `dropbox/files/types.go` stimmen mit der aufgelösten Modulcache-Version überein.
Die Klone sind lokale Forschungsreferenzen, keine eingecheckten Submodule und keine aktiven Ersatzabhängigkeiten.
`go.mod` verwendet weiterhin `github.com/dropbox/dropbox-sdk-go-unofficial/v6 v6.6.2` ohne `replace` auf den Referenzclone.
Eine Änderung allein im Clone würde das Verhalten der CLI daher nicht ändern.

Der SDK-interne Versionsstempel lautet `6.4.0`, der Spec-Stempel `6d985df`.
Diese Stempel sind nicht mit der aufgelösten Modulversion `v6.6.2` gleichzusetzen.
Vor Generatoränderungen muss die zugehörige Generator-/Spec-Provenienz separat geschlossen werden.

## 1. Aktueller Suchpfad und Kostenmodell

Der normale CLI-Prozess führt einen Cobra-Command aus und beendet sich anschließend.
Die Initialisierung lädt gegebenenfalls gespeicherte Credentials, erneuert sie bei Bedarf und ermittelt für `search` den Root-Namespace.
Danach folgen `search_v2`, cursorgebundene Fortsetzungen, lokale Filter, ID-Deduplizierung, Begrenzung, optionale Sortierung und Ausgabe.

Quellanker:

- `cmd/root.go`: `initDbx`, `makeDropboxConfig`, `withRootNamespace`, Command-Kontext und globale Zustandsvariablen.
- `cmd/auth.go`: `getAccessCredential`, `shouldRefresh`, `refreshStoredCredential`.
- `cmd/search.go:127-155,224-315,356-418`: Request, Pagination, Ergebnisprojektion und Ausgabe.
- `internal/output/output.go:48-59`: Renderer für Text beziehungsweise ein JSON-Dokument.

Für einen regulären erfolgreichen Suchlauf ohne konfigurierte SDK-Retries gilt als abgeleitetes Requestmodell:

```text
Gesamtrequests = 1 Root-Abfrage + P Suchseiten + R Token-Refresh
P >= 1
R = 0 oder 1 im untersuchten Initialisierungspfad
```

Der normale gespeicherte Token wird vorab erneuert, wenn die Refresh-Bedingung erfüllt ist, insbesondere im Fünf-Minuten-Fenster vor Ablauf.
Direkte Umgebungstokens durchlaufen diesen gespeicherten Refresh-Pfad nicht.
Das Modell ist keine Zusage für beliebige Fehler-, Redirect- oder zukünftige Retry-Pfade.

Die Suche besitzt bereits native Kategorie-/Endungsfilter, ID-Deduplizierung, erhaltene Teilergebnisse und explizite Abschlusszustände.
Diese Fähigkeiten sind keine neuen Optimierungsbefunde.
Es existiert kein zusätzlicher Metadatenrequest pro Suchtreffer, den man als N+1-Problem beseitigen könnte.

## 2. Größere Suchseiten

### Befund

`files.NewSearchOptions()` setzt `MaxResults=100`.
`newSearchV2Arg()` verkleinert diesen Wert bei kleinen ungefilterten Ergebnislimits, erhöht ihn aber nicht.
Die öffentliche Spezifikation erlaubt `max_results` zwischen 1 und 1.000, mit Default 100.

Quellanker: `cmd/search.go:224-241`, `references/dropbox-sdk-go/v6/dropbox/files/types.go:4113-4144` und `references/dropbox-api-spec/files.stone:2194-2210`.
Quelle: https://github.com/dropbox/dropbox-api-spec/blob/b63b5d5be0d1473b6bf237843fea977574c18fad/files.stone#L2194-L2210

Bei 10.000 tatsächlich gelieferten Rohmatches und durchgehend vollen Seiten ergeben sich rechnerisch 100 Suchrequests mit Seitengröße 100 oder zehn mit Seitengröße 1.000.
Das ist ein idealisiertes Requestmodell, keine garantierte Seitenfüllung oder gemessene Latenzersparnis.
Duplikate, kurze Seiten und lokale Filter beeinflussen die Zahl nutzbarer Ergebnisse.

### Vorschlag

- Bei endlichem Ergebnislimit ohne lokale Filter ist `min(limit, 1000)` ein plausibler Startwert.
- Bei Bulk-Abfragen ohne lokales Limit sind bis zu 1.000 angeforderte Treffer pro Seite sinnvoll zu prüfen.
- Bei lokalen Filtern muss die Seitengröße vom gewünschten akzeptierten Ergebnislimit entkoppelt bleiben.
- Die bestehende Prüfung des Seitenrests, Deduplizierung und Unterscheidung zwischen Limit und API-Ende bleiben erhalten.

Größere Seiten erhöhen Antwortgröße, Decoding-Aufwand und transienten Speicherbedarf.
Sie verbessern weder automatisch die erste Antwortzeit noch die serverseitige Indexabdeckung.
Der Continue-Request enthält nur den Cursor; die Seitengröße lässt sich innerhalb derselben Kette nicht einfach nachregeln.

### Verifikation

Requestgröße, Cursorfolge und Stopgründe zunächst offline für kleine Limits, lokale Filter, Duplikate und volle Endseiten prüfen.
Reale Seitengrößen 100, 500 und 1.000 anschließend nur in einer gesondert autorisierten Live-Messung vergleichen.
Dabei Requests, Bytes, erste nutzbare Ergebnisse und Gesamtdauer getrennt erfassen.

## 3. Suchspezifischer Custom-Decoder

### Wiederholte Arbeit im SDK

Der generierte Decoder liest die äußere Suchantwort und delegiert jedes Metadatenobjekt an weitere Decoder.
`MetadataV2.UnmarshalJSON` dekodiert einen Wrapper und kopiert das innere Objekt als `json.RawMessage`.
`IsMetadataFromJSON` und `metadataUnion.UnmarshalJSON` bestimmen danach den inneren Tag und dekodieren dasselbe Objekt erneut als Datei oder Ordner.
Das bedeutet wiederholtes Scannen verschachtelter Objekte und eine innere Kopie, nicht mehrere Kopien des gesamten Response-Bodys.

Quellanker:

- `references/dropbox-sdk-go/v6/dropbox/files/types.go:964-1012,2616-2635`: Metadaten-Uniondecoder.
- `references/dropbox-sdk-go/v6/dropbox/files/client.go:3365-3394,3412-3441`: Search-RPC und anschließendes `json.Unmarshal`.
- `references/dropbox-sdk-go/v6/dropbox/sdk.go:217-232`: vollständiges Einlesen des RPC-Bodys.
- `references/dropbox-sdk-go/generator/go_types.stoneg.py`: Generatorbesitzer der Unionmuster.

### Gemessener Offline-Vergleich

Umgebung: Linux x64, Go `1.27.3`; die Projekt-Mindestversion bleibt Go `1.26.0`.
Die Fixture enthält ausschließlich synthetische Dateien mit Name, Pfaden, ID, Revision, Größe, beiden Änderungszeiten, Downloadbarkeit und Matchtyp.
Die Antwortgröße beträgt 37.330 Bytes für 100 Treffer und 373.030 Bytes für 1.000 Treffer.

Der Vergleichskandidat dekodiert eine flachere Wire-Struktur und baut daraus dieselben SDK-Ergebnisobjekte auf.
`reflect.DeepEqual` zwischen SDK-Ergebnis und Kandidat bestand für beide Fixtures.
Die Messung verwendete `testing.Benchmark` mit `b.Loop`; Fixture-Erzeugung und Gleichheitsprüfung lagen außerhalb des Timings.

| Treffer | Decoder | Zeit pro Operation | Bytes pro Operation | Allokationen pro Operation |
| --- | --- | ---: | ---: | ---: |
| 100 | SDK | 0,967820 ms | 92.316 | 1.338 |
| 100 | flacher Decoder mit SDK-Wiederaufbau | 0,332509 ms | 112.603 | 740 |
| 1.000 | SDK | 10,375929 ms | 941.368 | 14.073 |
| 1.000 | flacher Decoder mit SDK-Wiederaufbau | 3,304886 ms | 948.942 | 8.076 |

Der 1.000er-Vergleich zeigt rund 3,1-fach schnelleres Decoding, nicht eine 3,1-fach schnellere Dropbox-Suche.
Die Allokationsanzahl sinkt; das allokierte Bytevolumen sinkt in diesen Messungen nicht.
Allokierte Bytes sind weder retained Heap noch gemessene Speicherhöchstlast.
Es handelt sich um einen lokalen Microbenchmark-Durchlauf, nicht um wiederholte Produktionsmessungen mit Konfidenzintervallen.

### Geeignete Integrationsgrenze

Ein schmaler Search-Adapter kann eigene Wire-Typen und einen eigenen Decoder verwenden und den vorhandenen SDK-Transport behalten.
`dropbox.Context.ExecuteContext` übernimmt weiterhin Requestaufbau, Authentifizierung, Root-/Member-Header, Proxyverhalten und Kontext.
Das Projekt besitzt mit `cmd/share_link_raw.go:109-124` bereits ein entsprechendes Raw-Request-Muster.

Ein gezielter SDK-Fork oder eine zusätzliche SDK-Decoderoberfläche ist ebenfalls legitim, wenn Messungen und benötigte Verträge das rechtfertigen.
Eine globale Änderung des Generators betrifft wesentlich mehr Endpoints und benötigt eine entsprechend breitere Paritätsprüfung.
Der öffentliche SDK-Decoder darf nicht still auf reduzierte CLI-Felder umgestellt werden, wenn andere Verbraucher vollständige Metadaten erwarten.
Ein ausschließlich CLI-interner Projektionspfad kann dagegen die tatsächlich benötigten Felder explizit besitzen.

Zu erhalten sind insbesondere Datei-/Ordnerunterscheidung, IDs, Filterfelder, Zeitdarstellung, optionale Highlights und offene Matchtyp-Strings.
Unbekannte Metadatenvarianten dürfen nicht still als Datei interpretiert oder als angeblich vollständiger Erfolg übersprungen werden.
Die aktuellen CLI-Fehler für nicht unterstützte Metadaten und die ursprünglichen SDK-Fehlerklassen bleiben Teil der Kompatibilitätsprüfung.

### Offene Parität

Nicht abgedeckt sind vollständige Ordner-/Dateimischungen, optionale Media-/Property-/Sharing-Felder, Highlights und unbekannte Varianten.
Auch fehlerhaftes JSON, Abbruch, Rate Limits und vollständige CLI-Ausgabe sind nicht durch den file-only Decodervergleich verifiziert.
Vor Umsetzung ist ein gezielter Differential-Probe dieser Fälle gegen den unveränderten SDK-Pfad erforderlich.
Ein pauschaler Austausch der JSON-Bibliothek beseitigt die inneren fest eingebauten `encoding/json`-Aufrufe nicht automatisch.

## 4. Namenssortierschlüssel einmal berechnen

`cmd/format.go:63-80` führt bei jedem Namensvergleich zweimal `strings.ToLower(entryPath(...))` aus.
Bei gemischt geschriebenen Pfaden entstehen dadurch wiederholt neue Strings.

Ein Offline-Probe verglich die bestehende stabile Sortierung mit einmalig vorberechneten Kleinschreibungsschlüsseln.
Verwendet wurden 10.000 synthetische Pfade mit gemischter Groß-/Kleinschreibung und deterministischer PCG-Zufallsfolge mit Seeds 1 und 2.
Beide Varianten kopierten beziehungsweise bauten ihre Arbeitsdaten pro Operation neu auf.
Die sortierten Ergebnislisten waren per `reflect.DeepEqual` identisch.

| Stabile Namenssortierung | Zeit pro Operation | Bytes pro Operation | Allokationen pro Operation |
| --- | ---: | ---: | ---: |
| Kleinschreibung je Vergleich | 46,518093 ms | 10.892.507 | 335.271 |
| Kleinschreibung einmal je Pfad | 12,256852 ms | 811.626 | 10.005 |

Dies ist ein konkreter kleiner Optimierungskandidat für große Namenssortierungen.
Die Wirkung betrifft nicht Suchläufe ohne Namenssortierung und ist kein gemessener End-to-End-Speedup.
Die neue Schlüsselbildung müsste exakt dieselbe `strings.ToLower(PathDisplay)`-Semantik und stabile Gleichstandsbehandlung erhalten.
Das Providerfeld `path_lower` ist ohne Gleichwertigkeitsnachweis kein austauschbarer Ersatz.

## 5. Root-Kontext und OAuth-Lifecycle

Die Suche ruft vor ihrer ersten Ergebnisseite weiterhin `users/get_current_account` auf.
Ein korrekt gebundener Root-Cache oder wiederverwendeter Sessionkontext kann diesen seriellen Request amortisieren.
Pauschales Weglassen verändert potenziell Teamroot und sichtbaren Suchraum.

Ein Root-Cache benötigt eine stabile Bindung an Auth-Profil, Domain, Account, App-/Zugriffskontext und ausgewähltes Mitglied.
Gespeicherte Credentials enthalten derzeit keine fertige Account-/Root-ID; diese Identität muss zunächst korrekt aufgelöst werden.
Unklare Identität erzwingt einen Cache-Miss statt Wiederverwendung fremden Kontexts.

Die `root`-Variante von `Dropbox-API-Path-Root` meldet bei geändertem Root `invalid_root` einschließlich aktueller RootInfo.
Dieser Vertrag eignet sich als konkrete Recovery-/Invalidierungsgrenze, ersetzt aber keine vollständige Cache-Identität.
Quellanker: `references/dropbox-api-spec/common.stone:21-32,69-75`.
Quelle: https://github.com/dropbox/dropbox-api-spec/blob/b63b5d5be0d1473b6bf237843fea977574c18fad/common.stone#L21-L75

Der aktuelle CLI-Pfad erneuert gespeicherte Tokens vorab und übergibt anschließend einen statischen Access-Token an das SDK.
Ein sehr langer Suchlauf oder ein neuer langlebiger Dienst benötigt daher einen bewusst erweiterten Refresh-Lifecycle.
Ein gemeinsam verwendeter TokenSource muss Rotation persistieren und parallele Refreshes desselben Principals koordinieren.
Ein Prozessmutex allein schützt nicht mehrere CLI-Prozesse mit derselben Auth-Datei.
Die bestehenden Credential-Persistenzbefunde bleiben eigenständige Korrektheitsvoraussetzungen für aggressive Parallelität.

Die Dropbox-OAuth-Hilfe setzt den Auth-Stil ausdrücklich auf `AuthStyleInParams`.
Ein vermeintlicher zweiter OAuth-Request durch Auth-Stil-Probing ist auf diesem Pfad kein belegter Optimierungsfall.

## 6. Getrennte Cache- und Snapshot-Arten

### Ergebnissnapshot und Resume

Ein Snapshot bewahrt eine bereits geladene Suchgeneration für weitere Ausschnitte, lokale Sortierung oder zusätzliche lokale Filter.
Er ist eine explizite historische Sicht, keine Behauptung über den aktuellen Dropbox-Zustand.
Bereits geladene Daten können ohne erneute Suchrequests ausgegeben werden; noch nicht geladene Seiten benötigen weiterhin Fortsetzungen.

Ein Dropbox-Cursor allein reicht nicht für verlustfreie Fortsetzung.
Wenn eine Seite 1.000 Matches enthält und der Agent zunächst 50 sieht, müssen die übrigen Matches erhalten bleiben.
Der Cursor verweist bereits hinter die geladene Seite.

Ein Resume-Zustand benötigt:

- Bereits geladene, noch nicht ausgegebene Matches und lokale Position.
- Folgeseiten-Cursor und Suchgeneration.
- Effektive Serveroptionen sowie Account-/Root-/Member-/Domain-Kontext.
- Deduplizierungszustand und bei Bedarf lokale Auswahlparameter.
- Atomare Veröffentlichung von Seitendaten und zugehörigem Fortschritt.
- Ablauf-, Lösch- und Größenregeln für die gesamte Snapshotgeneration.

Cursor-Lebensdauer und Wiederverwendbarkeit nach längerer Unterbrechung sind nicht hinreichend geschlossen.
Ein ungültiger Cursor muss sichtbar behandelt werden; stiller Neustart ist kein identischer Snapshot.
Ein fortgesetzter Serverabruf garantiert auch keinen transaktional unveränderten Dropbox-Datenbestand.

### Exakter Query-Cache

Ein gültiger Hit kann sämtliche Search-/Continue-Requests einer identischen Abfrage vermeiden.
Das verhindert nicht automatisch die vorgelagerte Root-Abfrage oder einen erforderlichen OAuth-Refresh.

Der Schlüssel muss Principal und effektive Serveroptionen umfassen, einschließlich Query, Scope, Content-Modus, Kategorien, Endungen, Reihenfolge und Highlights.
Für fertige CLI-Antworten gehören zusätzlich lokale Filter, explizite Nullwerte, Limit, Sortierung und Ausgabe-/Schemaversion dazu.
Rohseiten und lokale Projektionen können getrennte Schlüssel besitzen, wenn die gespeicherte Abdeckung dafür ausreicht.
Querystrings dürfen nicht ohne Semantiknachweis kleingeschrieben oder heuristisch zusammengelegt werden.

Eine gecachte erste Seite ist keine ausreichende Antwort für ein höheres Limit oder einen anderen lokalen Filter.
Teilresultate und `stop_reason=limit` dürfen nicht als abgeschlossene Gesamtmenge wiederverwendet werden.
Auch ein abgeschlossener Search-Snapshot beweist keine vollständige Kontoabdeckung.

Ein sinnvoller Vertrag benötigt explizite Frischepolitik, Abrufzeit/Alter, Bytebegrenzung, Eviction und Invalidierung.
Es gibt keine evidenzbasierte Universal-TTL, die aktuelle Dropbox-Inhalte garantiert.
Ein transparenter stale-if-error-Fallback darf Auth-, Berechtigungs- oder Root-Fehler nicht verdecken.
CLI-eigene Schreiboperationen können konservativ invalidieren; Änderungen anderer Clients bleiben ohne Abgleich unbeobachtet.

### Laufende identische Abfragen zusammenlegen

In einem gemeinsam betriebenen Abrufprozess kann ein Leader einen Suchlauf ausführen und identische Waiter bedienen.
Der Vorteil ist weniger doppelte Arbeit bei gleichzeitigen Agenten, nicht automatisch eine frühere Antwort für jeden Waiter.

Gemeinsam sind nur identische Principal-/Requestkontexte.
Ein einzelner abbrechender Waiter darf andere Verbraucher nicht abbrechen.
Wenn alle Verbraucher weg sind, kann der gemeinsame Abruf beendet werden.
In-flight-Einträge, Ergebnisbytes und parallele Leader benötigen Grenzen; Fehler dürfen nicht unbeabsichtigt dauerhaft gecacht werden.

Eine globale In-Memory-Map hilft nicht zwischen separaten One-shot-CLI-Prozessen.
Persistente Caches und Indexabfragen benötigen dagegen nicht zwangsläufig einen Daemon.
Ein langlebiger Prozess wird vor allem für gemeinsame aktive Abrufe, Verbindungspools und kontinuierliche Überwachung interessant.

### Inhalts- und Extraktionscache

Die aktuelle `--content`-Suche lässt Dropbox suchen und lädt keine Dokumente zur lokalen Extraktion herunter.
Ein Inhaltscache spart deshalb nicht unmittelbar die bisherigen Search-Requests.
Er kann in einem erweiterten Agent-Workflow wiederholte Downloads, Exports und Textextraktion vermeiden.

Der Schlüssel benötigt mindestens Principal, Datei-ID, Revision, Export-/Extraktionsformat und Extractor-Version.
Datei-ID oder Pfad allein würden nach einer Bearbeitung alten Inhalt liefern.
Größe, Eviction und Sichtbarkeits-/Frischeaussagen sind getrennt vom Metadatencache zu besitzen.
Bulk-Downloads, OCR und lokale Inhaltsindexierung wurden nicht ausgeführt.

## 7. Metadaten-Deltaindex

Für häufige Namen-, Pfad-, Größen-, Typ- oder Datumsabfragen ist eine cursorgepflegte Inventur eine starke Architekturoption.
Nach einer Erstinventur können solche Abfragen lokal ohne erneuten Search-Request beantwortet werden.
Das ist eine eigene Suchfähigkeit, kein stiller Ersatz für die öffentliche Inhaltssuche.

`list_folder` erlaubt ein ungefähres Seitenlimit bis 2.000 und cursorbasierte Fortsetzung für Bestand und Änderungen.
Die Spezifikation warnt bei großen rekursiven Bäumen vor möglichen Performanceproblemen und nennt nichtrekursive Traversierung als Ausweichweg.
Quelle: https://github.com/dropbox/dropbox-api-spec/blob/b63b5d5be0d1473b6bf237843fea977574c18fad/files.stone#L1814-L1847

Wesentliche Invarianten:

- Erstinventur vollständig abarbeiten, bevor der gewählte Scope als vollständig inventarisiert markiert wird.
- Dateiersetzung, Ordneränderung und Löschmetadaten einschließlich betroffener Unterbäume korrekt anwenden.
- Seiteneinträge und Cursorcheckpoint atomar persistieren; eine Wiederholung nach Crash muss idempotent sein.
- Bei Cursorreset oder unvollständiger Indexgeneration den Frische-/Abdeckungsstatus zurücksetzen und neu aufbauen.
- Root, Mitglied, gemountete Ordner und Listingoptionen in die Indexidentität aufnehmen.
- Einen teilweise evicteten Index niemals als vollständige Inventur weiterverwenden.

Die Merge-Regeln und die Warnung vor gleichzeitigen identischen Listings stehen im offiziellen Vertrag.
Quelle: https://github.com/dropbox/dropbox-api-spec/blob/b63b5d5be0d1473b6bf237843fea977574c18fad/files.stone#L2845-L2877

`get_latest_cursor` liefert keine bestehenden Einträge und ist daher kein Ersatz für die Erstinventur.
`longpoll` meldet Änderungen, liefert aber keine Deltadaten; dafür bleibt `list_folder/continue` notwendig.
Longpoll-Backoff, Cursorreset und Indexverzug müssen sichtbar behandelt werden.
Ein gesunder Cursor garantiert nicht jede Darstellungsänderung: reine Pfad-Groß-/Kleinschreibungsänderungen werden laut SDK-Metadatenkommentar nicht über Continue gemeldet.

Der Index enthält Metadaten, keine Dropbox-Volltext-, OCR- oder Relevanzinformationen.
Lokale Dokumentinhaltssuche benötigt zusätzliche Beschaffung und Extraktion und besitzt eigene Kosten und Abdeckungsgrenzen.

## 8. Mehr Recall durch unabhängige Suchpartitionen

Die öffentliche Search-/Continue-Kette liefert höchstens 10.000 Matches.
Dropbox dokumentiert zusätzlich Indexverzug, Duplikate und mögliche Auslassungen.
Ein größeres lokales Limit oder eine größere einzelne `max_results`-Zahl hebt diese Grenze nicht auf.
Quelle: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/files/search-v-2

Mehrere unabhängige Suchketten können insgesamt mehr Treffer erschließen, beispielsweise über bekannte Ordnerbereiche.
Begrenzte Parallelität kann die Wandzeit dieser unabhängigen Ketten reduzieren.
Sie kann gleichzeitig Gesamtrequests, Payload und Rate-Limit-Druck erhöhen.

Eine belastbare Partitionierung muss folgende Grenzen berücksichtigen:

- Root-Dateien, Unterbäume, Mounts und Teamspaces dürfen nicht zwischen Teilbereichen verloren gehen.
- Kategorien und Endungen sind nicht automatisch disjunkte oder vollständige Partitionen.
- Dateien ohne Endung und nicht eindeutig klassifizierte Dateien brauchen eine explizite Abdeckung.
- IDs entfernen Überschneidungen, beweisen aber keine Vollständigkeit.
- Jede Partition behält ihren eigenen Cursor, Abbruchgrund, Fehler und Abdeckungsstatus.
- Zusammengeführte Serverrangfolgen ergeben ohne vergleichbare Scores keine belegte globale Relevanzordnung.

Die API liefert keinen verlässlichen Gesamtzähler oder `cap_reached`-Marker.
Natürliches Ende und Obergrenze lassen sich bei genau 10.000 Matches nicht zuverlässig unterscheiden.
Auch `has_more=false` oder `api_exhausted` beweist keine vollständige Kontosuche.

### Parallelität und Prefetch

Cursorseiten derselben Suche können nicht unabhängig parallel angefordert werden, weil die nächste Anfrage den Cursor der vorherigen Antwort benötigt.
Nach Empfang einer Seite kann der nächste Netzwerkabruf jedoch mit lokaler Verarbeitung oder Ausgabe der vorherigen Seite überlappt werden.
Ein begrenzter Ein-Seiten-Prefetch ist damit denkbar, aber keine Requestreduktion.
Er kann bei erreichtem Limit oder Abbruch unnötige Requests erzeugen und passt eher zu bewusst gewählten Bulk-Läufen.

Beliebiger Query-/Synonym-Fan-out besitzt keinen belegten Vollständigkeits- oder Gleichwertigkeitsvertrag.
Die öffentliche Query ist kein dokumentierter Ersatz für eine universelle Boolean-/Regex-Sprache.
Der deprecated v1-Search mit Offset ist kein belastbarer Ausweg aus den modernen Suchgrenzen.

## 9. Serverfilter und alternative Suchflächen

### Öffentliche Filter

Pfadscope, Kategorien und Endungen sind bereits native Einschränkungen.
Sie reduzieren Kandidaten nur für eine tatsächlich gewünschte Teilmenge; enger filtern ist nicht kostenloser Recall-Gewinn.
Die kombinierte Semantik und Klassifikationsabdeckung darf nicht stärker behauptet werden als belegt.

Namens-, Größen-, Typ- und Zeitfilter der aktuellen CLI wirken dagegen lokal auf bereits empfangene Matches.
Ein kleines akzeptiertes Ergebnislimit kann deshalb trotzdem viele Seiten verbrauchen.
Lokale `--sort`-Sortierung erfolgt nach Auswahl und ist kein globales Server-Top-N.
`--content` erweitert die Suchfläche und kann Trefferzahl sowie Backend-Arbeit erhöhen.

### Android-FSS

Der bestehende Reverse-Bericht belegt zusätzliche Android-Requestfelder für Datumsuntergrenzen, umgekehrte Serverordnung, Tags und angeforderte Snippets.
Ein serverseitiger Datumsfilter könnte lokale Verwerfungen reduzieren.
Echte Inhaltssnippets könnten nachgelagerte Dokumentabrufe durch Agenten vermeiden.

Die statische App-Konstruktion beweist weder Akzeptanz mit CLI-App-Key/User-Grant noch die tatsächliche Filterwirkung.
Eine erfolgreiche HTTP-Antwort allein würde ebenfalls nicht beweisen, dass ein zusätzliches Feld ausgewertet statt ignoriert wurde.
Datumsgrenzen, Zeitzonen, Trefferkorpus und Snippetsemantik benötigen kontrollierte Gegenbeispiele.
Ein nachgewiesenes unteres Datum ist keine belegte beliebige Zeitintervallsuche.

### DCS und Ask

`/2/dcs/search` ist ein eigener interner Suchvertrag mit anderen Filtern, Gates und Ergebnisstrukturen.
Der untersuchte Android-Consumer liefert keinen geschlossenen All-results-Fortsetzungsvertrag; sein Mapping setzt `hasMore=false` und `cursor=null`.
DCS ist daher kein nachgewiesener Ersatz für umfangreiche paginierte FSS-Suche.
Ask/Dash-Antwortstreams sind ebenfalls keine austauschbaren Dateitrefferlisten.

### Property Search

`file_properties/properties/search` ist für bereits gepflegte Custom-Property-Templates interessant.
Es ist keine allgemeine Namen-, OCR-, Inhalts- oder Metadatensuche.
App-/Ownerbindung, Templatezugriff und eigener Cursorvertrag bleiben entscheidend.
Ein fehlendes dokumentiertes Seiten- oder Gesamtlimit ist kein Beweis unbegrenzter Ergebnisse.
Quelle: https://docs.dropboxapi.com/dropbox-api/api-reference/user-endpoints/file-properties/properties-search

### Quellen und Aktualitätsgrenze

Die Android-Evidenz liegt in `~/projects/REVERSE/project-reverse-android-dropbox/context/`:

- `findings-search-comparison.md`: vergleichende Sicht und Authentifizierungsgrenzen.
- `findings-search-files.md`: FSS-/DCS-Request- und Pagingbefunde.
- `findings-search-dash.md`: zusätzliche Dash-/Ask-Verträge.

Der CLI-Baselineabschnitt des Vergleichsreports beschreibt einen älteren Forkstand.
Seine Aussagen über fehlende Kategorien, Highlights, Matchtypen, Deduplizierung oder verlorene Teilergebnisse gelten nicht mehr für den aktuellen Fork.
Die statischen Android-Befunde bleiben getrennte Forschungsevidenz, keine Live-Verifikation.

## 10. Transport: vorhandene Vorteile bewahren

Neue SDK-Kontexte erzeugen neue Client-/OAuth-Wrapper, aber im untersuchten Standardpfad keine getrennten Verbindungspools.
Die OAuth-Transports verwenden letztlich `http.DefaultTransport`.
Die SDK-RPC-Pfade lesen Response-Bodies vollständig und schließen sie, was Verbindungswiederverwendung ermöglicht.

Ein ausgeführter Loopback-Probe erzeugte drei frische SDK-Kontexte und führte drei serielle RPCs aus.
Der lokale Server beobachtete über `ConnState` genau eine TCP-Verbindung.
Dies belegt Wiederverwendung für lokales HTTP, nicht Dropbox-TLS, HTTP/2 oder reale Netzwerklatenz.

Go stellt bereits Connection-Pooling sowie HTTP/2- und gzip-Unterstützung bereit.
Ein Custom-Transport muss diese Eigenschaften, Proxykonfiguration und Authentifizierung bewusst erhalten.
Ein neuer Transport je Client würde den bestehenden gemeinsamen Pool aufspalten.
`Config.Client` hat Vorrang vor `Token` und `TokenSource`; ein unbedacht eingesetzter nackter HTTP-Client umgeht die SDK-Tokenbehandlung.

Ein eigener HTTP-Pfad ist legitim, wenn er einen belegten Engpass beseitigt, etwa eine notwendige Streaminggrenze.
Nur denselben JSON-Body über selbst geschriebenes HTTP zu laden spart keine Roundtrips und übernimmt zusätzliche Vertragsverantwortung.

## 11. Retry- und Fehlermetadaten

`makeDropboxConfig` setzt derzeit keine SDK-`RetryPolicy`.
Die SDK-Policy kann grundsätzlich 429, 503 und konfigurierte 409-Tags behandeln, Wartezeiten beachten und kontextabhängig abbrechen.
Retries reduzieren keine normalen Requests; sie sind eine Recovery-Strategie gegen verlorenen Fortschritt und komplette Neustarts.

### Offline nachgewiesene Informationsverluste

Ein temporärer kontrollierter `RoundTripper` gab dem echten Search-SDK Fehlerantworten ohne Netzwerkzugriff zurück.
Die Retry-Policy blieb wie im normalen CLI-Pfad ungesetzt.

```text
http_status=429 attempts=1 header_retry_after=17 returned_body_retry_after=0
http_status=503 attempts=1 returned_status_code=0
```

Beim 429-Fall enthielt nur der Header `Retry-After: 17`; ein entsprechender Wert im JSON-Body fehlte.
Die SDK-Fehlerrückgabe bewahrte diesen Header nicht, sodass der dekodierte Bodywert null blieb.
Beim 503-Fall erzeugte `auth.ParseError` einen `ServerError`, ohne dessen vorhandenes Feld `StatusCode` zu setzen.

Quellanker:

- `references/dropbox-sdk-go/v6/dropbox/sdk.go:239-273`: Body-/Retry-Verarbeitung und Fehler ohne Header.
- `references/dropbox-sdk-go/v6/dropbox/auth/sdk.go:35-51,75-81`: Serverfehler und Rate-Limit-Fehlerübersetzung.
- `references/dropbox-sdk-go/v6/dropbox/retry/retry.go:104-140`: Retry-After-Auswertung vor der Fehlerübersetzung.

### Konsequenz

Vor aggressiver Parallelisierung sollten Retryinformationen, Statuscodes und Fortschritt zuverlässig verfügbar sein.
Ein schmaler SDK-Fix oder eine bewusst eigene Fehlergrenze ist dafür eine legitime Custom-Implementierung.
Bei Aktivierung vorhandener SDK-Retries kann der SDK-Transport Headerwartezeiten bereits vor der Fehlerübersetzung berücksichtigen.

Retry-Owner, Versuchszahl, Gesamtdeadline und zulässige HTTP-/Endpointfehler müssen explizit sein.
Globale Retries für sämtliche Schreib- und Leseoperationen sind keine sichere pauschale Optimierung.
Mehrere Worker dürfen denselben Rate-Limit-Zustand nicht durch unabhängige sofortige Wiederholungen verstärken.
Teilergebnisse und Cursorfortschritt sollten einen vollständigen Neustart vermeiden, soweit der jeweilige Cursorvertrag das erlaubt.

## 12. Ausgabe, Streaming und Speicher

Aktuell sammelt `collectSearchMatches` akzeptierte Treffer bis zum Ende oder Abbruch und rendert danach.
Globale lokale Sortierung benötigt den ausgewählten Bestand; unsortierte progressive Ausgabe besitzt dagegen eine andere mögliche Grenze.

Für Agenten sollten geladene Daten und sichtbarer Modellkontext getrennt werden:

- Große Ergebnismengen effizient abrufen und als Snapshot oder Artefakt verfügbar halten.
- Eine kompakte auswählbare Projektion mit benötigten Feldern ausgeben.
- Weitere Ausschnitte ohne Wiederholung derselben Suche anbieten.
- Zähler, Abbruchgründe, Quelle, Alter und Abdeckung unabhängig von der Trefferprojektion sichtbar halten.

Weniger CLI-Ausgabefelder sparen Modellkontext und lokale Verarbeitung, aber nicht automatisch Dropbox-Netzwerkbytes.
Die öffentliche Search-API besitzt keine allgemeine Feldprojektionsoption; Titel-Highlights sind eine gesonderte optionale Erweiterung.
Highlights sind bereits standardmäßig ausgeschaltet und keine neue Requestoptimierung.

### Streaminggrenzen

Textstreaming oder ein neuer NDJSON-Vertrag könnten erste Ergebnisse früher ausgeben.
Streaming spart ohne frühen Abbruch keine API-Requests, kann aber bei geschlossenem Verbraucher weitere Seiten vermeiden.
Eine langsame Ausgabe benötigt Rückstau statt unbegrenzt wachsender Zwischenpuffer.

JSON v1 garantiert aktuell ein einziges Erfolgs- oder Fehlerdokument mit erhaltenen Teilergebnissen.
Nach bereits ausgegebenen Erfolgsteilen kann ein späterer Fehler nicht einfach durch eine zweite Fehlerhülle repariert werden.
Progressive Ereignisse benötigen deshalb einen ausdrücklichen neuen Ausgabevertrag; sie dürfen nicht still JSON v1 umdeuten.

Das SDK puffert jede RPC-Antwort vollständig, bevor der generierte Decoder läuft.
Ein echter elementweiser Streamingdecoder braucht eine Transportgrenze mit klarer Body-Ownership, Cleanup und Fehlerbehandlung.
`json.Decoder.Decode(&SearchV2Result)` allein ist kein Nachweis inkrementeller Verarbeitung einzelner Matches.
Auch ein Streamingdecoder entfernt nicht automatisch die Retention für spätere Sortierung und Ausgabe.

## 13. Sessionbetrieb und Parallelität

Ein langlebiger Prozess kann Verbindungspools über mehrere Agent-Anfragen, gemeinsamen Root-Kontext und laufende identische Abrufe teilen.
Das ist eine mögliche Betriebsform, keine Voraussetzung für jeden Cache oder Index.

Der aktuelle CLI-Prozess ist nicht durch paralleles Aufrufen von `RootCmd.Execute` in einen sicheren Dienst verwandelbar.
`config`, Command-Kontext, Cancel-Funktion und Auth-Kontext sind globale Zustände; Cobra-Flags sind ebenfalls veränderlich.
Vor parallelen Anfragen müssen Konfiguration, Auth-Bindung, Deadline und Ergebniszustand auf den einzelnen Aufruf begrenzt werden.
Geteilt werden dürfen nur ausdrücklich gemeinsam besessene Ressourcen wie kompatible Transports und koordinierte Tokenquellen.

Unabhängige Queries benötigen begrenzte Parallelität, Warteschlangen-/Speichergrenzen und einen eindeutigen Shutdown-Owner.
Prozessübergreifende Credential-Schreibvorgänge benötigen eigene Koordination und sichere Persistenz.
Ein Sessionprozess darf keine fremde Account-/Root-Konfiguration oder Deadline zwischen Requests übernehmen.

## 14. Verworfene oder eingeschränkte Optimierungsannahmen

- Neue SDK-Clients bedeuten nicht automatisch neue DNS-/TLS-/TCP-Verbindungen; die lokale Wiederverwendung wurde beobachtet.
- Ein neuer eigener Transport je Client wäre gegenüber dem gemeinsamen Standardpool eher ein Rückschritt.
- HTTP/2 und gzip sind nicht pauschal fehlende Fähigkeiten; die konkrete Dropbox-Nutzung ist noch ungemessen.
- SDK-Body-Cleanup ist im untersuchten RPC-Pfad kein belegtes Leak.
- Search-Logging serialisiert nicht heimlich ganze Antworten bei ausgeschaltetem Logging.
- OAuth-Auth-Stil-Probing erzeugt im untersuchten Dropbox-Refresh-Pfad keinen belegten Zusatzrequest.
- Ein JSON-Bibliothekswechsel beseitigt die generierten inneren Decoder nicht automatisch.
- Eine universelle persistente Cache-Schicht besitzt keine einheitliche korrekte Frische- und Invalidierungsregel.
- Ein Metadatenindex ersetzt weder Dropbox-Volltextsuche noch OCR oder Serverrelevanz.
- `get_latest_cursor` ersetzt keine Inventur vorhandener Dateien.
- Parallelität hebt Cursorabhängigkeiten und Rate Limits nicht auf.
- DCS ist ohne belegte Pagination und Auth-/Feature-Akzeptanz kein nachgewiesenes Upgrade für Massenergebnisse.
- Zusätzliche Wrapper, MCP oder ein Daemon beseitigen die eigentlichen Abrufkosten nicht von selbst.
- Ein vollständiger SDK-, OAuth- oder Cobra-Rewrite hat derzeit keinen belegten Vorteil gegenüber gezielten Eingriffen.

## 15. Priorisierung und nächste Nachweise

### Reihenfolge einer möglichen Umsetzung

1. Größere bedarfsgerechte Suchseiten und korrekte Wiederverwendung des Root-Kontexts.
2. Verlustfreie Snapshots/Resume mit getrenntem API-Abruf- und Agent-Ausgabelimit.
3. Vorberechnete Namenssortierschlüssel und ein gezielt verifizierter Custom-Search-Decoder.
4. Erhaltene Fehlermetadaten, begrenzte Recovery und transparente Arbeitsbudgets.
5. Wiederholungsorientierte Query-/Inhaltscaches und ein optionaler Metadaten-Deltaindex.
6. Bei belegtem Bedarf Sessionbetrieb, gemeinsame laufende Abrufe und begrenzte unabhängige Suchparallelität.
7. Kontrollierte Partitionierung und gesonderte Verifikation zusätzlicher Android-FSS-Felder.

Diese Reihenfolge ist eine Empfehlung, keine Implementierungsfreigabe oder universelle Aussage für alle Workloads.
Ein stark wiederholungsorientierter Metadatenworkflow kann den Deltaindex früher rechtfertigen.
Ein einmaliger Bulk-Abruf profitiert dagegen zuerst von weniger seriellen Seitenrequests.

### Erforderliche Messgrößen

- Tatsächliche HTTP-Requests getrennt nach Root, OAuth, Suchstart, Fortsetzung und Retry.
- Übertragene Bytes, Rohmatches, akzeptierte eindeutige Treffer und lokale Verwerfungen.
- Zeit bis zum ersten nutzbaren Treffer und vollständige Laufzeit.
- CPU, Allokationsvolumen und separat gemessene Speicherhöchstlast.
- Cachequelle, Alter, Trefferquote, Evictions und gespeicherte Bytes.
- Anzahl Partitionen, partielle Fehler, Budgetabbrüche und jeweils erreichter Fortschritt.

Ein Ergebnislimit ist kein Requestbudget.
Optionale Request-, Byte- und Zeitbudgets sollten Abbrüche ehrlich ausweisen und erhaltene Ergebnisse fortsetzbar machen.
Sie dürfen keine versteckten zusätzlichen Vollständigkeitsgrenzen einführen.

### Offene Verifikation

- Reale Seitenfüllung, Netzwerkwartezeit und Gesamtlatenz bei 100, 500 und 1.000 angeforderten Treffern.
- Decoderparität über alle tatsächlich benötigten Varianten und Fehlerklassen.
- Sortiergleichheit bei Gleichständen, Unicodepfaden und umgekehrter Reihenfolge.
- Cacheidentität, Ablauf, Rechte-/Rootwechsel und Wiederaufnahme nach Prozessabbruch.
- Deltaindex-Korrektheit bei Löschung, Umbenennung, Mountwechsel und Cursorreset.
- Gemeinsame Abbruch-/Refresh-Semantik für gleichzeitige Agenten.
- Tatsächliche Akzeptanz und Wirkung interner Android-FSS-Felder.
- Reale Trefferabdeckung und Kosten unabhängiger Suchpartitionen.

## Ausgeführte Prüfungen und verbleibende Artefakte

Abgeschlossen sind die gezielte Quellcodeanalyse, der Abgleich der Referenzversionen und die beschriebenen Offline-Probes.
Die Decoder- und Sortiervergleiche prüften ihre synthetischen Ergebnisse auf Gleichheit.
Der Transportprobe beobachtete eine TCP-Verbindung für drei frische SDK-Kontexte.
Der Fehlerprobe bestätigte Headerverlust beim ausgewählten 429-Fall und `StatusCode=0` bei HTTP 503.

Die temporären Programme `references/performance-probe.go` und `references/error-probe.go` wurden nach Ausführung entfernt.
Die fünf Referenzklone bleiben lokal erhalten.
Es wurden keine permanenten Tests hinzugefügt, verändert, abgeschwächt, übersprungen oder gelöscht.
Es wurden für diese Untersuchung keine Produktionsfunktionen, aktiven Abhängigkeiten oder Laufzeitkonfigurationen geändert.

Die Messprogramme sind nicht Bestandteil dieses Dokuments und nicht als ausführbare Benchmark-Suite verfügbar.
Eine Wiederholung benötigt einen neu erstellten temporären Probe mit den beschriebenen Fixtures und Vergleichsgrenzen.
Die vorhandenen Ergebnisse tragen gezielte Optimierungskandidaten, aber keine behauptete Gesamtbeschleunigung, vollständige API-Parität oder vollständige Kontosuche.
