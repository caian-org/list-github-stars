.DEFAULT_GOAL := run

ARTIFACT = list-github-stars


build:
	@go build

release:
	@go build -ldflags "-w -s"
	@upx --best --lzma $(ARTIFACT)

run: build
	@./$(ARTIFACT)
