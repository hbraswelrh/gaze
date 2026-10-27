## Why

`gaze init` writes a `<!-- scaffolded by gaze <version> -->` provenance
marker into each scaffolded file. Re-running `gaze init` with a different
version (e.g. `dev` to `v1.8.0`) left a second marker instead of replacing
the existing one, because `insertMarkerAfterFrontmatter` inserted
unconditionally (issue #279).

Additionally, the MegaLinter auto-commit on this PR's branch reformatted
only the `internal/` copies of three embedded assets (blank-line insertions),
breaking the byte-identical invariant between `.opencode/` sources and
`internal/scaffold/assets/` mirrors and failing
`TestEmbeddedAssetsMatchSource`. The root cause is structural: one side of a
byte-identical pair was linted while the other side was excluded, so any
linter reformat could silently drift the pair and re-break CI.

This change is a retroactive record for PR #280 (implemented and reviewed
before the spec was written). Closes #279.

## What Changes

### Idempotent Marker Insertion

`insertMarkerAfterFrontmatter` now runs `removeExistingMarker` (a new
unexported helper) as its first step, stripping all existing
`<!-- scaffolded by gaze ... -->` lines before insertion. Any number of
re-scaffolds with any version leaves exactly one marker carrying the
newest version.

### Stale Baked-in Markers Removed

Seven scaffold assets contained a hardcoded
`<!-- scaffolded by gaze v1.8.0 -->` line — the same stale-provenance data
at the root of the duplicate-marker bug. Removed from
`internal/scaffold/assets/`, mirrored to `.opencode/`, and applied to the
aireport embedded prompt copy.

### Embedded Assets Excluded from markdownlint

`internal/scaffold/assets/` and `internal/aireport/assets/` are added to
`MARKDOWN_MARKDOWNLINT_FILTER_REGEX_EXCLUDE` in `.mega-linter.yml`.
Embedded copies are byte-identical mirrors of the `.opencode/` source of
truth and must never be reformatted independently. The 43 blank lines the
linter added to four asset files were reverted.

### No Other Behavioral Change

Marker placement (after the frontmatter closing delimiter, appended when
there is no frontmatter), file ownership rules (tool-owned vs user-owned),
and scaffold output are otherwise unchanged.

## Capabilities

### New Capabilities

- `scaffold-marker-idempotency` — re-scaffolding replaces the version
  marker in place; embedded assets are lint-invariant mirrors

### Modified Capabilities

- None

### Removed Capabilities

- None

## Impact

- **Modified files**: `internal/scaffold/scaffold.go` (new
  `removeExistingMarker`, one call site), `internal/scaffold/scaffold_test.go`
  (4 new test functions, 19 subtests), 7 `.opencode/` +
  `internal/scaffold/assets/` asset pairs (stale marker removal),
  `internal/aireport/assets/agents/gaze-reporter.md` (stale marker removal),
  `.mega-linter.yml` (markdownlint exclusion), `internal/crap/report.go`
  (gofmt tag alignment only, semantic no-op)
- **Risk**: Low. Marker stripping is line-based prefix/suffix matching (no
  regex, no backtracking surface); non-marker HTML comments are preserved.
  Excluding embedded assets from markdownlint reduces lint coverage of
  generated mirrors only.
- **Test impact**: New unit, table-driven, and end-to-end tests.
  `TestEmbeddedAssetsMatchSource` and `TestEmbeddedPromptMatchesScaffold`
  pass again.

## Constitution Alignment

Assessed against the Gaze project constitution
(`.specify/memory/constitution.md`).

### I. Accuracy

**Assessment**: PASS

Provenance markers now report the actual build version instead of a
stale baked-in one. Regression tests pin the single-marker invariant
across force and non-force re-scaffold scenarios.

### II. Minimal Assumptions

**Assessment**: PASS

No new host assumptions. Marker detection matches the exact marker
syntax Gaze itself emits (prefix/suffix on trimmed lines).

### III. Actionable Output

**Assessment**: N/A

No report output changes.

### IV. Testability

**Assessment**: PASS

`removeExistingMarker` is directly unit-testable (5-case table-driven
test). End-to-end re-scaffold tests use `t.TempDir`, are isolated, and
pass under `-race`.
