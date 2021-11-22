.DEFAULT_GOAL := run


build:
	@go build

run: build
	@./list-github-stars
