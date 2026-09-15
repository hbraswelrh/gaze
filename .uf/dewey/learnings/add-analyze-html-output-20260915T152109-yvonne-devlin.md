---
tag: add-analyze-html-output
author: yvonne-devlin
category: gotcha
created_at: 2026-09-15T15:21:09Z
identity: add-analyze-html-output-20260915T152109-yvonne-devlin
tier: draft
---

On branch opsx/add-analyze-html-output (2026-09-15), safe deterministic HTML reporting used html/template with an embedded self-contained template and ordinary string view-model fields. Opaque SideEffect.Detail maps were JSON-marshaled before rendering, which both sorts string keys deterministically and allows contextual HTML escaping; unsupported detail values fail before template execution rather than being silently dropped. Adversarial tests should assert the exact escaped detail representation, not permit an unescaped fallback.
