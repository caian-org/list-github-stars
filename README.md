[![CI][gh-ci-shield]][gh-ci-url]
[![Release][gh-release-shield]][gh-release-url]
[![GitHub tag][tag-shield]][tag-url]

# `list-github-stars`

<img src="icon.svg" height="240px" align="right"/>

Tiny CLI program capable of retrieving a GitHub user's starred repositories
and formatting the data into a nice Markdown document. You can see an
example of the program's output
[here](https://gist.github.com/upsetbit/ada2117bd8c73a1e94e49580fd5c7cf7).

[gh-ci-shield]: https://img.shields.io/github/actions/workflow/status/caian-org/list-github-stars/ci.yml?label=ci&logo=github&style=flat-square
[gh-ci-url]: https://github.com/caian-org/list-github-stars/actions/workflows/ci.yml

[gh-release-shield]: https://img.shields.io/github/actions/workflow/status/caian-org/list-github-stars/release.yml?label=release&logo=github&style=flat-square
[gh-release-url]: https://github.com/caian-org/list-github-stars/actions/workflows/release.yml

[tag-shield]: https://img.shields.io/github/tag/caian-org/list-github-stars.svg?logo=git&logoColor=FFF&style=flat-square
[tag-url]: https://github.com/caian-org/list-github-stars/releases


### Usage

Authentication is made via [personal access tokens][pat]. Create a token,
export it to the environment variable `GITHUB_TOKEN`, and run. Output is
written to `STDOUT`:

```sh
./bin/lgs > my-stars.md
```

Flags:

- `--token`, `-t` — GitHub PAT (defaults to `$GITHUB_TOKEN`).
- `--user`, `-u` — list a different user's stars (defaults to the
  authenticated user).
- `--version`, `-v` — print version, commit, and build timestamp.

You can use external programs such as [`pandoc`][pandoc] to convert the
output to other formats:

```sh
# MS Word document
./bin/lgs | pandoc -o stars.docx

# HTML page
./bin/lgs | pandoc -o stars.html

# HTML page with custom stylesheet
./bin/lgs | pandoc -o stars.html --self-contained --css=style.css
```

[pat]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
[pandoc]: https://pandoc.org


### Run with Docker

Images are published to GitHub Container Registry:

```sh
docker run --rm \
    -e GITHUB_TOKEN="your-github-auth-token" \
    ghcr.io/caian-org/list-github-stars:latest > my-stars.md
```


### Build from source

This repo uses [`devbox`][devbox] + [`just`][just]. Inside `devbox shell`:

```sh
just build lgs              # produces ./bin/lgs
./bin/lgs --version
```

Common targets:

```sh
just test                   # go test ./...
just test-race              # go test ./... -race
just lint                   # go vet ./...
just release-check          # goreleaser check
just release-snapshot       # local goreleaser dry-run into dist/
```

[devbox]: https://www.jetify.com/devbox
[just]: https://github.com/casey/just


### Releases

Releases are tag-driven. Pushing a `v*` tag runs
`.github/workflows/release.yml`, which uses GoReleaser to publish:

- Multi-platform archives (`linux/darwin/windows × amd64/arm64`) attached
  to a GitHub Release.
- A Docker image to `ghcr.io/caian-org/list-github-stars:<tag>` and
  `:latest`.

A `daily-run` workflow rewrites a public gist with fresh output every day
at 03:00 UTC, consuming the `:latest` image.


## License

To the extent possible under law, [Caian Ertl][me] has waived __all
copyright and related or neighboring rights to this work__. In the spirit
of _freedom of information_, I encourage you to fork, modify, change,
share, or do whatever you like with this project! [`^C ^V`][kopimi]

[![License][cc-shield]][cc-url]

[me]: https://github.com/upsetbit
[cc-shield]: https://forthebadge.com/images/badges/cc-0.svg
[cc-url]: http://creativecommons.org/publicdomain/zero/1.0

[kopimi]: https://kopimi.com
