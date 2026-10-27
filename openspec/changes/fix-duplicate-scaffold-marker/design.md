## Context

`gaze init` scaffolds nine files into a target project. Each file receives
a `<!-- scaffolded by gaze <version> -->` marker inserted after the YAML
frontmatter closing delimiter (or appended when there is no frontmatter)
by `insertMarkerAfterFrontmatter` (`internal/scaffold/scaffold.go`).
Re-runs were expected to refresh the marker in place, but insertion was
unconditional, so a second `gaze init` with a different version produced
two markers (issue #279).

While PR #280 was under review, the org MegaLinter auto-commit added blank
lines to four `internal/` asset files because `.opencode/` is excluded in
`.mega-linter.yml` while `internal/` is linted.
`TestEmbeddedAssetsMatchSource` (which enforces the byte-identical
invariant documented in AGENTS.md) failed on the PR head.

## Goals / Non-Goals

### Goals
- Idempotent marker insertion: exactly one marker, always the newest version
- Remove the stale hardcoded `v1.8.0` marker from all embedded assets
- Make the byte-identical invariant structurally safe under MegaLinter
  auto-commits
- Regression coverage for both force and non-force re-scaffold paths

### Non-Goals
- Changing marker placement rules (after frontmatter / append)
- Changing file ownership rules (tool-owned vs user-owned)
- Replacing the sync-invariant tests with a different mechanism

## Decisions

### D1: Strip-then-insert inside `insertMarkerAfterFrontmatter`

A new unexported helper `removeExistingMarker(content []byte) []byte`
strips every line whose trimmed form starts with `<!-- scaffolded by gaze `
and ends with `-->`, and is called as the first step of
`insertMarkerAfterFrontmatter`.

Alternatives considered:

- Regex-based removal: more expressive, but adds catastrophic-backtracking
  surface to a fixed-syntax match. Prefix/suffix matching is O(n) and
  sufficient.
- Version-matched replacement (replace only the marker of the old
  version): leaves duplicates if two different stale versions are already
  present. Stripping all marker lines is strictly more robust.

Line-based trimming (rather than byte-exact matching) makes indented and
CRLF-terminated markers removable. Malformed lines (prefix present but no
closing `-->`) are conservatively retained.

### D2: Exclude embedded assets from markdownlint

`internal/scaffold/assets/` and `internal/aireport/assets/` were added to
`MARKDOWN_MARKDOWNLINT_FILTER_REGEX_EXCLUDE`, and the 43 blank lines the
linter added to four asset files were reverted.

The `.opencode/` files are the source of truth and the `internal/` copies
are byte-identical mirrors enforced by `TestEmbeddedAssetsMatchSource` /
`TestEmbeddedPromptMatchesScaffold`. Linting only one side of the pair
means any future markdownlint rule change silently drifts the pair and
re-breaks CI — exactly the failure that motivated this decision.
Propagating the linter's output into `.opencode/` (the alternative) was
rejected: it bakes mechanical formatting into the source of truth and
leaves the hazard in place.

### D3: Non-force re-scaffold covered by a dedicated e2e case

`TestRun_ReScaffoldNonForceSingleMarker` covers the issue's actual repro:
a second `gaze init` without `--force` updates tool-owned files on diff
(replacing the marker) and leaves user-owned files untouched.