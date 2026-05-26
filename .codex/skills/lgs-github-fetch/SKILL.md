---
name: lgs-github-fetch
description: GitHub API fetcher guidance for list-github-stars. Use when changing the paginated starred-repos client, token handling, error wrapping, or the httptest-based test harness.
---

# LGS GitHub Fetch

Use this skill before touching `internal/github/`.

## Client Shape

`internal/github` wraps `github.com/google/go-github/v76` with a minimal,
test-friendly surface:

- `New(token, user string) *Client` — authenticated client (token may not
  be empty). Empty `user` means "the authenticated user".
- `SetBaseURL(baseURL string) error` — points the underlying go-github
  client at a different host. The only legitimate caller is the test
  harness using `httptest.Server`. The helper appends a trailing slash
  if missing, because go-github's `BaseURL` requires one.
- `AuthenticatedUser(ctx) (login, name string, err error)` — returns the
  login of the token owner and an optional display name.
- `FetchStarred(ctx) ([]Starred, error)` — paginates `Activity.ListStarred`
  with `PerPage=50` until an empty page is returned.

## Wire Format

`Starred` flattens the relevant repository fields:

```go
type Starred struct {
    Owner       string
    Name        string
    Description string
    Language    string  // "(NA)" when the API returns null
    Stars       int
}
```

GitHub returns `null` for the language field on repositories without a
detected language. The fetcher normalises that to the sentinel string
`(NA)` so the renderer's grouping logic always has a key.

## Pagination Contract

The loop in `FetchStarred` terminates on the first empty page. Pages are
1-indexed, mirroring the GitHub API. Page size is hard-coded at 50 to
match the historical program behaviour.

If `FetchStarred` ever needs to support partial fetches (e.g. resume after
rate limit), do it in this package — keep the orchestration layer
(`internal/starlist`) unaware of pagination.

## Error Wrapping

Errors use a `github:` package prefix:

```go
fmt.Errorf("github: failed to list starred (page %d): %w", page, err)
fmt.Errorf("github: failed to get authenticated user: %w", err)
fmt.Errorf("github: invalid base URL: %w", err)
```

Keep the prefix when adding new error sites.

## Testing Pattern

Tests build a `httptest.Server` and configure the client through
`SetBaseURL`:

```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // route on r.URL.Path / r.URL.Query()
}))
defer srv.Close()

c := github.New("test-token", "")
require.NoError(t, c.SetBaseURL(srv.URL))
```

For starred fixtures, return the wire format used by the
`application/vnd.github.v3.star+json` Accept header — an array of
`{ "starred_at": "...", "repo": { ... } }` envelopes. See
`internal/github/client_test.go` for the helpers (`repoFixture`,
`arrayJSON`).

Coverage should include at minimum:

- Authenticated user lookup with and without `name`.
- A multi-page run that terminates on the first empty page.
- A repository with `language: null` mapping to `(NA)`.
- A non-2xx response surfacing as a wrapped `github:` error.

## Common Pitfalls

- Forgetting the trailing slash in `BaseURL` → go-github silently appends
  paths to the wrong directory. `SetBaseURL` already handles this; do not
  bypass it.
- Calling the real GitHub API in tests → flaky CI and rate-limit failures.
  Always use `httptest.Server`.
- Adding new repo fields without also extending the renderer → the data
  is fetched but never surfaces in the output.
