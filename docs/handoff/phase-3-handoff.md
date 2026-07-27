# Handoff Prompt — Continue after Phase 3

Use this mini-prompt to start a new session after Phase 3 narrative site.

```text
I'm working on documentation for the astrate project in ~/astrate.

Before making any changes, read:
- ~/astrate/docs/AI-DOCUMENTATION-WORKFLOW.md
- ~/astrate/docs/handoff/README.md
- ~/astrate/docs/handoff/phase-3-memory.md
- ~/astrate/docs/handoff/phase-3-handoff.md
- g-mind/08 - Progetti/Piano documentazione astrate.md

Known state:
- Phase 1 (OpenAPI YAML) is complete: 5 specs in ~/astrate/docs/api/.
- Phase 2 (Swagger UI) is complete: ~/astrate/docs/swagger-ui/index.html.
- Phase 3 (Narrative site) is complete: 15 Markdown pages in ~/astrate/docs/site/.
- To view Swagger UI locally, serve from ~/astrate/docs, then open http://localhost:9090/swagger-ui/.
- Do not use local port 8080 while the Lenovo Legion Go SSH tunnel is active.
- npx --yes @redocly/cli lint docs/api/*.yaml passes with exit code 0.
- Phase 2.2 (go:embed + /swagger route) was not done.

Next phase:
- Phase 4: wire docs/site/ into a static site generator (optional), or begin Phase 3 code generation milestones (M0-M9 per docs/ROADMAP.md).

Handoff rule:
- At end of phase, update or create docs/handoff/phase-N-memory.md with what was done, verifications, risks, and next steps.
- At end of phase, update or create docs/handoff/phase-N-handoff.md with a mini-prompt for the next session, explicitly referencing relevant handoff files.
```
