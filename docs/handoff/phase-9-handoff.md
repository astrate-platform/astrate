# Handoff Prompt — Start Phase 9

Use this prompt to start the next documentation session.

````text
I'm working on the astrate project documentation in ~/astrate.

Before changing files, read:
- ~/astrate/docs/AI-DOCUMENTATION-WORKFLOW.md
- ~/astrate/docs/handoff/README.md
- ~/astrate/docs/handoff/phase-8-memory.md

Known state:
- Phases 1-8 are complete and committed to main.
- Live site: https://astrate-platform.github.io/astrate/
- MkDocs nav has 20 entries (index + 19 pages).
- 4 new pages added in Phase 8: Configuration Reference, Contributing, Troubleshooting, Migration from Astarte.
- OPERATIONS.md and JSON-PAYLOAD-PROFILE.md exist as standalone docs but are not in the MkDocs nav.
- Build: cd docs && make sync && make build

Task:
- Pick next documentation work based on gaps. Candidates:
  1. Add OPERATIONS.md to MkDocs nav (copy to docs/site/ or create a site/ version).
  2. Add JSON-PAYLOAD-PROFILE.md to MkDocs nav.
  3. Review and refresh existing narrative pages for accuracy against current source.
  4. Create any remaining pages from the ROADMAP milestones.

Rules:
- Read source/tests before documenting behavior.
- Keep each page focused and scannable (tables, code blocks, short paragraphs).
- Mark uncertain or preliminary content explicitly.
- Verify build before declaring complete.
- At the end, provide: files changed, verification results, risks, and tell the user which handoff file to read for the next session.
````
