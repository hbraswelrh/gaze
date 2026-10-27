## ADDED Requirements

### Requirement: Idempotent version marker insertion

`gaze init` MUST leave exactly one `<!-- scaffolded by gaze <version> -->`
marker in each scaffolded tool-owned file after any number of re-runs, and
the marker MUST carry the version of the most recent run.

#### Scenario: Re-scaffold with a different version (force)

- **GIVEN** a file scaffolded by `gaze init` with version `dev`
- **WHEN** `gaze init` is re-run with `--force` and version `v1.8.0`
- **THEN** the file MUST contain exactly one marker: `<!-- scaffolded by gaze v1.8.0 -->`

#### Scenario: Re-scaffold without --force

- **GIVEN** a scaffolded project whose tool-owned files carry a `dev` marker
- **WHEN** `gaze init` is re-run without `--force` with version `v1.8.0`
- **THEN** tool-owned files that differ MUST be updated with exactly one `v1.8.0` marker
- **AND** user-owned files MUST be left untouched

#### Scenario: Duplicate markers are collapsed

- **GIVEN** a file containing two or more `<!-- scaffolded by gaze ... -->` lines
- **WHEN** `insertMarkerAfterFrontmatter` is called
- **THEN** all existing marker lines MUST be removed before the new marker is inserted

#### Scenario: Non-marker HTML comments are preserved

- **GIVEN** a file containing HTML comments that do not match the marker syntax
- **WHEN** `insertMarkerAfterFrontmatter` is called
- **THEN** those comments MUST remain unchanged in the output

### Requirement: Embedded assets are lint-invariant mirrors

Files under `internal/scaffold/assets/` and `internal/aireport/assets/`
are byte-identical mirrors of their `.opencode/` sources and MUST NOT be
independently reformatted by lint tooling.

#### Scenario: markdownlint skips embedded assets

- **GIVEN** `.mega-linter.yml` excludes `internal/scaffold/assets/` and `internal/aireport/assets/` from `MARKDOWN_MARKDOWNLINT`
- **WHEN** MegaLinter runs over the repository
- **THEN** markdownlint MUST NOT modify any file under those directories
- **AND** `TestEmbeddedAssetsMatchSource` and `TestEmbeddedPromptMatchesScaffold` MUST pass