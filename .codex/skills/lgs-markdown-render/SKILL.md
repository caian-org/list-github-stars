---
name: lgs-markdown-render
description: Markdown rendering guidance for list-github-stars. Use when changing the document layout, language grouping, slug disambiguation, sorting, or golden fixtures.
---

# LGS Markdown Render

Use this skill before touching `internal/render/` or `testdata/golden/`.

## Output Contract

`render.Markdown(user, starred)` returns a stable Markdown document with
this structure:

```
# GitHub Stars

Starred by [@<login>[ (<name>)]](https://github.com/<login>).


## Summary

  - [<lang>](#<slug>)
  ...

## <lang>

[🔝 back to top](#summary)

### `<repo>`

  - [<owner>/<repo>](https://github.com/<owner>/<repo>) | ⭐ <stars>
  - <description>     # bullet omitted when description is empty
```

The daily-run workflow rewrites a public gist with this output every day,
so the layout is effectively a public contract. Changing it requires
regenerating the golden fixture and is a visible change to consumers.

## Sort Order

- Languages are sorted case-insensitive ascending by raw name. Special
  characters such as `(` in the `(NA)` sentinel sort before letters by
  ASCII codepoint, which puts `(NA)` first whenever it is present.
- Repositories within a language are sorted case-insensitive ascending by
  name.
- Sorts use `sort.SliceStable` so iteration order of equal items is
  deterministic relative to insertion order. Insertion order itself comes
  from the input slice (`starred`), which is whatever
  `internal/github.Client.FetchStarred` produced — that, in turn, is
  page-then-API order.

## Slug Disambiguation

Anchor slugs come from `slug.Make(lang)`. Languages like `C`, `C++`, and
`C#` all slugify to `c`, so the renderer assigns suffixes on collision:

| Iteration | Slug   |
|-----------|--------|
| 1st       | `c`    |
| 2nd       | `c-1`  |
| 3rd       | `c-2`  |

The summary links and the section anchors must always be in sync — the
slug for a language is computed once per render via `assignLangSlugs`.

## Purity

`internal/render/render.go` is I/O-free:

- Returns a `string`. Never accepts a `Writer`.
- No logging, no globals, no environment access.
- Time and randomness are out of scope.

This makes it trivial to golden-test deterministically.

## Golden Workflow

The golden fixture lives at `testdata/golden/starred.md` (path is
`../../testdata/golden/starred.md` from the test file).

```bash
# Regenerate after an intentional layout change:
just update-goldens

# Verify the diff in the golden, then run the full suite:
just test-race
```

The renderer change and the regenerated golden must land in the same
commit — never split them.

## Common Pitfalls

- Reordering struct fields in `Starred` and forgetting to update
  `internal/render` — adding new fields requires also rendering them.
- Reverting golden bytes mechanically when the diff looks "ugly" — every
  byte change should be intentional and reviewed.
- Adding fields to the per-repo line without thinking about how it
  interacts with the trailing-space-and-newline format
  (`| ⭐ <stars> \n`). The trailing space is intentional in the historical
  output.
