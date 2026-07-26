# Phase 1 Memory — Fase 1 completata (2026-07-27)

## Cosa è stato fatto

### Pre-work (sessione precedente)
1. **Esplorazione completa del repo astrate** — struttura, route registration, golden fixtures, doc esistenti.
2. **Aggiornamento del piano** in `g-mind/08 - Progetti/Piano documentazione astrate.md` con dati reali.
3. **Setup directory** `docs/handoff/` per la memoria iterativa.

### Fase 1 — OpenAPI YAML (questa sessione)

Tutte le 5 YAML sono state scritte in `docs/api/`:

| File | Operazioni | Righi | Note |
|---|---|---|---|
| `astarte_housekeeping_api.yaml` | 4 | 333 | Realm CRUD, `a_ha` JWT |
| `astarte_pairing_api.yaml` | 5 | 558 | Flows A/B/C, `a_pa` JWT + bearer secret, 15 fixture usati |
| `astarte_realm_management_api.yaml` | 19 | 813 | Interfaces CRUD, triggers CRUD, config, policies, device delete, `a_rma` JWT |
| `astarte_appengine_api.yaml` | 15 | 892 | Device list/detail/patch, data CRUD, groups CRUD, `a_aea` JWT, pagination links |
| `astrate_native_api.yaml` | 8 | 336 | Health/readiness/metrics (no auth), 3 compat health, 2 WebSocket (native stream + Phoenix V2) |

**Totale: 51 operazioni** (49 REST + 2 WebSocket).

#### Approccio adottato
- **Ogni YAML è self-contained** — definisce i propri `components.schemas` e `components.responses` senza `$ref` esterni. Questo massima compatibilità con tool OpenAPI (Swagger UI, Redocly, code generatori).
- **Schema condivisi replicati** — `ErrorDetail`, `ErrorFields`, `HealthStatus` sono definiti in ogni YAML dove servono. Non c'è un file comune separato (OpenAPI 3.0 non supporta `$ref` a file esterni senza tooling specifico).
- **Examples dai golden fixtures** — ogni response 200 ha un example derivato dai testdata reali. I fixture pairing (15 file) sono stati usati direttamente. I fixture envelope (`data_object.json`, `error_bad_request.json`, ecc.) hanno informato le shape di tutti gli error envelope.
- **Auth schemes** — 4 definiti: `a_ha` (apiKey in header), `a_pa` (apiKey), `a_rma` (apiKey), `a_aea` (apiKey), più `bearerSecret` (http bearer) nel pairing.
- **WebSocket documentati come REST operations** — i 2 WS endpoint hanno `operationId` e descrizione del protocollo, anche se Swagger UI non li visualizza interattivamente.

#### Go struct tradotte in schema
- `realmBody` → `Realm`, `RealmCreate` (housekeeping)
- `registerRequest/Response` → `RegisterRequest`, `CredentialsSecret` (pairing)
- `credentialsRequest/Response` → `CredentialsRequest`, `ClientCert` (pairing)
- `infoResponse` → `PairingInfo` (pairing)
- `verifyValidResponse/Invalid` → `VerifyValid`, `VerifyInvalid` (pairing)
- `devicePatchBody` → `DevicePatch` (appengine)
- `groupBody` → `GroupCreate` (appengine)
- `authConfig` → `AuthConfig` (realm management)
- `wireEvent` → documentato inline nella description del WS endpoint (native)

#### Non documentato (out of scope Fase 1)
- **Request body per `installInterface`** — usato `type: object` generico. Lo schema JSON interface è complesso (mappings, types, retention policies) e va documentato nella Fase 3 (sito narrativo) oppure con un `$ref` a un JSON Schema separato.
- **Request body per `createTrigger`** — stessa situazione: `type: object` generico.
- **Request body per `createPolicy`** — stessa situazione.

## Dati chiave per la fase successiva

### File creati
```
docs/api/
├── astarte_housekeeping_api.yaml    (4 op, 333 righe)
├── astarte_pairing_api.yaml         (5 op, 558 righe)
├── astarte_realm_management_api.yaml (19 op, 813 righe)
├── astarte_appengine_api.yaml       (15 op, 892 righe)
└── astrate_native_api.yaml          (8 op, 336 righe)
```

### Cosa deve fare la prossima sessione

1. Leggere questo file (`docs/handoff/phase-1-memory.md`).
2. **Fase 2 — Swagger UI**: creare `docs/swagger-ui/index.html` con CDN refs e configurazione per le 5 YAML.
3. **Opzionale: Fase 2.2** — embed `go:embed` + route `/swagger` nel binary.
4. **Fase 3 — Sito narrativo**: 15 pagine Markdown in `docs/site/`.
5. **Step 1.7** — lint con `@redocly/cli lint docs/api/*.yaml` (da installare prima).

### Note per il lint
- `redocly/cli` va installato via npm: `npm install -g @redocly/cli`
- Le YAML usano `oneOf` in due punti (pairing verify responses, appengine device list) — redocly li gestisce.
- I WebSocket non hanno un visualizzatore interattivo in Swagger UI — documentati come GET endpoints con描述 completa del protocollo.
