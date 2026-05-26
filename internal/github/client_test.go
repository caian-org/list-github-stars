package github_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/caian-org/list-github-stars/internal/github"
)

// repoFixture renders a single starred-repo JSON entry as returned by the
// GitHub REST API when called with Accept: application/vnd.github.v3.star+json.
type repoFixture struct {
	ID          int
	Name        string
	Owner       string
	Language    string
	Description string
	Stars       int
}

func (r repoFixture) JSON() string {
	lang := "null"
	if r.Language != "" {
		lang = fmt.Sprintf("%q", r.Language)
	}
	return fmt.Sprintf(`{
  "starred_at": "2024-01-01T00:00:00Z",
  "repo": {
    "id": %d,
    "name": %q,
    "full_name": "%s/%s",
    "owner": {"login": %q},
    "language": %s,
    "description": %q,
    "stargazers_count": %d
  }
}`, r.ID, r.Name, r.Owner, r.Name, r.Owner, lang, r.Description, r.Stars)
}

func arrayJSON(repos []repoFixture) string {
	if len(repos) == 0 {
		return "[]"
	}
	parts := make([]string, 0, len(repos))
	for _, r := range repos {
		parts = append(parts, r.JSON())
	}
	return "[\n" + strings.Join(parts, ",\n") + "\n]"
}

func newTestClient(t *testing.T, handler http.HandlerFunc) *github.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	c := github.New("test-token", "")
	require.NoError(t, c.SetBaseURL(srv.URL))
	return c
}

func TestAuthenticatedUser_ReturnsLoginAndName(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/user", r.URL.Path)
		require.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
		fmt.Fprint(w, `{"login":"testuser","name":"Test User"}`)
	})

	login, name, err := c.AuthenticatedUser(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "testuser", login)
	assert.Equal(t, "Test User", name)
}

func TestAuthenticatedUser_OmitsNameWhenMissing(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"login":"justlogin"}`)
	})

	login, name, err := c.AuthenticatedUser(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "justlogin", login)
	assert.Empty(t, name)
}

func TestFetchStarred_SinglePage(t *testing.T) {
	repos := []repoFixture{
		{ID: 1, Name: "go", Owner: "golang", Language: "Go", Description: "The Go language", Stars: 100000},
		{ID: 2, Name: "claude", Owner: "anthropic", Language: "Python", Description: "Claude AI", Stars: 50000},
	}

	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/user/starred", r.URL.Path)
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 1 {
			fmt.Fprint(w, arrayJSON(repos))
			return
		}
		fmt.Fprint(w, "[]")
	})

	got, err := c.FetchStarred(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 2)

	assert.Equal(t, "golang", got[0].Owner)
	assert.Equal(t, "go", got[0].Name)
	assert.Equal(t, "Go", got[0].Language)
	assert.Equal(t, 100000, got[0].Stars)
	assert.Equal(t, "The Go language", got[0].Description)
}

func TestFetchStarred_PaginationTerminatesOnEmpty(t *testing.T) {
	pages := map[int][]repoFixture{
		1: makeFixtures(50, "go", "Go"),
		2: makeFixtures(5, "py", "Python"),
		3: nil,
	}

	hits := 0
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		hits++
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		fmt.Fprint(w, arrayJSON(pages[page]))
	})

	got, err := c.FetchStarred(context.Background())
	require.NoError(t, err)
	assert.Len(t, got, 55)
	assert.Equal(t, 3, hits, "should stop after the first empty page")
}

func TestFetchStarred_MissingLanguageFallsBackToNA(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 1 {
			fmt.Fprint(w, arrayJSON([]repoFixture{
				{ID: 99, Name: "untyped", Owner: "u", Language: "", Description: "", Stars: 0},
			}))
			return
		}
		fmt.Fprint(w, "[]")
	})

	got, err := c.FetchStarred(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "(NA)", got[0].Language)
}

func TestFetchStarred_PropagatesAPIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"message":"Bad credentials"}`, http.StatusUnauthorized)
	})

	_, err := c.FetchStarred(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "github: failed to list starred")
}

func makeFixtures(n int, prefix, language string) []repoFixture {
	out := make([]repoFixture, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, repoFixture{
			ID:       i + 1,
			Name:     fmt.Sprintf("%s-%d", prefix, i),
			Owner:    "u",
			Language: language,
			Stars:    i,
		})
	}
	return out
}
