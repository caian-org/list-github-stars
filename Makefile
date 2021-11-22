.DEFAULT_GOAL := run

ARTIFACT = list-github-stars

ifeq ($(OS),Windows_NT)
	ARTIFACT = list-github-stars.exe
endif


build:
	@go build

release:
	@go build -ldflags "-w -s"
	@upx --best --lzma $(ARTIFACT)

run: build
	@./$(ARTIFACT)
