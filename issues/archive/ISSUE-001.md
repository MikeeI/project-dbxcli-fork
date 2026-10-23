# ISSUE-001 — Uploads: vermeidbare Zwischenpuffer für vollständige Chunks

State: Archived
Authorized-Work: Not-Selected
Publication-Target: Not-Selected
External-Reference: Not published.
Contribution-Priority: Medium
Root-Cause-Confidence: High
Finding-Category: Performance
Created: 2026-10-22
Updated: 2026-10-23
Source: `upstream/master@8fc3850baecaaa751e85b1673f619d29b68f3b76`

## Root-Cause

Root-Cause [S]: `uploadChunked` nutzt `io.ReadAll`, obwohl die Größe vollständiger Chunks bereits bekannt ist.
`io.ReadAll` erzeugt Zwischenpuffer und kopiert deren Inhalt anschließend in einen weiteren Ergebnispuffer.
Gesprächsreferenz: `P1`; der anschließende Review begrenzt den Vorschlag ausdrücklich auf vollständige Chunks.

## Reach-and-Impact

Reach [S]: `putFileWithResult` aktiviert Chunk-Uploads für lokale Dateien über 32 MiB.
Der Standardchunk umfasst 16 MiB; standardmäßig laufen vier Upload-Worker.
Impact [S]: Jeder vollständig gelesene Standardchunk verursacht mindestens 16 MiB vermeidbare Zusammenkopierarbeit.
Die vorgeschlagene Änderung senkt das Pufferallokationsvolumen, nicht nachweislich die Übertragungsdauer oder den RSS.
Die Optimierung greift nicht, wenn die konfigurierte Chunkgröße die gesamte Quelldatei übersteigt.

## Evidence

- [S] `cmd/put.go:149-159` kennt Gesamtgröße, Chunkgröße, Offset und erwartete Länge vor der Übergabe an Worker.
- [S] `cmd/put.go:597-598` aktiviert den Pfad oberhalb 32 MiB; `cmd/put.go:1099-1100` definiert die Defaults.
- [S] `cmd/put.go:496-503` erlaubt Chunkgrößen von 4 MiB bis 128 MiB in Vielfachen von 4 MiB.
- [S] `cmd/put.go:58-65` erstellt für jeden Retry einen neuen Reader über dieselbe unveränderliche Payload.
- [S] Go 1.26.6 sammelt Zwischenpuffer und erstellt danach eine Ergebniskopie: https://github.com/golang/go/blob/go1.26.6/src/io/io.go#L715-L753.
- [S] Die untersuchte Go-1.27.3-Laufzeit verwendet denselben Ansatz in `src/io/io.go:709-747`.
- [O] Der ausgeführte Vergleich mit dem frisch abgerufenen `upstream/master` zeigt keine Abweichung in `cmd/put.go`.

## Performance-Evidence

Cost model [S]: Für `q` vollständig gelesene Chunks der Größe `C` entstehen mindestens `2*q*C` Pufferallokationsbytes.
Davon entfallen mindestens `q*C` Bytes auf Zwischenpuffer und `q*C` Bytes auf endgültige Ergebnispuffer.
Zusätzlich werden `q*C` Nutzdatenbytes in die Ergebnispuffer kopiert.
Bei 16-MiB-Standardchunks sind das mindestens 32 MiB Pufferallokationsvolumen und 16 MiB Kopie je vollständigem Chunk.
Allocator-Rundung, Metadaten, Teilchunk, HTTP, Hashing und GC sind in diesen Untergrenzen nicht enthalten.
Measurement status: Quellcodeabgeleitete Untergrenzen; keine gemessenen Allokationen oder Ende-zu-Ende-Gewinne.

## Prior-Art

Coverage: Lokaler Index und `issues/` am 2026-10-22 geprüft; keine vorhandenen Finding-Datensätze.
Gaps: Upstream-Issues, PRs, Diskussionen und Releases wurden für diesen Befund noch nicht recherchiert.
Contribution fit: Offen; messbarer Nutzen und bestehende Upstream-Arbeiten sind vor einer Empfehlung zu klären.

## Proposed-Change

Vollständige Chunks direkt in einmalig auf die Chunkgröße dimensionierte Puffer lesen.
Den bisherigen Teilchunk-Pfad einschließlich seiner Längenprüfung unverändert lassen.
Nicht-EOF-Lesefehler auch bei gleichzeitig vollständiger Datenlieferung weitergeben.
EOF und Kurzlesen müssen weiterhin die bisherige Behandlung beziehungsweise Fehlermeldung auslösen.
Der Nutzer autorisierte die private Umsetzung dieses Review-Befunds `P1` am 2026-10-23 mit „fix all issues“.
Diese Freigabe wählt weder einen Upstream-Arbeitsmodus noch ein Veröffentlichungsziel.

