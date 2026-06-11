---
name: feedback-no-ai-attribution
description: Never add Co-Authored-By/AI attribution to commits or files in this project — user explicitly forbids it
metadata:
  type: feedback
---

Never add a `Co-Authored-By: Claude ...` trailer, "Generated with Claude Code" line, or any other AI attribution to commit messages, PR bodies, README, or any file in this project.

**Why:** The user explicitly forbids mentioning Claude as a co-author of this project (stated in the S1/M0 Phase 3 instructions, 2026-06-11). This is a direct user instruction that overrides the default harness guidance to append such trailers.

**How to apply:** Every `git commit` and any generated doc/file in [[project-astrate-phases]] — write clean conventional commit messages with no attribution lines whatsoever.
