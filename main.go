package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"unicode"

	"github.com/google/go-github/v43/github"
	"github.com/gosimple/slug"
	"golang.org/x/oauth2"
)

type starredRepository struct {
	owner       string
	name        string
	description string
	language    string
	stars       int
}

func toLower(v string) string {
	return strings.Map(unicode.ToLower, v)
}

func getGitHubTokenFromEnv() string {
	token_var := "GITHUB_TOKEN"
	token := os.Getenv(token_var)
	if len(token) == 0 {
		panic(fmt.Sprintf("required environment variable %s is undefined", token_var))
	}

	return token
}

func getOAuthClient(token string) (*context.Context, *http.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	return &ctx, oauth2.NewClient(ctx, ts)
}

func getAuthenticatedUser(ctx *context.Context, client *github.Client) *github.User {
	user, _, err := client.Users.Get(*ctx, "")
	if err != nil {
		panic(err)
	}

	return user
}

func getAuthenticatedUserStarredRepos(ctx *context.Context, client *github.Client, login string, page int) []starredRepository {
	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: 50,
		},
	}

	sts, _, err := client.Activity.ListStarred(*ctx, login, opts)
	if err != nil {
		panic(err)
	}

	s := []starredRepository{}
	for _, starred := range sts {
		r := starred.GetRepository()
		if r == nil {
			continue
		}

		lang := r.GetLanguage()
		if len(lang) == 0 {
			lang = "(NA)"
		}

		s = append(s, starredRepository{
			owner:       r.GetOwner().GetLogin(),
			name:        r.GetName(),
			description: r.GetDescription(),
			language:    lang,
			stars:       r.GetStargazersCount(),
		})
	}

	return s
}

func organizeReposByLanguage(starred *[]starredRepository) ([]string, map[string][]starredRepository) {
	organizedRepos := make(map[string][]starredRepository)

	for _, s := range *starred {
		_, found := organizedRepos[s.language]
		if !found {
			organizedRepos[s.language] = []starredRepository{}
		}

		organizedRepos[s.language] = append(organizedRepos[s.language], s)
	}

	langs := []string{}
	for k, v := range organizedRepos {
		langs = append(langs, k)
		sort.Slice(organizedRepos[k], func(i, j int) bool { return toLower(v[i].name) < toLower(v[j].name) })
	}

	sort.Slice(langs, func(i, j int) bool { return toLower(langs[i]) < toLower(langs[j]) })

	return langs, organizedRepos
}

func main() {
	ctx, oauthClient := getOAuthClient(getGitHubTokenFromEnv())
	client := github.NewClient(oauthClient)

	authenticatedUser := getAuthenticatedUser(ctx, client)

	page := 1
	starred := []starredRepository{}

	// loop until we got all the starred repos
	for {
		s := getAuthenticatedUserStarredRepos(ctx, client, *authenticatedUser.Login, page)
		if len(s) == 0 {
			break
		}

		starred = append(starred, s...)
		page++
	}

	langs, organizedRepos := organizeReposByLanguage(&starred)

	fmt.Printf("# GitHub Stars\n\n")
	fmt.Printf("Starred by [@%s (%s)](%s).\n\n\n", *authenticatedUser.Login, *authenticatedUser.Name, *authenticatedUser.URL)

	fmt.Printf("## Summary\n\n")
	for _, lang := range langs {
		fmt.Printf("  - [%s](#%s)\n", lang, slug.Make(lang))
	}

	fmt.Printf("\n")
	for _, lang := range langs {
		fmt.Printf("\n## %s\n\n", lang)
		fmt.Printf("[🔝 back to top](#summary)\n\n")

		for _, repo := range organizedRepos[lang] {
			repoAndOwner := fmt.Sprintf("%s/%s", repo.owner, repo.name)
			repoUrl := fmt.Sprintf("https://github.com/%s", repoAndOwner)

			fmt.Printf("### `%s`\n\n", repo.name)
			fmt.Printf("  - [%s](%s) | ⭐ %d \n", repoAndOwner, repoUrl, repo.stars)

			if len(repo.description) > 0 {
				fmt.Printf("  - %s\n", repo.description)
			}

			fmt.Printf("\n")
		}
	}
}
