# Handoff Prompt — Start Phase 11

Use this prompt to start the next documentation session.

````text
I'm working on the astrate project documentation in ~/astrate.

Before changing files, read:
- ~/astrate/docs/AI-DOCUMENTATION-WORKFLOW.md
- ~/astrate/docs/handoff/README.md
- ~/astrate/docs/handoff/phase-10-memory.md

Known state:
- Phases 1-10 are complete and committed to main.
- Live site: https://astrate-platform.github.io/astrate/
- MkDocs nav has 22 entries (index + 21 pages).
- Phase 10 fixed observability.md metric names and added cross-references to 4 pages.
- Build: cd docs && make sync && make build

Task:
- Pick next documentation work based on gaps. Candidates:
  1. Review and refresh existing narrative pages for accuracy against current source code (data-modeling, mqtt-protocol, interface-schema, payload-formats, triggers, appengine-api, realm-management-api, housekeeping-api, quickstart, migration-from-astarte, contributing).
  2. Create remaining pages from the ROADMAP milestones (check ROADMAP.md for undocumented features).
  3. Add ROADMAP.md to the MkDocs nav so users can find it.

Rules:
- Read source/tests before documenting behavior.
- Keep each page focused and scannable (tables, code blocks, short paragraphs).
- Mark uncertain or preliminary content explicitly.
- Verify build before declaring complete.
- At the end, provide: files changed, verification results, risks, and tell the user which handoff file to read for the next session.
````
