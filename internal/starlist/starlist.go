package starlist

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/caian-org/list-github-stars/internal/github"
	"github.com/caian-org/list-github-stars/internal/render"
)

// Options configures a Run invocation.
type Options struct {
	Token string
	User  string // empty means the authenticated user
	Out   io.Writer

	// BaseURL overrides the GitHub API endpoint. Used by tests.
	BaseURL string
}

// Run fetches starred repositories and writes a Markdown document to Out.
func Run(ctx context.Context, opts Options) error {
	if opts.Token == "" {
		return fmt.Errorf("starlist: missing token")
	}
	if opts.Out == nil {
		opts.Out = os.Stdout
	}

	client := github.New(opts.Token, opts.User)
	if opts.BaseURL != "" {
		if err := client.SetBaseURL(opts.BaseURL); err != nil {
			return err
		}
	}

	var user render.User
	if opts.User != "" {
		user.Login = opts.User
	} else {
		login, name, err := client.AuthenticatedUser(ctx)
		if err != nil {
			return err
		}
		user.Login = login
		user.Name = name
	}

	starred, err := client.FetchStarred(ctx)
	if err != nil {
		return err
	}

	if _, err := io.WriteString(opts.Out, render.Markdown(user, starred)); err != nil {
		return fmt.Errorf("starlist: write output: %w", err)
	}

	return nil
}
