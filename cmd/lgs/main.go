package main

import (
	"context"
	"fmt"
	"os"

	"github.com/caian-org/list-github-stars/cmd/lgs/cli"
)

func main() {
	if err := cli.New().Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
