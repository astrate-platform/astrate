# Handoff Prompt — Continue after Phase 7

Use this prompt to start a new documentation session.

```text
I'm working on documentation for the astrate project in ~/astrate.

Before making any changes, read:
- ~/astrate/docs/AI-DOCUMENTATION-WORKFLOW.md
- ~/astrate/docs/handoff/README.md
- ~/astrate/docs/handoff/phase-7-memory.md

Known state:
- Phases 1-7 are complete.
- Phase 1 (OpenAPI YAML): 5 specs in ~/astrate/docs/api/.
- Phase 2 (Swagger UI): ~/astrate/docs/swagger-ui/index.html.
- Phase 3 (Narrative site): 15 Markdown pages in ~/astrate/docs/site/.
- Phase 4 (Static site generator): MkDocs + Material theme configured.
- Phase 5 (GitHub Pages deployment): .github/workflows/docs.yml created.
- Phase 6 (Verification): Local build verified — `make sync && make build` passes.
- Phase 7 (Push to main): Docs committed and pushed (commit cc34e6a). GitHub Pages enabled in repo Settings.
- Local preview: cd docs && make serve, then open http://localhost:8000/.
- Live site: https://astrate-platform.github.io/astrate/ (after next push to main).

Documentation gaps to fill (see phase-7-memory.md for details):
- Configuration Reference: TOML options, ASTRATE_* env vars, defaults.
- Operations Runbook: production TLS, CA re-keying, backup, upgrade.
- JSON Payload Profile: normative spec for AtomVM/bare-metal JSON clients.
- Contributing / Dev Guide: test tiers, dev setup, how to contribute.
- Troubleshooting: common failure modes and fixes.
- Migration from Astarte: moving an existing deployment.

Handoff rule:
- At end of phase, update or create docs/handoff/phase-N-memory.md with what was done, verifications, risks, and next steps.
- At end of phase, update or create docs/handoff/phase-N-handoff.md with a mini-prompt for the next session, explicitly referencing relevant handoff files.
```
