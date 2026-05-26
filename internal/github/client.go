package github

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	gh "github.com/google/go-github/v76/github"
)

// Starred is a flat representation of a starred GitHub repository.
type Starred struct {
	Owner       string
	Name        string
	Description string
	Language    string
	Stars       int
}

// Client fetches starred repositories from the GitHub REST API.
type Client struct {
	api     *gh.Client
	user    string
	perPage int
}

// New returns a Client authenticated with the given personal access token.
// `user` selects the target account; empty means "the authenticated user".
func New(token, user string) *Client {
	return &Client{
		api:     gh.NewClient(nil).WithAuthToken(token),
		user:    user,
		perPage: 50,
	}
}

// SetBaseURL replaces the GitHub API base URL. Intended for tests that point
// at a httptest server. A trailing slash is added if missing.
func (c *Client) SetBaseURL(baseURL string) error {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}
	parsed, err := url.Parse(baseURL)
	if err != nil {
		return fmt.Errorf("github: invalid base URL: %w", err)
	}
	c.api.BaseURL = parsed
	c.api.UploadURL = parsed
	return nil
}

// AuthenticatedUser returns the login and (optional) display name of the
// account that owns the configured token.
func (c *Client) AuthenticatedUser(ctx context.Context) (login, name string, err error) {
	u, _, err := c.api.Users.Get(ctx, "")
	if err != nil {
		return "", "", fmt.Errorf("github: failed to get authenticated user: %w", err)
	}
	return u.GetLogin(), u.GetName(), nil
}

// FetchStarred paginates through every starred repository for the configured
// user. If `c.user` is empty the authenticated user is used.
func (c *Client) FetchStarred(ctx context.Context) ([]Starred, error) {
	var all []Starred

	for page := 1; ; page++ {
		opts := &gh.ActivityListStarredOptions{
			ListOptions: gh.ListOptions{Page: page, PerPage: c.perPage},
		}

		batch, _, err := c.api.Activity.ListStarred(ctx, c.user, opts)
		if err != nil {
			return nil, fmt.Errorf("github: failed to list starred (page %d): %w", page, err)
		}
		if len(batch) == 0 {
			break
		}

		for _, s := range batch {
			r := s.GetRepository()
			if r == nil {
				continue
			}

			lang := r.GetLanguage()
			if lang == "" {
				lang = "(NA)"
			}

			all = append(all, Starred{
				Owner:       r.GetOwner().GetLogin(),
				Name:        r.GetName(),
				Description: r.GetDescription(),
				Language:    lang,
				Stars:       r.GetStargazersCount(),
			})
		}
	}

	return all, nil
}
