# Handoff Prompt — Continue after Phase 5

Use this prompt to start a new session after Phase 5 GitHub Pages deployment setup.

```text
I'm working on documentation for the astrate project in ~/astrate.

Before making any changes, read:
- ~/astrate/docs/AI-DOCUMENTATION-WORKFLOW.md
- ~/astrate/docs/handoff/README.md
- ~/astrate/docs/handoff/phase-5-memory.md
- ~/astrate/docs/handoff/phase-5-handoff.md

Known state:
- Phase 1 (OpenAPI YAML) is complete: 5 specs in ~/astrate/docs/api/.
- Phase 2 (Swagger UI) is complete: ~/astrate/docs/swagger-ui/index.html.
- Phase 3 (Narrative site) is complete: 15 Markdown pages in ~/astrate/docs/site/.
- Phase 4 (Static site generator) is complete: MkDocs + Material theme configured.
- Phase 5 (GitHub Pages deployment) is complete: .github/workflows/docs.yml created.
- Local preview: cd docs && make serve, then open http://localhost:8000/.
- GitHub Pages deployment: push to main or manual trigger deploys to GitHub Pages.

Next phase:
- Phase 6: Verify GitHub Pages deployment works, or proceed with code generation milestones (M0-M9 per docs/ROADMAP.md).

Handoff rule:
- At end of phase, update or create docs/handoff/phase-N-memory.md with what was done, verifications, risks, and next steps.
- At end of phase, update or create docs/handoff/phase-N-handoff.md with a mini-prompt for the next session, explicitly referencing relevant handoff files.
```
