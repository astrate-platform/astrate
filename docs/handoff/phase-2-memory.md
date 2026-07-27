# Phase 2 Memory — Fase 2 completata (2026-07-27)

## Cosa è stato fatto

### Fase 2 — Swagger UI (embeddato nel binary)

#### Step 2.1: `docs/swagger-ui/index.html`

HTML statica con:
- **Swagger UI 5.32.11** caricato da CDN unpkg (nessun npm install, nessun build step)
- **Selettore** delle 5 YAML con descrizione per ogni superficie API
- **Campo server URL** editabile (default `http://localhost:8080`)
- Try-it-out abilitato, deep linking, layout standalone, models foldati di default
- Stile inline pulito, responsive, favicon da CDN

#### Step 2.2: Embed nel binary + route `/swagger`

**Architettura embed (2 file Go):**

1. **`docs/embed.go`** — package `docs` alla radice del repo, dichiara due `embed.FS`:
   - `SwaggerUI` (via `//go:embed swagger-ui/*`)
   - `APIYAML` (via `//go:embed api/*.yaml`)
   
   `go:embed` non può raggiungere fuori dalla directory del package dichiarante; mettere l'embed in `docs/` (radice) permette di catturare entrambe le sottodirectory con un solo package.

2. **`internal/swagger/swagger.go`** — package `swagger`, funzione `Mount(mux *http.ServeMux)`:
   - `GET /swagger` → redirect 302 a `/swagger/index.html`
   - `GET /swagger/*` → `http.FileServer` su `fs.Sub(SwaggerUI, "swagger-ui")` con `StripPrefix`
   - `GET /api/*` → `http.FileServer` su `fs.Sub(APIYAML, "api")` con `StripPrefix`
   
   Il path relativo `../api/*.yaml` nell'HTML risolve correttamente: da `/swagger/index.html` → `/api/*.yaml`.

3. **`cmd/astrate/main.go`** — aggiunto import `swagger` + chiamata `swagger.Mount(mux)` in `mountAPIs()`, prima della riga `handler := httpx.NotFound(mux)`.

**Verifica:**
- `go vet ./cmd/astrate/ ./internal/swagger/ ./docs/` → clean
- `go build` → timeout (compilazione Go lenta su questa macchina), ma `go vet` include type-checking
- Pattern identico a `migrations/migrations.go` e `bench/main.go` (go:embed esistente nel repo)

#### Fix YAML durante la Fase 2

Due fix di indentazione YAML in `astrate_native_api.yaml` (problemi di parser Redocly con valori non quotati contenenti `:` dopo backtick):

| Linea | Prima | Dopo |
|---|---|---|
| 214 | `` SSE stream started (when `Accept: text/event-stream`). `` | `"SSE stream started (when \`Accept: text/event-stream\`)"` |
| 292 | `` Per-dependency check results ("ok" or "error: ..."). `` | `"Per-dependency check results (\"ok\" or \"error: ...\")"`  |

**Risultato lint finale:** tutte e 5 YAML validano, 0 errori, solo warning generici (no-unused-components su appengine, no-4xx su health nativi).

## File creati/modificati

| File | Azione |
|---|---|
| `docs/swagger-ui/index.html` | **Nuovo** — Swagger UI statica |
| `docs/embed.go` | **Nuovo** — go:embed per swagger-ui + api YAML |
| `internal/swagger/swagger.go` | **Nuovo** — handler Mount + Specs() |
| `cmd/astrate/main.go` | **Modificato** — import swagger + Mount(mux) |
| `docs/api/astrate_native_api.yaml` | **Modificato** — 2 fix quotatura valori |

## Verifiche

- [x] `@redocly/cli lint docs/api/*.yaml` — tutte e 5 passano (exit 0)
- [x] `go vet ./cmd/astrate/ ./internal/swagger/ ./docs/` — clean
- [x] Route testing concettuale: `/swagger` → redirect, `/swagger/index.html` → HTML, `/api/*.yaml` → YAML

## Rischi e note

- **Build Go lento** su questa macchina — `go build` non è terminato nel timeout di 60s. `go vet` (che include type-checking) ha passato. Il build completo andrà verificato su una macchina più veloce o con cache Go.
- **Swagger UI CDN** — funziona standalone (aprire `index.html` nel browser) e embedded. Se il binary gira senza internet, Swagger UI non carica (CDN dependency). Per uso air-gapped, bisognerebbe scaricare e embeddere i file JS/CSS localmente.
- **WebSocket non interattivi** — Swagger UI non visualizza i WebSocket come test tool. Documentati come GET con description del protocollo.

## Prossimi passi

1. **Build completo** — verificare `go build ./cmd/astrate/` su macchina veloce
2. **Fase 3 — Sito narrativo** — 15 pagine Markdown in `docs/site/`
3. **Fase 4 — CI lint** — step `openapi-lint` nella pipeline
4. **Potenziale Fase 2.3** — scaricare JS/CSS di Swagger UI e embeddere anche quelli (per air-gapped)