## Scope-and-Constraints

- Preserve: Chunkgrenzen, Offsets, `Close`, Queue-Kapazität und Abschluss erst nach allen Appends erhalten.
- Preserve: Payloads bis zum Abschluss aller Upload-Retries exklusiv halten und nicht verändern.
- Preserve: Der Teilchunk muss zwischenzeitliches Dateiwachstum weiterhin durch seine Längenprüfung erkennen.
- Preserve: SDK-Integritätsprüfung und wiederholbares Lesen derselben Chunkdaten erhalten.
- Exclude: Kein mechanischer Ersatz durch `io.ReadFull`, da dieser bestimmte abschließende Lesefehler unterdrückt.
- Exclude: Kein `sync.Pool`, keine neue Parallelität und kein erneutes Lesen der Quelldatei während Retries.
- Exclude: Bestehende Worker- und Fehlerkanalprobleme sind eigenständige Korrektheitsfragen.
- Cost: Ein begrenzter Füllpfad benötigt eine explizite Fehlerbehandlung statt eines einzelnen Bibliotheksaufrufs.

## Verification

### Historischer Untersuchungsstand

Abgeschlossen [O]: Source-Abgleich mit dem frisch abgerufenen kanonischen Upstream ohne relevante Quelldifferenz.
Nicht ausgeführt: Implementierung, Tests, Transfer-Szenarien und Laufzeitmessungen.
Nach autorisierter Implementierung ist folgende Verifikation erforderlich:

- `go test ./cmd -run '^TestUploadChunked_'` muss die bestehenden Upload-Verträge weiterhin erfüllen.
- Ein temporäres Smoke-Szenario muss Kurzlesen, kombinierte Daten-/Fehlerrückgabe und Retry-Payloads prüfen.
- Ein Vergleich identischer vollständiger Chunks muss das Allokationsvolumen vor und nach der Änderung erfassen.
- Übertragungsdauer und Spitzenspeicher sind separat zu messen und nicht aus Allokationsbytes abzuleiten.

### Privater Abschluss am 2026-10-23

[O] Commit `5843c2f` auf `personal` wurde nach `origin/personal` gepusht.
[S] `cmd/put.go:117-136,181-192` liest vollständige Chunks direkt und erhält den bisherigen Teilchunk-Pfad samt Längenprüfung.
[O] Der isolierte 16-MiB-Lesevergleich unter Go 1.27.3 auf Linux/amd64 ergibt 36224962 → 16777264 Allokationsbytes und 33 → 2 Allokationen pro Operation.
[O] Gemessen wurden 6584643 → 2103644 ns pro isolierter Leseoperation; das ist keine Ende-zu-Ende-Transfermessung.
[O] Acht temporäre Lesefälle bestätigen volle und kurze Eingaben, EOF, Quelldateiwachstum sowie Daten zusammen mit Nicht-EOF-Fehlern.
[O] Vorhandene Chunk-/Retry-Prüfungen sowie `go test ./...`, `go vet ./...`, `go build ./...` und `go test -race ./...` bestehen.
[O] `golangci-lint run --fix ./...` meldet ausschließlich den unveränderten ungenutzten Helfer `cmd/json_metadata.go:64`.
Test decision: none; keine permanente Testfunktion für die Pufferoptimierung ergänzt.
Temporäre Mess- und Smoke-Dateien wurden entfernt; keine Live-Transfers, globale Installation oder Upstream-Veröffentlichung.
`ISSUE-007` und `ISSUE-008` bleiben separate Korrektheitsbefunde.

## Publication-Blockers

- Live-Transferdauer und Spitzenspeicher sind nicht gemessen; die isolierte Allokationsmessung belegt keine Ende-zu-Ende-Beschleunigung.
- Prior Art und aktuelle Upstream-Beitragsvorgaben sind vor einer Veröffentlichung zu prüfen.
- `Authorized-Work` und `Publication-Target` sind vom Nutzer noch nicht ausgewählt.
- Ein exakter Veröffentlichungsentwurf und dessen Freigabe liegen nicht vor.

## Next-Action

Summary: —
Action: None.
Done-When: None.

## Archive

Archive-Reason: Other
Detail: Privat autorisierte Pufferoptimierung implementiert, geprüft und auf origin/personal gepusht; keine Upstream-Erledigung behauptet.
Evidence: `5843c2f`; isolierter 16-MiB-Vergleich, Lesefehler-Smokes und vorgeschriebene Go-Prüfungen.
Checked: 2026-10-23
