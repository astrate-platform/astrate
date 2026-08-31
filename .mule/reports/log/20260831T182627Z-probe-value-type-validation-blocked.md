slug: probe-value-type-validation
verdict: blocked
at:  db8048e
ran: 2026-08-31T18:26:27Z on DietPi in 48s


> build · big-pickle

→ Read .mule/MULE.md
→ Read .mule/task.md
✱ Grep "validate_value_type|ValidateValueType|value.?type.?valid" in . · 23 matches
This is a `[readonly]` probe: compare Astrate's per-mapping value-type validation against upstream Astarte's restored `validate_value_type` from v1.4.0-rc.3. I must report, not patch.

Let me find the relevant code in both codebases.
! permission requested: external_directory (/root/*); auto-rejecting
✗ Grep "validate_value_type|ValidateValueType|value.?type.?valid" failed in /root
Error: The user rejected permission to use this specific tool call.
