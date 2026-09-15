---
tag: add-analyze-html-output
author: yvonne-devlin
category: pattern
created_at: 2026-09-15T15:21:09Z
identity: add-analyze-html-output-20260915T152109-yvonne-devlin-2
tier: draft
---

On branch opsx/add-analyze-html-output (2026-09-15), command-scoped output format expansion was safely implemented by making cliutil.ValidateFormat accept optional extra formats. Only runAnalyze passes "html"; unchanged commands pass no extras and preserve their legacy text/json allowlist and error. This prevents unsupported commands from accepting a new format and silently falling through to their text default. Invalid-format guidance must dynamically include extras for the opted-in caller while remaining byte-compatible for unchanged callers.
