# ISSUE-002 — CLI-Start: unnötige Namespace-Abfrage bei account und du

State: Investigating
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Low
Root-Cause-Confidence: High
Finding-Category: Performance
Created: 2026-10-22
Updated: 2026-10-22
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: `initDbx` ermittelt für `account` und `du` einen Root-Namespace, den deren Nutzabfragen nicht benötigen.
Dafür wird seriell `GetCurrentAccountContext` aufgerufen, bevor das eigentliche Kommando beginnt.
Gesprächsreferenzen: `P2` und `CW1`; beide Reviews bestätigen denselben eng begrenzten Vorschlag.

## Reach-and-Impact

Reach [S]: Reguläre authentifizierte Aufrufe von `account [<account-id>]` und `du` durchlaufen diese Initialisierung.
`account` fragt danach den aktuellen oder den angegebenen Account ab; `du` fragt den Speicherverbrauch ab.
Impact [S]: Pro regulärem Aufruf entsteht ein vermeidbarer serieller Dropbox-RPC.
Eine konkrete Latenzersparnis ist nicht gemessen.
Die RPC-Zählung umfasst weder einen möglichen OAuth-Refresh noch lokale Credential-Operationen.

## Evidence

- [S] `cmd/root.go:242` und `cmd/root.go:258` rufen die Namespace-Ermittlung in beiden Auth-Pfaden auf.
- [S] `cmd/root.go:263-279` überspringt nur Team-Management-Tokens und fragt andernfalls den aktuellen Account ab.
- [S] `cmd/root.go:270-272` behandelt Fehler dieser Abfrage als nichtfatalen Hinweis, nicht als Auth-Gate.
- [S] `cmd/account.go:150-164` sendet anschließend die eigentliche Account-Abfrage.
- [S] `cmd/du.go:47-56` benötigt ausschließlich `GetSpaceUsageContext` für sein Ergebnis.
- [S] `cmd/auth.go:264-276` erledigt erforderlichen Credential-Refresh unabhängig von der Namespace-Ermittlung.
- [S] Dropbox dokumentiert `PathRoot` als Kontext für Dateipfade: https://docs.dropboxapi.com/dropbox-api/docs/team-files.
- [O] Der Vergleich mit dem frisch abgerufenen `upstream/master` zeigt keine relevante Quelldifferenz.

## Performance-Evidence

Cost model [S]: Die reguläre Ausführung umfasst einen Namespace-RPC und einen Nutz-RPC statt nur des Nutz-RPCs.
Bei `account` ohne ID wird `GetCurrentAccountContext` zweimal nacheinander aufgerufen.
Bei `account` mit ID und bei `du` sind Namespace-Abfrage und Nutzabfrage unterschiedliche Endpunkte.
Der Befund betrifft einen CLI-Kaltpfad mit einem festen zusätzlichen Netzwerkaufruf, nicht einen gemessenen Hotspot.
Measurement status: RPC-Anzahl aus dem Kontrollfluss abgeleitet; keine Live-Konto- oder Latenzmessung.

## Prior-Art

Coverage: Der aktuelle Index und die Akten `ISSUE-001` bis `ISSUE-003` wurden auf gleiche Ursachen geprüft.
`CW1` bestätigt den vorhandenen Namespace-Befund und wird nicht als zweite Finding angelegt.
Gaps: Upstream-Issues, PRs, Diskussionen und Releases wurden für diesen Befund noch nicht recherchiert.
Contribution fit: Offen; bestehende Arbeiten und der Nutzen eines kleinen separaten Beitrags sind zu prüfen.

## Proposed-Change

Die Namespace-Ermittlung ausschließlich bei `account` und `du` überspringen.
Credential-Laden, OAuth-Refresh, `currentAuthContext`, Domain und Member-Auswahl bleiben erhalten.
Der eigentliche Nutzaufruf bleibt frisch und behält seine Fehlerbehandlung.
Keinen Account-Snapshot wiederverwenden und keinen persistenten Cache einführen.
Die Änderung ist ein Vorschlag und nicht zur Implementierung autorisiert.

## Scope-and-Constraints

- Preserve: Text- und JSON-Ergebnisse sowie Fehler des eigentlichen Nutzendpunkts erhalten.
- Preserve: Direkte Tokens, gespeicherte Credentials, Refresh und `--as-member` weiterhin unterstützen.
- Preserve: Namespace-Ermittlung für dateipfadabhängige Kommandos und das bisherige Team-Verhalten erhalten.
- Exclude: Die Ausnahme nicht auf weitere ungeprüfte Kommandos ausweiten.
- Exclude: Keine neue Lazy-Initialisierungsarchitektur und kein Namespace- oder Account-Cache.
- Cost: Die überflüssige Namespace-Warnung und deren Request-Logging entfallen bewusst bei diesen beiden Kommandos.
- Compatibility: Der initiale Namespace-RPC entfällt auch als unbeabsichtigter zusätzlicher Versuch bei transienten Fehlern.
- Preserve: Keine zugesicherte Retry-Funktion ersetzen; die Fehlerbehandlung des eigentlichen Nutzaufrufs bleibt erhalten.
- Cost: Die Cheap-Win-Einstufung ist Medium; die bestehende Beitragspriorität bleibt wegen ungemessener Wirkung Low.
- Rollback: Nur die Ausnahme für diese beiden Befehle zurücknehmen; persistierter Zustand ist nicht betroffen.

## Verification

Abgeschlossen [O]: Source-Abgleich mit dem frisch abgerufenen kanonischen Upstream ohne relevante Quelldifferenz.
Nicht ausgeführt: Implementierung, Tests, Kontoabfragen und Laufzeitmessungen.
Nach autorisierter Implementierung ist folgende Verifikation erforderlich:

- Die Initialisierung samt Handler offline für beide Auth-Quellen und alle drei Nutzabfragen ausführen.
- `account` ohne ID muss genau eine Current-Account-Operation statt zwei ausführen.
- `account <id>` und `du` dürfen keine vorgeschaltete Current-Account-Operation mehr ausführen.
- Text-/JSON-Ausgabe, direkte Tokens, gespeicherte Credentials, Refresh und `--as-member` sind zu prüfen.
- Fehler des Nutzendpunkts müssen erhalten bleiben; der entfallene Namespace-Hinweis ist kein Ausgaberegressionsfehler.
- Ein dateipfadabhängiges Kommando muss seinen bisherigen Root-Namespace weiterhin erhalten.

## Publication-Blockers

- Eine kontrollierte Ablaufverifikation ohne Live-Kontozugriff fehlt.
- Prior Art und aktuelle Upstream-Beitragsvorgaben sind vor einer Veröffentlichung zu prüfen.
- `Authorized-Work` und `Publication-Target` sind vom Nutzer noch nicht ausgewählt.
- Ein exakter Veröffentlichungsentwurf und dessen Freigabe liegen nicht vor.

## Next-Action

Summary: Prior Art vollständig recherchieren
Action: Upstream-Issues, PRs, Diskussionen und Releases zum unnötigen Namespace-RPC bei Account-Abfragen prüfen.
Done-When: Suchabdeckung, relevante Treffer, Dublettenbewertung und verbleibende Recherchelücken sind dokumentiert.
