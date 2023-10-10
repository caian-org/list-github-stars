[![Build & Check][gh-bnc-shield]][gh-bnc-url]
[![Tests][gh-test-shield]][gh-test-url]
[![GitHub tag][tag-shield]][tag-url]

# list-github-stars

<img src="icon.svg" height="240px" align="right"/>

Tiny CLI program capable of retrieving a GitHub user's starred repositories and
formatting the data into a nice markdown document. You can see an example of
the program's output [here](https://gist.github.com/upsetbit/ada2117bd8c73a1e94e49580fd5c7cf7).

[gh-bnc-shield]: https://img.shields.io/github/actions/workflow/status/caian-org/list-github-stars/build-many.yml?label=build&logo=github&style=for-the-badge

[gh-bnc-url]: https://github.com/caian-org/list-github-stars/actions/workflows/build-many.yml

[gh-test-shield]: https://img.shields.io/github/actions/workflow/status/caian-org/list-github-stars/test-many.yml?label=test&logo=github&style=for-the-badge
[gh-test-url]: https://github.com/caian-org/list-github-stars/actions/workflows/test-many.yml

[tag-shield]: https://img.shields.io/github/tag/caian-org/list-github-stars.svg?logo=git&logoColor=FFF&style=for-the-badge
[tag-url]: https://github.com/caian-org/list-github-stars/releases


## Usage

Authentication is made via [personal access tokens][pat]. Create a token,
export to the environment variable `GITHUB_TOKEN` and run. The output is sent
to `STDOUT`, so to create a file, just do:

```sh
./list-github-stars >> my-stars.md
```

You can use external programs such as [`pandoc`][pandoc] to convert the output
to other formats. E.g.:

```sh
# MS Word document
./list-github-stars | pandoc -o stars.docx

# HTML page
./list-github-stars | pandoc -o stars.html

# HTML page with custom stylesheet
./list-github-stars | pandoc -o stars.html --self-contained --css=style.css
```

[pat]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
[pandoc]: https://pandoc.org

### Run with Docker

```sh
docker run --rm -e GITHUB_TOKEN="your-github-auth-token" caian/list-github-stars >> stars.md
```


## License

To the extent possible under law, [Caian Ertl][me] has waived __all copyright
and related or neighboring rights to this work__. In the spirit of _freedom of
information_, I encourage you to fork, modify, change, share, or do whatever
you like with this project! [`^C ^V`][kopimi]

[![License][cc-shield]][cc-url]

[me]: https://github.com/upsetbit
[cc-shield]: https://forthebadge.com/images/badges/cc-0.svg
[cc-url]: http://creativecommons.org/publicdomain/zero/1.0

[kopimi]: https://kopimi.com
