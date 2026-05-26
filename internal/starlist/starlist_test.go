package starlist_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/caian-org/list-github-stars/internal/starlist"
)

func TestRun_AuthenticatedUser_WritesMarkdown(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/user":
			fmt.Fprint(w, `{"login":"testuser","name":"Test User"}`)
		case r.URL.Path == "/user/starred" && r.URL.Query().Get("page") == "1":
			fmt.Fprint(w, `[
                {"starred_at":"2024-01-01T00:00:00Z","repo":{"id":1,"name":"go","full_name":"golang/go","owner":{"login":"golang"},"language":"Go","description":"The Go language","stargazers_count":100000}}
            ]`)
		default:
			fmt.Fprint(w, "[]")
		}
	}))
	t.Cleanup(srv.Close)

	var buf bytes.Buffer
	err := starlist.Run(context.Background(), starlist.Options{
		Token:   "test-token",
		Out:     &buf,
		BaseURL: srv.URL,
	})
	require.NoError(t, err)

	got := buf.String()
	assert.True(t, strings.HasPrefix(got, "# GitHub Stars\n\n"))
	assert.Contains(t, got, "Starred by [@testuser (Test User)](https://github.com/testuser)")
	assert.Contains(t, got, "## Go")
	assert.Contains(t, got, "### `go`")
	assert.Contains(t, got, "  - [golang/go](https://github.com/golang/go) | ⭐ 100000")
}

func TestRun_ExplicitUser_SkipsAuthenticatedLookup(t *testing.T) {
	authenticatedHit := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/user":
			authenticatedHit = true
			fmt.Fprint(w, `{"login":"shouldnotsee"}`)
		case "/users/explicit/starred":
			fmt.Fprint(w, "[]")
		default:
			fmt.Fprint(w, "[]")
		}
	}))
	t.Cleanup(srv.Close)

	var buf bytes.Buffer
	err := starlist.Run(context.Background(), starlist.Options{
		Token:   "test-token",
		User:    "explicit",
		Out:     &buf,
		BaseURL: srv.URL,
	})
	require.NoError(t, err)
	assert.False(t, authenticatedHit, "explicit --user must skip the /user lookup")
	assert.Contains(t, buf.String(), "Starred by [@explicit](https://github.com/explicit)")
}

func TestRun_RequiresToken(t *testing.T) {
	err := starlist.Run(context.Background(), starlist.Options{Out: &bytes.Buffer{}})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing token")
}
