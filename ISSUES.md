# Issue and Pull Request Tracking

Read this index at the start of every agent session before repository work.
`FORMAT.md` owns research, lifecycle, drafting, implementation, and publication rules.
Each linked `issues/ISSUE-NNN.md` is the complete authoritative record for one root cause.
This file owns `Next finding ID` and projects current issue-file state.
`Next-Action` is the 2–6 word `Next-Action/Summary` projection from the issue record.
When a row disagrees with its issue file, correct the row from the issue file in the same task.

Next finding ID: ISSUE-040

Die Akten `ISSUE-006` bis `ISSUE-022` übernehmen das semantische Audit einschließlich Review-Korrekturen.
Gesprächs-IDs gelten nur innerhalb ihrer benannten Analyse; gleiche E-IDs verschiedener Audits sind keine Dubletten.
`E10` des semantischen Audits wurde mangels belegten Validierungsvertrags zurückgezogen; siehe `issues/ISSUE-006.md`.

Der Swarm-Explore-Review ergänzt 13 Akten.
Vorhandene Ursachen bleiben in `ISSUE-002`, `ISSUE-003`, `ISSUE-004` und `ISSUE-007`.
Seine `E10` bezeichnet die eigenständige Pipe-Spooling-Belastung in `ISSUE-025`, keinen zurückgezogenen Exportbefund.
Seine `E17` wird mangels aktuellem Änderungsdruck nicht angelegt.
Der Dot-Segment-Fall bleibt ohne belegte Authority blockiert.

## Open-Findings

