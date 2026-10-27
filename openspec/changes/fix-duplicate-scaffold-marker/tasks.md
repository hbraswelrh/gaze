<!--
  Retroactive record for PR #280 (issue #279). Implementation predates
  this artifact; all tasks are complete as written.
-->

## 1. Idempotent Marker Insertion

- [x] 1.1 Add `removeExistingMarker(content []byte) []byte` to `internal/scaffold/scaffold.go` — strip every line whose trimmed form starts with `<!-- scaffolded by gaze ` and ends with `-->`. O(n) line scan, no regex.
- [x] 1.2 Call `removeExistingMarker` as the first step of `insertMarkerAfterFrontmatter`; update the GoDoc comment to document strip-before-insert.

## 2. Stale Marker Removal

- [x] 2.1 Remove the hardcoded `<!-- scaffolded by gaze v1.8.0 -->` line from the seven `internal/scaffold/assets/` assets and mirror to `.opencode/` (TestEmbeddedAssetsMatchSource).
- [x] 2.2 Apply the same removal to the `internal/aireport/assets/agents/gaze-reporter.md` embedded prompt copy (TestEmbeddedPromptMatchesScaffold).

## 3. Regression Tests

- [x] 3.1 `TestRemoveExistingMarker` — table-driven: single marker, multiple markers, no marker, marker at end of file, non-marker HTML comments preserved.
- [x] 3.2 `TestInsertMarkerAfterFrontmatter_ReplacesExistingMarker` — table-driven: replacement after frontmatter, replacement without frontmatter, duplicate removal, no existing marker.
- [x] 3.3 `TestRun_ReScaffoldSingleMarker` — end-to-end: `gaze init` twice (`dev` then `v1.8.0`) with `--force`; assert exactly one marker per file.
- [x] 3.4 `TestRun_ReScaffoldNonForceSingleMarker` — end-to-end non-force re-init: tool-owned files updated with one marker, user-owned files untouched.

## 4. Embedded Asset Lint Invariance

- [x] 4.1 Add `internal/scaffold/assets/` and `internal/aireport/assets/` to `MARKDOWN_MARKDOWNLINT_FILTER_REGEX_EXCLUDE` in `.mega-linter.yml`.
- [x] 4.2 Revert the 43 blank lines the MegaLinter auto-commit added to the four affected asset files, restoring the byte-identical pairs.

## 5. Verification

- [x] 5.1 `go build ./cmd/gaze` — zero errors.
- [x] 5.2 `go test -race -count=1 -short ./...` — all tests pass, including `TestEmbeddedAssetsMatchSource` and `TestEmbeddedPromptMatchesScaffold`.
- [x] 5.3 `go vet ./...` — clean.
- [x] 5.4 `golangci-lint run` — zero findings.

## 6. Documentation

- [x] 6.1 Update `AGENTS.md` Recent Changes with a `fix-duplicate-scaffold-marker` entry.
- [x] 6.2 No website issue needed — bug fix to existing scaffold behavior; no new command, flag, or feature.