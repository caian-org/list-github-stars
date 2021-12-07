.DEFAULT_GOAL := build

ARTIFACT = list-github-stars
LDFLAGS = -w -s

ifeq ($(OS),Windows_NT)
	ARTIFACT = list-github-stars.exe
endif


build:
	@go build

release:
	@go build -ldflags "$(LDFLAGS)"
	@upx --best --lzma $(ARTIFACT)

run: build
	@./$(ARTIFACT)

format:
	gofmt -s -w .
