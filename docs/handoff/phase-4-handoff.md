# Handoff Prompt — Continue after Phase 4

Use this prompt to start a new session after Phase 4 static site generator setup.

```text
I'm working on documentation for the astrate project in ~/astrate.

Before making any changes, read:
- ~/astrate/docs/AI-DOCUMENTATION-WORKFLOW.md
- ~/astrate/docs/handoff/README.md
- ~/astrate/docs/handoff/phase-4-memory.md
- ~/astrate/docs/handoff/phase-4-handoff.md

Known state:
- Phase 1 (OpenAPI YAML) is complete: 5 specs in ~/astrate/docs/api/.
- Phase 2 (Swagger UI) is complete: ~/astrate/docs/swagger-ui/index.html.
- Phase 3 (Narrative site) is complete: 15 Markdown pages in ~/astrate/docs/site/.
- Phase 4 (Static site generator) is complete: MkDocs + Material theme configured.
- Local preview: cd docs && make serve, then open http://localhost:8000/.
- DESIGN.md, ROADMAP.md, api/, swagger-ui/ are copied into docs/site/ for MkDocs. Run make sync to refresh.
- docs/site/index.md links were updated (removed ../ prefix).

Next phase:
- Phase 5: Add site-dist/ to .gitignore, set up GitHub Pages deployment, or begin code generation milestones (M0-M9 per docs/ROADMAP.md).

Handoff rule:
- At end of phase, update or create docs/handoff/phase-N-memory.md with what was done, verifications, risks, and next steps.
- At end of phase, update or create docs/handoff/phase-N-handoff.md with a mini-prompt for the next session, explicitly referencing relevant handoff files.
```
