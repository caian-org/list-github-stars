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

	starreds, _, err := client.Activity.ListStarred(ctx, login, opts)
	if err != nil {
		panic(err)
	}

	s := []starredRepository{}
	for _, starred := range starreds {
		s = append(s, starredRepository{
			owner:       *starred.Repository.Owner.Login,
			name:        *starred.Repository.Name,
			description: *starred.Repository.Description,
			stars:       *starred.Repository.StargazersCount,
		})
	}

	return s
}

func main() {
	ctx, oauthClient := getOAuthClient(getGitHubTokenFromEnv())
	client := github.NewClient(oauthClient)

	authenticatedUser := getAuthenticatedUser(ctx, client)
	fmt.Printf("authenticated as @%s (%s)\n", *authenticatedUser.Login, *authenticatedUser.Name)

	starred := getAuthenticatedUserStarredRepos(ctx, client, *authenticatedUser.Login, 1)
	for _, s := range starred {
		fmt.Println(s.name)
	}
}
