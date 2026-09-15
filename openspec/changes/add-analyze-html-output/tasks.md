<!--
  [P] marks tasks eligible for parallel execution.
  Add [P] when a task: (a) touches different files from
  other [P] tasks in the group, (b) has no dependency
  on prior tasks in the group, (c) can safely execute
  without ordering constraints.
  Do NOT add [P] when tasks modify the same file —
  parallel workers will cause merge conflicts.
  Tasks without [P] run sequentially first, then [P]
  tasks run in parallel.
-->

## Execution Checklist

- [x] Step 0: Startup Cleanup
- [x] Step 1: Branch Safety Gate
- [x] Step 2: Resumability Detection
- [x] Step 3: Clarify (Step 1)
- [x] Step 4: Plan (Step 2)
- [x] Step 5: Tasks (Step 3)
- [x] Step 6: Spec Review (Step 4) -- iteration: 3/3
- [ ] Step 7: Implement (Step 5) -- phase: 0/5, batch: 0/1, workers: 0/2
- [ ] Step 8: Code Review (Step 6) -- iteration: 0/3
- [ ] Step 9: Retrospective (Step 7)
- [ ] Step 10: Demo (Step 8)

## 0. Pre-Implementation Gate

- [x] 0.1 Confirm `proposal.md`, `design.md`, `specs/analyze-html-output/spec.md`, and `tasks.md` are committed and pushed on `opsx/add-analyze-html-output` before modifying production or test code, as required by the repository spec commit gate.

## 1. HTML Formatter

- [ ] 1.1 Add a narrowly embedded analyze HTML template under `internal/report` containing a complete document, compact inline CSS, native `<details>`/`<summary>` function sections, an empty-result state, version provenance, and no external resources.
- [ ] 1.2 Implement `report.WriteHTML(io.Writer, []taxonomy.AnalysisResult, string) error` with a deterministic presentation view model, `dev` version fallback, stable JSON serialization of effect detail metadata, contextual `html/template` escaping, and wrapped preparation/execution errors.
- [ ] 1.3 Add report formatter tests named for the specification scenarios: complete document and version, empty results, all functions and effects, classification and detail metadata, native collapsible sections, no external resources, adversarial escaping, deterministic bytes, and unsupported detail serialization failure. Assert specific rendered values rather than presence alone where applicable.

## 2. Analyze CLI Integration

- [ ] 2.1 Extend `cliutil.ValidateFormat` with optional additional allowed formats while preserving the existing `text`/`json` behavior and error for all unchanged callers; update table-driven tests to prove `html` is accepted only when explicitly supplied.
- [ ] 2.2 Update `runAnalyze` to validate with the additional `html` format and dispatch an explicit HTML case to `report.WriteHTML`, leaving interactive, text, and JSON paths unchanged.
- [ ] 2.3 Add CLI integration coverage proving `gaze analyze --format=html` emits a complete report for an existing lightweight fixture and that an out-of-scope command such as `gaze quality --format=html` still returns an invalid-format error instead of falling through to text.

## 3. Documentation

- [ ] 3.1 [P] Update `README.md` command examples and output-format summary to include `gaze analyze --format=html` and describe the single-file, offline-safe report.
- [ ] 3.2 [P] Update `docs/reference/cli/analyze.md` with the `html` format value, behavior, constraints, and an example that writes the report to a file.
- [ ] 3.3 Assess `AGENTS.md` and exported GoDoc impact; update GoDoc for all new exported identifiers and record that no AGENTS convention or architecture update is required unless implementation introduces an unplanned pattern.
- [ ] 3.4 Prepare the title and complete GitHub-flavored Markdown body for the required `unbound-force/website` documentation issue covering the new analyze HTML workflow; the user must create it before the implementing PR is merged because active workspace policy prohibits agent-created GitHub issues.

## 4. Verification And Governance

- [ ] 4.1 Read the current `.github/workflows/` files and record the exact local CI-equivalent build, test, and lint commands rather than relying on a memorized command list.
- [ ] 4.2 Format changed Go files with `gofmt` and `goimports`, then run targeted race-enabled formatter, CLI utility, and command tests with `-count=1`.
- [ ] 4.3 Run every CI-equivalent check discovered in task 4.1, including both race-enabled test suites and linting; treat any failure as blocking and do not modify protected gate values.
- [ ] 4.4 Verify the final implementation against every scenario in `specs/analyze-html-output/spec.md`, including byte-for-byte text/JSON non-regression and rejection of HTML by unsupported commands.
- [ ] 4.5 Re-check the proposal's constitution alignment: the report remains a self-contained collaboration artifact, introduces no mandatory dependency, preserves machine-readable JSON provenance, and is testable without external services.
- [ ] 4.6 Run the required review council on the final implementation and resolve all REQUEST CHANGES findings before any PR submission; make minimal to no code changes after all four reviewers approve.
<!-- scaffolded by uf vdev -->
<!-- spec-review: passed -->
