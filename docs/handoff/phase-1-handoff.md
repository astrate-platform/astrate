# Handoff Prompt — Fase 1 OpenAPI YAML

Usa questo prompt per iniziare la prossima sessione di lavoro sulla documentazione astrate.

---

## Prompt

```
Sto lavorando sulla documentazione OpenAPI 3.0 del progetto astrate (~/astrate).

Leggi questi due file prima di iniziare:
1. ~/astrate/docs/handoff/phase-1-memory.md — memoria di cosa è stato fatto e dati chiave
2. g-mind/08 - Progetti/Piano documentazione astrate.md — piano completo aggiornato

Il lavoro consiste nella Fase 1 del piano: generare 5 file YAML OpenAPI 3.0 per le 5 superfici API di astrate, partendo dal codice sorgente.

Inizia dallo step 1.1 (setup scaffolding docs/api/ + schema condivisi) e prosegui con 1.2 (astarte_housekeeping_api.yaml, 4 endpoint, il più semplice).

Per ogni YAML:
- Leggi il file http.go corrispondente per estrarre route, parameters, request/response
- Leggi i golden fixtures per popolare gli examples
- Leggi i Go types per definire gli schemas
- Scrivi la YAML in docs/api/

Verifica con @redocly/cli lint se disponibile, altrimenti validazione manuale della sintassi.

Alla fine di ogni step completato, committa nel repo con messaggio descrittivo del lavoro fatto (es. "docs: add housekeeping OpenAPI spec"). Push sul remote.
```
