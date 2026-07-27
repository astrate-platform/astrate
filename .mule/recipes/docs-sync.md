# Recipe — docs sync

The docs site (`docs/site/`, MkDocs, plus OpenAPI specs in `docs/api/` and Swagger UI) drifts
from the code. This recipe finds the drift. **Prose in `docs/site/` is on the never-touch
list — it is Giulio's voice.** You propose corrections; you do not write them.

The exception, and it is the useful half of this recipe: `docs/api/*.yaml` are *generated
artefacts describing an interface*, not prose. Fixing a wrong path, a wrong status code or a
missing field there is a normal task and can be proposed as one.

## Find drift, cheaply

Pick one surface per run (there are five: appengine, housekeeping, pairing, realm-management,
astrate-native). For that one:

```sh
rg -n 'r\.(Get|Post|Put|Delete|Patch)\(|HandleFunc' internal/ --glob '*<surface>*'
rg -n '^\s{2}/' docs/api/astrate_<surface>_api.yaml   # the documented paths
```

Compare the two lists. Then for three or four endpoints — not all of them — check the
documented status codes and response fields against the handler.

## What to propose

- **Documented but absent**, or **present but undocumented**: a `docs/api/` fix task, naming
  the path and the file.
- **Wrong status code, wrong field name, wrong required-ness**: same, and say what the code
  actually does and where you read it.
- **A `docs/site/` page that contradicts the code**: append to `.mule/for-giulio.md`,
  quoting the sentence and the source line that contradicts it. Never edit the page.
- **A config key documented in `docs/site/configuration-reference.md` that no longer exists**
  (or the reverse) — the same escalation. Find them with:
  `rg -o '\bASTRATE_[A-Z_]+' -N internal/ | sort -u`

## Verify the site still builds

If any task in this family touches `docs/api/` or `docs/mkdocs.yml`, the executing task must
run the docs build (see `docs/Makefile`) and confirm every YAML the Swagger UI references
still loads. A docs change that breaks the build is worse than the drift it fixed.

Five proposals maximum. An accurate spec is a good result; say so and stop.
