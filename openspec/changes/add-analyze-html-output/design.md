## Context

`gaze analyze` currently validates `--format` through `internal/cliutil` and dispatches results to `internal/report.WriteTextOptions` or `internal/report.WriteJSON`. The report package has no HTML implementation despite HTML being a longstanding tracked capability in issue #8. Issue #260 deliberately limits the first delivery to analyze output; issue #261 covers quality, CRAP, and combined reports.

The proposal establishes that the output must be portable, self-contained, safely escaped, deterministic, and implemented without new dependencies. This preserves the proposal's constitution alignment: the HTML file is an independently exchangeable artifact, existing formats remain composable, JSON remains machine-parseable, and rendering remains isolated and testable.

## Goals / Non-Goals

### Goals
- Add `gaze analyze --format=html` without changing text or JSON behavior.
- Render all analyzed functions and side effects in a complete single-file document.
- Use contextual escaping for all source-derived and analyzer-derived content.
- Provide native collapsible sections with useful inline styling and no runtime dependencies.
- Establish a reusable HTML template and testing pattern for issue #261.
- Keep output deterministic for reliable tests and artifact comparison.

### Non-Goals
- Add HTML support to `gaze quality`, `gaze crap`, `gaze report`, or `gaze self-check`.
- Add a web server, JavaScript application, asset pipeline, CSS framework, or third-party dependency.
- Change analysis, classification, scoring, taxonomy, text output, or JSON schemas.
- Define a universal HTML schema shared by report types that have materially different data models.

## Decisions

### D1: Add a report-layer HTML formatter

Add `report.WriteHTML(w io.Writer, results []taxonomy.AnalysisResult, version string) error`. The signature mirrors `WriteJSON`, including the `dev` version fallback, while retaining the existing `io.Writer`-based testability pattern.

The formatter will convert taxonomy values into a presentation view model before template execution. The view model may add deterministic counts and serialized detail text, but it will not alter analysis semantics. Template parsing or execution failures will be returned with operation-specific wrapped errors.

### D2: Embed a dedicated analyze template

Store the document template as a narrowly targeted embedded asset in `internal/report` and load it with `embed.FS`, as required by convention AP-004. The template contains the complete document structure and inline CSS. It uses native `<details>` and `<summary>` elements and contains no JavaScript or external URLs.

Keeping the template separate from Go source makes the document structure reviewable while embedding it preserves standalone binary behavior. Issue #261 may reuse layout conventions and CSS, but this change will not introduce speculative cross-report abstractions.

### D3: Rely on html/template contextual escaping

Parse and execute the asset with `html/template`. Dynamic report values remain ordinary strings or scalar values; the implementation will not cast untrusted values to `template.HTML` or other trusted content types.

Opaque `SideEffect.Detail` metadata will be deterministically JSON-encoded into a string before template execution and displayed as escaped text. Go's JSON encoder sorts string map keys, preventing map iteration from making output nondeterministic. Serialization errors will fail rendering rather than dropping metadata silently.

### D4: Extend format validation without enabling unsupported commands

Evolve `cliutil.ValidateFormat` to accept optional additional allowed formats. Existing callers continue calling it with only the format and therefore retain the exact `text`/`json` allowlist. `runAnalyze` supplies `html` as its additional allowed value.

This avoids globally accepting HTML for commands whose dispatch switches would otherwise fall through to text. It also gives issue #261 a small, concrete extension point without introducing command-specific validator functions.

### D5: Dispatch HTML explicitly

Add an explicit `case "html"` in `runAnalyze` that calls `report.WriteHTML`. The default branch remains the text path after validation, preserving current behavior. Interactive analyze behavior remains unchanged because interactive mode takes precedence over formatter dispatch today.

### D6: Keep HTML deterministic and self-contained

The template includes no current time, generated IDs, environment-derived values, or remote resources. Input function and effect ordering is preserved. Any derived detail representation is stable. Repeated calls with identical ordered inputs and a version must therefore produce byte-identical output.

### D7: Update user-facing documentation

Update the analyze CLI reference and concise README examples to list HTML. The implementation task must also prepare the full title and body for the website documentation issue required by the repository gate; active workspace policy prohibits the agent from creating that GitHub issue directly.

## Coverage Strategy

- Unit-test `WriteHTML` with synthetic `taxonomy.AnalysisResult` values for complete structure, every function/effect, classifications, detail metadata, and empty input.
- Add adversarial escaping cases for function names, paths, descriptions, classification excerpts, and detail values containing markup and attribute-breaking characters.
- Parse or inspect output to prove the absence of external resource elements and URLs and the presence of native collapsible sections.
- Render identical input twice and compare complete byte sequences.
- Extend `ValidateFormat` tests to prove optional HTML acceptance and unchanged rejection when HTML is not supplied.
- Add a `runAnalyze` integration test against an existing lightweight fixture to prove CLI format selection and complete HTML output.
- Retain existing text and JSON tests as non-regression coverage.
- Run the exact CI-equivalent build, race-enabled test suites, and lint commands derived from `.github/workflows/` before implementation is declared complete.

## Risks / Trade-offs

### R1: Shared format validation could enable unsupported output

If HTML were added to the unconditional allowlist, quality and CRAP commands could accept it and silently emit text. The optional-allowlist design prevents this, and tests will cover both accepted and rejected paths.

### R2: Embedded CSS increases binary and output size

Every report contains the same CSS. This is an accepted trade-off for a portable single-file artifact; the stylesheet will remain intentionally small and dependency-free.

### R3: html/template can be bypassed with trusted content types

The implementation will not use trusted raw template types for report data. Adversarial escaping tests provide the regression gate.

### R4: Arbitrary detail values may not serialize

External analyzer detail normally originates in JSON and is serializable, but in-process callers can construct unsupported Go values. The formatter will return a wrapped error rather than omit or partially render the detail.

### R5: HTML visual quality is difficult to assert exhaustively

Tests will verify semantic structure, accessibility through native elements, self-containment, and deterministic content rather than brittle pixel-level styling. Manual browser inspection may supplement but will not replace automated acceptance tests.
<!-- scaffolded by uf vdev -->
