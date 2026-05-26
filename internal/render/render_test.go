package render_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/caian-org/list-github-stars/internal/github"
	"github.com/caian-org/list-github-stars/internal/render"
)

var updateGolden = flag.Bool("update", false, "regenerate golden Markdown fixtures")

func sampleStarred() []github.Starred {
	return []github.Starred{
		{Owner: "golang", Name: "go", Description: "The Go programming language", Language: "Go", Stars: 100000},
		{Owner: "anthropic", Name: "claude", Description: "Claude AI", Language: "Python", Stars: 50000},
		{Owner: "rust-lang", Name: "rust", Description: "Empowering everyone", Language: "Rust", Stars: 80000},
		{Owner: "example", Name: "no-lang", Description: "", Language: "(NA)", Stars: 10},
		{Owner: "zorg", Name: "zlib", Description: "A library", Language: "Python", Stars: 200},
	}
}

func TestMarkdown_Golden(t *testing.T) {
	user := render.User{Login: "testuser", Name: "Test User"}
	got := render.Markdown(user, sampleStarred())

	goldenPath := filepath.Join("..", "..", "testdata", "golden", "starred.md")

	if *updateGolden {
		require.NoError(t, os.MkdirAll(filepath.Dir(goldenPath), 0o755))
		require.NoError(t, os.WriteFile(goldenPath, []byte(got), 0o644))
	}

	want, err := os.ReadFile(goldenPath)
	require.NoError(t, err, "missing golden file; run with -update to regenerate")
	assert.Equal(t, string(want), got)
}

func TestMarkdown_OmitsDisplayNameWhenEmpty(t *testing.T) {
	user := render.User{Login: "testuser"}
	got := render.Markdown(user, sampleStarred())

	assert.Contains(t, got, "Starred by [@testuser](https://github.com/testuser).")
	assert.NotContains(t, got, "(Test User)")
}

func TestMarkdown_LanguageSortIsAlphabeticCaseInsensitive(t *testing.T) {
	starred := []github.Starred{
		{Owner: "u", Name: "a", Language: "Zoom", Stars: 1},
		{Owner: "u", Name: "b", Language: "apple", Stars: 1},
		{Owner: "u", Name: "c", Language: "Mango", Stars: 1},
	}

	got := render.Markdown(render.User{Login: "u"}, starred)

	apple := strings.Index(got, "## apple\n")
	mango := strings.Index(got, "## Mango\n")
	zoom := strings.Index(got, "## Zoom\n")
	require.NotEqual(t, -1, apple)
	require.NotEqual(t, -1, mango)
	require.NotEqual(t, -1, zoom)
	assert.True(t, apple < mango && mango < zoom, "expected languages sorted case-insensitive ascending")
}

func TestMarkdown_RepoNameOrderIsCaseInsensitive(t *testing.T) {
	starred := []github.Starred{
		{Owner: "u", Name: "Zeta", Language: "X", Stars: 1},
		{Owner: "u", Name: "apple", Language: "X", Stars: 1},
		{Owner: "u", Name: "Banana", Language: "X", Stars: 1},
	}

	got := render.Markdown(render.User{Login: "u"}, starred)

	apple := strings.Index(got, "### `apple`")
	banana := strings.Index(got, "### `Banana`")
	zeta := strings.Index(got, "### `Zeta`")
	require.True(t, apple < banana && banana < zeta, "expected case-insensitive repo ordering")
}

func TestMarkdown_HandlesMissingDescription(t *testing.T) {
	starred := []github.Starred{
		{Owner: "u", Name: "r", Language: "Go", Stars: 1, Description: ""},
	}

	got := render.Markdown(render.User{Login: "u"}, starred)

	assert.Contains(t, got, "### `r`")
	assert.Contains(t, got, "⭐ 1")
	// The description bullet must be omitted when empty.
	assert.NotContains(t, got, "  - \n")
}

func TestMarkdown_DisambiguatesCollidingLanguageSlugs(t *testing.T) {
	// "C", "C++" and "C#" all slugify to "c". The renderer must emit unique
	// anchors so the summary links stay valid.
	starred := []github.Starred{
		{Owner: "u", Name: "a", Language: "C", Stars: 1},
		{Owner: "u", Name: "b", Language: "C++", Stars: 1},
		{Owner: "u", Name: "d", Language: "C#", Stars: 1},
	}

	got := render.Markdown(render.User{Login: "u"}, starred)

	assert.Contains(t, got, "(#c)")
	assert.Contains(t, got, "(#c-1)")
	assert.Contains(t, got, "(#c-2)")
}
