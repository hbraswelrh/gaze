## Why

Gaze currently emits terminal text and JSON from `gaze analyze`. Text is concise but awkward to share or inspect interactively, while JSON is intended for machines. A self-contained HTML report gives users a portable, human-readable artifact for local review and CI uploads without requiring a web server, network access, or additional dependencies.

This change implements the `gaze analyze` portion of tracking issue #8. HTML output for `gaze quality`, `gaze crap`, and `gaze report` remains separately scoped to issue #261.

## What Changes

- Add a shared HTML rendering path under `internal/report` using Go's `html/template` package.
- Add `html` to the supported output formats for `gaze analyze`.
- Render a complete, self-contained HTML document with inline CSS and native `<details>` and `<summary>` navigation.
- Represent every analyzed function and detected side effect while safely escaping source-derived content.
- Add deterministic formatter and CLI integration tests covering rendering, escaping, self-containment, and format selection.
- Update user-facing CLI documentation for the new format.

## Capabilities

### New Capabilities
- `analyze-html-output`: Produces a complete, self-contained HTML representation of `gaze analyze` results.

### Modified Capabilities
- `analyze-output-formats`: Extends the accepted `gaze analyze --format` values from `text` and `json` to include `html` without changing existing output contracts.

### Removed Capabilities
- None.

## Impact

- `internal/report/`: new HTML formatter and tests.
- `internal/cliutil/`: format validation recognizes `html`.
- `cmd/gaze/`: analyze command dispatches HTML output and gains CLI integration coverage.
- CLI documentation: supported formats and examples include HTML.
- No changes to side-effect analysis, classification, JSON schemas, text rendering, quality reports, CRAP reports, or AI reports.
- No new external dependencies, runtime services, or remote assets.
- Because this is a user-facing CLI capability, the repository's website documentation gate requires a website tracking issue before the implementing PR is merged.

## Constitution Alignment

Assessed against the Gaze project constitution.

### I. Accuracy

**Assessment**: PASS

The generated report is a self-contained artifact that can be exchanged between local development, CI, and reviewers without runtime coupling or an external rendering service.

### II. Minimal Assumptions

**Assessment**: PASS

HTML is an optional output format implemented with the Go standard library. Existing text and JSON modes remain independently usable, and Gaze gains no mandatory service or third-party dependency.

### III. Actionable Output

**Assessment**: PASS

JSON remains the stable machine-parseable output with existing provenance metadata. HTML adds a deterministic human-readable view of the same analysis results and does not replace or weaken the JSON contract.

### IV. Testability

**Assessment**: PASS

The formatter accepts an `io.Writer` and analysis results, allowing isolated deterministic tests. Tests will verify complete rendering, escaping, self-containment, CLI selection, and non-regression of existing formats without external services.
<!-- scaffolded by uf vdev -->
