# ISSUE-019 — Freigabe-Downloads: Traversierung verlangt optionalen Anzeigepfad

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

Root-Cause [S]: Rekursive Shared-Link-Traversierung bildet die Identität ausschließlich aus optionalem `PathDisplay`.
Gesprächsreferenz: `E15`; der Review bestätigt den Widerspruch zum nullable SDK-Feldvertrag.

## Reach-and-Impact

Reach [S]: Ordner-Freigabelinks mit gültigen, im aufrufenden Account nicht eingebundenen Dateimetadaten.
Impact [O]: Der Offline-Smoke listete einen gültigen Eintrag, brach jedoch vor jedem Inhaltsdownload ab.
Tatsächliche ListFolder-Antworten eines solchen Live-Links wurden nicht aufgezeichnet.

## Evidence

- [S] `cmd/share_link_download.go:220-224` besitzt bereits den aktuellen relativen Ordner als Traversierungskontext.
- [S] `cmd/share_link_download.go:240-263` verwendet für Datei und Ordner trotzdem nur `PathDisplay`.
- [S] `cmd/share_link_download.go:337-353` lehnt dessen zulässigen leeren Wert ab.
- [S] SDK v6.6.2, `dropbox/files/types.go:909-922`: `Name` ist erforderlich; `PathDisplay` darf bei unmounted Inhalt fehlen.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierter Freigabelink-Gegenpart.
Aufruf: `share-link download https://www.dropbox.com/scl/fo/isolated --recursive <temporäres-ziel> --output=json`.
Root-Metadaten beschreiben einen Ordner; ListFolder liefert eine Datei mit `name=report.pdf` ohne Pfadfelder.
Beobachtet [O]: `command_error=share-link download: 1 error(s)` und `invalid shared link entry path ""`.
Der Inhaltsendpunkt wurde nicht aufgerufen; der temporäre Harness wurde entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Traversierungsannahme.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Der vorhandene relative Traversierungskontext und der erforderliche Name reichen zur Identitätsbildung.

## Proposed-Change

Relative Pfade aus dem aktuellen `relFolder` und dem erforderlichen einzelnen `Name` bilden.
Den Namensbestandteil validieren und vorhandene lokale Pfadprüfungen beibehalten.

## Scope-and-Constraints

- Preserve: Link-Passwort, Cursor-Pagination und sichere lokale Tempfile-Commits erhalten.
- Preserve: Root-Dateien und verschachtelte Inhalte relativ zum Freigabelink adressieren.
- Exclude: Keine Pfadidentität aus geratenen Root-Präfixen oder optionalen Anzeigedaten ableiten.
- Exclude: Keine neue Symlink-Isolationszusage oder allgemeine Dateisystem-Sandbox einführen.

## Verification

Abgeschlossen [O]: Der Offline-Smoke reproduziert die Ablehnung SDK-konformer Metadaten ohne Anzeigepfad.
Erforderlich: Root-Datei und verschachtelte Datei ohne `path_display` auflisten.
Erwartet: Beide Inhalte werden am richtigen relativen Ziel heruntergeladen.
Test decision: none; keine Source- oder Teständerung in diesem Auftrag.

## Publication-Blockers

- Reproduktionsantworten und zulässigen SDK-Feldvertrag dauerhaft zusammenführen.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Freigabe-Traversierung gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zu Shared-Link-Downloads ohne eingebundene Anzeigepfade recherchieren.
Done-When: Suchabdeckung, relevante Treffer und Abgrenzung von Anzeige- und Traversierungsidentität sind dokumentiert.
