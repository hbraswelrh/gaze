## ADDED Requirements

### Requirement: Analyze HTML Format Selection

The `gaze analyze` command MUST accept `html` as a value for `--format` and MUST route successful analysis results to the HTML formatter. Commands outside this change's scope MUST continue to reject `html` until they implement their own HTML formatter.

#### Scenario: User selects HTML output
- **GIVEN** a package containing at least one analyzable function
- **WHEN** the user runs `gaze analyze --format=html <package>`
- **THEN** the command writes a complete HTML document to standard output
- **AND** the command exits successfully

#### Scenario: Another command selects HTML output
- **GIVEN** HTML support has not yet been implemented for `gaze quality`
- **WHEN** the user runs `gaze quality --format=html <package>`
- **THEN** the command returns an invalid-format error
- **AND** it MUST NOT silently emit text or JSON

### Requirement: Complete HTML Document

The HTML formatter MUST write a complete document containing a doctype, an `<html>` root, document metadata, an identifying title, and a body. The document MUST identify the Gaze version, defaulting to `dev` when no version is supplied.

#### Scenario: Formatter receives analysis results
- **GIVEN** one or more analysis results and a Gaze version
- **WHEN** the HTML formatter renders them
- **THEN** the output contains a complete HTML document
- **AND** the supplied version is represented in the document

#### Scenario: Formatter receives no results
- **GIVEN** an empty analysis result set
- **WHEN** the HTML formatter renders it directly
- **THEN** the output remains a valid complete HTML document
- **AND** the document clearly states that no functions were analyzed

### Requirement: Complete Analysis Representation

The HTML report MUST represent every input function and every side effect associated with each function. For each function, the report MUST show its identity and source location. For each side effect, the report MUST show its type, priority tier, description, and source location when those values are present. Classification and detail metadata MUST be shown when present.

#### Scenario: Multiple functions and effects are rendered
- **GIVEN** analysis results containing multiple functions with multiple side effects
- **WHEN** the HTML formatter renders the results
- **THEN** every function appears exactly once as a function section
- **AND** every side effect appears under its associated function
- **AND** no function or side effect is omitted

#### Scenario: Classified effect is rendered
- **GIVEN** a side effect with classification and detail metadata
- **WHEN** the HTML formatter renders the effect
- **THEN** the classification is visible with the effect
- **AND** the detail metadata is represented without interpretation or data loss

### Requirement: Native Collapsible Navigation

The HTML report MUST use native `<details>` and `<summary>` elements to make function sections collapsible. Core report navigation and content visibility MUST NOT depend on JavaScript.

#### Scenario: Function section is collapsible offline
- **GIVEN** a generated HTML report
- **WHEN** it is opened in a browser without network access or JavaScript
- **THEN** each function can be expanded or collapsed using native browser behavior
- **AND** all report content remains accessible

### Requirement: Self-Contained Output

The HTML report MUST contain all required styling inline and MUST NOT reference remote or runtime filesystem assets. It MUST NOT emit external `<link>` elements, `<script src>` elements, remote images, CSS imports, or other URLs that cause network resource loading.

#### Scenario: Report is inspected for external resources
- **GIVEN** a generated HTML report
- **WHEN** its elements and CSS are inspected
- **THEN** no external resource reference is present
- **AND** the report remains styled and usable as a single file

### Requirement: Context-Aware Escaping

The formatter MUST use Go's `html/template` package for HTML generation. All source-derived and analyzer-derived values MUST pass through contextual escaping and MUST NOT be converted to trusted `template.HTML`, `template.JS`, `template.CSS`, or equivalent raw content types.

#### Scenario: Source data contains HTML markup
- **GIVEN** function names, file paths, descriptions, or detail values containing `<script>`, quotes, ampersands, and HTML tags
- **WHEN** the HTML formatter renders the data
- **THEN** those values appear as inert escaped text
- **AND** no source-derived element or executable script is introduced into the document

### Requirement: Deterministic Rendering

The formatter MUST produce byte-identical output when called repeatedly with the same ordered input and version. It MUST NOT include timestamps, random identifiers, environment-dependent paths, or nondeterministic map iteration.

#### Scenario: Identical input is rendered twice
- **GIVEN** identical ordered analysis results and version values
- **WHEN** the formatter renders them in two independent calls
- **THEN** the resulting byte sequences are identical

### Requirement: Existing Format Compatibility

Adding HTML output MUST NOT alter the accepted behavior or rendered bytes of existing `text` and `json` output paths.

#### Scenario: Existing formats are selected
- **GIVEN** analysis input supported before this change
- **WHEN** the user selects `--format=text` or `--format=json`
- **THEN** the existing formatter and output contract are used unchanged

## MODIFIED Requirements

None.

## REMOVED Requirements

None.
<!-- scaffolded by uf vdev -->