| ID | Finding | State | Authorized-Work | Publication-Target | Contribution-Priority | Next-Action | External-Reference |
| --- | --- | --- | --- | --- | --- | --- | --- |
| [ISSUE-001](issues/ISSUE-001.md) | Uploads: vermeidbare Zwischenpuffer für vollständige Chunks | Investigating | Not-Selected | Not-Selected | Medium | Allokationsdruck ohne Produktivänderungen messen | Not published. |
| [ISSUE-002](issues/ISSUE-002.md) | CLI-Start: unnötige Namespace-Abfrage bei account und du | Investigating | Not-Selected | Not-Selected | Low | Prior Art vollständig recherchieren | Not published. |
| [ISSUE-003](issues/ISSUE-003.md) | Rückwärtssortierung: instabile Reihenfolge bei gleichen Schlüsseln | Investigating | Not-Selected | Not-Selected | Medium | Private Korrektur im Ledger abschließen | Not published. |
| [ISSUE-004](issues/ISSUE-004.md) | Uploads: doppelte Traversierung für Planung und Ausführung | Investigating | Not-Selected | Not-Selected | Low | Traversierungsrefaktorierung gegen Prior Art prüfen | Not published. |
| [ISSUE-005](issues/ISSUE-005.md) | CLI-Metadaten: getrennte Register und redundanter Auditstatus | Investigating | Not-Selected | Not-Selected | Low | Registervereinfachung gegen Prior Art prüfen | Not published. |
| [ISSUE-006](issues/ISSUE-006.md) | Downloads: Exportnamen überschreiben andere rekursive Quellen | Investigating | Not-Selected | Not-Selected | High | Exportkollision gegen Prior Art prüfen | Not published. |
| [ISSUE-007](issues/ISSUE-007.md) | Uploads: Workerfehler blockieren den Abschluss dauerhaft | Investigating | Not-Selected | Not-Selected | High | Upload-Deadlock gegen Prior Art prüfen | Not published. |
| [ISSUE-008](issues/ISSUE-008.md) | Uploads: finaler Append schließt vor ausstehenden Chunks | Investigating | Not-Selected | Not-Selected | High | Sessionabschluss gegen Prior Art prüfen | Not published. |
| [ISSUE-009](issues/ISSUE-009.md) | Credentials: Schreibfehler zerstört die gültige Auth-Datei | Investigating | Not-Selected | Not-Selected | High | Credential-Persistenz gegen Prior Art prüfen | Not published. |
| [ISSUE-010](issues/ISSUE-010.md) | Credentials: veralteter Snapshot verliert zwischenzeitliche Anmeldungen | Investigating | Not-Selected | Not-Selected | High | Credential-Rennen gegen Prior Art prüfen | Not published. |
| [ISSUE-011](issues/ISSUE-011.md) | Teamaufnahme: individuelles Fehlerergebnis wird als Erfolg bestätigt | Investigating | Not-Selected | Not-Selected | High | Teamaufnahme gegen Prior Art prüfen | Not published. |
| [ISSUE-012](issues/ISSUE-012.md) | Zielauflösung: fehlgeschlagene Typprüfung erzeugt falschen Skip | Investigating | Not-Selected | Not-Selected | Medium | Zielauflösung gegen Prior Art prüfen | Not published. |
| [ISSUE-013](issues/ISSUE-013.md) | Uploads: Autorename-Retry erzeugt zusätzliche Remote-Dateien | Investigating | Not-Selected | Not-Selected | Medium | Autorename-Retries gegen Prior Art prüfen | Not published. |
| [ISSUE-014](issues/ISSUE-014.md) | Uploads: Close-Recovery behandelt zulässiges closed nicht | Investigating | Not-Selected | Not-Selected | Medium | Close-Recovery am Protokoll eingrenzen | Not published. |
| [ISSUE-015](issues/ISSUE-015.md) | Downloads: rekursive Wurzel scheitert an Root-Metadaten | Investigating | Not-Selected | Not-Selected | Medium | Root-Download gegen Prior Art prüfen | Not published. |
| [ISSUE-016](issues/ISSUE-016.md) | Relocation: roher Quell-Basename erzeugt falsches Ordnerziel | Investigating | Not-Selected | Not-Selected | Medium | Pfadnormalisierung gegen Prior Art prüfen | Not published. |
| [ISSUE-017](issues/ISSUE-017.md) | Relocation: Dry-run ignoriert die clientseitige Skip-Policy | Investigating | Not-Selected | Not-Selected | Medium | Skip-Vorschau gegen Prior Art prüfen | Not published. |
| [ISSUE-018](issues/ISSUE-018.md) | Freigabelinks: Raw-Änderung liefert veraltete Metadaten zurück | Investigating | Not-Selected | Not-Selected | Medium | Linkmetadaten gegen Prior Art prüfen | Not published. |
| [ISSUE-019](issues/ISSUE-019.md) | Freigabe-Downloads: Traversierung verlangt optionalen Anzeigepfad | Investigating | Not-Selected | Not-Selected | Medium | Freigabe-Traversierung gegen Prior Art prüfen | Not published. |
| [ISSUE-020](issues/ISSUE-020.md) | Fehleraggregation: einzelne SDK-Ursache wird generisch umcodiert | Investigating | Not-Selected | Not-Selected | Medium | Einzelursachen gegen Prior Art prüfen | Not published. |
| [ISSUE-021](issues/ISSUE-021.md) | Listings: Löschmetadaten werden vor ihrer Anreicherung sortiert | Investigating | Not-Selected | Not-Selected | Medium | Löschmetadaten-Sortierung gegen Prior Art prüfen | Not published. |
| [ISSUE-022](issues/ISSUE-022.md) | Listings: erforderlicher Recoveryfehler wird leerer Erfolg | Investigating | Not-Selected | Not-Selected | Low | Listing-Recovery gegen Prior Art prüfen | Not published. |
| [ISSUE-023](issues/ISSUE-023.md) | Rekursiver Upload: Textmodus verschweigt übersprungene Symlinks | Investigating | Not-Selected | Not-Selected | Medium | Symlink-Diagnose offline vergleichen | Not published. |
| [ISSUE-024](issues/ISSUE-024.md) | Transferdiagnostik: dynamischer Fortschritt wird ohne Terminal ausgegeben | Investigating | Not-Selected | Not-Selected | Low | Terminalpolitik am Diagnosewriter prüfen | Not published. |
| [ISSUE-025](issues/ISSUE-025.md) | Pipe-Uploads: vollständige Eingabe benötigt temporären Speicher | Investigating | Not-Selected | Not-Selected | Medium | TMPDIR-Ausweg offline bestätigen | Not published. |
| [ISSUE-026](issues/ISSUE-026.md) | Listing: Löschzustand wird über Textmarkierungen transportiert | Investigating | Not-Selected | Not-Selected | Low | Listing-Zustandsgrenze gezielt prüfen | Not published. |
| [ISSUE-027](issues/ISSUE-027.md) | Rekursives get: Textpfad verwirft den gewählten Fehlerwriter | Investigating | Not-Selected | Not-Selected | Low | Rekursiven Fehlerwriter offline nachweisen | Not published. |
| [ISSUE-028](issues/ISSUE-028.md) | CLI-Hilfe: Begrenzung vor lokaler Sortierung bleibt unklar | Investigating | Not-Selected | Not-Selected | Medium | Verbleibende Listing-Hilfe offline abgleichen | Not published. |
| [ISSUE-029](issues/ISSUE-029.md) | Shell-Completion: bekannte Konfliktstrategien werden nicht angeboten | Investigating | Not-Selected | Not-Selected | Low | Enum-Completion offline abgrenzen | Not published. |
| [ISSUE-030](issues/ISSUE-030.md) | Gelöschte Dateien: neun ungenutzte Revisionen werden mitabgerufen | Investigating | Not-Selected | Not-Selected | Low | Revisionslimit offline vergleichen | Not published. |
| [ISSUE-031](issues/ISSUE-031.md) | Sammeloperationen: Gesamtfehler verwirft bestätigte Einzelresultate | Investigating | Not-Selected | Not-Selected | Medium | Teilergebnisvertrag offline eingrenzen | Not published. |
| [ISSUE-032](issues/ISSUE-032.md) | CLI-Validierung: fehlende Operanden werden erst nach Authentifizierung erkannt | Investigating | Not-Selected | Not-Selected | Medium | Validierungsreihenfolge offline nachweisen | Not published. |
| [ISSUE-033](issues/ISSUE-033.md) | Uploads: revisionsgebundenes Zurückschreiben fehlt im Bearbeitungszyklus | Investigating | Not-Selected | Not-Selected | Low | Bedarf für Revisionsbedingung bestätigen | Not published. |
| [ISSUE-034](issues/ISSUE-034.md) | Ergebnislisten: begrenzte Ausschnitte bleiben unmarkiert | Investigating | Not-Selected | Not-Selected | Medium | Verbleibende Listing-Kürzung offline abgrenzen | Not published. |
| [ISSUE-035](issues/ISSUE-035.md) | Uploadstatus: Lesefortschritt wirkt wie bestätigte Übertragung | Investigating | Not-Selected | Not-Selected | Low | Uploadstatus bei Verzögerung prüfen | Not published. |
| [ISSUE-036](issues/ISSUE-036.md) | Search: öffentliche Kategorie- und Erweiterungsfilter fehlen | Investigating | Not-Selected | Not-Selected | Medium | Privaten Abschluss revisionsfest erfassen | Not published. |
| [ISSUE-037](issues/ISSUE-037.md) | Search: Matchtyp und Titel-Highlights gehen verloren | Investigating | Not-Selected | Not-Selected | Medium | Privaten Abschluss revisionsfest erfassen | Not published. |
| [ISSUE-038](issues/ISSUE-038.md) | Search: doppelte IDs verbrauchen das Ergebnislimit | Investigating | Not-Selected | Not-Selected | Medium | Privaten Abschluss revisionsfest erfassen | Not published. |
| [ISSUE-039](issues/ISSUE-039.md) | Search: lokale Metadatenfilter fehlen vor der Ergebnisauswahl | Investigating | Not-Selected | Not-Selected | Medium | Privaten Abschluss revisionsfest erfassen | Not published. |

## Archived-Findings

| ID | Finding | Authorized-Work | Publication-Target | Contribution-Priority | Archive-Reason | External-Reference |
| --- | --- | --- | --- | --- | --- | --- |
