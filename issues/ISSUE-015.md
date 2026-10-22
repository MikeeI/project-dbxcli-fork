# ISSUE-015 — Downloads: rekursive Wurzel scheitert an Root-Metadaten

State: Investigating
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

Root-Cause [S]: `get` verlangt vor jeder Rekursion eine Metadatenabfrage, die für die Dropbox-Wurzel nicht unterstützt wird.
Gesprächsreferenz: `E11`; der Review ergänzt ein gültiges strukturiertes Root-Ergebnis als Teil der Korrektur.

## Reach-and-Impact

Reach [S]: `get / <lokales-ziel> --recursive` nach Normalisierung auf den SDK-Root-Sentinel `""`.
Impact [O]: Der vertragskonforme Root-Gegenpart ließ den Command vor jedem ListFolder-Aufruf scheitern.
Ein Kontozugriff oder vollständiger Live-Backup-Lauf wurde nicht ausgeführt.

## Evidence

- [S] `cmd/get.go:72-73,97-101` fragt vor der Rekursion zwingend Metadaten ab und gibt deren Fehler zurück.
- [S] `cmd/dropbox_path.go:25-30` normalisiert `/` auf den leeren Root-Sentinel.
- [S] SDK v6.6.2, `dropbox/files/client.go:1975-1976`: Root-Metadaten sind ausdrücklich nicht unterstützt.
- [S] `cmd/get.go:199-205,213-220` liefert ohne Metadaten keinen strukturierten Metadatentyp.
- [S] `internal/jsonschema/definition_schemas.go:184-187` verlangt `type` im Metadatenobjekt.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierter HTTP-Gegenpart.
Aufruf: `get / <temporäres-ziel> --recursive --output=json` über `RootCmd.ExecuteC`.
Der Gegenpart verweigert Root-GetMetadata gemäß SDK-Vertrag und bietet Root-ListFolder an.
Beobachtet [O]: `command_error=get metadata for : ...` und `list_calls=0`.
Der temporäre Harness wurde entfernt; keine Live-Dateien wurden gelesen oder übertragen.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Root-Download-Ursache.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Ein bekannter Root-Sonderfall genügt; reguläre Datei- und Referenzpfade müssen unverändert bleiben.

## Proposed-Change

Root ohne `GetMetadata` als bekannten Ordner behandeln und direkt im rekursiven Auflistungspfad verarbeiten.
Ein gültiges Root-Ordnerergebnis für JSON erzeugen, ohne nicht vorhandene Objektidentitäten zu erfinden.

## Scope-and-Constraints

- Preserve: Explizite lokale Ziele, relative Unterpfade, Pagination und sichere Dateicom­mits erhalten.
- Preserve: Gewöhnliche Ordner sowie `id:`, `rev:` und `ns:` bleiben auf ihren bestehenden Vertragswegen.
- Review: Bloßes Entfernen der Root-Abfrage ohne gültiges strukturiertes Ergebnis ist unvollständig.
- Exclude: Keine allgemeine Metadatenabfrage- oder Downloadarchitektur einführen.

## Verification

Abgeschlossen [O]: Root-GetMetadata-Verweigerung und unerreichtes ListFolder wurden offline beobachtet.
Erforderlich: Root-Metadaten verweigern und eine gültige Root-Auflistung mit Inhalt anbieten.
Erwartet: Textdownload und JSON-Ausgabe funktionieren einschließlich gültigem Root-Ordnerergebnis.
Test decision: none; keine Source- oder Teständerung in diesem Auftrag.

## Publication-Blockers

- Reproduktionsablauf samt Root-JSON-Vertrag dauerhaft nachvollziehbar dokumentieren.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Root-Download gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zu rekursiven Root-Downloads und Root-Ergebnisdarstellung recherchieren.
Done-When: Suchabdeckung, relevante Treffer und der vollständige Text-/JSON-Sonderfall sind dokumentiert.
