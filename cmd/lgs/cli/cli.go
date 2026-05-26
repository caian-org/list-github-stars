package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/caian-org/list-github-stars/internal/starlist"
)

const (
	flagToken = "token"
	flagUser  = "user"
)

// New returns the root CLI command for `lgs`.
func New() *cli.Command {
	return &cli.Command{
		Name:    "lgs",
		Usage:   "list a GitHub user's starred repositories as Markdown",
		Version: programVersion(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    flagToken,
				Aliases: []string{"t"},
				Usage:   "GitHub personal access token (read from $GITHUB_TOKEN by default)",
				Sources: cli.EnvVars("GITHUB_TOKEN"),
			},
			&cli.StringFlag{
				Name:    flagUser,
				Aliases: []string{"u"},
				Usage:   "GitHub username to list stars for (defaults to the authenticated user)",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			token := cmd.String(flagToken)
			if token == "" {
				return fmt.Errorf("missing GitHub token: pass --token or set GITHUB_TOKEN")
			}

			out := cmd.Writer
			if out == nil {
				out = os.Stdout
			}

			return starlist.Run(ctx, starlist.Options{
				Token: token,
				User:  cmd.String(flagUser),
				Out:   out,
			})
		},
	}
}
