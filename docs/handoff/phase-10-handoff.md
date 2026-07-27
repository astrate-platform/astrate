# Handoff Prompt — Start Phase 10

Use this prompt to start the next documentation session.

````text
I'm working on the astrate project documentation in ~/astrate.

Before changing files, read:
- ~/astrate/docs/AI-DOCUMENTATION-WORKFLOW.md
- ~/astrate/docs/handoff/README.md
- ~/astrate/docs/handoff/phase-9-memory.md

Known state:
- Phases 1-9 are complete and committed to main.
- Live site: https://astrate-platform.github.io/astrate/
- MkDocs nav has 22 entries (index + 21 pages).
- 2 new pages added in Phase 9: Operations, JSON Payload Profile.
- Both are adapted from standalone docs/ root files with fixed internal links.
- Build: cd docs && make sync && make build

Task:
- Pick next documentation work based on gaps. Candidates:
  1. Review and refresh existing narrative pages for accuracy against current source code.
  2. Create remaining pages from the ROADMAP milestones (check ROADMAP.md for undocumented features).
  3. Add cross-references between the new operations page and related pages (deployment, configuration-reference, troubleshooting).

Rules:
- Read source/tests before documenting behavior.
- Keep each page focused and scannable (tables, code blocks, short paragraphs).
- Mark uncertain or preliminary content explicitly.
- Verify build before declaring complete.
- At the end, provide: files changed, verification results, risks, and tell the user which handoff file to read for the next session.
````
