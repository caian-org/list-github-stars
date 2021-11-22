package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

type starredRepository struct {
	owner       string
	name        string
	description string
	language    string
	stars       int
}

func getGitHubTokenFromEnv() string {
	token_var := "GITHUB_AUTH_TOKEN"
	token := os.Getenv(token_var)
	if len(token) == 0 {
		panic(fmt.Sprintf("required environment variable %s is undefined", token_var))
	}

	return token
}

func getOAuthClient(token string) (context.Context, *http.Client) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	return ctx, oauth2.NewClient(ctx, ts)
}

func getAuthenticatedUser(ctx context.Context, client *github.Client) *github.User {
	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		panic(err)
	}

	return user
}

func getAuthenticatedUserStarredRepos(ctx context.Context, client *github.Client, login string, page int) []starredRepository {
	opts := &github.ActivityListStarredOptions{
		ListOptions: github.ListOptions{
			Page:    page,
			PerPage: 50,
		},
	}

	sts, _, err := client.Activity.ListStarred(ctx, login, opts)
	if err != nil {
		panic(err)
	}

	s := []starredRepository{}
	for _, starred := range sts {
		r := starred.GetRepository()
		if r == nil {
			continue
		}

		s = append(s, starredRepository{
			owner:       r.GetOwner().GetLogin(),
			name:        r.GetName(),
			description: r.GetDescription(),
			language:    r.GetLanguage(),
			stars:       r.GetStargazersCount(),
		})
	}

	return s
}

func main() {
	ctx, oauthClient := getOAuthClient(getGitHubTokenFromEnv())
	client := github.NewClient(oauthClient)

	authenticatedUser := getAuthenticatedUser(ctx, client)
	fmt.Printf("authenticated as @%s (%s)\n", *authenticatedUser.Login, *authenticatedUser.Name)

	page := 1
	starred := []starredRepository{}

	fmt.Println("fetching starred repositories")

	// loop until we got all the starred repos
	for {
		s := getAuthenticatedUserStarredRepos(ctx, client, *authenticatedUser.Login, page)
		if len(s) == 0 {
			break
		}

		starred = append(starred, s...)
		page++
	}

	fmt.Printf("got %d results\n", len(starred))

	file, err := os.Create("out.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println("writing file")

	defer file.Close()
	for _, s := range starred {
		file.WriteString(fmt.Sprintf("%s/%s (%s, %d stars)\n", s.owner, s.name, s.language, s.stars))
		file.WriteString(fmt.Sprintf("%s\n\n", s.description))
	}

	fmt.Println("done")
}
