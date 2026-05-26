package render

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/gosimple/slug"

	"github.com/caian-org/list-github-stars/internal/github"
)

// User identifies the GitHub user whose stars are being rendered.
type User struct {
	Login string
	Name  string
}

// Markdown renders the starred repositories grouped by language as a Markdown
// document. The output layout matches the historical program output.
func Markdown(user User, starred []github.Starred) string {
	langs, byLang := groupByLanguage(starred)

	var b strings.Builder

	accountURL := fmt.Sprintf("https://github.com/%s", user.Login)
	display := ""
	if user.Name != "" {
		display = fmt.Sprintf(" (%s)", user.Name)
	}

	fmt.Fprintf(&b, "# GitHub Stars\n\n")
	fmt.Fprintf(&b, "Starred by [@%s%s](%s).\n\n\n", user.Login, display, accountURL)
	fmt.Fprintf(&b, "## Summary\n\n")

	langSlugs := assignLangSlugs(langs)
	for _, lang := range langs {
		fmt.Fprintf(&b, "  - [%s](#%s)\n", lang, langSlugs[lang])
	}
	fmt.Fprintf(&b, "\n")

	for _, lang := range langs {
		fmt.Fprintf(&b, "\n## %s\n\n", lang)
		fmt.Fprintf(&b, "[🔝 back to top](#summary)\n\n")

		for _, repo := range byLang[lang] {
			owned := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
			fmt.Fprintf(&b, "### `%s`\n\n", repo.Name)
			fmt.Fprintf(&b, "  - [%s](https://github.com/%s) | ⭐ %d \n", owned, owned, repo.Stars)
			if repo.Description != "" {
				fmt.Fprintf(&b, "  - %s\n", repo.Description)
			}
			fmt.Fprintf(&b, "\n")
		}
	}

	return b.String()
}

func groupByLanguage(starred []github.Starred) ([]string, map[string][]github.Starred) {
	by := make(map[string][]github.Starred)

	for _, s := range starred {
		by[s.Language] = append(by[s.Language], s)
	}

	langs := make([]string, 0, len(by))
	for lang, repos := range by {
		langs = append(langs, lang)
		sort.SliceStable(repos, func(i, j int) bool {
			return lowercase(repos[i].Name) < lowercase(repos[j].Name)
		})
	}
	sort.SliceStable(langs, func(i, j int) bool {
		return lowercase(langs[i]) < lowercase(langs[j])
	})

	return langs, by
}

func assignLangSlugs(langs []string) map[string]string {
	out := make(map[string]string, len(langs))
	counts := map[string]int{}

	for _, lang := range langs {
		base := slug.Make(lang)
		idx := counts[base]
		counts[base] = idx + 1
		if idx == 0 {
			out[lang] = base
		} else {
			out[lang] = fmt.Sprintf("%s-%d", base, idx)
		}
	}

	return out
}

func lowercase(s string) string {
	return strings.Map(unicode.ToLower, s)
}
