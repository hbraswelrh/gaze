# Website Documentation Issue

> **Status: NOT YET CREATED** -- This issue must be created in
> `unbound-force/website` before the implementing PR is merged.
> The repository's Website Documentation Gate (AGENTS.md) requires it.

## Command

```bash
gh issue create --repo unbound-force/website \
  --title "docs: add gaze analyze --format=html output documentation" \
  --body-file openspec/changes/add-analyze-html-output/website-issue.md
```

## Issue Body

`gaze analyze` now supports `--format=html`, producing a self-contained
single-file HTML report with inline CSS, native `<details>`/`<summary>`
collapsible sections, and no JavaScript or external resources.

### What changed

- New output format: `gaze analyze --format=html <package>`
- The report renders every analyzed function and side effect with
  classification labels, confidence percentages, and detail metadata.
- Output is deterministic, offline-safe, and contextually escaped via
  Go's `html/template`.
- HTML is accepted only by `gaze analyze`; other commands (`quality`,
  `crap`, `report`) continue to support `text` and `json` only.

### Why it matters

Users can now generate portable, human-readable analysis reports for
local review, CI artifact uploads, and team sharing without requiring
a web server or network access.

### Affected website pages

- **CLI Reference -- `gaze analyze`**: Add `html` to the `--format`
  flag values, describe behavior, and include a usage example
  (`gaze analyze --format=html ./pkg > report.html`).
- **Output Formats overview**: Note that `analyze` additionally
  supports `html` alongside the universal `text`/`json` formats.
- **Getting Started / Quickstart** (if applicable): Consider
  mentioning HTML output as an option for first-time users exploring
  analysis results.
