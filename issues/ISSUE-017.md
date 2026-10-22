# ISSUE-017 — Relocation: Dry-run ignoriert die clientseitige Skip-Policy

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

Root-Cause [S]: Der Dry-run-Zweig verlässt die Iteration vor der bereits vorhandenen clientseitigen Skip-Entscheidung.
Gesprächsreferenz: `E13`; der Review verlangt keine vollständige serverseitige Ausführungssimulation.

## Reach-and-Impact

Reach [S]: `cp` und `mv` mit `--if-exists=skip`, bestehendem Ziel und `--dry-run`.
Impact [O]: Dieselbe unveränderte Fixture erzeugte im Dry-run `planned`, in der normalen Ausführung dagegen `skipped`.
Impact [S]: Bei fehlender Quelle kann nur die Vorschau scheitern, obwohl die echte Ausführung erfolgreich überspringt.

## Evidence

- [S] `cmd/dry_run.go:33-34` beschreibt eine Vorschau beabsichtigter Schreiboperationen.
- [S] `cmd/relocation.go:98-109` plant anhand der Quelle und setzt die Iteration sofort fort.
- [S] `cmd/relocation.go:111-120` erreicht die clientseitige Skip-Prüfung ausschließlich ohne Dry-run.
- [O] Der aktuelle Source-Abgleich zum frisch abgerufenen angegebenen Upstream zeigt keine Produktionsdifferenz.

## Bug-Reproduction

Umgebung: Ubuntu 24.04.5, Linux x64, Go 1.27.3, SDK v6.6.2; isolierte, unveränderte Source- und Zielmetadaten.
Aufrufe: `cp /src.txt /dst.txt --if-exists=skip --dry-run --output=json` und derselbe Aufruf ohne `--dry-run`.
Der Gegenpart liefert für das bestehende Ziel jeweils dieselben Dateimetadaten.
Beobachtet [O]: Beide Commands kehren erfolgreich zurück, der erste mit `planned`, der zweite mit `skipped`.
Keine Copy-Operation wurde benötigt; der temporäre Harness wurde entfernt.

## Prior-Art

Coverage: Index und Akten `ISSUE-001` bis `ISSUE-005` enthalten keine gleiche Skip-Entscheidungsursache.
`ISSUE-004` betrifft rekursive Upload-Traversierung, nicht die Vorschau von Copy und Move.
Gaps: Externe Issues, PRs, Diskussionen und Releases sind noch nicht recherchiert.
Contribution fit: Eine lokale gemeinsame Vorentscheidung genügt; Serverzustände müssen nicht vollständig simuliert werden.

## Proposed-Change

Die bestehende clientseitige Skip-Entscheidung vor den Split zwischen Vorschau und Ausführung legen.
Übersprungene Ziele nicht als beabsichtigte Schreiboperationen ankündigen.

## Scope-and-Constraints

- Preserve: Dry-run führt keine Remote-Schreiboperation aus.
- Preserve: Bestehende Skip-Semantik bleibt maßgeblich, auch bei nicht vorhandener Quelle.
- Exclude: Keine Garantie für serverseitige Autorename-Prognose, Berechtigungsprüfung oder Snapshot-Atomarität hinzufügen.
- Cost: Geplante No-ops müssen in Text und JSON eindeutig repräsentiert werden.

## Verification

Abgeschlossen [O]: Der Modusvergleich reproduziert den Vorschauwiderspruch bei gleichem Remote-Zustand.
Erforderlich: Bestehendes Ziel jeweils mit vorhandener und fehlender Quelle vergleichen.
Erwartet: Vorschau und Ausführung treffen dieselbe Skip-Entscheidung ohne zusätzliche Schreiboperation.
Test decision: none; kein Source- oder Testeingriff in diesem Auftrag.

## Publication-Blockers

- Vorschauvertrag für übersprungene No-ops und vorhandene JSON-Statuswerte konkret abgleichen.
- Prior Art und aktuelle Upstream-Beitragsvorgaben prüfen.
- Arbeitsmodus, Veröffentlichungsziel und exakte Freigabe fehlen.

## Next-Action

Summary: Skip-Vorschau gegen Prior Art prüfen
Action: Bestehende Upstream-Arbeiten zur Parität clientseitiger Skip-Entscheidungen in Copy-/Move-Vorschauen recherchieren.
Done-When: Suchabdeckung, relevante Treffer und die begrenzte Vorschaukorrektur sind dokumentiert.
