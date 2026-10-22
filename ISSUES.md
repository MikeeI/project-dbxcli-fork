# Issue and Pull Request Tracking

Read this index at the start of every agent session before repository work.
`FORMAT.md` owns research, lifecycle, drafting, implementation, and publication rules.
Each linked `issues/ISSUE-NNN.md` is the complete authoritative record for one root cause.
This file owns `Next finding ID` and projects current issue-file state.
`Next-Action` is the 2–6 word `Next-Action/Summary` projection from the issue record.
When a row disagrees with its issue file, correct the row from the issue file in the same task.

Next finding ID: ISSUE-023

Die Akten ab `ISSUE-006` übernehmen die bestätigten Befunde des semantischen Audits einschließlich Review-Korrekturen.
Die ursprünglichen E-IDs bleiben in den Akten als Gesprächsreferenzen erhalten.
`E10` wurde mangels belegten Validierungsvertrags zurückgezogen; die Abgrenzung ist in `issues/ISSUE-006.md` dokumentiert.

## Open-Findings

| ID | Finding | State | Authorized-Work | Publication-Target | Contribution-Priority | Next-Action | External-Reference |
| --- | --- | --- | --- | --- | --- | --- | --- |
| [ISSUE-001](issues/ISSUE-001.md) | Uploads: vermeidbare Zwischenpuffer für vollständige Chunks | Investigating | Not-Selected | Not-Selected | Medium | Allokationsdruck ohne Produktivänderungen messen | Not published. |
| [ISSUE-002](issues/ISSUE-002.md) | CLI-Start: unnötige Namespace-Abfrage bei account und du | Investigating | Not-Selected | Not-Selected | Low | Prior Art vollständig recherchieren | Not published. |
| [ISSUE-003](issues/ISSUE-003.md) | Rückwärtssortierung: instabile Reihenfolge bei gleichen Schlüsseln | Investigating | Not-Selected | Not-Selected | Medium | Gleichstandsfehler offline reproduzieren | Not published. |
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

## Archived-Findings

| ID | Finding | Authorized-Work | Publication-Target | Contribution-Priority | Archive-Reason | External-Reference |
| --- | --- | --- | --- | --- | --- | --- |
