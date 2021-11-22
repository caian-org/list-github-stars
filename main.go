package main

import (
	"context"
    "fmt"
    "net/http"
    "os"

	"github.com/google/go-github/v40/github"
	"golang.org/x/oauth2"
)

func getGitHubTokenFromEnv() string {
    token_var := "GITHUB_AUTH_TOKEN"
    token := os.Getenv(token_var)
    if (len(token) == 0) {
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

func main() {
    ctx, oauthClient := getOAuthClient(getGitHubTokenFromEnv())
    client := github.NewClient(oauthClient)

    user, _, err := client.Users.Get(ctx, "")
    if err != nil {
        panic(err)
    }

    login := *user.Login
    fmt.Printf("authenticated as \"%s\"\n", login)
}
