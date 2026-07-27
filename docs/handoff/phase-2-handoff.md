# Handoff Prompt — Fase 2 Swagger UI completata

Usa questo prompt per iniziare la prossima sessione di lavoro sulla documentazione astrate.

---

## Prompt

```
Sto lavorando sulla documentazione del progetto astrate (~/astrate).

Leggi questi file prima di iniziare:
1. ~/astrate/docs/handoff/phase-2-memory.md — memoria di cosa è stato fatto nella Fase 2 (Swagger UI)
2. ~/astrate/docs/handoff/phase-1-memory.md — memoria della Fase 1 (OpenAPI YAML)
3. g-mind/08 - Progetti/Piano documentazione astrate.md — piano completo aggiornato

Situazione attuale:
- Fase 1 (OpenAPI YAML): completata — 5 YAML in docs/api/, tutte lintate con @redocly/cli
- Fase 2 (Swagger UI): completata — docs/swagger-ui/index.html + go:embed in docs/embed.go + route /swagger via internal/swagger/swagger.go
- Fix YAML: astrate_native_api.yaml ha 2 valori quotati (backtick+colon)

Prossima fase: Fase 3 — Sito narrativo (15 pagine Markdown in docs/site/).

Per la Fase 3:
- Creare docs/site/ con 15 pagine come da piano
- Ogni pagina ha una fonte primaria (DESIGN.md, README.md, codice sorgente)
- Contenuto narrativo, non solo API reference
- Iniziare da index.md e 001-intro.md che dipendono solo da README.md

Prima di iniziare, verifica che go build ./cmd/astrate/ compili (la sessione precedente non è riuscita a completare il build per timeout).
```
